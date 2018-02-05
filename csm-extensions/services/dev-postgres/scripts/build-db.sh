#!/bin/bash
topdir=$(dirname $(dirname "$0"))
docker build \
    --tag ${SIDECAR_EXTENSION_SVC_IMAGE_NAME}:${SIDECAR_EXTENSION_SVC_IMAGE_TAG} \
    --rm \
    --file  "${topdir}/Dockerfile-db" .
