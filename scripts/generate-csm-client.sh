#!/bin/sh

. ${CSM_ROOT}/scripts/colors.sh

if [ -d generated/CatalogServiceManager-client ]
then
	echo "${WARN_MAGENTA}==> Generated code found @ generated/CatalogServiceManager-client ${NO_COLOR}"
else
	echo "${OK_GREEN_COLOR}==> Calling swagger generate service @ generated/CatalogServiceManager-client ${NO_COLOR}"
	swagger generate client -A CatlogServiceManager -t generated/CatalogServiceManager-client -f docs/swagger-spec/api.yml
fi
