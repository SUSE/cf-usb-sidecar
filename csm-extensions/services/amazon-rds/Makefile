NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
OK_GREEN_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_CYN_COLOR=\033[33;01m

ifndef DOCKER_REPOSITORY
	export DOCKER_REPOSITORY:=catalog-service-manager
endif

export CSM_ROOT:=${GOPATH}/src/github.com/hpcloud/catalog-service-manager
export CSM_RDS_EXTENSION_ROOT:=${CSM_ROOT}/csm-extensions/services/amazon-rds

ifeq ($(strip $(VERSION)),)
	export VERSION := $(shell ${CSM_ROOT}/scripts/build_version.sh "VERSION")
endif

ifeq ($(strip $(APP_VERSION_TAG)),)
	export APP_VERSION_TAG := $(shell VERSION=$(VERSION) ${CSM_ROOT}/scripts/build_version.sh "APP_VERSION_TAG")
endif

export CSM_EXTENSION_IMAGE_NAME:=csm-rds-mysql
export CSM_EXTENSION_IMAGE_TAG:=latest
export CSM_EXTENSION_BUILD_IMAGE_NAME:=csm-rds-mysql-build


.PHONY: clean-all clean build-image run-image tools

default: help

help:
	@echo "These 'make' targets are available."
	@echo
	@echo "  all                  cleans existing container and images and then"
	@echo "                       build docker image and runs the csm-mysql container "
	@echo "  build                build the amazon rds mysql binary"
	@echo "  run                  run the csm docker container for amazon rds"
	@echo "  tools                Installs tools needed to run"
	@echo "  test-format          Run the formatting tests"
	@echo "  test                 Run the formatting tests"
	@echo "  clean-containers     Remove all docker containers amazon-rds-mysql"
	@echo "  clean-images         Remove all docker images for amazon-rds-mysql"
	@echo "  clean-all            Remove docker container and images"
	@echo "  build-image          Build amazon-rds-mysql docker image "
	@echo "  publish-image        Publish the amazon-rds-mysql docker image to registry"
	@echo

all:	clean-all build-image tools run

build: 
	go install ./src/amazon-rds-mysql.go

run:
	@echo "$(OK_COLOR)==> Run image $(NO_COLOR)"
	${CSM_RDS_EXTENSION_ROOT}/scripts/docker-run-csm-mysql.sh

tools:
	go get github.com/aws/aws-sdk-go
	go get github.com/go-sql-driver/mysql
	go get github.com/go-ini/ini
	go get github.com/jmespath/go-jmespath

test-format:
	@echo "$(OK_COLOR)==> Running gofmt $(NO_COLOR)"
	${CSM_ROOT}/scripts/testFmt.sh src
		
# (required) run tests
test:	test-format

test-integration:
	echo "$(WARN_CYN_COLOR)==>No Integration tests $(NO_COLOR) ";n=0;\


# (required) clean containers
clean-containers:
	${CSM_ROOT}/scripts/docker/remove-docker-container.sh ${CSM_EXTENSION_IMAGE_NAME}
	${CSM_ROOT}/scripts/docker/remove-docker-container.sh ${CSM_EXTENSION_BUILD_IMAGE_NAME}

# (required) clean docker images
clean-images:
	${CSM_ROOT}/scripts/docker/remove-docker-image.sh ${CSM_EXTENSION_IMAGE_NAME}
	${CSM_ROOT}/scripts/docker/remove-docker-image.sh ${CSM_EXTENSION_BUILD_IMAGE_NAME}

# (required) clean docker containers and images
clean-all:	clean-containers clean-images

# (required) builds the service image
build-service-image:
	echo "$(WARN_CYN_COLOR)==>No Service image available $(NO_COLOR) ";n=0;\

# (required) build docker image for service
build-image:	
	@echo "$(OK_COLOR)==> Building Docker image $(NO_COLOR)"
	${CSM_RDS_EXTENSION_ROOT}/scripts/build-csm-amazon-rds-binary.sh
	${CSM_ROOT}/scripts/docker/remove-docker-container.sh ${CSM_EXTENSION_IMAGE_NAME}
	${CSM_ROOT}/scripts/docker/remove-docker-image.sh ${CSM_EXTENSION_BUILD_IMAGE_NAME}
	${CSM_RDS_EXTENSION_ROOT}/scripts/build-docker-image.sh


# (required) push image to docker registry
publish-image:
	IMAGE_NAME=${CSM_EXTENSION_IMAGE_NAME} IMAGE_TAG=${CSM_EXTENSION_IMAGE_TAG} ${CSM_ROOT}/scripts/docker/publish-image.sh
	