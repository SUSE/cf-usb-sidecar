#!/bin/sh

current_dir=$(pwd)
SIDECAR_EXTENSION_BIN_DIR=${current_dir}/SIDECAR_HOME
rm -rf ${SIDECAR_EXTENSION_BIN_DIR}
mkdir -p ${SIDECAR_EXTENSION_BIN_DIR}


docker build -t ${SIDECAR_EXTENSION_BUILD_IMAGE_NAME} --rm -f Dockerfile-build . 
 
docker images | grep ${SIDECAR_EXTENSION_BUILD_IMAGE_NAME} | grep build > /dev/null 2>&1
if [ $? -eq 0 ]
then
	echo "${OK_GREEN_COLOR}==> Copying CMS-EXTENSION binary to the host ${NO_COLOR}"
	docker run \
		--name ${SIDECAR_EXTENSION_BUILD_IMAGE_NAME} \
		-v ${SIDECAR_EXTENSION_BIN_DIR}:/out \
		${SIDECAR_EXTENSION_BUILD_IMAGE_NAME} \
		'cp -v -r /SIDECAR_HOME/* /out'
fi