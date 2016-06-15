#!/bin/sh

if [ -z ${POSTGRES_HOST} ]
then
    export POSTGRES_HOST="postgres-int.${HCP_SERVICE_DOMAIN_SUFFIX}"
fi

echo "Starting catalog-service-manager ..."

/catalog-service-manager/bin/catalog-service-manager