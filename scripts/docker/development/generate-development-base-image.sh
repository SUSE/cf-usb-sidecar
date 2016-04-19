#!/bin/sh

OK_COLOR='\033[1;32m'
OK_GREEN_COLOR='\033[0;32m'
OK_BG_COLOR='\033[42m'
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

if [ "${script_dir}" != "scripts/docker/development" ]
then
	echo "${ERROR_COLOR}==> Script directory is not correct, please run script from \${PROJECT_ROOT}/scripts/docker/development ${NO_COLOR}"
	exit 1	
fi

docker images | grep catalog-service-manager | grep development > /dev/null 2>&1
if [ $? -eq 0 ]
then
	if [ "$force_rebuild" != "rebuild-image" ]
	then
		echo "catalog-service-manager:development already exists!"
		exit 0
	fi
	
	if [ "$force_rebuild" != "rebuild-image" ]
	then
		docker images | grep catalog-service-manager | grep development | awk '{print $3}' | xargs -L 1 docker rmi -f > /dev/null 2>&1
		sleep 5
	fi	
fi

docker build -t catalog-service-manager:development -f ${script_dir}/Dockerfile-development .
docker tag catalog-service-manager:development catalog-service-manager:base

echo "${OK_BG_COLOR}==> catalog-service-manager:base tag is published ${NO_COLOR}"
