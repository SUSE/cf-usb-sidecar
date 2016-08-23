#!/bin/sh

MYSQL_ROOT_PASSWORD="root123"
SIDECAR_LOG_LEVEL="debug"
MYSQL_SERVICE_PORT_MYSQL="3306"
SIDECAR_DEV_MODE="true"
SIDECAR_API_KEY="sidecar-auth-token"



if [ ! -z ${DOCKER_HOST} ]
then
    export DOCKER_HOST_IP=`echo ${DOCKER_HOST} | cut -d "/" -f 3 | cut -d ":" -f 1`
else
    export DOCKER_HOST_IP=`ip route get 8.8.8.8 | awk 'NR==1 {print $NF}'`
fi


docker run --name csm-mysql \
       -p 8081:8081 \
       -e MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD} \
       -e SIDECAR_LOG_LEVEL=${SIDECAR_LOG_LEVEL} \
       -e MYSQL_SERVICE_HOST=${DOCKER_HOST_IP} \
       -e MYSQL_SERVICE_PORT_MYSQL=${MYSQL_SERVICE_PORT_MYSQL} \
       -e SIDECAR_API_KEY=${SIDECAR_API_KEY} \
       -e SIDECAR_DEV_MODE=${SIDECAR_DEV_MODE} \
       -d ${SIDECAR_MYSQL_IMAGE_NAME}:${SIDECAR_MYSQL_IMAGE_TAG}
