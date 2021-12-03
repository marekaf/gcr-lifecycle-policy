# GCR Retention policy
This simple golang CLI tool is handling retention policy of images in Google Container Registry (and Google Artifact Registry for Docker) in a similar manner as Elastic Container Registry does (original [blogpost talking about the details here](https://blog.marekbartik.com/posts/2019-09-26_google-container-registry-lifecycle-policy-for-images-retention/))

It scans all the Docker registry images (it supports paths like `eu.gcr.io/my-project/foo/bar/my-service:123`) and fetches their tags using Docker v2 API. Then it deletes some of them based on parametrized filters and input GKE clusters.

![](/assets/render.gif)

## Retention policy configuration

Using `--cluster path/to/kubeconfig` it queries all pods and replicasets in the GKE cluster and prevents these image tags before deleting.

```shell
$ gcr cleanup --help
Usage:
  gcr cleanup [flags]

Flags:
      --dry-run         dry-run for images cleaning (default true)
  -h, --help            help for cleanup
      --keep-tags int   number of tags to keep per image (default 10)
      --retention int   number of days of retention to keep images (default 365)

Global Flags:
      --cluster string      kubeconfig path (default "/Users/marek/.kube/config")
      --creds string        credential file (default "./creds/serviceaccount.json")
      --log-level string    log level (default "ERROR")
      --registry string     GCR url to use (default "eu.gcr.io")
      --repos stringArray   list of repos you want to work with
      --sort-by string      field to sort images by (default "timeCreatedMs")
```

`TIP`: For Google Artifact Registry use `--sort-by="timeUploadedMs"`

`WARNING`: This is an alpha version! Be very careful using this in production.

## How to build it
```shell
make dep
make build
```

## How to run it

### Provision GCP IAM service account
`serviceaccount.json` is a json key of IAM service account called `gcr-retention-policy` that has these roles

```
Kubernetes Engine Cluster Viewer
Kubernetes Engine Developer
Kubernetes Engine Service Agent
Storage Object Admin
Artifact Registry Repository Administrator
```

### Run it

```shell
make build && ./bin/gcr --help
```

## How to Build and Run using Docker
* get a copy of the Dockerfile (no need to clone the full repo)
* build the docker in the usual way
    ```shell
    docker build -t gcr:latest ./
    ```
    _note: if you are using macOS with the new M1 cpu(arm64), you need to force the
    usage of amd64 images, as some dependencies of `gcr` are not yet available for
    arm64_
    ```shell
    docker buildx build --platform linux/amd64 -t gcr:latest ./
    ```

* create a IAM service account and get a json key for it (see
  [Provision GCP IAM service account][0])

[0]: #provision-gcp-iam-service-account

* create `kubeconfig.yaml` as per [this post from ahmet][1] on gcp headless authentication  
  on the below code snippet, replace `[CLUSTER]` and `[ZONE]` with the
  appropriate values from your GKE cluster
    ```shell
    GET_CMD="gcloud container clusters describe [CLUSTER] --zone=[ZONE]"
    cat > kubeconfig.yaml <<EOF
    apiVersion: v1
    kind: Config
    current-context: my-cluster
    contexts: [{name: my-cluster, context: {cluster: cluster-1, user: user-1}}]
    users: [{name: user-1, user: {auth-provider: {name: gcp}}}]
    clusters:
    - name: cluster-1
      cluster:
        server: "https://$(eval "$GET_CMD --format='value(endpoint)'")"
        certificate-authority-data: "$(eval "$GET_CMD
    --format='value(masterAuth.clusterCaCertificate)'")"
    EOF
    ```

[1]: https://ahmet.im/blog/authenticating-to-gke-without-gcloud/

* place both `serviceaccount.json` and `kubeconfig.yaml` in `~/.gcr`

* use `gcr` from your CLI
    ```shell
    docker run -it --rm \
      -v $HOME/.gcr/serviceaccount.json:/root/creds/serviceaccount.json \
      -v $HOME/.gcr/kubeconfig.yaml:/root/.kube/config \
      gcr <gcr commands and flags>
    ```

* or, alternatively, get a shell inside the gcr docker and run `gcr` commands
  from there
    ```shell
    docker run -it --rm \
      -v $HOME/.gcr/serviceaccount.json:/root/creds/serviceaccount.json \
      -v $HOME/.gcr/kubeconfig.yaml:/root/.kube/config \
      --entrypoint /bin/bash \
      gcr
    ```

## `gcr cleanup` without kubeconfig

If you want to nuke old images without checking if they're currently used in a cluster you can do so by setting `--cluster=""`. When using this "dummy mode" your IAM service account doesn't need any `Kubernetes Engine` roles listed above.

## TODO:
- release binary to Github Release
- release docker image to Docker Hub
- try to make use of `https://github.com/google/go-containerregistry/tree/master`
- support Cloud Run in `--cluster` (not it supports only GKE)
- prepare example deployment for Cloud Scheduler + Cloud Run (or Cloud Functions)
- make the horrible code a bit less horrible and more reusable

# deprecation info
this tool used bash/jq before and the code is still kept [here](https://github.com/marekaf/gcr-lifecycle-policy/releases/tag/0.1). I will not continue maintaining the bash/jq version but focus purely in the new one in Go.
