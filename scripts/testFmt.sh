#!/bin/sh

NO_COLOR="\033[0m"
ERROR_COLOR="\033[31;01m"

if [ $# = 0 ]; then
    printf "${ERROR_COLOR}Usage: testFmt.sh <files-to-check>${NO_COLOR}\n"
    exit 1
fi
# check go format on files
unformatted_files=$(gofmt -l $1)
[ -z "$unformatted_files" ] && exit 0

# show how to fix the unformatted files.
echo Run
for file in $unformatted_files; do
    echo "  gofmt -d $file"
done
echo to see details.
printf "${ERROR_COLOR}Files are not go fmt compliant.${NO_COLOR}\n"
exit 1
