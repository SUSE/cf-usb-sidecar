#!/bin/sh
TOPDIR=$(cd "$(dirname "$0")/.." && pwd)

# Manipulate the GOPATH to include both go-swagger and go-openapi from
# the sub-moduled swagger.
GOPATH="${TOPDIR}:${TOPDIR}/v:${GOPATH}" godep go run cmd/catalog-service-manager/catalog-service-manager.go
