#!/bin/sh

current_dir=$(pwd)
CSM_RDS_BIN=${current_dir}/CSM_RDS_BIN
mkdir -p ${CSM_RDS_BIN}
rm -rf ${CSM_RDS_BIN}/catalog-service-manager

docker build -t ${CSM_EXTENSION_BUILD_IMAGE_NAME} --rm -f Dockerfile-build . 
 
docker images | grep ${CSM_EXTENSION_BUILD_IMAGE_NAME} | grep build > /dev/null 2>&1
if [ $? -eq 0 ]
then
	echo "${OK_GREEN_COLOR}==> Copying CMS-Amazon-RDS binary to the host ${NO_COLOR}"
	DOCKER_CONTAINER_ID=`docker run -d --name ${CSM_EXTENSION_BUILD_IMAGE_NAME} -v ${CSM_RDS_BIN}:/csm-amazon-rds ${CSM_EXTENSION_BUILD_IMAGE_NAME}`
	
	sleep 5

	docker rm ${DOCKER_CONTAINER_ID}

fi

