#!/bin/bash

mkdir -p ${SIDECAR_EXTENSION_ROOT}/db/image/
docker pull redis:${SIDECAR_EXTENSION_SVC_VERSION}
docker tag redis:${SIDECAR_EXTENSION_SVC_VERSION} redis:hcf
docker rmi -f redis:${SIDECAR_EXTENSION_SVC_VERSION}
docker save -o ${SIDECAR_EXTENSION_ROOT}/db/image/redis.tgz redis:hcf
docker rmi -f redis:hcf
pushd ${SIDECAR_EXTENSION_ROOT}/db
docker build -t ${SIDECAR_EXTENSION_SVC_IMAGE_NAME}:${SIDECAR_EXTENSION_SVC_IMAGE_TAG} .
popd
rm -rf ${SIDECAR_EXTENSION_ROOT}/db/image/
