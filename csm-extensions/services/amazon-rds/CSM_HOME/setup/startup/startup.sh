#!/bin/sh -x

OUTPUT_FILE=$1

#create mysql database instance 
/catalog-service-manager/bin/amazon-rds-mysql createdb ${AWS_RDS_REGION} ${MYSQL_RDS_INSTANCE} ${MYSQL_DB_CLASS} ${DB_SIZE} ${MULTIAZ} ${MYSQL_ROOT_USER} ${MYSQL_ROOT_PASSWORD} ${OUTPUT_FILE}