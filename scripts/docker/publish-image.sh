#!/bin/sh
. ${SIDECAR_ROOT}/scripts/colors.sh

if [ -z ${DOCKER_REPOSITORY} ]; then
	echo "Cannot push images as DOCKER_REPOSITORY is not set"
	echo "if you want to push to your local docker registry use"
	printf "${WARN_MAGENTA} export DOCKER_REPOSITORY=localhost:5000${NO_COLOR}\n"
	echo ""
	echo "if you want to push to dockerhub use"
	printf "${WARN_MAGENTA} export DOCKER_REPOSITORY=docker.io${NO_COLOR}\n"
	exit 1
fi

if [ -z ${IMAGE_NAME} ]; then
	printf "${ERROR_COLOR}Error${NO_COLOR}: Please set environment variable SIDECAR_IMAGE_NAME\n"
	exit 1
fi

if [ -z ${IMAGE_TAG} ]; then
	printf "${ERROR_COLOR}Error${NO_COLOR}: Please set environment variable SIDECAR_IMAGE_TAG\n"
	exit 1
fi

if [ -z ${DOCKER_ORGANIZATION} ]; then
	printf "${ERROR_COLOR}Error${NO_COLOR}: Please set environment variable DOCKER_ORGANIZATION\n"
	exit 1
fi

if [ -z ${APP_VERSION_TAG} ]; then
	printf "${ERROR_COLOR}Error${NO_COLOR}: Please set environment variable APP_VERSION_TAG\n"
	exit 1
fi

printf "${NOTE_COLOR}APP_VERSION_TAG${NO_COLOR} ... = ${APP_VERSION_TAG}\n"
printf "${NOTE_COLOR}IMAGE_NAME${NO_COLOR} ........ = ${IMAGE_NAME}\n"
printf "${NOTE_COLOR}IMAGE_TAG${NO_COLOR} ......... = ${IMAGE_TAG}\n"
printf "${NOTE_COLOR}DOCKER_ORGANIZATION${NO_COLOR} = ${DOCKER_ORGANIZATION}\n"
printf "${NOTE_COLOR}DOCKER_REPOSITORY${NO_COLOR} . = ${DOCKER_REPOSITORY}\n"

#exit 0

docker images ${SIDECAR_IMAGE_NAME} | grep ${IMAGE_TAG} > /dev/null
if [ $? -eq 0 ]; then
	docker tag ${IMAGE_NAME}:${IMAGE_TAG} ${DOCKER_REPOSITORY}/${DOCKER_ORGANIZATION}/${IMAGE_NAME}:${IMAGE_TAG}
	docker tag ${IMAGE_NAME}:${IMAGE_TAG} ${DOCKER_REPOSITORY}/${DOCKER_ORGANIZATION}/${IMAGE_NAME}:${APP_VERSION_TAG}
	docker push ${DOCKER_REPOSITORY}/${DOCKER_ORGANIZATION}/${IMAGE_NAME}:${IMAGE_TAG}
	docker push ${DOCKER_REPOSITORY}/${DOCKER_ORGANIZATION}/${IMAGE_NAME}:${APP_VERSION_TAG}
else
	printf "${ERROR_COLOR}Error${NO_COLOR}: Docker image ${SIDECAR_IMAGE_NAME}:${IMAGE_TAG} not found\n"
	echo "Before running publish-image, please use 'make build-image' to build the docker image."
	exit 1
fi
