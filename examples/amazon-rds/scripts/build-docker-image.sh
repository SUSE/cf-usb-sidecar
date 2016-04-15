#!/bin/sh -x

CSM_IMAGE_NAME="amazon-rds-mysql:release"
CSM_LOG_LEVEL="debug"
docker build -t ${CSM_IMAGE_NAME} -f Dockerfile .
