#!/bin/sh

docker build -t ${CSM_RDS_IMAGE_NAME}:${CSM_RDS_IMAGE_TAG} --rm -f Dockerfile .
