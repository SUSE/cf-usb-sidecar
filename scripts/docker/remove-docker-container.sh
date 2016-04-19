#!/bin/sh

CONTAINER_NAME=$1

docker ps -a | grep ${CONTAINER_NAME}

if [ $? -eq 0 ]
then
	docker ps -a | grep ${CONTAINER_NAME} | awk '{print $1}' | xargs -L 1 docker rm -f 
fi

