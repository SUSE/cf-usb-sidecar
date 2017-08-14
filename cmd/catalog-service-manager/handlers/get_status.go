package handlers

import (
	"github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/SUSE/cf-usb-sidecar/src/csm_manager"

	"github.com/SUSE/cf-usb-sidecar/generated/CatalogServiceManager/restapi/operations/status"
)

func GetStatus() middleware.Responder {

	internalWorkspaces := csm_manager.GetStatus()
	statusResponse, err := internalWorkspaces.GetStatus()
	if err != nil {
		return status.NewStatusDefault(int(*err.Code)).WithPayload(err)
	}
	return status.NewStatusOK().WithPayload(statusResponse)
}
