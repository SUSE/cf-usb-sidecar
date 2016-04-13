#!/bin/sh -x

OUTPUT_FILE=$1

# creates output file for successful execution
write_success_output () { 
	cat <<EOF > ${OUTPUT_FILE}
{
	"status": "successful"
}
EOF
}

# creates output file for failed execution
write_failed_output(){
	cat <<EOF > ${OUTPUT_FILE}
{
	"status": "failed"
}
EOF
}


# this is NOOP extension
write_success_output