package handlers

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/SUSE/cf-usb-sidecar/generated/CatalogServiceManager/models"
	"github.com/SUSE/cf-usb-sidecar/generated/CatalogServiceManager/restapi/operations/connection"
	"github.com/SUSE/cf-usb-sidecar/src/csm_manager"
)

func CreateConnection(workspaceID string, connectionRequest *models.ServiceManagerConnectionCreateRequest) middleware.Responder {
	internalConnection := csm_manager.GetConnection()
	conn, err := internalConnection.CreateConnection(workspaceID, *connectionRequest.ConnectionID, connectionRequest.Details)
	if err != nil {
		return connection.NewCreateConnectionDefault(int(err.Code)).WithPayload(err)
	}

	return connection.NewCreateConnectionCreated().WithPayload(conn)
}
