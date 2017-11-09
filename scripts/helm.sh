#!/bin/sh
TOPDIR=$(cd "$(dirname "$0")/.." && pwd)

. ${TOPDIR}/scripts/colors.sh

if [ -z ${APP_VERSION_TAG} ]; then
    printf "${ERROR_COLOR}Error${NO_COLOR}: Please set environment variable APP_VERSION_TAG\n"
    exit 1
fi

rm -rf   output
mkdir -p output/helm

# First remove the dev-internal database role from the chart and add
# the post-deployment task in its stead.

# The p-d-s is configured with the location of an SCF to talk to and
# registers the mysql-sidecar with it as a service, via the universal
# service broker.

cp -rf chart/* output/helm/
mv     output/helm/other/* output/helm/templates/
rmdir  output/helm/other
rm     output/helm/templates/db.yaml
rm     output/helm/*.sh

# Fix the version information in the image references

for path in output/helm/templates/*
do
    sed < $path > $$ -e "s/:latest/:${APP_VERSION_TAG}/"
    cp $$ $path
done
