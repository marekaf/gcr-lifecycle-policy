# GCR Retention policy
This simple golang CLI tool is handling retention policy of images in Google Container Registry.

It scans all the GCR images (it supports paths like `eu.gcr.io/my-project/foo/bar/my-service:123`) and fetches their tags using Docker v2 API. Then it deletes some of them based on parametrized filters and input GKE clusters.

![](/assets/render.gif)

## Retention policy configuration

Using `--cluster path/to/kubeconfig` it queries all pods and replicasets in the GKE cluster and prevents these image tags before deleting.

```
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
```

`WARNING`: This is an alpha version! Be very careful using this in production.

## How to build it
```
make dep
make build
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
make build && ./bin/gcr --help
```

## TODO:
- release binary to Github Release
- release docker image to Docker Hub
- try to make use of `https://github.com/google/go-containerregistry/tree/master`
- support Cloud Run in `--cluster` (not it supports only GKE)
- prepare example deployment for Cloud Scheduler + Cloud Run (or Cloud Functions)
- make the horrible code a bit less horrible and more reusable
