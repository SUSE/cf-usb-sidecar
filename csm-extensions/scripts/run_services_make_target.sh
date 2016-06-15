#!/bin/sh

TARGET=$1
CSM_ROOT=$GOPATH/src/github.com/hpcloud/catalog-service-manager
CSM_SERVICES=$GOPATH/src/github.com/hpcloud/catalog-service-manager/csm-extensions/services

. ${CSM_ROOT}/scripts/colors.sh

BUILT_SERVICES=""
FAILED_SERVICES=""
command_status=0

cd  $GOPATH/src/github.com/hpcloud/catalog-service-manager/csm-extensions/services
for serviceDir in `find . -maxdepth 1 -mindepth 1 -type d `
do
    oldDir=`pwd`
    cd ${serviceDir}
    echo "${OK_COLOR} ---> Running 'make ${TARGET}' on ${PWD##*/} ${NO_COLOR} "
    if [ -f Makefile ]
    then
        make $TARGET
        service_command_status=$?
        if [ ${service_command_status} -ne 0 ];
        then
            command_status=${service_command_status}
            FAILED_SERVICES="${FAILED_SERVICES}${serviceDir} "
        else
            BUILT_SERVICES="${BUILT_SERVICES}${serviceDir} "
        fi
    else
        echo "Error: Makefile not found for Service ${PWD##*/}"
        exit 1
    fi
    cd $oldDir
done

echo "${OK_COLOR} ---> Finished building ${BUILT_SERVICES} ${NO_COLOR} "

if [ ${command_status} -ne 0 ]; then
echo "${ERROR_COLOR} ---> Failed Services ${FAILED_SERVICES} ${NO_COLOR} "

fi