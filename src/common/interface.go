package common

import (
	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager/models"
)

type CSMSetupInterface interface {
	CheckExtensions()
	Startup() bool
}

type CSMWorkspaceInterface interface {
	CheckExtensions()
	GetWorkspace(workspaceID string) *models.ServiceManagerWorkspaceResponse
	CreateWorkspace(workspaceCreate *models.ServiceManagerWorkspaceCreateRequest) *models.ServiceManagerWorkspaceResponse
	DeleteWorkspace(workspaceID string) *models.ServiceManagerWorkspaceResponse
}

type CSMConnectionInterface interface {
	CheckExtensions()
	GetConnection(workspaceID string, connectionID string) *models.ServiceManagerConnectionResponse
	CreateConnection(workspaceID string, createConnection *models.ServiceManagerConnectionCreateRequest) *models.ServiceManagerConnectionResponse
	DeleteConnection(workspaceID string, connectionID string) *models.ServiceManagerConnectionResponse
}
