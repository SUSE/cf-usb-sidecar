package redis

import (
	"testing"

	"github.com/SUSE/cf-usb-sidecar/csm-extensions/services/dev-redis/config"
	"github.com/SUSE/cf-usb-sidecar/csm-extensions/services/dev-redis/provisioner/provisionerfakes"
	"github.com/hpcloud/go-csm-lib/extension"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("redis-test")

func getFakeProvisioner() (*provisionerfakes.FakeRedisProvisionerInterface, extension.Extension) {
	fakeProv := new(provisionerfakes.FakeRedisProvisionerInterface)
	redis := NewRedisExtension(fakeProv, config.RedisConfig{}, logger)

	return fakeProv, redis
}

func Test_CreateWorkspace(t *testing.T) {
	assert := assert.New(t)

	fakeProv, redis := getFakeProvisioner()
	fakeProv.CreateContainerReturns(nil)

	workspaceID := "testId"

	response, err := redis.CreateWorkspace(workspaceID)

	assert.NotNil(response)
	assert.Equal("successful", response.Status)
	assert.NoError(err)
}

func Test_GetWorkspace(t *testing.T) {
	assert := assert.New(t)

	fakeProv, redis := getFakeProvisioner()
	fakeProv.ContainerExistsReturns(true, nil)

	workspaceID := "testId"

	response, err := redis.GetWorkspace(workspaceID)

	assert.NotNil(response)
	assert.Equal("successful", response.Status)
	assert.NoError(err)
}

func Test_GetConnection(t *testing.T) {
	assert := assert.New(t)
	_, redis := getFakeProvisioner()

	workspaceID := "testId"
	credentialsID := "credentialId"

	response, err := redis.GetConnection(workspaceID, credentialsID)

	assert.NotNil(response)
	assert.Equal("failed", response.Status)
	assert.NoError(err)
}

func Test_CreateConnection(t *testing.T) {
	assert := assert.New(t)

	_, redis := getFakeProvisioner()

	workspaceID := "testId"
	credentialsID := "credentialId"

	response, err := redis.CreateConnection(workspaceID, credentialsID)

	assert.NotNil(response)
	assert.Equal("successful", response.Status)
	assert.NoError(err)
}

func Test_DeleteConnection(t *testing.T) {
	assert := assert.New(t)

	_, redis := getFakeProvisioner()

	workspaceID := "testId"
	credentialsID := "credentialId"

	response, err := redis.DeleteConnection(workspaceID, credentialsID)

	assert.NotNil(response)
	assert.Equal("successful", response.Status)
	assert.NoError(err)
}

func Test_DeleteWorkspace(t *testing.T) {
	assert := assert.New(t)
	fakeProv, redis := getFakeProvisioner()
	fakeProv.DeleteContainerReturns(nil)

	workspaceID := "testId"

	response, err := redis.DeleteWorkspace(workspaceID)

	assert.NotNil(response)
	assert.Equal("successful", response.Status)
	assert.NoError(err)
}
