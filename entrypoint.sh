#!/bin/bash
set -euo pipefail

# todo: rewrite this to python (pykube-ng, docker-py)

# auth
gcloud config set project "$PROJECT_ID"
[[ "$GCLOUD_SERVICE_KEY" != "" ]] && ( echo "$GCLOUD_SERVICE_KEY" | base64 -d > /gcloud-service-key.json )
gcloud auth activate-service-account --key-file /gcloud-service-key.json
echo "successfully authenticated"

# fetch the GKE's kubeconfig
gcloud container clusters get-credentials "$CLUSTER_NAME" --zone "$ZONE" --project "$PROJECT_ID"

# arg handling
KEEP=${KEEP_TAGS:-10}
RETENTION=${RETENTION_DAYS:-365}

## convert RETENTION_DAYS to milliseconds (that is what docker api is using in "timeCreatedMs")
STAMP=$(date --date="$RETENTION days ago" +%s%3N)
echo "stamp=$STAMP"

## auth
ACCESS_TOKEN=$(gcloud auth print-access-token)

## use docker v2 api directly
DOCKER_REGISTRY_PROTO='https://'
DOCKER_REGISTRY_HOST='eu.gcr.io'

## fetch all images from the registry
### I'm ignoring pagination as there are not so many images in _catalog
curl --silent --show-error -u_token:"$ACCESS_TOKEN" -X GET "${DOCKER_REGISTRY_PROTO}${DOCKER_REGISTRY_HOST}/v2/_catalog" |
  # format it and get the raw output we need (--raw-output)
  jq --raw-output '.[][]' |
  # make sure we only work with our GCR
  grep "^${PROJECT_ID}" > images.txt

## fetch in-use images:tags from cluster
### make sure to list all pods and replicasets (even the old ones that are scaled to zero, because we may want to rollback to them)
kubectl get rs,po --all-namespaces -o jsonpath={..image} |
  # work with it line-by-line
  tr ' ' '\n' |
  # make sure we work only with the ones that are in GCR and not other docker registries
  grep "${DOCKER_REGISTRY_HOST}/${PROJECT_ID}" |
  # sort them and get rid of the duplicities
  sort -u >in-use.txt

### separate "image:tag" into "image tag"
tr ':' ' ' > in-use-spaces.txt <in-use.txt

## "normalize" the images and tags per line to one line per image with all its tags following
## image1 tag1  ->  image1 tag1 tag2
## image1 tag2
awk -F' ' -v OFS=' ' '{x=$1;$1="";a[x]=a[x]$0}END{for(x in a)print x,a[x]}' > in-use-merged.txt <in-use-spaces.txt

## make a json object out of it
jq -Rn '''

  {
    "usedtags": [inputs | split("\\s+"; "g") |
       select(length > 0 and .[0] != "") |
       {(.[0]): .[1:]} ]  |
       add
  }
''' >used-tags.json <in-use-merged.txt

echo "--------------------------------"

## loop through all "images" in GCR
grep -v '^ *#' < images.txt | while IFS= read -r image
do

  # debug only
  #echo "image=$image"

  # fetch all tags for this particular image
  # TODO: solve pagination?
  # https://docs.docker.com/registry/spec/api/
  curl --silent --show-error -u_token:"$ACCESS_TOKEN" -X GET  "https://eu.gcr.io/v2/$image/tags/list" 2>/dev/null |
    # and dump the whole object formatted to a temp file (this file is each loop overwritten)
    jq '.' > tmp.json

  #echo "DEBUG:"
  #head tmp.json

  # filter out only the image we want and don't write "null" ( // empty) to the output file if the image is not used in the cluster
  jq '.usedtags."'eu.gcr.io/"$image"'"  // empty ' > used-tags-tmp.json <used-tags.json

  #echo "DEBUG:"
  #head used-tags-tmp.json

  ### ALL TOGETHER
  jq --sort-keys '.' tmp.json |
    # make the digest value (object's name) part of the object and reduce the data to the timestamp and tags array, delete everything else
    jq '[.manifest | to_entries[] | { "digest":.key, "tag":.value.tag, "timeCreatedMs": .value.timeCreatedMs}] ' |
    # add checkMostRecentTagsPassed="false" to all of the objects, later we will check if they passed this check and change to "true"
    jq '.[].checkMostRecentTagsPassed="false"' |
    # sort the digests historically
    jq '. | sort_by(.timeCreatedMs | tonumber)' |
    # reverse the array (to have the newest first == [0]), and the first N mark to pass the check, reverse the sorting
    jq 'reverse  | .[0:'"${KEEP}"'][].checkMostRecentTagsPassed="true" | reverse ' |
    # add checkInClusterUsePassed="false" to all of the objects, later we will check if they passed this check and change to "true"
    jq '.[].checkInClusterUsePassed="false"' |
    # load the tags used in cluster to a json array and
    # check all the digests and their tags, if any of digest's tag array contain a tag that is used in the cluster, change the check bool
    # the IN() jq function is available in jq version >1.5. That is why is jq installed from binary and not via apt.
    jq --slurpfile usedTags used-tags-tmp.json '[.[] |  select(any(.tag[] ; . |  IN($usedTags[][]))).checkInClusterUsePassed="true"]' |
    # add checkDatePassed="false" to all of the objects, later we will check if they passed this check and change to "true"
    jq '.[].checkDatePassed="true" ' |
    # select the digests that are older than our RETENTION_DAYS timestamp
    jq '[.[] | select(.timeCreatedMs < "'"$STAMP"'").checkDatePassed="false"]' |
    # select the ones that failed all three checks, those are the candidates to be deleted
    jq '.[] | select(.checkDatePassed == "true" or .checkInClusterUsePassed == "true" or .checkMostRecentTagsPassed == "true" | not)' \
    >final.json

  TO_DELETE=$(jq --raw-output '.digest' final.json) #| head
  if [[ "$TO_DELETE" != "" ]]
  then
    echo "image=$image:"
    echo "$( wc -l <used-tags-tmp.json ) tags found in cluster, $(echo "$TO_DELETE" | wc -l) digests to delete:"
    echo "$TO_DELETE" | head || true
    echo "... and others"
    # TOOD: how many failed the check Date?

    # TODO: now this is only a dry-run, later you can delete these digests
    # foreach digest in TO_DELETE
    # gcloud container images delete -q --force-delete-tags "${IMAGE}@${digest}"
    # + gcloud container images delete -q --force-delete-tags "eu.gcr.io/my-project/foo/bar@sha256:296e2378f7a14695b2f53101a3bd443f656f823c46d13bf6406b91e9e9950ef0"

  else
    echo "image=$image:"
    echo "$( wc -l < used-tags-tmp.json) tags found in cluster"
    echo "no digests to delete"
  fi

  echo ""

done
