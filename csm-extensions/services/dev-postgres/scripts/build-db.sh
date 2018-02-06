#!/bin/bash
# The SIDECAR variables are supplied by the Makefile.
docker build \
    --tag ${SIDECAR_EXTENSION_SVC_IMAGE_NAME}:${SIDECAR_EXTENSION_SVC_IMAGE_TAG} \
    --rm \
    --file  "${SIDECAR_EXTENSION_ROOT}/Dockerfile-db" .
