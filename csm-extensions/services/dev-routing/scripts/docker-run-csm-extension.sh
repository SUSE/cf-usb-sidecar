#!/bin/sh

ROUTING_URL="https://`ifconfig eth0 | awk '/inet addr/{split($2,a,":"); print a[2]}'`:50000"
SIDECAR_LOG_LEVEL="debug"
SIDECAR_DEV_MODE="true"
SIDECAR_API_KEY="sidecar-auth-token"


docker run --name csm-routing \
       -p 8100:8081 \
       -e ROUTING_URL=${ROUTING_URL} \
       -e SIDECAR_LOG_LEVEL=${SIDECAR_LOG_LEVEL} \
	   -e SIDECAR_API_KEY=${SIDECAR_API_KEY} \
       -e SIDECAR_DEV_MODE=${SIDECAR_DEV_MODE} \
       -d ${SIDECAR_EXTENSION_IMAGE_NAME}:${SIDECAR_EXTENSION_IMAGE_TAG}
