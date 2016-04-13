package handlers

import (
	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager/restapi/operations/connection"
	"github.com/hpcloud/catalog-service-manager/src/csm_manager"
)

func GetConnection(workspaceID string, connectionID string) middleware.Responder {
	internalConnection := csm_manager.GetConnection()
	return connection.NewGetConnectionOK().WithPayload(internalConnection.GetConnection(workspaceID, connectionID))
}
