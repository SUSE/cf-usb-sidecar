NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
OK_GREEN_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_CYN_COLOR=\033[33;01m

ifndef DOCKER_ORGANIZATION
	export DOCKER_ORGANIZATION:=stackatodev
endif

export SIDECAR_ROOT:=${GOPATH}/src/github.com/hpcloud/catalog-service-manager
export SIDECAR_MYSQL_EXTENSION_ROOT:=${SIDECAR_ROOT}/csm-extensions/services/mysql

ifeq ($(strip $(VERSION)),)
	export VERSION := $(shell ${SIDECAR_ROOT}/scripts/build_version.sh "VERSION")
endif

ifeq ($(strip $(APP_VERSION_TAG)),)
	export APP_VERSION_TAG := $(shell VERSION=$(VERSION) ${SIDECAR_ROOT}/scripts/build_version.sh "APP_VERSION_TAG")
endif

export SIDECAR_EXTENSION_IMAGE_NAME:=hsm-sidecar-mysql
export SIDECAR_EXTENSION_IMAGE_TAG:=latest

.PHONY: all clean build test release

default: help

help:
	@echo "These 'make' targets are available."
	@echo
	@echo "  all                  cleans existing container and images and then"
	@echo "                       build docker image and runs the csm-mysql container "
	@echo "  run                  run csm-mysql docker container using docker"
	@echo "  tools                run csm-mysql-db docker container using docker"
	@echo "  test-format          Run the formatting tests for csm-mysql"
	@echo "  test-integration     Run the integration tests for csm-mysql"
	@echo "  test                 Run the formatting and integration tests for csm-mysql"
	@echo "  clean-containers     Remove all docker containers for csm-mysql"
	@echo "  clean-images         Remove all docker images for csm-mysql"
	@echo "  clean-all            Remove docker container and images for csm-mysql"
	@echo "  build-image          Build csm-mysql docker image "
	@echo "  publish-image        Publish csm-mysql docker image to registry"
	@echo


all: 	clean-all build-image tools run

run:
	@echo "$(OK_COLOR)==> Run image $(NO_COLOR)"
	SIDECAR_API_KEY=sidecar-auth-token scripts/docker-run-csm-mysql.sh

tools:
	@echo "$(OK_COLOR)==> Start MySQL DB container $(NO_COLOR)"
	${SIDECAR_MYSQL_EXTENSION_ROOT}/scripts/docker-run-csm-mysql-db.sh

test-format:
	@echo "$(OK_COLOR)==> Running gofmt $(NO_COLOR)"
	${SIDECAR_ROOT}/scripts/testFmt.sh tests

# (required)
test:	test-format

# (required) build docker image for service
build-image:
	@echo "$(OK_COLOR)==> Building Docker image $(NO_COLOR)"
	docker build -t ${SIDECAR_EXTENSION_IMAGE_NAME}:${SIDECAR_EXTENSION_IMAGE_TAG} --rm -f ${SIDECAR_MYSQL_EXTENSION_ROOT}/Dockerfile .

test-integration: build-image
	${SIDECAR_MYSQL_EXTENSION_ROOT}/scripts/run_int_test.sh

# (required) clean containers
clean-containers:
	${SIDECAR_ROOT}/scripts/docker/remove-docker-container.sh csm-mysql-db
	${SIDECAR_ROOT}/scripts/docker/remove-docker-container.sh ${SIDECAR_EXTENSION_IMAGE_NAME}

# (required) clean docker images
clean-images:
	${SIDECAR_ROOT}/scripts/docker/remove-docker-image.sh ${SIDECAR_EXTENSION_IMAGE_NAME}

# (required) clean docker containers and images
clean-all: clean-containers clean-images

# (required) builds the service image
build-service-image:
	echo "$(WARN_CYN_COLOR)==>No Service image available $(NO_COLOR) ";n=0;\


# (required) push image to docker registry
publish-image:
	IMAGE_NAME=${SIDECAR_EXTENSION_IMAGE_NAME} IMAGE_TAG=${SIDECAR_EXTENSION_IMAGE_TAG} ${SIDECAR_ROOT}/scripts/docker/publish-image.sh
