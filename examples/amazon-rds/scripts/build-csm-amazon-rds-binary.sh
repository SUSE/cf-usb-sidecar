#!/bin/sh -x

current_dir=$(pwd)
CSM_RDS_BIN=${current_dir}/CSM_RDS_BIN
mkdir -p ${CSM_RDS_BIN}
rm -rf ${CSM_RDS_BIN}/catalog-service-manager

docker build -t csm-rds:build -f Dockerfile-build . 
 
docker images | grep csm-rds | grep build > /dev/null 2>&1
if [ $? -eq 0 ]
then
	echo -e "${OK_GREEN_COLOR}==> Copying CMS-Amazon-RDS binary to the host ${NO_COLOR}"
	docker run --name csm-rds-build -v ${CSM_RDS_BIN}:/csm-amazon-rds csm-rds:build
fi