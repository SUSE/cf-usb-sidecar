#!/bin/sh

go build -v -o ${SIDECAR_EXTENSION_ROOT}/SIDECAR_HOME/connection/create/create	./cmd/createconnection
go build -v -o ${SIDECAR_EXTENSION_ROOT}/SIDECAR_HOME/connection/delete/delete	./cmd/deleteconnection
go build -v -o ${SIDECAR_EXTENSION_ROOT}/SIDECAR_HOME/connection/get/get		./cmd/getconnection
go build -v -o ${SIDECAR_EXTENSION_ROOT}/SIDECAR_HOME/workspace/create/create	./cmd/createworkspace
go build -v -o ${SIDECAR_EXTENSION_ROOT}/SIDECAR_HOME/workspace/delete/delete	./cmd/deleteworkspace
go build -v -o ${SIDECAR_EXTENSION_ROOT}/SIDECAR_HOME/workspace/get/get		./cmd/getworkspace
go build -v -o ${SIDECAR_EXTENSION_ROOT}/SIDECAR_HOME/status/status             ./cmd/status

mkdir -p ${SIDECAR_EXTENSION_ROOT}/SIDECAR_HOME/bin/ && cp ./scripts/run.sh ${SIDECAR_EXTENSION_ROOT}/SIDECAR_HOME/bin/
