#!/bin/sh

MYSQL_ROOT_PASSWORD="root123"
MYSQL_SERVICE_PORT_MYSQL="3306"

docker run \
    --name ${SIDECAR_EXTENSION_SVC_IMAGE_NAME} \
    -e MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD} \
    -p 0.0.0.0:${MYSQL_SERVICE_PORT_MYSQL}:${MYSQL_SERVICE_PORT_MYSQL} \
    -d ${SIDECAR_EXTENSION_SVC_IMAGE_NAME}:${SIDECAR_EXTENSION_SVC_IMAGE_TAG}
