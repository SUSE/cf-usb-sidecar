#!/bin/sh

export GO15VENDOREXPERIMENT=1

docker images | grep ${SIDECAR_MYSQL_IMAGE_NAME} | grep ${SIDECAR_MYSQL_IMAGE_TAG}
if [ $? -ne 0 ]
then
    echo "Error: Please run 'make build-image' to build docker image for ${SIDECAR_MYSQL_IMAGE_NAME}"
    exit 1
fi

if [ ! -z ${DOCKER_HOST} ]
then
    export DOCKER_HOST_IP=`echo ${DOCKER_HOST} | cut -d "/" -f 3 | cut -d ":" -f 1`
else
    export DOCKER_HOST_IP=`ip route get 8.8.8.8 | awk 'NR==1 {print $NF}'`
fi

cd ${SIDECAR_ROOT}
make generate

cd ${SIDECAR_MYSQL_EXTENSION_ROOT}
# remove existing containers
make clean-containers
# start mysql container
make tools
make run

checkServer(){
    n=0
    until [ $n -ge 6 ]
    do
        nc -w 2 ${DOCKER_HOST_IP} 8081 && break  # substitute your command here
        #n=$[$n+1]
        n=$(expr $n + 1)
        sleep 10
    done

}
testStatus=1
checkServer

nc -w 2 ${DOCKER_HOST_IP} 8081

if [ $? -eq 0 ]
then
    cd ${SIDECAR_MYSQL_EXTENSION_ROOT}/tests
    go test ./... -integration=true -host=${DOCKER_HOST_IP} -port="8081" -v
    testStatus=$?
    cd ${SIDECAR_MYSQL_EXTENSION_ROOT}
    make clean-containers
    sleep 5 #Wait for server to die
fi
exit $testStatus

