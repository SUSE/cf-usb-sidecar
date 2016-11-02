#!/bin/bash

mkdir -p ${SIDECAR_EXTENSION_ROOT}/db/image/
docker pull rabbitmq:${SIDECAR_EXTENSION_SVC_VERSION}
docker tag rabbitmq:${SIDECAR_EXTENSION_SVC_VERSION} rabbitmq:hsm
docker rmi -f rabbitmq:${SIDECAR_EXTENSION_SVC_VERSION}
docker save -o ${SIDECAR_EXTENSION_ROOT}/db/image/rabbit.tgz rabbitmq:hsm
docker rmi -f rabbitmq:hsm
pushd ${SIDECAR_EXTENSION_ROOT}/db
docker build -t ${SIDECAR_EXTENSION_SVC_IMAGE_NAME}:${SIDECAR_EXTENSION_SVC_IMAGE_TAG} .
popd
rm -rf ${SIDECAR_EXTENSION_ROOT}/db/image/
