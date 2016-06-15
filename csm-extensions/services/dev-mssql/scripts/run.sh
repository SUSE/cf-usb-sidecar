#!/bin/sh

if [ -z ${MSSQL_HOST} ]
then
    export MSSQL_HOST="mysql-int.${HCP_SERVICE_DOMAIN_SUFFIX}"
fi

echo "Starting catalog-service-manager ..."

/catalog-service-manager/bin/catalog-service-manager