#!/bin/sh -x


OUTPUT_FILE=$1
INSTANCE_ID=$2

#delete workspace
/catalog-service-manager/bin/amazon-rds-mysql deleteworkspace ${AWS_RDS_REGION} ${MYSQL_RDS_INSTANCE} ${MYSQL_ROOT_USER} ${MYSQL_ROOT_PASSWORD} d${INSTANCE_ID} ${OUTPUT_FILE}
