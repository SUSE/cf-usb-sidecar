#!/bin/sh

POSTGRES_PASSWORD="password"
POSTGRES_PORT="5432"

docker run --name csm-dev-postgres-db -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -p ${POSTGRES_PORT}:${POSTGRES_PORT} -d postgres:9.4
