#!/bin/sh

MSSQL_USER="sa"
MSSQL_PASS="password"
MSSQL_HOST="sqlserver"
MSSQL_PORT="1433"
SIDECAR_LOG_LEVEL="debug"
SIDECAR_DEV_MODE="true"
SIDECAR_API_KEY="sidecar-auth-token"



if [ ! -z ${DOCKER_HOST} ]
then
    export DOCKER_HOST_IP=`echo ${DOCKER_HOST} | cut -d "/" -f 3 | cut -d ":" -f 1`
else
    export DOCKER_HOST_IP=`ip route get 8.8.8.8 | awk 'NR==1 {print $NF}'`
fi


docker run --name ${SIDECAR_EXTENSION_IMAGE_NAME} \
	-p 0.0.0.0:8091:8081 \
	-e MSSQL_USER=${MSSQL_USER} \
	-e MSSQL_PASS=${MSSQL_PASS} \
	-e MSSQL_HOST=${MSSQL_HOST} \
	-e MSSQL_PORT=${MSSQL_PORT} \
	-e SIDECAR_LOG_LEVEL=${SIDECAR_LOG_LEVEL} \
	-e SIDECAR_API_KEY=${SIDECAR_API_KEY} \
	-e SIDECAR_DEV_MODE=${SIDECAR_DEV_MODE} \
	-d ${SIDECAR_EXTENSION_IMAGE_NAME}:${SIDECAR_EXTENSION_IMAGE_TAG}
