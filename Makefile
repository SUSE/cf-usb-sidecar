NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
OK_GREEN_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_CYN_COLOR=\033[33;01m

# Set environment variables if they are not before starting.
ifndef CSM_API_KEY
	export CSM_API_KEY:=csm-auth-token
endif

# List of files to be tested
TESTLIST=$(shell go list ./... | grep -v examples)

.PHONY: all clean build test release

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
	@echo "  build              generates swagger code and rebuilds the service only"
	@echo "  test               Run the unit tests"
	@echo "  coverage           Run the unit tests and produces a coverage report"
	@echo "  tools              Installs tools needed to run"
	@echo "  dev-base           Builds docker image for dev/test"
	@echo "  release-base       Builds docker image for release"
	@echo


run:	generate
	godep go run cmd/catalog-service-manager/catalog-service-manager.go

all: 	clean-all build test

clean:
	@echo "$(OK_COLOR)==> Removing build artifacts$(NO_COLOR)"
	rm -rf ${GOBIN}/catalog-service-manager
	rm -rf bin

clean-all: clean
	@echo "$(OK_COLOR)==> Removing generated code$(NO_COLOR)"
	rm -rf generated

clean-docker:
	scripts/docker/remove-docker-container.sh csm
	scripts/docker/remove-docker-container.sh catalog-service-manager
	scripts/docker/remove-docker-image.sh csm
	scripts/docker/remove-docker-image.sh catalog-service-manager

generate:
	@echo "$(OK_COLOR)==> Generating code $(NO_COLOR)"
	scripts/generate-server.sh

coverage:
	@echo "$(OK_COLOR)==> Running tests with coverage tool$(NO_COLOR)"
	./scripts/testCoverage.sh

build:	generate
	@echo "$(OK_COLOR)==> Building Catalog Service Manager code $(NO_COLOR)"
	cd cmd/catalog-service-manager;\
        godep go install .

test:
	@echo "$(OK_COLOR)==> Running tests $(NO_COLOR)"
	godep go test $(TESTLIST) | grep -v generated | grep -v cmd/catalog-service-manager/handlers | grep -v scripts 

tools:
	@echo "$(OK_COLOR)==> Installing tools and go dependancies $(NO_COLOR)"
	go get golang.org/x/tools/cmd/cover
	go get github.com/tools/godep

	./scripts/tools/codegen.sh

dev-base: clean-all
	@echo "$(OK_COLOR)==> Building dev/test docker image for Catalog Service Manager $(NO_COLOR)"
	scripts/docker/development/generate-development-base-image.sh

release-base: clean-all
	@echo "$(OK_COLOR)==> Building release docker image for Catalog Service Manager $(NO_COLOR)"
	scripts/docker/release/generate-release-base-image.sh

