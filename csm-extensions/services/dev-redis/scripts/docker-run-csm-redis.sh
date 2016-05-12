#!/bin/sh

PORT="4446"

docker run --privileged --name csm-dev-redis -e PORT=${PORT} -p ${PORT}:${PORT} -d jpetazzo/dind