#!/bin/sh

if [ -z "${SERVICE_MYSQL_HOST}" ]
then
    if [ -z "${MYSQL_INT_SERVICE_PORT}" ]
    then
      export SERVICE_MYSQL_HOST="mysql.${KUBE_SERVICE_DOMAIN_SUFFIX}"
	else
	  export SERVICE_MYSQL_HOST="mysql-int.${KUBE_SERVICE_DOMAIN_SUFFIX}"
	fi
fi

if [ -z "${SERVICE_MYSQL_PORT}" ]
  then
    if [ -z "${MYSQL_INT_SERVICE_PORT}" ]
      then
       export SERVICE_MYSQL_PORT=${MYSQL_SERVICE_PORT}
      else
       export SERVICE_MYSQL_PORT=${MYSQL_INT_SERVICE_PORT}
    fi
fi

if [ -z "${SERVICE_MYSQL_PASS}" ]
  then
    export SERVICE_MYSQL_PASS=${MYSQL_ROOT_PASSWORD}
fi

echo "Starting catalog-service-manager ..."

/catalog-service-manager/bin/catalog-service-manager
