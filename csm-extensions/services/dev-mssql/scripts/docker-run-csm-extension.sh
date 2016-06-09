#!/bin/sh

MSSQL_USER="sa"
MSSQL_PASS="password"
MSSQL_HOST="sqlserver"
MSSQL_PORT="1433"
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
	-p 0.0.0.0:8091:8081 \
	-e MSSQL_USER=${MSSQL_USER} \
	-e MSSQL_PASS=${MSSQL_PASS} \
	-e MSSQL_HOST=${MSSQL_HOST} \
	-e MSSQL_PORT=${MSSQL_PORT} \
	-e CSM_LOG_LEVEL=${CSM_LOG_LEVEL} \
	-e CSM_API_KEY=${CSM_API_KEY} \
	-e CSM_DEV_MODE=${CSM_DEV_MODE} \
	-d ${CSM_EXTENSION_IMAGE_NAME}:${CSM_EXTENSION_IMAGE_TAG}
