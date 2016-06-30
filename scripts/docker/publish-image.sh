#!/bin/sh
. ${CSM_ROOT}/scripts/colors.sh

if [ -z ${REGISTRY_LOCATION} ]; then
	echo "Cannot push images as REGISTRY_LOCATION is not set"
	echo "if you want to push this to local docker registry use"
	echo "${WARN_MAGENTA} export REGISTRY_LOCATION=localhost:5000${NO_COLOR}"
	echo ""
	echo "if you want to push this to cnap shared docker registry use"
	echo "${WARN_MAGENTA} export REGISTRY_LOCATION=docker-registry.helion.space:443${NO_COLOR}"
	exit 1
fi

if [ -z ${IMAGE_NAME} ]; then
	echo "Error: Please set value for environemtn variable CSM_IMAGE_NAME"
	exit 1
fi

if [ -z ${IMAGE_TAG} ]; then
	echo "Error: Please set value for environemtn variable CSM_IMAGE_TAG"
	exit 1
fi

if [ -z ${DOCKER_REPOSITORY} ]; then
	echo "Error: Please set value for environemtn variable DOCKER_REPOSITORY"
	exit 1
fi

if [ -z ${APP_VERSION_TAG} ]; then
	echo "Error: Please set value for environemtn variable APP_VERSION_TAG"
	exit 1
fi

docker images ${CSM_IMAGE_NAME} | grep ${IMAGE_TAG} > /dev/null
if [ $? -eq 0 ]; then
	docker tag ${IMAGE_NAME}:${IMAGE_TAG} ${REGISTRY_LOCATION}/${DOCKER_REPOSITORY}/${IMAGE_NAME}:${IMAGE_TAG}
	docker tag ${IMAGE_NAME}:${IMAGE_TAG} ${REGISTRY_LOCATION}/${DOCKER_REPOSITORY}/${IMAGE_NAME}:${APP_VERSION_TAG}
	docker push ${REGISTRY_LOCATION}/${DOCKER_REPOSITORY}/${IMAGE_NAME}:${IMAGE_TAG}
	docker push ${REGISTRY_LOCATION}/${DOCKER_REPOSITORY}/${IMAGE_NAME}:${APP_VERSION_TAG}
else
	echo "Error: Docker image ${CSM_IMAGE_NAME}:${IMAGE_TAG} not found"
	echo "Before running publish-image, please use 'make build-image' to build the docker image."
	exit 1
fi
