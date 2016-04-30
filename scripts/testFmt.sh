#!/bin/sh

NO_COLOR="\033[0m"
ERROR_COLOR="\033[31;01m"

if [ $# = 0 ]; then
    echo "${ERROR_COLOR}Usage: testFmt.sh <directory-to-check>${NO_COLOR}"
    exit 1
fi
# check go format on files
unformatted_files=$(gofmt -l $1)
[ -z "$unformatted_files" ] && exit 0

# show how to fix the unformatted files.
for file in $unformatted_files; do
    echo "  go fmt $file"
done
echo "${ERROR_COLOR}Files are not go fmt compliant.${NO_COLOR}"
exit 1
