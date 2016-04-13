package handlers

import (
	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager/restapi/operations/workspace"
	"github.com/hpcloud/catalog-service-manager/src/csm_manager"
)

func DeleteWorkspace(workspaceID string) middleware.Responder {
	internalWorkspaces := csm_manager.GetWorkspace()
	internalWorkspaces.DeleteWorkspace(workspaceID)
	return workspace.NewDeleteWorkspaceOK()
}
