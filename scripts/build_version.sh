#!/bin/sh 

set +o errexit +o nounset

test -n "${XTRACE}" && set -o xtrace

set -o errexit -o nounset

env_var_to_echo=$1

git fetch --tags
export APP_VERSION_TAG=$(git describe --tags)

set +o errexit +o nounset +o xtrace

if [ "${env_var_to_echo}" = "APP_VERSION_TAG" ]; then
  echo ${APP_VERSION_TAG}
fi

