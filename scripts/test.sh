#!/bin/sh
TOPDIR=$(cd "$(dirname "$0")/.." && pwd)

# List of files to be tested
TESTLIST=$(go list ./cmd/... ./src/... | grep -v github.com/go-swagger)

# Manipulate the GOPATH to include both go-swagger and go-openapi from
# the sub-moduled swagger. Unfortunately, go(1) gets _really_ confused when
# symlinks are involved; copying the tree is more likely to do the correct thing
swaggerdir="$(mktemp -d)"
trap "rm -rf '${swaggerdir}'" EXIT
cp -r "${TOPDIR}/src/github.com/go-swagger/go-swagger/vendor" "${swaggerdir}/src"
export GOPATH="${swaggerdir}:${GOPATH}"
go test ${TESTLIST}
