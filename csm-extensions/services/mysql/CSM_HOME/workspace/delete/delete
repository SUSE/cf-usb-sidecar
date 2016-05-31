#!/bin/sh -x


OUTPUT_FILE=$1
INSTANCE_ID=$2

Username="root"
Password=${MYSQL_ROOT_PASSWORD}

# creates output file for successful execution
write_success_output () { 
	cat <<EOF > ${OUTPUT_FILE}
{
	"status": "successful",
	"details":{
		"result" : "database deleted"
	}
}
EOF
}

# creates output file for failed execution
write_failed_output(){
	cat <<EOF > ${OUTPUT_FILE}
{
	"error_code" : 500, 
	"error_message" : "$1", 
	"status": "failed"
}
EOF
}

mysql -h ${MYSQL_SERVICE_HOST} -P ${MYSQL_SERVICE_PORT_MYSQL} -u ${Username} -p${Password} -e "show databases" | grep "d${INSTANCE_ID}" > /dev/null 2>&1

if [ $? -ne 0 ]; then
	write_failed_output "Database does not exist"
	return
fi

# delete workspace/database
mysql -h ${MYSQL_SERVICE_HOST} -P ${MYSQL_SERVICE_PORT_MYSQL} -u ${Username} -p${Password} -e "drop database d${INSTANCE_ID}" > /dev/null 2>&1
	
sleep 1
	
# make sure database is deleted
mysql -h ${MYSQL_SERVICE_HOST} -P ${MYSQL_SERVICE_PORT_MYSQL} -u ${Username} -p${Password} -e "show databases" | grep "d${INSTANCE_ID}" > /dev/null 2>&1
	
if [ $? -ne 0 ]; then
	# delete was successful as database is not found
	write_success_output
else
	# database still exists, something went wrong
	write_failed_output "Could not delete the database"
fi
