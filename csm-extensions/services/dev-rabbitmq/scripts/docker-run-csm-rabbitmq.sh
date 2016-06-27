#!/bin/sh

PORT="4445"

docker run --privileged \
	--name ${CSM_EXTENSION_SVC_CONTAINER_NAME} \
	-p ${PORT}:2375 \
	-p ${CSM_EXTENSION_SVC_PORTS_START}-${CSM_EXTENSION_SVC_PORTS_END}:${CSM_EXTENSION_SVC_PORTS_START}-${CSM_EXTENSION_SVC_PORTS_END} \
	-d ${CSM_EXTENSION_SVC_IMAGE_NAME}:${CSM_EXTENSION_SVC_IMAGE_TAG}
