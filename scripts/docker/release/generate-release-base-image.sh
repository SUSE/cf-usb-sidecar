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
	printf "${ERROR_COLOR}==> Please execute script from catalog-service-manager's project root directory ${NO_COLOR}\n"
	exit 1
fi

if [ "${script_dir}" != "scripts/docker/release" ]
then
	printf "${ERROR_COLOR}==> Script directory is not correct, please run script from \${PROJECT_ROOT}/scripts/docker/release ${NO_COLOR}\n"
	exit 1
fi

SIDECAR_BUILD_IMAGE_NAME=cf-usb-sidecar-build-binary
SIDECAR_BUILD_IMAGE_TAG=latest
SIDECAR_BUILD_CONTAINER_NAME=sidecar-build

SIDECAR_BIN=${current_dir}/SIDECAR_BIN
mkdir -p ${SIDECAR_BIN}
rm -rf ${SIDECAR_BIN}/catalog-service-manager

docker images | grep ${SIDECAR_BUILD_IMAGE_NAME} | grep ${SIDECAR_BUILD_IMAGE_TAG} > /dev/null 2>&1
if [ $? -eq 0 ]
then
	printf "${OK_GREEN_COLOR}==> Remove ${SIDECAR_BUILD_IMAGE_NAME}:${SIDECAR_BUILD_IMAGE_TAG} image from docker before building new image ${NO_COLOR}\n"
	docker rmi ${SIDECAR_BUILD_IMAGE_NAME}:${SIDECAR_BUILD_IMAGE_TAG}
fi

printf "${OK_BG_COLOR}==> Building ${SIDECAR_BUILD_IMAGE_NAME}:${SIDECAR_BUILD_IMAGE_TAG} image (to build catalog-service-manager binary).. ${NO_COLOR}\n"

docker build -t ${SIDECAR_BUILD_IMAGE_NAME}:${SIDECAR_BUILD_IMAGE_TAG} --rm -f scripts/docker/release/Dockerfile-release-build .

docker images | grep ${SIDECAR_BUILD_IMAGE_NAME} | grep ${SIDECAR_BUILD_IMAGE_TAG} > /dev/null 2>&1
if [ $? -eq 0 ]
then
	printf "${OK_GREEN_COLOR}==> Copying CMS binary to the host ${NO_COLOR}\n"
	docker run --name ${SIDECAR_BUILD_CONTAINER_NAME} -v ${SIDECAR_BIN}:/catalog-service-manager-bin/ ${SIDECAR_BUILD_IMAGE_NAME}:${SIDECAR_BUILD_IMAGE_TAG}
fi

printf "${OK_GREEN_COLOR}==> Removing ${SIDECAR_BUILD_IMAGE_NAME} container  ${NO_COLOR}\n"
docker ps -a | grep ${SIDECAR_BUILD_CONTAINER_NAME} | awk '{print $1}' | xargs -n 1  docker rm


if [ -f ${SIDECAR_BIN}/catalog-service-manager ]
then
	printf "${OK_BG_COLOR}==> CSM binary is built successfully ${NO_COLOR}\n"
	echo ""
	printf "${OK_GREEN_COLOR}==> Removing ${SIDECAR_BASE_IMAGE_NAME}:${SIDECAR_BASE_IMAGE_TAG} images  ${NO_COLOR}\n"
	docker images | grep ${SIDECAR_BASE_IMAGE_NAME} | grep -v ${SIDECAR_BUILD_BASE_IMAGE_NAME}  | grep ${SIDECAR_BASE_IMAGE_TAG} | awk '{print $3}' | xargs -n 1 docker rmi -f

	sleep 5
	echo ""
	printf "${OK_GREEN_COLOR}==> Building ${SIDECAR_BASE_IMAGE_NAME}:${SIDECAR_BASE_IMAGE_TAG} image ..  ${NO_COLOR}\n"
	docker build -t ${SIDECAR_BASE_IMAGE_NAME}:${SIDECAR_BASE_IMAGE_TAG} --rm -f scripts/docker/release/Dockerfile-release .

	echo ""
	echo ""

	printf "${OK_BG_COLOR}==> ${SIDECAR_BASE_IMAGE_NAME}:${SIDECAR_BASE_IMAGE_TAG} built successfully ${NO_COLOR}\n"
fi
