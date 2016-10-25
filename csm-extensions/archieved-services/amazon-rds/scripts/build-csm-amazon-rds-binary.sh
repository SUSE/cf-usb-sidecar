#!/bin/sh

current_dir=$(pwd)
SIDECAR_RDS_BIN=${current_dir}/SIDECAR_RDS_BIN
mkdir -p ${SIDECAR_RDS_BIN}
rm -rf ${SIDECAR_RDS_BIN}/catalog-service-manager

docker build -t ${SIDECAR_EXTENSION_BUILD_IMAGE_NAME} --rm -f Dockerfile-build . 
 
docker images | grep ${SIDECAR_EXTENSION_BUILD_IMAGE_NAME} | grep build > /dev/null 2>&1
if [ $? -eq 0 ]
then
	echo "${OK_GREEN_COLOR}==> Copying CSM-Amazon-RDS binary to the host ${NO_COLOR}"
	DOCKER_CONTAINER_ID=`docker run -d --name ${SIDECAR_EXTENSION_BUILD_IMAGE_NAME} -v ${SIDECAR_RDS_BIN}:/csm-amazon-rds ${SIDECAR_EXTENSION_BUILD_IMAGE_NAME}`
	
	sleep 5

	docker rm ${DOCKER_CONTAINER_ID}

fi

