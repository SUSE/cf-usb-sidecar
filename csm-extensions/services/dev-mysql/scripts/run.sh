#!/bin/sh

if [ -z ${MYSQL_HOST} ]
then
    export MYSQL_HOST="mysql-int.${HCP_SERVICE_DOMAIN_SUFFIX}"
fi

echo "Starting catalog-service-manager ..."

/catalog-service-manager/bin/catalog-service-manager