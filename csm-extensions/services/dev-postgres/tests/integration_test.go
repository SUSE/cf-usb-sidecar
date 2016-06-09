// +build integration

package main

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/fsouza/go-dockerclient"
	swaggerClient "github.com/go-swagger/go-swagger/client"
	httpClient "github.com/go-swagger/go-swagger/httpkit/client"
	"github.com/go-swagger/go-swagger/strfmt"
	csmClient "github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager-client/client"
	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager-client/client/connection"
	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager-client/client/workspace"
	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager-client/models"
	"github.com/stretchr/testify/assert"
)

const (
	DockerName   = "csm-dev-postgres:latest"
	DockerPort   = 8093
	DockerIP     = "127.0.0.1"
	WorkspaceID  = "test-onnllyy123"
	ConnectionID = "testconnonnllyy123"
	Token        = "csm-auth-token"
)

var (
	transportHost string
	transport     *httpClient.Runtime
	client        *csmClient.CatlogServiceManager
	authFunc      swaggerClient.AuthInfoWriter
)

func initializeTestAssets(t *testing.T) bool {

	err := checkPrerequisites()

	if err != nil {
		return assert.Fail(t, err.Error())
	}
	testServerIP := os.Getenv("TEST_SERVER_IP")
	testServerPort := os.Getenv("TEST_SERVER_PORT")
	testServerToken := os.Getenv("TEST_SERVER_TOKEN")

	host := DockerIP
	port := DockerPort
	token := Token
	if testServerIP != "" {
		host = testServerIP
	}
	if testServerPort != "" {
		sPort, err := strconv.Atoi(testServerPort)
		if err == nil {
			port = sPort
		}
	}

	if testServerToken != "" {
		token = testServerToken
	}

	transportHost = host + ":" + strconv.Itoa(port)
	transport = httpClient.New(transportHost, "", []string{"http"})
	client = csmClient.New(transport, strfmt.Default)
	authFunc = httpClient.APIKeyAuth("x-csm-token", "header", token)
	return true
}

func dockerImgExists(client *docker.Client, dockerName string) (bool, error) {
	imgs, err := client.ListImages(docker.ListImagesOptions{All: false})

	if err != nil {
		return false, err
	}

	for _, img := range imgs {
		for _, imgName := range img.RepoTags {
			if imgName == dockerName {
				return true, nil
			}
		}
	}

	return false, nil

}

func checkPrerequisites() error {

	client, err := docker.NewClientFromEnv()
	if err != nil {
		return err
	}

	exists, err := dockerImgExists(client, DockerName)

	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("The %s docker image was not found. Please run make run and make run-db or just make all first!", DockerName)
	}

	return nil
}

func TestDeleteWorkspaceShouldFail(t *testing.T) {
	assert := assert.New(t)

	if !initializeTestAssets(t) {
		return
	}

	params := workspace.NewDeleteWorkspaceParams().WithWorkspaceID(WorkspaceID)
	resp, err := client.Workspace.DeleteWorkspace(params, authFunc)
	if err != nil {
		t.Logf("Delete workspace error: %s", err.Error())
	}
	if resp != nil {
		t.Logf("Delete workspace resp: %s", resp.Error())
	}
	if assert.Error(err, "There should be an error while deleting a non existing workspace") {
		assert.Contains(err.Error(), "pq: database \"dtestonnllyy123\" does not exist", "Incorrect answer when deleting a database that does not exist")
	}
	assert.Nil(resp, "There should be no correct unswer when deleting a non existing workspace")
}

func TestGetWorkspaceShouldFail(t *testing.T) {
	assert := assert.New(t)
	if !initializeTestAssets(t) {
		return
	}

	params := workspace.NewGetWorkspaceParams().WithWorkspaceID(WorkspaceID)
	resp, err := client.Workspace.GetWorkspace(params, authFunc)
	if err != nil {
		t.Logf("Get workspace error: %s", err.Error())
	}
	if resp != nil {
		t.Logf("Get workspace resp: Status = %s, ProcessingType = %s", resp.Payload.Status, resp.Payload.ProcessingType)
	}
	if assert.Error(err, "Expected error since workspace does not exist") {
		assert.Contains(err.Error(), "Workspace does not exist", "The error message is incorrect for getting an inexistent workspace")
	}
	assert.Nil(resp, "response should be nil since there was an error")
}

func TestCreateWorkspaceShouldSucced(t *testing.T) {
	assert := assert.New(t)
	if !initializeTestAssets(t) {
		return
	}
	createWorkspaceRequest := models.ServiceManagerWorkspaceCreateRequest{
		WorkspaceID: WorkspaceID,
	}
	params := workspace.NewCreateWorkspaceParams().WithCreateWorkspaceRequest(&createWorkspaceRequest)
	resp, err := client.Workspace.CreateWorkspace(params, authFunc)
	if err != nil {
		t.Logf("Create workspace error: %s", err.Error())
	}
	if resp != nil {
		t.Logf("Create workspace resp: Status = %s, ProcessingType = %s", resp.Payload.Status, resp.Payload.ProcessingType)
	}
	assert.NoError(err, "There was an unexpected error while creating workspace")
	if assert.NotNil(resp, "There should be no error when creating a workspace") {
		assert.Equal("Extension", resp.Payload.ProcessingType, "Unexpected processing_type")
		assert.Equal("successful", resp.Payload.Status, "Unexpected status")
	}

}

func TestGetConnectionShouldFail(t *testing.T) {
	assert := assert.New(t)
	if !initializeTestAssets(t) {
		return
	}

	params := connection.NewGetConnectionParams().WithConnectionID(ConnectionID).WithWorkspaceID(WorkspaceID)
	resp, err := client.Connection.GetConnection(params, authFunc)

	if err != nil {
		t.Logf("Get connection error: %s", err.Error())
	}
	if resp != nil {
		t.Logf("Get connection resp: Status = %s, ProcessingType = %s", resp.Payload.Status, resp.Payload.ProcessingType)
	}
	assert.Error(err, "Expected error since the connection has not been created yet")
	assert.Nil(resp, "response should be nil as no connection with this name was yet created")
}

func TestDeleteConnectionShouldFail(t *testing.T) {
	assert := assert.New(t)
	if !initializeTestAssets(t) {
		return
	}

	params := connection.NewDeleteConnectionParams().WithConnectionID(ConnectionID).WithWorkspaceID(WorkspaceID)
	resp, err := client.Connection.DeleteConnection(params, authFunc)

	if err != nil {
		t.Logf("Delete connection error: %s", err.Error())
	}
	if resp != nil {
		t.Logf("Delete connection resp: %s", resp.Error())
	}
	assert.Error(err, "Expected error since the connection has not been created yet")
	assert.Nil(resp, "response should be nil as no connection with this name was yet created")
}

func TestCreateConnectionShouldSucced(t *testing.T) {
	assert := assert.New(t)
	if !initializeTestAssets(t) {
		return
	}

	createConnectionRequest := models.ServiceManagerConnectionCreateRequest{
		ConnectionID: ConnectionID,
	}
	params := connection.NewCreateConnectionParams().WithWorkspaceID(WorkspaceID).WithConnectionCreateRequest(&createConnectionRequest)
	resp, err := client.Connection.CreateConnection(params, authFunc)

	if err != nil {
		t.Logf("Create connection error: %s", err.Error())
	}
	if resp != nil {
		t.Logf("Create connection resp: Status = %s, ProcessingType = %s, Error=%s", resp.Payload.Status, resp.Payload.ProcessingType, resp.Error())
	}
	assert.NoError(err, "No error expected since the connection has not been created yet")
	if assert.NotNil(resp, "response should not be nil as no connection with this name was yet created") {
		assert.Equal("Extension", resp.Payload.ProcessingType, "Incorrect extension received")
		assert.Equal("successful", resp.Payload.Status, "Invalid status received")
		assert.NotNil(resp.Payload.Details, "The details should contain connection info")
	}
}

func TestGetConnectionShouldSucceed(t *testing.T) {
	assert := assert.New(t)
	if !initializeTestAssets(t) {
		return
	}

	params := connection.NewGetConnectionParams().WithWorkspaceID(WorkspaceID).WithConnectionID(ConnectionID)
	resp, err := client.Connection.GetConnection(params, authFunc)

	if err != nil {
		t.Logf("Get connection error: %s", err.Error())
	}
	if resp != nil {
		t.Logf("Get connection resp: Status = %s, ProcessingType = %s, Error=%s", resp.Payload.Status, resp.Payload.ProcessingType, resp.Error())
	}
	assert.NoError(err, "No error expected since the connection has already been created yet")
	if assert.NotNil(resp, "response should not be nil as a connection with this name was already created") {
		assert.Equal("Extension", resp.Payload.ProcessingType, "Incorrect extension received")
		assert.Equal("successful", resp.Payload.Status, "Invalid status received")
	}
}

func TestDeleteConnectionShouldSucceed(t *testing.T) {
	assert := assert.New(t)
	if !initializeTestAssets(t) {
		return
	}

	params := connection.NewDeleteConnectionParams().WithWorkspaceID(WorkspaceID).WithConnectionID(ConnectionID)
	resp, err := client.Connection.DeleteConnection(params, authFunc)

	if err != nil {
		t.Logf("Delete connection error: %s", err.Error())
	}
	if resp != nil {
		t.Logf("Delete connection resp: %s", resp.Error())
	}
	assert.NoError(err, "No error expected since the connection has already been created yet")
	assert.NotNil(resp, "response should not be nil as a connection with this name was already created")
}

func TestCreateWorkspaceShouldFail(t *testing.T) {
	assert := assert.New(t)
	if !initializeTestAssets(t) {
		return
	}

	createWorkspaceRequest := models.ServiceManagerWorkspaceCreateRequest{
		WorkspaceID: WorkspaceID,
	}
	params := workspace.NewCreateWorkspaceParams().WithCreateWorkspaceRequest(&createWorkspaceRequest)
	resp, err := client.Workspace.CreateWorkspace(params, authFunc)
	if err != nil {
		t.Logf("Create workspace error: %s", err.Error())
	}
	if resp != nil {
		t.Logf("Create workspace resp: Status = %s, ProcessingType = %s", resp.Payload.Status, resp.Payload.ProcessingType)
	}
	if assert.Error(err, "There should be an error when creating an workspace that allready exists") {
		assert.Contains(err.Error(), "pq: database \"dtestonnllyy123\" already exists", "There should be an error message stating that this db allready exists when attempting to create an existing database")
	}
	assert.Nil(resp, "There should be no correct unswer when creating a workspace that allready exists")

}

func TestGetWorkspacesShouldSucced(t *testing.T) {
	assert := assert.New(t)
	if !initializeTestAssets(t) {
		return
	}

	params := workspace.NewGetWorkspaceParams().WithWorkspaceID(WorkspaceID)
	resp, err := client.Workspace.GetWorkspace(params, authFunc)
	if err != nil {
		t.Logf("Get workspace error: %s", err.Error())
	}
	if resp != nil {
		t.Logf("Get workspace resp: Status: %s, ProcessingType: %s", resp.Payload.Status, resp.Payload.ProcessingType)
	}
	if assert.NoError(err, "There was an unexpected error while getting existing workspace.") {
		assert.Equal("Extension", resp.Payload.ProcessingType, "Unexpected processing_type")
		assert.Equal("successful", resp.Payload.Status, "Unexpected status")
	}
}

func TestDeleteWorkspaceShouldSucced(t *testing.T) {
	assert := assert.New(t)
	if !initializeTestAssets(t) {
		return
	}

	params := workspace.NewDeleteWorkspaceParams().WithWorkspaceID(WorkspaceID)
	resp, err := client.Workspace.DeleteWorkspace(params, authFunc)
	if err != nil {
		t.Logf("Delete workspace error: %s", err.Error())
	}
	if resp != nil {
		t.Logf("Delete workspace resp: Status %s", resp.Error())
	}
	assert.NoError(err, "There was an unexpected error while deleting the workspace")
	assert.NotNil(resp, "Unexpected err occured")
}
