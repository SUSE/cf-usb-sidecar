package handlers

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/SUSE/cf-usb-sidecar/generated/CatalogServiceManager/restapi/operations/connection"
	"github.com/SUSE/cf-usb-sidecar/src/csm_manager"
)

func DeleteConnection(workspaceID string, connectionID string) middleware.Responder {
	internalConnection := csm_manager.GetConnection()
	_, err := internalConnection.DeleteConnection(workspaceID, connectionID)
	if err != nil {
		return connection.NewDeleteConnectionDefault(int(err.Code)).WithPayload(err)
	}

	return connection.NewDeleteConnectionOK()
}
