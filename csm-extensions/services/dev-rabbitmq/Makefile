NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
OK_GREEN_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_CYN_COLOR=\033[33;01m

ifndef DOCKER_REPOSITORY
	export DOCKER_REPOSITORY:=hsm-sidecar
endif

export SIDECAR_ROOT:=${GOPATH}/src/github.com/hpcloud/catalog-service-manager
export SIDECAR_EXTENSION_ROOT:=${SIDECAR_ROOT}/csm-extensions/services/dev-rabbitmq

ifeq ($(strip $(VERSION)),)
	export VERSION := $(shell ${SIDECAR_ROOT}/scripts/build_version.sh "VERSION")
endif

ifeq ($(strip $(APP_VERSION_TAG)),)
	export APP_VERSION_TAG := $(shell VERSION=$(VERSION) ${SIDECAR_ROOT}/scripts/build_version.sh "APP_VERSION_TAG")
endif

export SIDECAR_EXTENSION_IMAGE_NAME:=sidecar-dev-rabbitmq
export SIDECAR_EXTENSION_IMAGE_TAG:=latest
export SIDECAR_EXTENSION_BUILD_IMAGE_NAME:=sidecar-dev-rabbitmq-build
export SIDECAR_EXTENSION_SVC_IMAGE_NAME:=sidecar-dev-rabbitmq-db
export SIDECAR_EXTENSION_SVC_VERSION:=3.6.0-management
export SIDECAR_EXTENSION_SVC_IMAGE_TAG:=latest
export SIDECAR_EXTENSION_SVC_CONTAINER_NAME:=sidecar-dev-rabbitmq-docker-host
export SIDECAR_EXTENSION_SVC_PORTS_START:=32000
export SIDECAR_EXTENSION_SVC_PORTS_END:=32020
export SIDECAR_EXTENSION_TOKEN:=sidecar-auth-token
export SIDECAR_EXTENSION_PORT:=8094

.PHONY: clean-all clean build-image build-service-image tools run test test-integration

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
	@echo "  build-service-image  Build database docker image"
	@echo "  publish-image        Publish csm-dev-rabbitmq docker image to registry"
	@echo

all:	clean-all build-image build-service-image run-db run

build:
	@(export GO15VENDOREXPERIMENT=1; \
	${SIDECAR_EXTENSION_ROOT}/scripts/build.sh)

run-db:
	@echo "$(OK_COLOR)==> Start Rabbitmq container $(NO_COLOR)"
	${SIDECAR_EXTENSION_ROOT}/scripts/docker-run-csm-rabbitmq.sh

run: build-image
	@echo "$(OK_COLOR)==> Run image $(NO_COLOR)"
	${SIDECAR_EXTENSION_ROOT}/scripts/docker-run-csm-extension.sh

test-format:
	@echo "$(OK_COLOR)==> Running gofmt $(NO_COLOR)"
	go list ./... | grep -v vendor | grep -v Godeps | gofmt -l

build-service-image: 
	@echo "$(OK_COLOR)==> Building database image $(NO_COLOR)"
	${SIDECAR_EXTENSION_ROOT}/scripts/build-db.sh

# (required) run tests
test: test-format
	@(export GO15VENDOREXPERIMENT=1; \
	go list ./... | grep -v vendor | go test -v)
	
test-integration: clean-all build-image build-service-image run-db run
	@(if [ ! -z ${DOCKER_HOST} ]; then \
	SERVER_IP=`echo ${DOCKER_HOST} | cut -d "/" -f 3 | cut -d ":" -f 1`;\
	else \
	SERVER_IP=`ip route get 8.8.8.8 | cut -d" " -f8`;\
	fi; \
	export DOCKER_HOST_IP=$$SERVER_IP;\
	echo "Testing against $$SERVER_IP:$(SIDECAR_EXTENSION_PORT)...\n";\
	echo "$(OK_COLOR)==>Waiting for docker to come online:$(NO_COLOR) ";n=0;\
	rez=`curl -i http://$$SERVER_IP:$(SIDECAR_EXTENSION_PORT)/workspaces -X POST -d '{"workspace_id":"initial11"}' -H "Content-Type: application/json" -H "x-sidecar-token: $(SIDECAR_EXTENSION_TOKEN)" --silent | grep successful` ; \
	while [ -z "$$rez" ] ; \
	do n=$$(( $$n + 1 )) ;\
	printf ".";\
	rez=`curl -i http://$$SERVER_IP:$(SIDECAR_EXTENSION_PORT)/workspaces -X POST -d '{"workspace_id":"initial11"}' -H "Content-Type: application/json" -H "x-sidecar-token: $(SIDECAR_EXTENSION_TOKEN)" --silent | grep successful` ; \
	if [ $$n -ge 130 ] ;\
	then echo "\n$(ERROR_COLOR)==>Docker took to long to wakeup or incorect setup $(NO_COLOR) "; break ;\
	fi ;\
	sleep 2 ;\
	done; echo "";\
	if [ $$n -le 129 ]; then \
	echo "$(OK_COLOR)==>Running integration tests:$(NO_COLOR) ";\
	export GO15VENDOREXPERIMENT=1;\
	go test ./tests -tags integration -v;\
	else \
	echo "$(ERROR_COLOR)==>Not running integration tests:$(NO_COLOR) ";\
	fi)
	${MAKE} clean-all

# (required) clean containers
clean-containers:
	${SIDECAR_ROOT}/scripts/docker/remove-docker-container.sh ${SIDECAR_EXTENSION_BUILD_IMAGE_NAME}
	${SIDECAR_ROOT}/scripts/docker/remove-docker-container.sh ${SIDECAR_EXTENSION_IMAGE_NAME}
	${SIDECAR_ROOT}/scripts/docker/remove-docker-container.sh ${SIDECAR_EXTENSION_SVC_CONTAINER_NAME}

# (required) clean docker images
clean-images:
	${SIDECAR_ROOT}/scripts/docker/remove-docker-image.sh ${SIDECAR_EXTENSION_IMAGE_NAME}
	${SIDECAR_ROOT}/scripts/docker/remove-docker-image.sh ${SIDECAR_EXTENSION_BUILD_IMAGE_NAME}
	${SIDECAR_ROOT}/scripts/docker/remove-docker-image.sh ${SIDECAR_EXTENSION_SVC_IMAGE_NAME}

# (required) clean docker containers and images
clean-all:	clean-containers clean-images

# (required) build docker image for service
build-image:	
	@echo "$(OK_COLOR)==> Building Docker image $(NO_COLOR)"
	${SIDECAR_ROOT}/scripts/dev-extensions/build-csm-binaries.sh
	${SIDECAR_ROOT}/scripts/docker/remove-docker-container.sh ${SIDECAR_EXTENSION_BUILD_IMAGE_NAME}
	${SIDECAR_ROOT}/scripts/docker/remove-docker-image.sh ${SIDECAR_EXTENSION_BUILD_IMAGE_NAME}
	${SIDECAR_ROOT}/scripts/dev-extensions/build-docker-image.sh


# (required) push image to docker registry
publish-image:
	IMAGE_NAME=${SIDECAR_EXTENSION_IMAGE_NAME} IMAGE_TAG=${SIDECAR_EXTENSION_IMAGE_TAG} ${SIDECAR_ROOT}/scripts/docker/publish-image.sh
	IMAGE_NAME=${SIDECAR_EXTENSION_SVC_IMAGE_NAME} IMAGE_TAG=${SIDECAR_EXTENSION_SVC_IMAGE_TAG} ${SIDECAR_ROOT}/scripts/docker/publish-image.sh
