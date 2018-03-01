#!/bin/sh

set -o errexit

NO_COLOR="\\033\[0m"
OK_COLOR="\\033\[32\;01m"
ERROR_COLOR="\\033\[31\;01m"
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

printf "%s==>Waiting for docker to come online:%s\n" "${OK_COLOR}" "${NO_COLOR}"

n=0
until curl "http://${SERVER_IP}:${SERVER_PORT}/workspaces" \
        -X POST -d '{"workspace_id":"initial"}' \
        -H "Content-Type: application/json" \
        -H "x-sidecar-token: ${TOKEN}" \
        --fail --silent --output /dev/null
do n=$(( n + 1 ))
  printf "."
  if [ $n -ge 30 ] ; then
    printf "\n%s==>Docker took to long to wakeup or incorrect setup %s\n" "${ERROR_COLOR}" "${NO_COLOR}"
    break 
  fi 
  sleep 1 
done

printf "\n"

if [ $n -le 19 ]; then 
  printf "%s==>Running integration tests:%s\n" "${OK_COLOR}" "${NO_COLOR}"
  export GO15VENDOREXPERIMENT=1
  go test ./tests -tags integration -v
else 
  printf "%s==>Not running integration tests:%s\n" "${ERROR_COLOR}" "${NO_COLOR}"
fi
