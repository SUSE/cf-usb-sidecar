#!/bin/sh

if [ -z ${DOCKER_HOST} ]
then
  if [ -z "${RABBITMQ_INT_SERVICE_PORT}" ]
  then
   export DOCKER_HOST="rabbitmq.${HCP_SERVICE_DOMAIN_SUFFIX}"   
  else
   export DOCKER_HOST="rabbitmq-int.${HCP_SERVICE_DOMAIN_SUFFIX}"
  fi
fi

if [ -z "${DOCKER_PORT}" ]
then
    if [ -z "${RABBITMQ_INT_SERVICE_PORT}" ]
    then
     export DOCKER_PORT=${RABBITMQ_SERVICE_PORT}
    else
     export DOCKER_PORT=${RABBITMQ_INT_SERVICE_PORT}
    fi
fi

echo "Starting catalog-service-manager ..."

/catalog-service-manager/bin/catalog-service-manager
