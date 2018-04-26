package routing

import (
	"github.com/SUSE/cf-usb-sidecar/csm-extensions/services/dev-routing/config"
	"github.com/SUSE/go-csm-lib/csm"
	"github.com/SUSE/go-csm-lib/extension"
	"github.com/pivotal-golang/lager"
)

type routingExtension struct {
	conf   config.RoutingConfig
	logger lager.Logger
}

func NewRoutingExtension(conf config.RoutingConfig, logger lager.Logger) extension.Extension {
	return &routingExtension{conf: conf, logger: logger}
}

func (e *routingExtension) CreateConnection(workspaceID, connectionID string) (*csm.CSMResponse, error) {
	e.logger.Info("create-connection", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})

	binding := config.RouteBinding{
		RouteServiceURL: connectionID,
	}

	response := csm.CreateCSMResponse(binding)
	return &response, nil
}
func (e *routingExtension) CreateWorkspace(workspaceID string) (*csm.CSMResponse, error) {
	e.logger.Info("create-workspace", lager.Data{"workspaceID": workspaceID})

	response := csm.CreateCSMResponse("")

	return &response, nil
}
func (e *routingExtension) DeleteConnection(workspaceID, connectionID string) (*csm.CSMResponse, error) {
	e.logger.Info("delete-connection", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})

	response := csm.CreateCSMResponse("")

	return &response, nil
}
func (e *routingExtension) DeleteWorkspace(workspaceID string) (*csm.CSMResponse, error) {
	e.logger.Info("delete-workspace", lager.Data{"workspaceID": workspaceID})

	response := csm.CreateCSMResponse("")

	return &response, nil
}
func (e *routingExtension) GetConnection(workspaceID, connectionID string) (*csm.CSMResponse, error) {

	response := csm.CreateCSMErrorResponse(404, "Connection does not exist")

	return &response, nil
}
func (e *routingExtension) GetWorkspace(workspaceID string) (*csm.CSMResponse, error) {
	response := csm.CreateCSMErrorResponse(404, "Workspace does not exist")

	return &response, nil
}

func (e *routingExtension) GetStatus() (*csm.CSMResponse, error) {
	response := csm.CSMResponse{}
	response.Status = "successful"
	response.ServiceType = "routing"
	return &response, nil
}
