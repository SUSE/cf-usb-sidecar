package handlers

import (
	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/SUSE/cf-usb-sidecar/generated/CatalogServiceManager/restapi/operations/workspace"
	"github.com/SUSE/cf-usb-sidecar/src/csm_manager"
)

func DeleteWorkspace(workspaceID string) middleware.Responder {
	internalWorkspaces := csm_manager.GetWorkspace()
	_, err := internalWorkspaces.DeleteWorkspace(workspaceID)
	if err != nil {
		return workspace.NewDeleteWorkspaceDefault(int(*err.Code)).WithPayload(err)
	}
	return workspace.NewDeleteWorkspaceOK()
}
