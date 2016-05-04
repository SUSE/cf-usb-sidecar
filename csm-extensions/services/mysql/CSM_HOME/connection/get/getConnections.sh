#!/bin/sh -x


OUTPUT_FILE=$1
INSTANCE_ID=$2
CREDENTIALS_ID=$3

Username="root"
Password=${MYSQL_ROOT_PASSWORD}
NewUsername=`echo ${CREDENTIALS_ID} |  cut -c 1-15`

# creates output file for successful execution
write_success_output () { 
	cat <<EOF > ${OUTPUT_FILE}
{
	"http_code":200,
	"status": "successful",
	"details":"user exists",
	"processing_type":"Extension"
}
EOF
}

# creates output file for failed execution
write_failed_output(){
	cat <<EOF > ${OUTPUT_FILE}
{
	"http_code" : 500, 
	"details" : "$1", 
	"status": "failed",
	"processing_type" : "Extension"
}
EOF
}


# check if workspace/database exist
mysql -h ${MYSQL_SERVICE_HOST} -P ${MYSQL_SERVICE_PORT_MYSQL} -u ${Username} -p${Password} -e "show databases" | grep "d${INSTANCE_ID}" > /dev/null 2>&1

if [ $? -eq 0 ]; then
	# check if user exists
	mysql -h ${MYSQL_SERVICE_HOST} -P ${MYSQL_SERVICE_PORT_MYSQL} -u ${Username} -p${Password} -e "select * from mysql.user where User=\"d${NewUsername}\";" | grep ${NewUsername} > /dev/null 2>&1
	if [ $? -eq 0 ]; then
		write_success_output
		exit 0
	else
		write_failed_output "User not found"
		exit 0
	fi
fi
write_failed_output "Wrong database"
