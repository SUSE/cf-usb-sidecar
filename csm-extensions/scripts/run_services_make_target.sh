#!/bin/sh

TARGET=$1
CSM_ROOT=$GOPATH/src/github.com/hpcloud/catalog-service-manager
CSM_SERVICES=$GOPATH/src/github.com/hpcloud/catalog-service-manager/csm-extensions/services

. ${CSM_ROOT}/scripts/colors.sh


cd  $GOPATH/src/github.com/hpcloud/catalog-service-manager/csm-extensions/services
for serviceDir in `find . -maxdepth 1 -mindepth 1 -type d `
do
    oldDir=`pwd`
    cd ${serviceDir}
    echo "${OK_COLOR} ---> Running 'make ${TARGET}' on ${PWD##*/} ${NO_COLOR} "
    if [ -f Makefile ]
    then
        make $TARGET
    else
        echo "Error: Makefile not found for Service ${PWD##*/}"
        exit 1
    fi
    cd $oldDir
done