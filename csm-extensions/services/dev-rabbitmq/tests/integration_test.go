// +build integration

package rabbitmq

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/fsouza/go-dockerclient"
	swaggerClient "github.com/go-swagger/go-swagger/client"
	httpClient "github.com/go-swagger/go-swagger/httpkit/client"
	"github.com/go-swagger/go-swagger/strfmt"
	csmClient "github.com/SUSE/cf-usb-sidecar/generated/CatalogServiceManager-client/client"
	"github.com/SUSE/cf-usb-sidecar/generated/CatalogServiceManager-client/client/connection"
	"github.com/SUSE/cf-usb-sidecar/generated/CatalogServiceManager-client/client/workspace"
	"github.com/SUSE/cf-usb-sidecar/generated/CatalogServiceManager-client/models"
	"github.com/stretchr/testify/assert"
)

const WORKSPACE_ID = "test_workspace"
const CONNECTION_ID = "test_connection"

var transportHost string
var transport *httpClient.Runtime
var client *csmClient.CatlogServiceManager
var authFunc swaggerClient.AuthInfoWriter
var dockerContainerName string
var csmExtensionHost string
var csmExtensionToken string
var csmExtensionPort string

func initTest() {
	dockerHost := os.Getenv("DOCKER_HOST")
	if dockerHost == "" {
		dockerHost = "unix:///var/run/docker.sock"
		csmExtensionHost = "127.0.0.1"
	} else {
		csmExtensionHost = os.Getenv("DOCKER_HOST_IP")
	}
	csmExtensionToken = os.Getenv("SIDECAR_EXTENSION_TOKEN")
	dockerContainerName = os.Getenv("SIDECAR_EXTENSION_IMAGE_NAME")
	csmExtensionPort = os.Getenv("SIDECAR_EXTENSION_PORT")

	transportHost = fmt.Sprintf("%s:%s", csmExtensionHost, csmExtensionPort)
	transport = httpClient.New(transportHost, "", []string{"http"})
	client = csmClient.New(transport, strfmt.Default)
	authFunc = httpClient.APIKeyAuth("x-sidecar-token", "header", csmExtensionToken)
}

func dockerContainerExists(client *docker.Client, containerName string) (bool, error) {
	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
	if err != nil {
		return false, err
	}

	for _, c := range containers {
		for _, n := range c.Names {
			if strings.TrimPrefix(n, "/") == containerName {
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

	exists, err := dockerContainerExists(client, dockerContainerName)

	if err != nil {
		return err
	}
	if !exists {
		return errors.New(fmt.Sprintf("The %s docker container was not found. Please run make run and make run-db or just make all first!", dockerContainerName))
	}

	return nil
}

func Test_FailGetWorkspace(t *testing.T) {
	assert := assert.New(t)
	initTest()

	err := checkPrerequisites()
	if err != nil {
		log.Fatalf(err.Error())
	}

	params := workspace.NewGetWorkspaceParams().WithWorkspaceID(WORKSPACE_ID)
	response, err := client.Workspace.GetWorkspace(params, authFunc)

	if err != nil {
		log.Println(err.Error())
	}

	assert.NotNil(err)
	assert.Nil(response)
}

func Test_FailGetConnection(t *testing.T) {
	assert := assert.New(t)
	initTest()

	err := checkPrerequisites()
	if err != nil {
		log.Fatalf(err.Error())
	}

	params := connection.NewGetConnectionParams().WithConnectionID(CONNECTION_ID)
	response, err := client.Connection.GetConnection(params, authFunc)

	if err != nil {
		log.Println(err.Error())
	}

	assert.NotNil(err)
	assert.Nil(response)
}

func Test_CreateWorkspace(t *testing.T) {
	assert := assert.New(t)
	initTest()

	err := checkPrerequisites()
	if err != nil {
		log.Fatalf(err.Error())
	}

	createWorkspaceRequest := models.ServiceManagerWorkspaceCreateRequest{
		WorkspaceID: WORKSPACE_ID,
	}
	params := workspace.NewCreateWorkspaceParams().WithCreateWorkspaceRequest(&createWorkspaceRequest)
	response, err := client.Workspace.CreateWorkspace(params, authFunc)
	if err != nil {
		log.Println(err.Error())
	}

	assert.Nil(err)
	assert.NotNil(response)

	if response != nil {
		assert.Equal("Extension", response.Payload.ProcessingType)
		assert.Equal("successful", response.Payload.Status)
	}
}

func Test_GetWorkspace(t *testing.T) {
	assert := assert.New(t)
	initTest()

	err := checkPrerequisites()
	if err != nil {
		log.Fatalf(err.Error())
	}

	params := workspace.NewGetWorkspaceParams().WithWorkspaceID(WORKSPACE_ID)
	response, err := client.Workspace.GetWorkspace(params, authFunc)

	if err != nil {
		log.Println(err.Error())
	}

	assert.Nil(err)
	assert.NotNil(response)

	if response != nil {
		assert.Equal("Extension", response.Payload.ProcessingType)
		assert.Equal("successful", response.Payload.Status)
	}
}

func Test_CreateConnection(t *testing.T) {
	assert := assert.New(t)
	initTest()

	err := checkPrerequisites()
	if err != nil {
		log.Fatalf(err.Error())
	}

	// wait for server to start
	time.Sleep(10 * time.Second)

	createConnectionRequest := models.ServiceManagerConnectionCreateRequest{
		ConnectionID: CONNECTION_ID,
	}
	params := connection.NewCreateConnectionParams().WithWorkspaceID(WORKSPACE_ID).WithConnectionCreateRequest(&createConnectionRequest)
	response, err := client.Connection.CreateConnection(params, authFunc)

	assert.Nil(err)
	assert.NotNil(response)

	if response != nil {
		assert.Equal("Extension", response.Payload.ProcessingType)
		assert.Equal("successful", response.Payload.Status)
	}
}

func Test_GetConnection(t *testing.T) {
	assert := assert.New(t)
	initTest()

	err := checkPrerequisites()
	if err != nil {
		log.Fatalf(err.Error())
	}

	params := connection.NewGetConnectionParams().WithConnectionID(CONNECTION_ID).WithWorkspaceID(WORKSPACE_ID)
	response, err := client.Connection.GetConnection(params, authFunc)

	assert.Nil(err)
	assert.NotNil(response)

	if response != nil {
		assert.Equal("Extension", response.Payload.ProcessingType)
		assert.Equal("successful", response.Payload.Status)
	}
}

func Test_DeleteConnection(t *testing.T) {
	assert := assert.New(t)
	initTest()

	err := checkPrerequisites()
	if err != nil {
		log.Fatalf(err.Error())
	}

	params := connection.NewDeleteConnectionParams().WithConnectionID(CONNECTION_ID).WithWorkspaceID(WORKSPACE_ID)
	response, err := client.Connection.DeleteConnection(params, authFunc)

	assert.Nil(err)
	assert.NotNil(response)
}

func Test_GetConnectionAfterDelete(t *testing.T) {
	assert := assert.New(t)
	initTest()

	err := checkPrerequisites()
	if err != nil {
		log.Fatalf(err.Error())
	}

	params := connection.NewGetConnectionParams().WithConnectionID(CONNECTION_ID).WithWorkspaceID(WORKSPACE_ID)
	response, err := client.Connection.GetConnection(params, authFunc)

	assert.NotNil(err)
	assert.Nil(response)
}

func Test_DeleteWorkspace(t *testing.T) {
	assert := assert.New(t)
	initTest()

	err := checkPrerequisites()
	if err != nil {
		log.Fatalf(err.Error())
	}

	params := workspace.NewDeleteWorkspaceParams().WithWorkspaceID(WORKSPACE_ID)
	response, err := client.Workspace.DeleteWorkspace(params, authFunc)
	if err != nil {
		log.Println(err.Error())
	}

	assert.Nil(err)
	assert.NotNil(response)
}

func Test_GetWorkspaceAfterDelete(t *testing.T) {
	assert := assert.New(t)
	initTest()

	err := checkPrerequisites()
	if err != nil {
		log.Fatalf(err.Error())
	}

	params := workspace.NewGetWorkspaceParams().WithWorkspaceID(WORKSPACE_ID)
	response, err := client.Workspace.GetWorkspace(params, authFunc)

	if err != nil {
		log.Println(err.Error())
	}

	assert.NotNil(err)
	assert.Nil(response)
}
