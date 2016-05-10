#!/bin/sh

OUTPUT_FILE=$1
INSTANCE_ID=$2

Username="root"
Password=${MYSQL_ROOT_PASSWORD}

# creates output file for successful execution
write_success_output () { 
	cat <<EOF > ${OUTPUT_FILE}
{
	"status": "successful",
	"details": {
		"result":"database exists"	
	}
}
EOF
}

# creates output file for failed execution
write_failed_output(){
	cat <<EOF > ${OUTPUT_FILE}
{
	"error_code" : 500, 
	"error_message" : "no database with this name was found", 
	"status": "failed"
}
EOF
}

# check if database/worksapce exists
mysql -h ${MYSQL_SERVICE_HOST} -P ${MYSQL_SERVICE_PORT_MYSQL} -u ${Username} -p${Password} -e "show databases" | grep "d${INSTANCE_ID}" > /dev/null 2>&1

if [ $? -eq 0 ]; then
    write_success_output
else
    write_failed_output
fi
