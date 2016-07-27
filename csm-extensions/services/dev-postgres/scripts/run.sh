#!/bin/sh

if [ -z "${SERVICE_POSTGRES_HOST}" ]
then
    if [ -z "${POSTGRES_INT_SERVICE_PORT}" ]
    then
      export SERVICE_POSTGRES_HOST="postgres.${HCP_SERVICE_DOMAIN_SUFFIX}"
	else
	  export SERVICE_POSTGRES_HOST="postgres-int.${HCP_SERVICE_DOMAIN_SUFFIX}"
	fi
fi

if [ -z "${SERVICE_POSTGRES_PORT}" ]
 then
     if [ -z "${POSTGRES_INT_SERVICE_PORT}" ]
     then
      export SERVICE_POSTGRES_PORT=${POSTGRES_SERVICE_PORT}
     else
      export SERVICE_POSTGRES_PORT=${POSTGRES_INT_SERVICE_PORT}
     fi
 fi

if [ -z "${SERVICE_POSTGRES_PASSWORD}" ]
  then
    export SERVICE_POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
fi

echo "Starting catalog-service-manager ..."

/catalog-service-manager/bin/catalog-service-manager
