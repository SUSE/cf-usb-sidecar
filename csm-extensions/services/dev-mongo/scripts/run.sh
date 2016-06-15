#!/bin/sh

if [ -z ${MONGO_HOST} ]
then
    export MONGO_HOST="mongo-int.${HCP_SERVICE_DOMAIN_SUFFIX}"
fi

echo "Starting catalog-service-manager ..."

/catalog-service-manager/bin/catalog-service-manager