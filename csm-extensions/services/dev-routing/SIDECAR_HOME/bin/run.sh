#!/bin/sh

if [ -z "${ROUTING_URL}" ]
  then
    export ROUTING_URL=${ROUTING_URL}
fi

echo "Starting catalog-service-manager ..."

/catalog-service-manager/bin/catalog-service-manager
