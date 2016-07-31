#!/bin/sh

MYSQL_ROOT_PASSWORD="root123"
CSM_LOG_LEVEL="debug"
MYSQL_SERVICE_PORT_MYSQL="3306"
CSM_DEV_MODE="true"
CSM_API_KEY="csm-auth-token"



if [ ! -z ${DOCKER_HOST} ]
then
    export DOCKER_HOST_IP=`echo ${DOCKER_HOST} | cut -d "/" -f 3 | cut -d ":" -f 1`
else
    export DOCKER_HOST_IP=`ip route get 8.8.8.8 | awk 'NR==1 {print $NF}'`
fi


docker run --name csm-mysql \
       -p 8081:8081 \
       -e MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD} \
       -e CSM_LOG_LEVEL=${CSM_LOG_LEVEL} \
       -e MYSQL_SERVICE_HOST=${DOCKER_HOST_IP} \
       -e MYSQL_SERVICE_PORT_MYSQL=${MYSQL_SERVICE_PORT_MYSQL} \
       -e CSM_API_KEY=${CSM_API_KEY} \
       -e CSM_DEV_MODE=${CSM_DEV_MODE} \
       -d ${CSM_EXTENSION_IMAGE_NAME}:${CSM_EXTENSION_IMAGE_TAG}
