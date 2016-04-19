#!/bin/sh

. scripts/colors.sh

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

docker images | grep catalog-service-manager | grep build > /dev/null 2>&1
if [ $? -eq 0 ]
then
	if [ "$force_rebuild" != "rebuild-image" ]
	then
		echo "${WARN_MAGENTA}==> catalog-service-manager:build image already exists!${NO_COLOR}"
		exit 0
	fi
	
	if [ "$force_rebuild" == "rebuild-image" ]
	then
		echo "${OK_GREEN_COLOR}==> Removing old catalog-service-manager:build images ..${NO_COLOR}"
		docker images | grep catalog-service-manager | grep build | awk '{print $3}' | xargs -L 1 docker rmi -f > /dev/null 2>&1
		sleep 5
	fi	
fi

echo "${OK_GREEN_COLOR}==> Building catalog-service-manager:build image ..${NO_COLOR}"
docker build -t catalog-service-manager:build -f ${script_dir}/Dockerfile-build .
