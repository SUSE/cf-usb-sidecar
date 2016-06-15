#!/bin/sh

if [ -z ${DOCKER_HOST} ]
then
    export DOCKER_HOST="redis-int.${HCP_SERVICE_DOMAIN_SUFFIX}"
fi

echo "Starting catalog-service-manager ..."

/catalog-service-manager/bin/catalog-service-manager