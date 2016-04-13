package handlers

import (
	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager/models"
	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager/restapi/operations/connection"
	"github.com/hpcloud/catalog-service-manager/src/csm_manager"
)

func CreateConnection(workspaceID string, connectionRequest *models.ServiceManagerConnectionCreateRequest) middleware.Responder {
	internalConnection := csm_manager.GetConnection()
	return connection.NewCreateConnectionCreated().WithPayload(internalConnection.CreateConnection(workspaceID, connectionRequest))
}
