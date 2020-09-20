# GCR Retention policy
This `bash/jq` script packaged in Docker image is handling retention policy of images in Google Container Registry.

It scans all the GCR images (it supports paths like `eu.gcr.io/my-project/foo/bar/my-service:123`) and fetches their tags using Docker v2 API. Then it deletes some of them.

## Retention policy configuration

It queries all pods and replicasets in the GKE cluster and prevents these image tags before deleting.

It uses `RETENTION_DAYS` to delete images older than this value.

It uses `KEEP_TAGS` number to keep the `KEEP_TAGS` most recent tags

`TODO`: This is only a dry run. Check the `entrypoint.sh` to uncomment the section that does digest deletion. I'm planning on enabling and disabling dry-run via ENV vars but when I properly test this in production and get some confidence.

`WARNING`: This is an alpha version! Be very careful using this in production.

`TODO`: I want to rewrite this to python using `pykube-ng` and `docker-py`.

## How to build it
```
docker build -t gcr-retention:0.0.1 .
```

## How to run it

### provision GCP IAM service account
`serviceaccount.json` is a json key of IAM service account called `gcr-retention-policy` that has these roles

```
Kubernetes Engine Cluster Viewer
Kubernetes Engine Developer
Kubernetes Engine Service Agent
Storage Object Admin
```

### run it

```
docker run -it -e PROJECT_ID=my-project \
  -e GCLOUD_SERVICE_KEY="$(cat serviceaccount.json| base64)"  \
  -e KEEP_TAGS=10 -e RETENTION_DAYS=365 \
  -e REPOSITORY='eu.gcr.io/my-project' \
  -e ZONE="europe-west1-b" \
  -e CLUSTER_NAME="prod" \
    gcr-retention:0.0.1
```

### deploy it
I'm planning on using this in gitlab-ci scheduled pipeline - once a week to run this container with proper ENV vars, passing it the service account in base64 in gitlab secrets.


## TODO:
try to make use of `https://github.com/google/go-containerregistry/tree/master`
fix {"errors":[{"code":"DENIED","message":"Cloud Resource Manager API has not been used in project 50963927524 before or it is disabled. Enable it by visiting https://console.developers.google.com/apis/api/cloudresourcemanager.googleapis.com/overview?project=50963927524 then retry. If you enabled this API recently, wait a few minutes for the action to propagate to our systems and retry."}]}