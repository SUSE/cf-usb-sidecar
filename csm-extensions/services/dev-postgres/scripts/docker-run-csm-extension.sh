#!/bin/sh

POSTGRES_USER="postgres"
POSTGRES_PASS="password"
POSTGRES_PORT="5432"
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
	-p 8093:8081 \
	-e POSTGRES_USER=${POSTGRES_USER} \
	-e POSTGRES_PASS=${POSTGRES_PASS} \
	-e POSTGRES_HOST=${DOCKER_HOST_IP} \
	-e POSTGRES_PORT=${POSTGRES_PORT} \
	-e CSM_LOG_LEVEL=${CSM_LOG_LEVEL} \
	-e CSM_API_KEY=${CSM_API_KEY} \
	-e CSM_DEV_MODE=${CSM_DEV_MODE} \
	-d ${CSM_EXTENSION_IMAGE_NAME}:${CSM_EXTENSION_IMAGE_TAG}
