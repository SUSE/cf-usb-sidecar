#!/bin/bash

docker pull tutum/mongodb:3.0
docker tag tutum/mongodb:3.0 ${SIDECAR_EXTENSION_SVC_IMAGE_NAME}:${SIDECAR_EXTENSION_SVC_IMAGE_TAG}
