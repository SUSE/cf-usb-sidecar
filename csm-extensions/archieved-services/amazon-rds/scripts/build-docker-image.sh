#!/bin/sh

docker build -t ${SIDECAR_EXTENSION_IMAGE_NAME}:${SIDECAR_EXTENSION_IMAGE_TAG} --rm -f Dockerfile .
