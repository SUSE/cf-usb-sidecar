#!/bin/sh

. ${CSM_ROOT}/scripts/colors.sh

echo "${OK_GREEN_COLOR}==> Setting up go-swagger ${NO_COLOR}"


GOSWAGGER_VERSION=35bf94d48ffdd1eca2f287fb0950668c43650d52

# Installs go-swagger generation tool locked to a specified version
# adding buildutil as swagger will not build without it. It is not vendored
go get golang.org/x/tools/go/buildutil
go get -d -u github.com/go-swagger/go-swagger
cd ${GOPATH}/src/github.com/go-swagger/go-swagger && git reset --hard ${GOSWAGGER_VERSION}
GO15VENDOREXPERIMENT=1 go install github.com/go-swagger/go-swagger/cmd/swagger
