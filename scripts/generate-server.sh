#!/bin/sh

. ${SIDECAR_ROOT}/scripts/colors.sh

if [ -d generated/CatalogServiceManager ]
then
	printf "${WARN_MAGENTA}==> Generated code found @ generated/CatalogServiceManager ${NO_COLOR}"
	exit 0
fi

printf "${OK_GREEN_COLOR}==> Calling swagger generate service @ generated/CatalogServiceManager ${NO_COLOR}\n"
mkdir -p generated/CatalogServiceManager

# Locate ourselves
TOPDIR=$(cd "$(dirname "$0")/.." && pwd)

# And run the sub-moduled swagger (which comes with a vendor matching go-openapi).
GOPATH="${TOPDIR}:${GOPATH}" go run \
    "${TOPDIR}"/src/github.com/go-swagger/go-swagger/cmd/swagger/swagger.go \
    generate server \
    -A CatlogServiceManager \
    -t generated/CatalogServiceManager \
    -f docs/swagger-spec/api.yml
