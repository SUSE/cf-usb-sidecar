#!/bin/bash

mkdir -p ${CSM_EXTENSION_ROOT}/db/image/
docker pull rabbitmq:${CSM_EXTENSION_SVC_VERSION}
docker save -o ${CSM_EXTENSION_ROOT}/db/image/rabbit.tgz rabbitmq:${CSM_EXTENSION_SVC_VERSION}
docker rmi -f rabbitmq:${CSM_EXTENSION_SVC_VERSION}
pushd ${CSM_EXTENSION_ROOT}/db
docker build -t ${CSM_EXTENSION_SVC_IMAGE_NAME}:${CSM_EXTENSION_SVC_IMAGE_TAG} .
popd
rm -rf ${CSM_EXTENSION_ROOT}/db/image/
