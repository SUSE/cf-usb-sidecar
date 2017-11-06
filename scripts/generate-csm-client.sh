#!/bin/sh

. ${SIDECAR_ROOT}/scripts/colors.sh

if [ -d generated/CatalogServiceManager-client ]
then
	printf "${WARN_MAGENTA}==> Generated code found @ generated/CatalogServiceManager-client ${NO_COLOR}\n"
	exit 0
fi

printf "${OK_GREEN_COLOR}==> Calling swagger generate service @ generated/CatalogServiceManager-client ${NO_COLOR}\n"
mkdir -p generated/CatalogServiceManager-client

# Locate ourselves
TOPDIR=$(cd "$(dirname "$0")/.." && pwd)

# And run the sub-moduled swagger (which comes with a vendor matching go-openapi).
GOPATH="${TOPDIR}:${GOPATH}" go run \
    "${TOPDIR}"/src/github.com/go-swagger/go-swagger/cmd/swagger/swagger.go \
    generate client \
    -A CatlogServiceManager \
    -t generated/CatalogServiceManager-client \
    -f docs/swagger-spec/api.yml
