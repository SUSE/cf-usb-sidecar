#!/bin/bash

mkdir -p ${CSM_EXTENSION_ROOT}/db/image/
docker pull redis:${CSM_EXTENSION_SVC_VERSION}
docker save -o ${CSM_EXTENSION_ROOT}/db/image/redis.tgz redis:${CSM_EXTENSION_SVC_VERSION}
docker rmi -f redis:${CSM_EXTENSION_SVC_VERSION}
pushd ${CSM_EXTENSION_ROOT}/db
docker build -t ${CSM_EXTENSION_SVC_IMAGE_NAME}:${CSM_EXTENSION_SVC_IMAGE_TAG} .
popd
rm -rf ${CSM_EXTENSION_ROOT}/db/image/
