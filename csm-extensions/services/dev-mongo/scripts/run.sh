#!/bin/sh

if [ -z ${SERVICE_MONGO_HOST} ]
then
    if [ -z "${MONGO_INT_SERVICE_PORT}" ]
    then
      export SERVICE_MONGO_HOST="mongo.${HCP_SERVICE_DOMAIN_SUFFIX}"
	else
	  export SERVICE_MONGO_HOST="mongo-int.${HCP_SERVICE_DOMAIN_SUFFIX}"
	fi
fi

if [ -z "${SERVICE_MONGO_PORT}" ]
 then
     if [ -z "${MONGO_INT_SERVICE_PORT}" ]
     then
      export SERVICE_MONGO_PORT=${MONGO_SERVICE_PORT}
     else
      export SERVICE_MONGO_PORT=${MONGO_INT_SERVICE_PORT}
     fi
fi

if [ -z "${SERVICE_MONGO_PASS}" ]
  then
    export SERVICE_MONGO_PASS=${MONGODB_PASS}
fi

echo "Starting catalog-service-manager ..."

/catalog-service-manager/bin/catalog-service-manager
