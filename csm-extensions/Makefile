NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
OK_GREEN_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_CYN_COLOR=\033[33;01m

ifndef DOCKER_REPOSITORY
	export DOCKER_REPOSITORY:=hsm-sidecar
endif

export SIDECAR_ROOT:=${GOPATH}/src/github.com/hpcloud/catalog-service-manager

ifeq ($(strip $(VERSION)),)
	export VERSION := $(shell ${SIDECAR_ROOT}/scripts/build_version.sh "VERSION")
endif

ifeq ($(strip $(APP_VERSION_TAG)),)
	export APP_VERSION_TAG := $(shell VERSION=$(VERSION) ${SIDECAR_ROOT}/scripts/build_version.sh "APP_VERSION_TAG")
endif

.PHONY: clean-all test clean-containers clean-images clean-all build-docker-image push-docker-image

default: help

help:
	@echo "These 'make' targets are available."
	@echo
	@echo "  test                 Run the tests on all the service extensions"
	@echo "  clean-containers     Remove docker containers for all the service extensions"
	@echo "  clean-images         Remove docker images for all the service extensions"
	@echo "  clean-all            Remove docker container and images for all the service extensions"
	@echo "  build-image          Build docker image for all the service extensions"
	@echo "  publish-image        Publish docker images to registry for all the service extensions"
	@echo

test:
	scripts/run_services_make_target.sh  test

test-integration:
	scripts/run_services_make_target.sh  test-integration

clean-containers:
	scripts/run_services_make_target.sh  clean-containers
	
clean-images:
	scripts/run_services_make_target.sh  clean-images
	
clean-all:
	scripts/run_services_make_target.sh  clean-all
	
build-image:
	scripts/run_services_make_target.sh  build-image

build-service-image:
	scripts/run_services_make_target.sh build-service-image

publish-image:
	scripts/run_services_make_target.sh  publish-image
