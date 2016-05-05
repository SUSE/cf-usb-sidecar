package IntegrationTest

import (
	"testing"

	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager-client/client/workspace"
	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager-client/models"
	"github.com/stretchr/testify/assert"
)

func TestGetWorkspacesBeforeCreate(t *testing.T) {
	assert := assert.New(t)
	params := workspace.NewGetWorkspaceParams().WithWorkspaceID("123")
	resp, err := client.Workspace.GetWorkspace(params, authFunc)

	assert.NotNil(err, "Expected error since workspace does not exist")
	assert.Nil(resp, "response should be nil since there was an error")
}

func TestCreateWorkspaces(t *testing.T) {
	assert := assert.New(t)
	createWorkspaceRequest := models.ServiceManagerWorkspaceCreateRequest{
		WorkspaceID: "123",
	}
	params := workspace.NewCreateWorkspaceParams().WithCreateWorkspaceRequest(&createWorkspaceRequest)
	resp, err := client.Workspace.CreateWorkspace(params, authFunc)
	assert.Nil(err, "There was an unexpected error while creating workspace.")
	assert.Equal("Extension", resp.Payload.ProcessingType, "Unexpected processing_type")
	assert.Equal("successful", resp.Payload.Status, "Unexpected status")
}

func TestGetWorkspacesAfterCreate(t *testing.T) {
	assert := assert.New(t)
	params := workspace.NewGetWorkspaceParams().WithWorkspaceID("123")
	resp, err := client.Workspace.GetWorkspace(params, authFunc)
	assert.Nil(err, "There was an unexpected error while creating workspace.")
	assert.Equal("Extension", resp.Payload.ProcessingType, "Unexpected processing_type")
	assert.Equal("successful", resp.Payload.Status, "Unexpected status")
}
