#!/bin/sh -x


OUTPUT_FILE=$1
INSTANCE_ID=$2

Username="root"
Password=${MYSQL_ROOT_PASSWORD}

# creates output file for successful execution
write_success_output () { 
	cat <<EOF > ${OUTPUT_FILE}
{
	"http_code":201,
	"status": "successful",
	"details":"database created",
	"processing_type":"Extension"
}
EOF
}

# creates output file for failed execution
write_failed_output(){
	cat <<EOF > ${OUTPUT_FILE}
{
	"http_code":500,
	"status": "failed",
	"details":"database could not be created",
	"processing_type":"Extension"
}
EOF
}

# creates output file for failed execution
write_exists_output(){
	cat <<EOF > ${OUTPUT_FILE}
{
	"http_code":500,
	"status": "failed",
	"details":"database already exists",
	"processing_type":"Extension"
}
EOF
}

# check if database already exists
mysql -h ${MYSQL_SERVICE_HOST} -P ${MYSQL_SERVICE_PORT_MYSQL} -u ${Username} -p${Password} -e "show databases" | grep "d${INSTANCE_ID}" > /dev/null 2>&1

if [ $? -eq 0 ]; then
    write_exists_output
	return
fi

sleep 1

# create mysql workspace/database 
mysql -h ${MYSQL_SERVICE_HOST} -P ${MYSQL_SERVICE_PORT_MYSQL} -u ${Username} -p${Password} -e "create database d${INSTANCE_ID}" > /dev/null 2>&1

sleep 1

# make sure its created successfully
mysql -h ${MYSQL_SERVICE_HOST} -P ${MYSQL_SERVICE_PORT_MYSQL} -u ${Username} -p${Password} -e "show databases" | grep "d${INSTANCE_ID}" > /dev/null 2>&1

if [ $? -eq 0 ]; then
    write_success_output
else
    write_failed_output
fi
