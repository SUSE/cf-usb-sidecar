#!/bin/sh

. ${SIDECAR_ROOT}/scripts/colors.sh

if [ -d generated/CatalogServiceManager ]
then
	echo "${WARN_MAGENTA}==> Generated code found @ generated/CatalogServiceManager ${NO_COLOR}"
else
	echo "${OK_GREEN_COLOR}==> Calling swagger generate service @ generated/CatalogServiceManager ${NO_COLOR}"
	swagger generate server -A CatlogServiceManager -t generated/CatalogServiceManager -f docs/swagger-spec/api.yml
fi
