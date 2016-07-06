#!/bin/bash

docker pull postgres:9.4
docker tag postgres:9.4 ${CSM_EXTENSION_SVC_IMAGE_NAME}:${CSM_EXTENSION_SVC_IMAGE_TAG}
