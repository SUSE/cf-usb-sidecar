#!/bin/sh

. ${SIDECAR_ROOT}/scripts/colors.sh

force_rebuild=$1

current_dir=$(pwd)
script_dir=$(dirname "$0")

if ! [ -d ${current_dir}/.git ]
then
	echo "${ERROR_COLOR}==> Please execute script from catalog-service-manager's project root directory${NO_COLOR}"
	exit 1	
fi

if [ "${script_dir}" != "scripts/docker" ]
then
	echo "${ERROR_COLOR}==> Script directory is not correct, please run script from \${PROJECT_ROOT}/scripts/dockerv${NO_COLOR}"
	exit 1	
fi

docker images | grep ${SIDECAR_BUILD_BASE_IMAGE_NAME} | grep ${SIDECAR_BUILD_BASE_IMAGE_TAG} > /dev/null 2>&1
if [ $? -eq 0 ]
then
	if [ "$force_rebuild" != "rebuild-image" ]
	then
		echo "${WARN_MAGENTA}==> ${SIDECAR_BUILD_BASE_IMAGE_NAME}:${SIDECAR_BUILD_BASE_IMAGE_TAG} image already exists!${NO_COLOR}"
		exit 0
	fi
	
	if [ "$force_rebuild" == "rebuild-image" ]
	then
		echo "${OK_GREEN_COLOR}==> Removing old ${SIDECAR_BUILD_BASE_IMAGE_NAME}:${SIDECAR_BUILD_BASE_IMAGE_TAG} images ..${NO_COLOR}"
		docker images | grep ${SIDECAR_BUILD_BASE_IMAGE_NAME} | grep ${SIDECAR_BUILD_BASE_IMAGE_TAG} | awk '{print $3}' | xargs -n 1 docker rmi -f > /dev/null 2>&1
		sleep 5
	fi	
fi

echo "${OK_GREEN_COLOR}==> Building ${SIDECAR_BUILD_BASE_IMAGE_NAME}:build image ..${NO_COLOR}"
docker build -t ${SIDECAR_BUILD_BASE_IMAGE_NAME}:${SIDECAR_BUILD_BASE_IMAGE_TAG} --rm -f ${script_dir}/Dockerfile-build .