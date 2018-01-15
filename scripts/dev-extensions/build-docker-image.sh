#!/bin/sh

docker build -t ${SIDECAR_EXTENSION_IMAGE_NAME}:${SIDECAR_EXTENSION_IMAGE_TAG} --rm -f Dockerfile       .
if test -n "${SIDECAR_SETUP_IMAGE_NAME:-}" ; then
    docker build -t ${SIDECAR_SETUP_IMAGE_NAME}:${SIDECAR_SETUP_IMAGE_TAG}     --rm -f Dockerfile-setup .
fi
