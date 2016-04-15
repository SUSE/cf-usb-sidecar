NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
OK_GREEN_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_CYN_COLOR=\033[33;01m

.PHONY: clean-all clean build-image run-image tools

default: help

help:
	@echo "These 'make' targets are available."
	@echo
	@echo "  clean-docker	Remove all docker containers and images for amazon-rds-mysql"
	@echo "  clean-all	Remove all docker containers and images for amazon-rds-mysql"
	@echo "  build-image	Build amazon-rds-mysql docker image"
	@echo "  run		Run amazon-rds-mysql docker container using docker"
	@echo "  tools		Installs tools needed to run"
	@echo "  all		Cleans, build docker image and runs the amazon-rds-mysql container "


	@echo


all: 	clean-all build-image tools run

build: 
	go install ./src/amazon-rds-mysql.go

clean-docker:
	../../scripts/docker/remove-docker-container.sh amazon-rds-mysql
	../../scripts/docker/remove-docker-container.sh csm-rds
	../../scripts/docker/remove-docker-image.sh amazon-rds-mysql
	../../scripts/docker/remove-docker-image.sh csm-rds

clean-all: clean-docker

build-image:	
	@echo "$(OK_COLOR)==> Building Docker image $(NO_COLOR)"
	scripts/build-csm-amazon-rds-binary.sh
	scripts/build-docker-image.sh

run:
	@echo "$(OK_COLOR)==> Run image $(NO_COLOR)"
	scripts/docker-run-csm-mysql.sh

tools:
	go get github.com/aws/aws-sdk-go
	go get github.com/go-sql-driver/mysql
	go get github.com/go-ini/ini
	go get github.com/jmespath/go-jmespath
