package handlers

import (
	"github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/hpcloud/catalog-service-manager/src/csm_manager"

	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager/restapi/operations/workspace"
)

func GetWorkspace(workspaceID string) middleware.Responder {
	internalWorkspaces := csm_manager.GetWorkspace()
	return workspace.NewGetWorkspaceOK().WithPayload(internalWorkspaces.GetWorkspace(workspaceID))
}
