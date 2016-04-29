#!/bin/sh -x

OUTPUT_FILE=$1
INSTANCE_ID=$2
CREDENTIALS_ID=$3

NewUsername=`echo ${CREDENTIALS_ID} |  cut -c 1-15`

# get credentials
/catalog-service-manager/bin/amazon-rds-mysql getconnection ${AWS_RDS_REGION} ${MYSQL_RDS_INSTANCE} ${MYSQL_ROOT_USER} ${MYSQL_ROOT_PASSWORD} d${INSTANCE_ID} ${NewUsername} ${OUTPUT_FILE}
