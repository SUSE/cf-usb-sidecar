#!/bin/bash

TARGET=$1

. ${SIDECAR_ROOT}/scripts/colors.sh

BUILT_SERVICES=""
FAILED_SERVICES=""

cd  ${SIDECAR_ROOT}/csm-utils/utils
for serviceDir in `find . -maxdepth 1 -mindepth 1 -type d `
do
    pushd ${serviceDir}
    echo "${OK_COLOR} ---> Running 'make ${TARGET}' on ${PWD##*/} ${NO_COLOR} "
    if [ -f Makefile ]
    then
        make $TARGET
        service_command_status=$?
        if [ ${service_command_status} -ne 0 ];
        then
            FAILED_SERVICES="${FAILED_SERVICES}${serviceDir} "
        else
            BUILT_SERVICES="${BUILT_SERVICES}${serviceDir} "
        fi
    else
        echo "Error: Makefile not found for Service ${PWD##*/}"
        exit 1
    fi
    popd
done

echo "${OK_COLOR} ---> Finished building ${BUILT_SERVICES} ${NO_COLOR} "

if [ ! -z "${FAILED_SERVICES}" ]; then
    echo "${ERROR_COLOR} ---> Failed Services ${FAILED_SERVICES} ${NO_COLOR} "
    exit 1
fi
