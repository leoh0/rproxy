#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

git_commit="$(git describe --tags --always --dirty)"
build_date="$(date -u '+%Y%m%d')"
docker_tag="v${build_date}-${git_commit}"
# TODO(fejta): retire STABLE_PROW_REPO
cat <<EOF
STABLE_DOCKER_REPO ${DOCKER_REPO_OVERRIDE:-docker.io/leoh0}
CLUSTER ${CLUSTER_OVERRIDE:-al-cluster}
CONTEXT ${CONTEXT_OVERRIDE:-al-context}
NAMESPACE ${NAMESPACE_OVERRIDE:-rproxy}
DOCKER_TAG ${docker_tag}
EOF
