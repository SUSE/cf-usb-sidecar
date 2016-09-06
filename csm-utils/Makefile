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

.PHONY: clean-all build-image publish-image

default: help

help:
	@echo "These 'make' targets are available."
	@echo
	@echo "  clean-all            Remove docker container and images for all the service extensions"
	@echo "  build-image          Build docker image for all the service extensions"
	@echo "  publish-image        Publish docker images to registry for all the service extensions"
	@echo

clean-all:
	scripts/run_utils_make_target.sh  clean-all
	
build-image:
	scripts/run_utils_make_target.sh  build-image

publish-image:
	scripts/run_utils_make_target.sh  publish-image
