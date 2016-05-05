package IntegrationTest

import (
	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager-client/client/workspace"
	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager-client/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetWorkspaces_BeforeCreate(t *testing.T) {
	params := workspace.NewGetWorkspaceParams().WithWorkspaceID("123")
	resp, err := client.Workspace.GetWorkspace(params, authFunc)
	assert.Nil(t, err, "There was an unexpected error while creating workspace.")
	assert.Equal(t, "Extension", resp.Payload.ProcessingType, "Unexpected processing_type")
	assert.Equal(t, "failed", resp.Payload.Status, "Unexpected status")
}

func TestCreateWorkspaces(t *testing.T) {
	createWorkspaceRequest := models.ServiceManagerWorkspaceCreateRequest{
		WorkspaceID: "123",
	}
	params := workspace.NewCreateWorkspaceParams().WithCreateWorkspaceRequest(&createWorkspaceRequest)
	resp, err := client.Workspace.CreateWorkspace(params, authFunc)
	assert.Nil(t, err, "There was an unexpected error while creating workspace.")
	assert.Equal(t, "Extension", resp.Payload.ProcessingType, "Unexpected processing_type")
	assert.Equal(t, "successful", resp.Payload.Status, "Unexpected status")
}

func TestGetWorkspaces_AfterCreate(t *testing.T) {
	params := workspace.NewGetWorkspaceParams().WithWorkspaceID("123")
	resp, err := client.Workspace.GetWorkspace(params, authFunc)
	assert.Nil(t, err, "There was an unexpected error while creating workspace.")
	assert.Equal(t, "Extension", resp.Payload.ProcessingType, "Unexpected processing_type")
	assert.Equal(t, "successful", resp.Payload.Status, "Unexpected status")
}
