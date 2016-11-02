NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
OK_GREEN_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_CYN_COLOR=\033[33;01m
TOKEN=sidecar-auth-token
SERVER_PORT=8081

ifndef DOCKER_ORGANIZATION
	export DOCKER_ORGANIZATION:=stackatodev
endif

export SIDECAR_ROOT:=${GOPATH}/src/github.com/hpcloud/catalog-service-manager
export SIDECAR_UTILS_ROOT:=${SIDECAR_ROOT}/csm-utils/utils/postflight

ifeq ($(strip $(VERSION)),)
	export VERSION := $(shell ${SIDECAR_ROOT}/scripts/build_version.sh "VERSION")
endif

ifeq ($(strip $(APP_VERSION_TAG)),)
	export APP_VERSION_TAG := $(shell VERSION=$(VERSION) ${SIDECAR_ROOT}/scripts/build_version.sh "APP_VERSION_TAG")
endif

export SIDECAR_UTILS_IMAGE_NAME:=hsm-sidecar-postflight
export SIDECAR_UTILS_IMAGE_TAG:=latest

.PHONY: clean-all build-image publish-image

default: help

help:
	@echo "These 'make' targets are available."
	@echo
	@echo "  clean-all            Remove docker images"
	@echo "  build-image          Build extension docker image"
	@echo "  publish-image        Publish sidecar-postflight docker image to registry"
	@echo

# (required) clean docker containers and images
clean-all:
	${SIDECAR_ROOT}/scripts/docker/remove-docker-image.sh ${SIDECAR_UTILS_IMAGE_NAME}

# (required) build docker image for service
build-image:
	@echo "$(OK_COLOR)==> Building Docker image $(NO_COLOR)"
	docker build -t ${SIDECAR_UTILS_IMAGE_NAME}:${SIDECAR_UTILS_IMAGE_TAG} .

# (required) push image to docker registry
publish-image:
	IMAGE_NAME=${SIDECAR_UTILS_IMAGE_NAME} IMAGE_TAG=${SIDECAR_UTILS_IMAGE_TAG} ${SIDECAR_ROOT}/scripts/docker/publish-image.sh
