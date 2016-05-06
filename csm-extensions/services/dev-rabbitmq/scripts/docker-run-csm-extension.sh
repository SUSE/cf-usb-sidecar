#!/bin/sh

DOCKER_ENDPOINT="docker1"
DOCKER_IMAGE="rabbit"
DOCKER_IMAGE_VERSION="latest"
CSM_LOG_LEVEL="debug"
CSM_DEV_MODE="true"
CSM_API_KEY="csm-auth-token"



if [ ! -z ${DOCKER_HOST} ]
then
    export DOCKER_HOST_IP=`echo ${DOCKER_HOST} | cut -d "/" -f 3 | cut -d ":" -f 1`
else
    export DOCKER_HOST_IP=`ip route get 8.8.8.8 | awk 'NR==1 {print $NF}'`
fi


docker run --name ${CSM_EXTENSION_IMAGE_NAME} \
	-p 8094:8081 \
	-e DOCKER_ENDPOINT=${DOCKER_ENDPOINT} \
	-e DOCKER_IMAGE=${DOCKER_IMAGE} \
	-e DOCKER_IMAGE_VERSION=${DOCKER_IMAGE_VERSION} \
	-e CSM_LOG_LEVEL=${CSM_LOG_LEVEL} \
	-e CSM_API_KEY=${CSM_API_KEY} \
	-e CSM_DEV_MODE=${CSM_DEV_MODE} \
	-d ${CSM_EXTENSION_IMAGE_NAME}:${CSM_EXTENSION_IMAGE_TAG}
