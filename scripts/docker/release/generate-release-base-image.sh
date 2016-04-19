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

CSM_BIN=${current_dir}/CSM_BIN
mkdir -p ${CSM_BIN}
rm -rf ${CSM_BIN}/catalog-service-manager

echo "${OK_BG_COLOR}==> Building release-build image (to build catalog-service-manager binary).. ${NO_COLOR}"

docker build -t catalog-service-manager:release-build -f scripts/docker/release/Dockerfile-release-build .  

docker images | grep catalog-service-manager | grep release-build > /dev/null 2>&1
if [ $? -eq 0 ]
then
	echo "${OK_GREEN_COLOR}==> Copying CMS binary to the host ${NO_COLOR}"
	docker run --name csm-release-build -v ${CSM_BIN}:/catalog-service-manager-bin/ catalog-service-manager:release-build 
fi

if [ -f ${CSM_BIN}/catalog-service-manager ]
then
	echo "${OK_BG_COLOR}==> CSM binary is built successfuly ${NO_COLOR}"
	echo ""
	
	echo "${OK_GREEN_COLOR}==> Removing csm-release-build container  ${NO_COLOR}"
	docker ps -a | grep csm-release-build | awk '{print $1}' | xargs docker rm 
	
	echo "${OK_GREEN_COLOR}==> Removing catalog-service-manager:release-build images  ${NO_COLOR}"
	docker images | grep catalog-service-manager | grep release-build | awk '{print $3}' | xargs -L 1 docker rmi -f 
	
	echo "${OK_GREEN_COLOR}==> Removing old catalog-service-manager:release images  ${NO_COLOR}"
	docker images | grep catalog-service-manager | grep release | awk '{print $3}' | xargs -L 1 docker rmi -f 
	
	sleep 5
	echo ""
	echo "${OK_GREEN_COLOR}==> Building catalog-service-manager:release image ..  ${NO_COLOR}"
	docker build -t catalog-service-manager:release -f scripts/docker/release/Dockerfile-release . 
	
	echo ""
	echo ""
	
	echo "${OK_BG_COLOR}==> catalog-service-manager:release built successfuly ${NO_COLOR}"

	docker tag catalog-service-manager:release catalog-service-manager:base
	echo "${OK_BG_COLOR}==> catalog-service-manager:base tag is published ${NO_COLOR}"

fi
