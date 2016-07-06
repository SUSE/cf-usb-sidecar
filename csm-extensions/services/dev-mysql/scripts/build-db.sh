#!/bin/bash

docker pull mysql:5.5
docker tag mysql:5.5 ${CSM_EXTENSION_SVC_IMAGE_NAME}:${CSM_EXTENSION_SVC_IMAGE_TAG}
