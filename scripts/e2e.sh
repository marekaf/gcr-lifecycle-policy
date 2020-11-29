#!/usr/bin/env bash

set -eEuo pipefail

if [ -z "${REGISTRY_URL:-}" ]; then
  echo "✋ Missing REGISTRY_URL!" >&2
  exit 1
fi

echo "prepopulating GCR with some images"

docker pull nginx:latest
docker tag nginx:latest "$REGISTRY_URL"/nginx:latest
docker push "$REGISTRY_URL"/nginx:latest

docker tag nginx:latest "$REGISTRY_URL"/nginx:1
docker push "$REGISTRY_URL"/nginx:1

docker tag nginx:latest "$REGISTRY_URL"/nginx:master
docker push "$REGISTRY_URL"/nginx:master

docker tag nginx:latest "$REGISTRY_URL"/webapp:latest
docker push "$REGISTRY_URL"/webapp:latest

docker tag nginx:latest "$REGISTRY_URL"/webapp:dev
docker push "$REGISTRY_URL"/webapp:dev

docker tag nginx:latest "$REGISTRY_URL"/webapp:v1.0
docker push "$REGISTRY_URL"/webapp:v1.0

docker pull nginx:alpine
docker tag nginx:alpine "$REGISTRY_URL"/webapp:feature-improve-tracking

repush() {
  docker pull "$1"
  docker tag "$1" "${REGISTRY_URL}/${1}"
  docker push "${REGISTRY_URL}/${1}"
}

repush mysql:latest
repush mysql:5.7
repush mysql:5.6
repush redis:latest
repush redis:5.0
repush redis:6.0
repush traefik
repush postgres
repush busybox
repush mariadb
repush mariadb:beta

make build && ./bin/gcr list --log-level INFO

./bin/gcr cleanup --log-level DEBUG --retention 0 --keep-tags=0 --dry-run=true

./bin/gcr cleanup --log-level DEBUG --retention 0 --keep-tags=0 --dry-run=false

./bin/gcr list --log-level INFO
