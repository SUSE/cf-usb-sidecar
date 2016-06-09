#!/bin/sh

POSTGRES_PASSWORD="password"
POSTGRES_PORT="5432"

docker run --name ${CSM_EXTENSION_SVC_CONTAINER_NAME} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -p 0.0.0.0:${POSTGRES_PORT}:${POSTGRES_PORT} -d ${CSM_EXTENSION_SVC_IMAGE_NAME}:${CSM_EXTENSION_SVC_IMAGE_TAG}
