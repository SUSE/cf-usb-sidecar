package handlers

import (
	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"
	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager/models"
	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager/restapi/operations/workspace"
	"github.com/hpcloud/catalog-service-manager/src/csm_manager"
)

func CreateWorkspace(createRequest *models.ServiceManagerWorkspaceCreateRequest) middleware.Responder {
	internalWorkspaces := csm_manager.GetWorkspace()
	wksp, err := internalWorkspaces.CreateWorkspace(createRequest.WorkspaceID, createRequest.Details)
	if err != nil {
		return workspace.NewCreateWorkspaceDefault(int(*err.Code)).WithPayload(err)
	}
	return workspace.NewCreateWorkspaceCreated().WithPayload(wksp)
}
