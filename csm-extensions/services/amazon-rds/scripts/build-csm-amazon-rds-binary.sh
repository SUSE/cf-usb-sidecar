#!/bin/sh

current_dir=$(pwd)
CSM_RDS_BIN=${current_dir}/CSM_RDS_BIN
mkdir -p ${CSM_RDS_BIN}
rm -rf ${CSM_RDS_BIN}/catalog-service-manager

docker build -t ${CSM_RDS_BUILD_IMAGE_NAME} --rm -f Dockerfile-build . 
 
docker images | grep ${CSM_RDS_BUILD_IMAGE_NAME} | grep build > /dev/null 2>&1
if [ $? -eq 0 ]
then
	echo "${OK_GREEN_COLOR}==> Copying CMS-Amazon-RDS binary to the host ${NO_COLOR}"
	docker run --name ${CSM_RDS_BUILD_IMAGE_NAME} -v ${CSM_RDS_BIN}:/csm-amazon-rds ${CSM_RDS_BUILD_IMAGE_NAME}
fi