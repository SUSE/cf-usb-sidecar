#!/bin/sh
TOPDIR=$(cd "$(dirname "$0")/.." && pwd)

# List of files to be tested
TESTLIST=$(go list ./cmd/... ./src/... | grep -v github.com/go-swagger)

# Manipulate the GOPATH to include both go-swagger and go-openapi from
# the sub-moduled swagger for succesful compilation of test files.
GOPATH="${TOPDIR}:${TOPDIR}/go-swagger:${GOPATH}" godep go test ${TESTLIST}
