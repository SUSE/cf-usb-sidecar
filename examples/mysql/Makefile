NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
OK_GREEN_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_CYN_COLOR=\033[33;01m

.PHONY: all clean build test release

default: help

help:
	@echo "  tools 		run cms-mysql-db docker container using docker"

	@echo

tools:
	scripts/docker-run-csm-mysql-db.sh
