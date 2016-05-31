#!/bin/sh

MONGO_USER="admin"
MONGO_PASS="password"
MONGO_PORT="27017"
CSM_LOG_LEVEL="debug"
CSM_DEV_MODE="true"
CSM_API_KEY=${CSM_EXTENSION_TOKEN}

if [ ! -z ${DOCKER_HOST} ]
then
    export DOCKER_HOST_IP=`echo ${DOCKER_HOST} | cut -d "/" -f 3 | cut -d ":" -f 1`
else
    export DOCKER_HOST_IP=`ip route get 8.8.8.8 | awk 'NR==1 {print $NF}'`
fi


docker run --name ${CSM_EXTENSION_IMAGE_NAME} \
	-p 8092:8081 \
	-e MONGO_USER=${MONGO_USER} \
	-e MONGO_PASS=${MONGO_PASS} \
	-e MONGO_HOST=${DOCKER_HOST_IP} \
	-e MONGO_PORT=${MONGO_PORT} \
	-e CSM_LOG_LEVEL=${CSM_LOG_LEVEL} \
	-e CSM_API_KEY=${CSM_API_KEY} \
	-e CSM_DEV_MODE=${CSM_DEV_MODE} \
	-d ${CSM_EXTENSION_IMAGE_NAME}:${CSM_EXTENSION_IMAGE_TAG}
