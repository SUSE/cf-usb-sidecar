#!/bin/bash

pushd ${CSM_EXTENSION_ROOT}/db
 docker build -t ${CSM_EXTENSION_SVC_IMAGE_NAME}:${CSM_EXTENSION_SVC_IMAGE_TAG} .
popd