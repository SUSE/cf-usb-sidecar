#!/bin/sh -x

OUTPUT_FILE=$1

#delete mysql database instance
/catalog-service-manager/bin/amazon-rds-mysql deletedb ${AWS_RDS_REGION} ${MYSQL_RDS_INSTANCE} ${OUTPUT_FILE}
