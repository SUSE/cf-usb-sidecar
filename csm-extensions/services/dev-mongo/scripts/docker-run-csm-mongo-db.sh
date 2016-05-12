#!/bin/sh

MONGODB_PASS="password"
MONGODB_PORT="27017"

docker run --name csm-dev-mongodb-db -e MONGODB_PASS=${MONGODB_PASS} -p ${MONGODB_PORT}:${MONGODB_PORT} -d tutum/mongodb:3.0
