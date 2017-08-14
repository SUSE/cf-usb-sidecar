package handlers

import (
	"github.com/SUSE/cf-usb-sidecar/generated/CatalogServiceManager/models"
	"github.com/SUSE/cf-usb-sidecar/generated/CatalogServiceManager/restapi/operations/workspace"
	"github.com/SUSE/cf-usb-sidecar/src/csm_manager"
	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"
)

func CreateWorkspace(createRequest *models.ServiceManagerWorkspaceCreateRequest) middleware.Responder {
	internalWorkspaces := csm_manager.GetWorkspace()
	wksp, err := internalWorkspaces.CreateWorkspace(createRequest.WorkspaceID, createRequest.Details)
	if err != nil {
		return workspace.NewCreateWorkspaceDefault(int(*err.Code)).WithPayload(err)
	}
	return workspace.NewCreateWorkspaceCreated().WithPayload(wksp)
}
