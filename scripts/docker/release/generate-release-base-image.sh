#!/bin/sh

OK_COLOR='\033[1;32m'
OK_GREEN_COLOR='\033[0;32m'
OK_BG_COLOR='\033[42m'
WARN_CYN_COLOR='\033[36m'
ERROR_COLOR='\033[1;31m'
NO_COLOR='\033[0m'

force_rebuild=$1

. scripts/colors.sh
scripts/docker/generate-build-image.sh $force_rebuild

current_dir=$(pwd)
script_dir=$(dirname "$0")

if ! [ -d ${current_dir}/.git ]
then
	echo "${ERROR_COLOR}==> Please execute script from catalog-service-manager's project root directory ${NO_COLOR}"
	exit 1	
fi

if [ "${script_dir}" != "scripts/docker/release" ]
then
	echo "${ERROR_COLOR}==> Script directory is not correct, please run script from \${PROJECT_ROOT}/scripts/docker/release ${NO_COLOR}"
	exit 1	
fi

CSM_BUILD_IMAGE_NAME=csm-release-build
CSM_BUILD_IMAGE_TAG=latest
CSM_BUILD_CONTAINER_NAME=csm-build

CSM_BIN=${current_dir}/CSM_BIN
mkdir -p ${CSM_BIN}
rm -rf ${CSM_BIN}/catalog-service-manager

echo "${OK_BG_COLOR}==> Building release-build image (to build catalog-service-manager binary).. ${NO_COLOR}"

docker build -t ${CSM_BUILD_IMAGE_NAME}:${CSM_BUILD_IMAGE_TAG} --rm -f scripts/docker/release/Dockerfile-release-build .  

docker images | grep ${CSM_BUILD_IMAGE_NAME} | grep ${CSM_BUILD_IMAGE_TAG} > /dev/null 2>&1
if [ $? -eq 0 ]
then
	echo "${OK_GREEN_COLOR}==> Copying CMS binary to the host ${NO_COLOR}"
	docker run --name ${CSM_BUILD_CONTAINER_NAME} -v ${CSM_BIN}:/catalog-service-manager-bin/ ${CSM_BUILD_IMAGE_NAME}:${CSM_BUILD_IMAGE_TAG} 
fi

echo "${OK_GREEN_COLOR}==> Removing ${CSM_BUILD_IMAGE_NAME} container  ${NO_COLOR}"
docker ps -a | grep ${CSM_BUILD_CONTAINER_NAME} | awk '{print $1}' | xargs -n 1  docker rm 


if [ -f ${CSM_BIN}/catalog-service-manager ]
then
	echo "${OK_BG_COLOR}==> CSM binary is built successfuly ${NO_COLOR}"
	echo ""	
	echo "${OK_GREEN_COLOR}==> Removing ${CSM_BUILD_IMAGE_NAME}:${CSM_BUILD_IMAGE_TAG} images  ${NO_COLOR}"
	docker images | grep ${CSM_BASE_IMAGE_NAME} | grep -v ${CSM_BUILD_BASE_IMAGE_NAME}  | grep ${CSM_BASE_IMAGE_TAG} | awk '{print $3}' | xargs -n 1 docker rmi -f 
	
	sleep 5
	echo ""
	echo "${OK_GREEN_COLOR}==> Building ${CSM_BUILD_IMAGE_NAME}:release image ..  ${NO_COLOR}"
	docker build -t ${CSM_BASE_IMAGE_NAME}:${CSM_BASE_IMAGE_TAG} --rm -f scripts/docker/release/Dockerfile-release . 
	
	echo ""
	echo ""
	
	echo "${OK_BG_COLOR}==> ${CSM_BASE_IMAGE_NAME}:${CSM_BASE_IMAGE_TAG} built successfuly ${NO_COLOR}"
fi
