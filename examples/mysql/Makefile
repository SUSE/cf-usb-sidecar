NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
OK_GREEN_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_CYN_COLOR=\033[33;01m

.PHONY: all clean build test release

default: help

help:
	@echo "These 'make' targets are available."
	@echo
	@echo "  clean-docker		Remove all docker containers and images for csm-mysql"
	@echo "  clean-all			Remove all docker containers and images for csm-mysql"
	@echo "  build-image		build csm-mysql docker image"
	@echo "  run				run csm-mysql docker container using docker"
	@echo "  tools				run csm-mysql-db docker container using docker"
	@echo "  all				cleans, build docker image and runs the mysql and csm container "

	@echo


all: 	clean-all build-image tools run

clean-docker:
	../../scripts/docker/remove-docker-container.sh csm-mysql-db
	../../scripts/docker/remove-docker-container.sh csm-mysql
	../../scripts/docker/remove-docker-image.sh csm-mysql
	
clean-all: clean-docker

build-image:
	@echo "$(OK_COLOR)==> Building Docker image $(NO_COLOR)"
	scripts/build-docker-image.sh

run:
	@echo "$(OK_COLOR)==> Run image $(NO_COLOR)"
	scripts/docker-run-csm-mysql.sh

tools:
	scripts/docker-run-csm-mysql-db.sh
