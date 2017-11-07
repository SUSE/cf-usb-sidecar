package handlers

import (
	"github.com/go-openapi/runtime/middleware"

	"github.com/SUSE/cf-usb-sidecar/src/csm_manager"
	"github.com/SUSE/cf-usb-sidecar/generated/CatalogServiceManager/restapi/operations/workspace"
)

func GetWorkspace(workspaceID string) middleware.Responder {
	internalWorkspaces := csm_manager.GetWorkspace()
	wksp, err := internalWorkspaces.GetWorkspace(workspaceID)
	if err != nil {
		return workspace.NewGetWorkspaceDefault(int(err.Code)).WithPayload(err)
	}
	return workspace.NewGetWorkspaceOK().WithPayload(wksp)
}
