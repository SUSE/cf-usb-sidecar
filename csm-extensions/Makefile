NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
OK_GREEN_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_CYN_COLOR=\033[33;01m

ifndef DOCKER_REPOSITORY
	export DOCKER_REPOSITORY:=service-catalog
endif

ifndef REGISTRY_LOCATION
	export REGISTRY_LOCATION:=docker-registry.helion.space:443
endif

export CSM_ROOT:=${GOPATH}/src/github.com/hpcloud/catalog-service-manager

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
	@echo "  push-image           Push docker images to registry for all the service extensions"
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

push-image:
	scripts/run_services_make_target.sh  push-image
