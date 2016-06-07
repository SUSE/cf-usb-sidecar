#!/bin/sh

echo "Using ${CSM_EXTENSION_ROOT}"
go build -v -o ${CSM_EXTENSION_ROOT}/CSM_HOME/connection/create/create	./cmd/createconnection
go build -v -o ${CSM_EXTENSION_ROOT}/CSM_HOME/connection/delete/delete	./cmd/deleteconnection
go build -v -o ${CSM_EXTENSION_ROOT}/CSM_HOME/connection/get/get		./cmd/getconnection
go build -v -o ${CSM_EXTENSION_ROOT}/CSM_HOME/workspace/create/create	./cmd/createworkspace
go build -v -o ${CSM_EXTENSION_ROOT}/CSM_HOME/workspace/delete/delete	./cmd/deleteworkspace
go build -v -o ${CSM_EXTENSION_ROOT}/CSM_HOME/workspace/get/get		./cmd/getworkspace
go build -v -o ${CSM_EXTENSION_ROOT}/CSM_HOME/status/status             ./cmd/status

mkdir -p ${CSM_EXTENSION_ROOT}/CSM_HOME/bin/ && cp ./scripts/run.sh ${CSM_EXTENSION_ROOT}/CSM_HOME/bin/