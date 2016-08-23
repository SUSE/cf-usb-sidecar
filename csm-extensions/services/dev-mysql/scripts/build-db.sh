#!/bin/bash

docker pull mysql:5.5
docker tag mysql:5.5 ${SIDECAR_EXTENSION_SVC_IMAGE_NAME}:${SIDECAR_EXTENSION_SVC_IMAGE_TAG}
