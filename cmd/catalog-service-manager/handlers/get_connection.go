package handlers

import (
	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/SUSE/cf-usb-sidecar/generated/CatalogServiceManager/restapi/operations/connection"
	"github.com/SUSE/cf-usb-sidecar/src/csm_manager"
)

func GetConnection(workspaceID string, connectionID string) middleware.Responder {
	internalConnection := csm_manager.GetConnection()
	conn, err := internalConnection.GetConnection(workspaceID, connectionID)
	if err != nil {
		return connection.NewGetConnectionDefault(int(*err.Code)).WithPayload(err)
	}

	return connection.NewGetConnectionOK().WithPayload(conn)
}
