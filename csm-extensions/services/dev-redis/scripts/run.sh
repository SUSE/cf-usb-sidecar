#!/bin/sh

if [ -z "${DOCKER_HOST}" ]
then
  if [ -z "${REDIS_INT_SERVICE_PORT}" ]
  then
    export DOCKER_HOST="redis.${HCP_SERVICE_DOMAIN_SUFFIX}"
  else
    export DOCKER_HOST="redis-int.${HCP_SERVICE_DOMAIN_SUFFIX}"
  fi
fi


if [ -z "${DOCKER_PORT}" ]
then
    if [ -z "${REDIS_INT_SERVICE_PORT}" ]
    then
     export DOCKER_PORT=${REDIS_SERVICE_PORT}
    else
     export DOCKER_PORT=${REDIS_INT_SERVICE_PORT}
    fi
fi

echo "Starting catalog-service-manager ..."

/catalog-service-manager/bin/catalog-service-manager