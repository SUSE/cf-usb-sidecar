#!/bin/sh
set +e

DOCKER_IMAGE=$1
DOCKER_IMAGE_TAG=$2
if [ -z ${DOCKER_IMAGE_TAG} ] 
then
    docker images | grep ${DOCKER_IMAGE}
    if [ $? -eq 0 ]
    then
        echo "Deleting docker images for ${DOCKER_IMAGE}"
        docker images | grep ${DOCKER_IMAGE} | awk '{print $3}' | xargs -n 1 docker rmi -f
    else
        echo "No docker image found with name ${DOCKER_IMAGE}"
    fi
else
    docker images | grep ${DOCKER_IMAGE} | grep ${DOCKER_IMAGE_TAG}
    if [ $? -eq 0 ]
    then
        echo "Deleting docker images for ${DOCKER_IMAGE}:${DOCKER_IMAGE_TAG}"
        docker images | grep ${DOCKER_IMAGE} | grep ${DOCKER_IMAGE_TAG} | awk '{print $3}' | xargs -n 1 docker rmi -f
    else
        echo "No docker image found with name ${DOCKER_IMAGE} and tag ${DOCKER_IMAGE_TAG}"
    fi    
fi
