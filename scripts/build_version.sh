#!/bin/sh 
set -u
env_var_to_echo=$1

export WORKSPACE=${GOPATH}/src/github.com/hpcloud/catalog-service-manager
STARTDIR=${START_DIR:-}
CONCOURSEBUILD=${CONCOURSE_BUILD:-}

# if $VERSION is set from CI build run then use it
if [ -z "${VERSION-}" ]; then
  if [ -f ${STARTDIR}/version/number ]; then
      echo "export VERSION=$(cat ${STARTDIR}/version/number)">$WORKSPACE/version.mk
      echo "export MAJOR_MINOR=$(cat ${STARTDIR}/version/number| cut -d "." -f1,2)">>$WORKSPACE/version.mk
      export CONCOURSEBUILD="1"
  fi
  . $WORKSPACE/version.mk
fi

export BRANCH=$(git name-rev --name-only HEAD)

# since '/' in not allowed in docker image tag 
export BRANCH_TAG=$(echo $BRANCH | tr "/" "-" )

build_commit_hash=$(git rev-parse --short HEAD)
build_time=$(date -u +%Y%m%d%H%M%S)

if [ -n "${CONCOURSEBUILD}" ]; then
  # concourse build number
  if [ "${BRANCH}" = "master" ]; then
    export APP_VERSION=$(git describe --tags --long)
    export APP_VERSION_TAG=$(git describe --tags --long)
    export APP_LATEST_BRANCH_TAG="latest"
  else
    export APP_VERSION="${VERSION}+${BRANCH}.${build_commit_hash}.${build_time}"
    export APP_VERSION_TAG="${VERSION}-${BRANCH_TAG}"
    export APP_LATEST_BRANCH_TAG="latest-${BRANCH_TAG}"
  fi
else
  #dev build number
  user=$(echo $(whoami))
  export APP_VERSION="${VERSION}+${user}.${BRANCH}.${build_commit_hash}.${build_time}"
  export APP_VERSION_TAG="${VERSION}-${user}-${BRANCH_TAG}"
  export APP_LATEST_BRANCH_TAG="latest-${user}-${BRANCH_TAG}"
fi

if [ "${env_var_to_echo}" = "VERSION" ]; then
  echo ${VERSION}
fi

if [ "${env_var_to_echo}" = "APP_VERSION" ]; then
  echo ${APP_VERSION}
fi

if [ "${env_var_to_echo}" = "APP_VERSION_TAG" ]; then
  echo ${APP_VERSION_TAG}
fi

if [ "${env_var_to_echo}" = "APP_LATEST_BRANCH_TAG" ]; then
  echo ${APP_LATEST_BRANCH_TAG}
fi

if [ "${env_var_to_echo}" = "BRANCH" ]; then
  echo ${BRANCH}
fi
