NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
OK_GREEN_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_CYN_COLOR=\033[33;01m

include version.mk

ifeq ($(strip $(VERSION)),)
    export VERSION := $(shell scripts/build_version.sh "VERSION")
endif

ifeq ($(strip $(APP_VERSION)),)
	export APP_VERSION := $(shell VERSION=$(VERSION) scripts/build_version.sh "APP_VERSION")
endif

ifeq ($(strip $(APP_VERSION_TAG)),)
	export APP_VERSION_TAG := $(shell VERSION=$(VERSION) scripts/build_version.sh "APP_VERSION_TAG")
endif

# Set environment variables if they are not before starting.
ifndef SIDECAR_API_KEY
	export SIDECAR_API_KEY:=sidecar-auth-token
endif

ifndef DOCKER_ORGANIZATION
	export DOCKER_ORGANIZATION:=splatform
endif

export SIDECAR_ROOT:=${GOPATH}/src/github.com/SUSE/cf-usb-sidecar
export SIDECAR_BASE_IMAGE_NAME:=cf-usb-sidecar
export SIDECAR_BASE_IMAGE_TAG:=latest
export SIDECAR_BUILD_BASE_IMAGE_NAME:=cf-usb-sidecar-buildbase
export SIDECAR_BUILD_BASE_IMAGE_TAG:=latest

.PHONY: run all clean clean-all clean-docker generate build test coverage tools build-image publish-image

default: help

help:
	@echo "These 'make' targets are available."
	@echo
	@echo "  run                Generates, runs the service locally in go"
	@echo "  all                Cleans, builds, runs tests"
	@echo "  clean              Removes all build output"
	@echo "  clean-all          Remove all build output and generated code"
	@echo "  clean-docker       Remove all docker containers and images for catalog-service-manager"
	@echo "  generate           Generates both server and client"
	@echo "  build              Generates swagger code and rebuilds the service only"
	@echo "  test               Run the unit tests"
	@echo "  coverage           Run the unit tests and produces a coverage report"
	@echo "  tools              Installs tools needed to run"
	@echo "  release-base       Builds docker image for release"
	@echo "  build-image        Builds docker image for release"
	@echo "  publish-image      Publish csm docker image to registry"
	@echo


run:	generate
	./scripts/run.sh

all: 	clean-all build test

clean:
	@printf "$(OK_COLOR)==> Removing build artifacts$(NO_COLOR)\n"
	rm -rf ${GOBIN}/catalog-service-manager
	rm -rf bin
	rm -rf SIDECAR_BIN

clean-all: clean
	@printf "$(OK_COLOR)==> Removing generated code$(NO_COLOR)\n"
	rm -rf generated

clean-docker:
	scripts/docker/remove-docker-container.sh sidecar
	scripts/docker/remove-docker-image.sh sidecar

generate-server:
	@printf "$(OK_COLOR)==> Generating code: server$(NO_COLOR)\n"
	rm -rf generated/CatalogServiceManager
	scripts/generate-server.sh

generate-client:
	@printf "$(OK_COLOR)==> Generating code: client$(NO_COLOR)\n"
	rm -rf generated/CatalogServiceManager-client
	scripts/generate-csm-client.sh

generate: generate-server generate-client

coverage:
	@printf "$(OK_COLOR)==> Running tests with coverage tool$(NO_COLOR)\n"
	./scripts/testCoverage.sh

build:	generate
	@printf "$(OK_COLOR)==> Building Catalog Service Manager code $(NO_COLOR)\n"
	./scripts/build.sh

test-format:
	@printf "$(OK_COLOR)==> Running gofmt $(NO_COLOR)\n"
	FILES=`find cmd src -name "*.go" | grep -v github.com/go-swagger`;\
	./scripts/testFmt.sh "$$FILES"

test: test-format
	@printf "$(OK_COLOR)==> Running tests $(NO_COLOR)\n"
	./scripts/test.sh

tools:
	@printf "$(OK_COLOR)==> Installing tools and go dependencies $(NO_COLOR)\n"
	go get golang.org/x/tools/cmd/cover
	go get github.com/tools/godep
	go get github.com/fsouza/go-dockerclient

build-image: clean-all
	@printf "$(OK_COLOR)==> Building release docker image for Catalog Service Manager $(NO_COLOR)\n"
	scripts/docker/release/generate-release-base-image.sh

release-base: build-image

publish-image:
	IMAGE_NAME=${SIDECAR_BASE_IMAGE_NAME} IMAGE_TAG=${SIDECAR_BASE_IMAGE_TAG} scripts/docker/publish-image.sh
