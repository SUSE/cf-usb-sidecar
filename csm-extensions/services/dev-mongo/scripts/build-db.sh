#!/bin/bash

docker pull tutum/mongodb:3.0
docker tag tutum/mongodb:3.0 ${CSM_EXTENSION_SVC_IMAGE_NAME}:${CSM_EXTENSION_SVC_IMAGE_TAG}
