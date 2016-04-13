package handlers

import (
	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager/models"
	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager/restapi/operations/workspace"
	"github.com/hpcloud/catalog-service-manager/src/csm_manager"
)

func CreateWorkspace(createRequest *models.ServiceManagerWorkspaceCreateRequest) middleware.Responder {
	internalWorkspaces := csm_manager.GetWorkspace()
	return workspace.NewCreateWorkspaceCreated().WithPayload(internalWorkspaces.CreateWorkspace(createRequest))
}
