#!/bin/bash

set -o errexit

NO_COLOR="\033[0m"
OK_COLOR="\033[32m"
ERROR_COLOR="\033[31;01m"
TOKEN=sidecar-auth-token
SIDECAR_EXTENSION_PORT=8093

if [ ! -z "${DOCKER_HOST}" ]; then
    SERVER_IP=$(echo "${DOCKER_HOST}" | cut -d "/" -f 3 | cut -d ":" -f 1)
else
    SERVER_IP=$(ip route get 8.8.8.8 | cut -d" " -f8)
fi

SERVER_PORT=${SIDECAR_EXTENSION_PORT}

printf "Testing against %s:%s...\n" "${SERVER_IP}" "${SERVER_PORT}"

export TEST_SERVER_IP=${SERVER_IP}
export TEST_SERVER_PORT=${SERVER_PORT}
export TEST_SERVER_TOKEN=${TOKEN}

printf "${OK_COLOR}==> Waiting for docker to come online:${NO_COLOR} "

n=0
until curl "http://${SERVER_IP}:${SERVER_PORT}/workspaces" \
    -X POST -d '{"workspace_id":"initial"}' \
    -H "Content-Type: application/json" \
    -H "x-sidecar-token: ${TOKEN}" \
    --fail --silent --output /dev/null
do n=$(( n + 1 ))
    printf "."
    if [ $n -ge 30 ] ; then
	printf "\n${ERROR_COLOR}==> Docker took to long to wakeup or incorrect setup${NO_COLOR}\n"
	break
    fi
    sleep 1
done

printf "\n"

if [ $n -le 19 ]; then
    printf "${OK_COLOR}==> Running integration tests:${NO_COLOR}\n"
    export GO15VENDOREXPERIMENT=1
    go test ./tests -tags integration -v
else
    printf "${ERROR_COLOR}==> Not running integration tests${NO_COLOR}\n"
fi
