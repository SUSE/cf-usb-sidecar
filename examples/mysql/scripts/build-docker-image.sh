#!/bin/sh

CSM_IMAGE_NAME="csm-mysql:latest"
MYSQL_ROOT_PASSWORD="root123"
CSM_LOG_LEVEL="debug"
MYSQL_SERVICE_PORT_MYSQL="3306"

docker build -t ${CSM_IMAGE_NAME} -f Dockerfile .
