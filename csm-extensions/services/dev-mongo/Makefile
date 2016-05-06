NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
OK_GREEN_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_CYN_COLOR=\033[33;01m

.PHONY: clean-all clean build-image tools run

export CSM_ROOT := $(shell echo ${GOPATH}/src/github.com/hpcloud/catalog-service-manager)
export CSM_EXTENSION_ROOT := $(shell echo ${CSM_ROOT}/csm-extensions/services/dev-mongo)
export CSM_EXTENSION_IMAGE_NAME := $(shell echo "csm-dev-mongo")
export CSM_EXTENSION_IMAGE_TAG := $(shell echo "latest")

export CSM_EXTENSION_BUILD_IMAGE_NAME := $(shell echo "csm-dev-mongo-build")


default: help

help:
	@echo "These 'make' targets are available."
	@echo
	@echo "  all                  cleans existing container and images and then"
	@echo "                       build docker image and runs the extension container "
	@echo "  build                build the extension binary"
	@echo "  run                  run the csm docker container for extension"
	@echo "  tools                Installs tools needed to run"
	@echo "  test-format          Run the formatting tests"
	@echo "  test                 Run the formatting tests"
	@echo "  clean-containers     Remove all docker containers for extension"
	@echo "  clean-images         Remove all docker images for extension"
	@echo "  clean-all            Remove docker container and images"
	@echo "  build-image          Build extension docker image"
	@echo "  push-image           Push the extension docker image to registry"
	@echo

all:	clean-all build-image run

build:
	${CSM_EXTENSION_ROOT}/scripts/build.sh

run:
	@echo "$(OK_COLOR)==> Run image $(NO_COLOR)"
	${CSM_EXTENSION_ROOT}/scripts/docker-run-csm-extension.sh

test-format:
	@echo "$(OK_COLOR)==> Running gofmt $(NO_COLOR)"
	${CSM_ROOT}/scripts/testFmt.sh src
		
# (required) run tests
test:	test-format

# (required) clean containers
clean-containers:
	${CSM_ROOT}/scripts/docker/remove-docker-container.sh ${CSM_EXTENSION_BUILD_IMAGE_NAME}
	${CSM_ROOT}/scripts/docker/remove-docker-container.sh ${CSM_EXTENSION_IMAGE_NAME}
	${CSM_ROOT}/scripts/docker/remove-docker-container.sh csm-dev-mysql-db

# (required) clean docker images
clean-images:
	${CSM_ROOT}/scripts/docker/remove-docker-image.sh ${CSM_EXTENSION_IMAGE_NAME}
	${CSM_ROOT}/scripts/docker/remove-docker-image.sh ${CSM_EXTENSION_BUILD_IMAGE_NAME}

# (required) clean docker containers and images
clean-all:	clean-containers clean-images

# (required) build docker image for service
build-image:	
	@echo "$(OK_COLOR)==> Building Docker image $(NO_COLOR)"
	${CSM_ROOT}/scripts/dev-extensions/build-csm-binaries.sh
	${CSM_ROOT}/scripts/docker/remove-docker-container.sh ${CSM_EXTENSION_BUILD_IMAGE_NAME}
	${CSM_ROOT}/scripts/docker/remove-docker-image.sh ${CSM_EXTENSION_BUILD_IMAGE_NAME}
	${CSM_ROOT}/scripts/dev-extensions/build-docker-image.sh


# (required) push image to docker registry
push-image:
	docker tag -f dev-mongo:latest $(REGISTRY_LOCATION)/catalog/dev-mongo:$(APP_VERSION_TAG)
	docker tag -f dev-mongo:latest $(REGISTRY_LOCATION)/catalog/dev-mongo:latest
	docker push $(REGISTRY_LOCATION)/catalog/dev-mongo:$(APP_VERSION_TAG)
	docker push $(REGISTRY_LOCATION)/catalog/dev-mongo:latest
