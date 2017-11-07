#!/bin/sh
TOPDIR=$(cd "$(dirname "$0")/.." && pwd)

# List of files to be tested
TESTLIST=$(go list ./... | \
    grep -v examples | \
    grep -v services | \
    grep -v generated | \
    grep -v github.com/go-swagger | \
    grep -v v/src | \
    grep -v scripts | \
    grep -v SIDECAR_extensions)

# Manipulate the GOPATH to include both go-swagger and go-openapi from
# the sub-moduled swagger for succesful compilation of test files.
GOPATH="${TOPDIR}:${TOPDIR}/v:${GOPATH}" godep go test ${TESTLIST}
