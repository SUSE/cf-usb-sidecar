#!/bin/sh

docker build -t ${CSM_EXTENSION_IMAGE_NAME}:${CSM_EXTENSION_IMAGE_TAG} --rm -f Dockerfile .
