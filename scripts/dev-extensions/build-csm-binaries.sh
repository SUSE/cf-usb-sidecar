#!/bin/sh

current_dir=$(pwd)
CSM_EXTENSION_BIN_DIR=${current_dir}/CSM_HOME
rm -rf ${CSM_EXTENSION_BIN_DIR}
mkdir -p ${CSM_EXTENSION_BIN_DIR}


docker build -t ${CSM_EXTENSION_BUILD_IMAGE_NAME} --rm -f Dockerfile-build . 
 
docker images | grep ${CSM_EXTENSION_BUILD_IMAGE_NAME} | grep build > /dev/null 2>&1
if [ $? -eq 0 ]
then
	echo "${OK_GREEN_COLOR}==> Copying CMS-EXTENSION binary to the host ${NO_COLOR}"
	docker run \
		--name ${CSM_EXTENSION_BUILD_IMAGE_NAME} \
		-v ${CSM_EXTENSION_BIN_DIR}:/out \
		${CSM_EXTENSION_BUILD_IMAGE_NAME} \
		'cp -v -r /CSM_HOME/* /out'
fi