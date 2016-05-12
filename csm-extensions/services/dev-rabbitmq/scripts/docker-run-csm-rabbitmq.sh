#!/bin/sh

PORT="4445"

docker run --privileged --name csm-dev-rabbitmq -e PORT=${PORT} -p ${PORT}:${PORT} -d jpetazzo/dind
