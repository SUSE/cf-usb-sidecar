package common

import (
	"github.com/SUSE/cf-usb-sidecar/generated/CatalogServiceManager/models"
)

type CSMWorkspaceInterface interface {
	CheckExtensions()
	GetWorkspace(workspaceID string) *models.ServiceManagerWorkspaceResponse
	CreateWorkspace(workspaceCreate *models.ServiceManagerWorkspaceCreateRequest) *models.ServiceManagerWorkspaceResponse
	DeleteWorkspace(workspaceID string) *models.ServiceManagerWorkspaceResponse
}

type CSMConnectionInterface interface {
	CheckExtensions()
	GetConnection(workspaceID string, connectionID string) *models.ServiceManagerConnectionResponse
	CreateConnection(workspaceID string, createConnection *models.ServiceManagerConnectionCreateRequest, Details map[string]interface{}) *models.ServiceManagerConnectionResponse
	DeleteConnection(workspaceID string, connectionID string) *models.ServiceManagerConnectionResponse
}

type CSMStatusInterface interface {
	GetStatus() *models.StatusResponse
}
