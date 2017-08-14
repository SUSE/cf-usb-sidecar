package rabbitmq

import (
	"testing"

	"github.com/hpcloud/go-csm-lib/extension"
	"github.com/SUSE/cf-usb-sidecar/csm-extensions/services/dev-rabbitmq/config"
	"github.com/SUSE/cf-usb-sidecar/csm-extensions/services/dev-rabbitmq/provisioner/provisionerfakes"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("rabbitmq-test")

func getFakeProvisioner() (*provisionerfakes.FakeRabbitmqProvisionerInterface, extension.Extension) {
	fakeProv := new(provisionerfakes.FakeRabbitmqProvisionerInterface)
	rabbit := NewRabbitmqExtension(fakeProv, config.RabbitmqConfig{}, logger)

	return fakeProv, rabbit
}

func Test_CreateWorkspace(t *testing.T) {
	assert := assert.New(t)

	fakeProv, rabbit := getFakeProvisioner()
	fakeProv.CreateContainerReturns(nil)

	workspaceID := "testId"

	response, err := rabbit.CreateWorkspace(workspaceID)

	assert.NotNil(response)
	assert.Equal("successful", response.Status)
	assert.NoError(err)
}

func Test_GetWorkspace(t *testing.T) {
	assert := assert.New(t)

	fakeProv, rabbit := getFakeProvisioner()
	fakeProv.ContainerExistsReturns(true, nil)

	workspaceID := "testId"

	response, err := rabbit.GetWorkspace(workspaceID)

	assert.NotNil(response)
	assert.Equal("successful", response.Status)
	assert.NoError(err)
}

func Test_GetConnection(t *testing.T) {
	assert := assert.New(t)
	fakeProv, rabbit := getFakeProvisioner()
	fakeProv.UserExistsReturns(true, nil)

	workspaceID := "testId"
	credentialsID := "credentialId"

	response, err := rabbit.GetConnection(workspaceID, credentialsID)

	assert.NotNil(response)
	assert.Equal("successful", response.Status)
	assert.NoError(err)
}

func Test_CreateConnection(t *testing.T) {
	assert := assert.New(t)

	fakeProv, rabbit := getFakeProvisioner()
	fakeProv.CreateUserReturns(map[string]string{
		"host":     "127.0.0.1",
		"user":     "user",
		"password": "password",
		"port":     "1234",
	}, nil)

	workspaceID := "testId"
	credentialsID := "credentialId"

	response, err := rabbit.CreateConnection(workspaceID, credentialsID)

	assert.NotNil(response)
	assert.Equal("successful", response.Status)
	assert.NoError(err)
}

func Test_DeleteConnection(t *testing.T) {
	assert := assert.New(t)

	fakeProv, rabbit := getFakeProvisioner()
	fakeProv.DeleteUserReturns(nil)

	workspaceID := "testId"
	credentialsID := "credentialId"

	response, err := rabbit.DeleteConnection(workspaceID, credentialsID)

	assert.NotNil(response)
	assert.Equal("successful", response.Status)
	assert.NoError(err)
}

func Test_DeleteWorkspace(t *testing.T) {
	assert := assert.New(t)
	fakeProv, rabbit := getFakeProvisioner()
	fakeProv.DeleteContainerReturns(nil)

	workspaceID := "testId"

	response, err := rabbit.DeleteWorkspace(workspaceID)

	assert.NotNil(response)
	assert.Equal("successful", response.Status)
	assert.NoError(err)
}
