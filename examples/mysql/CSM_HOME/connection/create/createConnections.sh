#!/bin/sh

OUTPUT_FILE=$1
INSTANCE_ID=$2
CREDENTIALS_ID=$3

Username="root"
Password=${MYSQL_ROOT_PASSWORD}
if [ -f /sbin/md5 ]
then
	NewPassword=`date | md5 -r | head -c 20`
else
	NewPassword=`date | md5sum | head -c 20`
fi

write_success_output () {
	cat <<EOF > ${OUTPUT_FILE}
{
	"status": "successful",
	"details": {
		"host" : "${MYSQL_SERVICE_HOST}",
		"port" : "${MYSQL_SERVICE_PORT_MYSQL}",
		"username" : "d${NewUsername}",
		"password" : "${NewPassword}",
		"database" : "d${INSTANCE_ID}"
	}
}
EOF

}

write_failed_output(){
	cat <<EOF > ${OUTPUT_FILE}
{
	"status": "failed"
}
EOF

}

NewUsername=`echo ${CREDENTIALS_ID} |  cut -c 1-15`

mysql -h ${MYSQL_SERVICE_HOST} -P ${MYSQL_SERVICE_PORT_MYSQL} -u ${Username} -p${Password} -e "show databases" | grep "d${INSTANCE_ID}" > /dev/null 2>&1

if [ $? -eq 0 ]; then
	mysql -h ${MYSQL_SERVICE_HOST} -P ${MYSQL_SERVICE_PORT_MYSQL} -u ${Username} -p${Password} -e "CREATE USER \"d${NewUsername}\"@\"%\" IDENTIFIED BY \"${NewPassword}\";FLUSH PRIVILEGES" > /dev/null 2>&1

	if [ $? -ne 0 ]; then
		write_failed_output
		exit 0
	fi

	sleep 1

	mysql -h ${MYSQL_SERVICE_HOST} -P ${MYSQL_SERVICE_PORT_MYSQL} -u ${Username} -p${Password} -e "select * from mysql.user where User=\"d${NewUsername}\";" | grep ${NewUsername} > /dev/null 2>&1
	if [ $? -eq 0 ]; then
		mysql -h ${MYSQL_SERVICE_HOST} -P ${MYSQL_SERVICE_PORT_MYSQL} -u ${Username} -p${Password} -e "GRANT ALL PRIVILEGES ON d${INSTANCE_ID}.* TO \"d${NewUsername}\"@\"%\";FLUSH PRIVILEGES" > /dev/null 2>&1

		if [ $? -eq 0 ]; then
			write_success_output
		else
			write_failed_output
		fi
	else
		write_failed_output
	fi
else
    write_failed_output
fi
