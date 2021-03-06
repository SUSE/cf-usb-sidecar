#!/bin/sh

SERVICE_POSTGRES_USER="postgres"
SERVICE_POSTGRES_PASSWORD="password"
SERVICE_POSTGRES_PORT="5432"
SERVICE_POSTGRES_DBNAME="postgres"
SERVICE_POSTGRES_SSLMODE="disable"
SIDECAR_LOG_LEVEL="debug"
SIDECAR_DEV_MODE="true"


if [ ! -z ${DOCKER_HOST} ]
then
    export DOCKER_HOST_IP=`echo ${DOCKER_HOST} | cut -d "/" -f 3 | cut -d ":" -f 1`
else
    export DOCKER_HOST_IP=`ip route get 8.8.8.8 | awk 'NR==1 {print $NF}'`
fi


docker run --name ${SIDECAR_EXTENSION_IMAGE_NAME} \
	-p 0.0.0.0:${SIDECAR_EXTENSION_PORT}:8081 \
	-e SERVICE_POSTGRES_USER=${SERVICE_POSTGRES_USER} \
	-e SERVICE_POSTGRES_PASSWORD=${SERVICE_POSTGRES_PASSWORD} \
	-e SERVICE_POSTGRES_HOST=${DOCKER_HOST_IP} \
	-e SERVICE_POSTGRES_PORT=${SERVICE_POSTGRES_PORT} \
	-e SERVICE_POSTGRES_DBNAME=${SERVICE_POSTGRES_DBNAME} \
	-e SERVICE_POSTGRES_SSLMODE=${SERVICE_POSTGRES_SSLMODE} \
	-e SIDECAR_LOG_LEVEL=${SIDECAR_LOG_LEVEL} \
	-e SIDECAR_API_KEY=${SIDECAR_API_KEY} \
	-e SIDECAR_DEV_MODE=${SIDECAR_DEV_MODE} \
	-d ${SIDECAR_EXTENSION_IMAGE_NAME}:${SIDECAR_EXTENSION_IMAGE_TAG}
