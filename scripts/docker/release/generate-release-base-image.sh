#!/bin/sh

OK_COLOR='\033[1;32m'
OK_GREEN_COLOR='\033[0;32m'
OK_BG_COLOR='\033[42m'
WARN_CYN_COLOR='\033[36m'
ERROR_COLOR='\033[1;31m'
NO_COLOR='\033[0m'

force_rebuild="${1:-}"

. scripts/colors.sh
scripts/docker/generate-build-image.sh "${force_rebuild}"

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

printf "${OK_GREEN_COLOR}==> Building ${SIDECAR_BASE_IMAGE_NAME}:${SIDECAR_BASE_IMAGE_TAG} image ..  ${NO_COLOR}\n"
docker build -t ${SIDECAR_BASE_IMAGE_NAME}:${SIDECAR_BASE_IMAGE_TAG} --rm -f scripts/docker/release/Dockerfile-release .

echo ""
echo ""

printf "${OK_BG_COLOR}==> ${SIDECAR_BASE_IMAGE_NAME}:${SIDECAR_BASE_IMAGE_TAG} built successfully ${NO_COLOR}\n"
