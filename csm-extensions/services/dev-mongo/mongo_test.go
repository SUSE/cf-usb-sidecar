package mongo

import (
	"testing"

	"github.com/hpcloud/go-csm-lib/extension"

	"github.com/hpcloud/catalog-service-manager/csm-extensions/services/dev-mongo/config"

	"github.com/hpcloud/catalog-service-manager/csm-extensions/services/dev-mongo/provisioner/provisionerfakes"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("mongo-test")

func getFakeProvisioner() (*provisionerfakes.FakeMongoProvisionerInterface, extension.Extension) {
	fakeProv := new(provisionerfakes.FakeMongoProvisionerInterface)

	mongo := NewMongoExtension(fakeProv, config.MongoDriverConfig{}, logger)
	return fakeProv, mongo
}

func TestCreateWorkspace(t *testing.T) {
	assert := assert.New(t)

	fakeProv, mongo := getFakeProvisioner()
	fakeProv.CreateDatabaseReturns(nil)

	workspaceid := "8b490a70-a892-4eff-a495-81e905f3960f"

	response, err := mongo.CreateWorkspace(workspaceid)
	assert.NotNil(response)
	assert.Equal("successful", response.Status)
	assert.NoError(err)
}

func TestGetWorkspace(t *testing.T) {
	assert := assert.New(t)

	fakeProv, mongo := getFakeProvisioner()
	fakeProv.IsDatabaseCreatedReturns(true, nil)

	workspaceid := "8b490a70-a892-4eff-a495-81e905f3960f"

	response, err := mongo.GetWorkspace(workspaceid)
	assert.NotNil(response)
	assert.Equal("successful", response.Status)
	assert.NoError(err)
}

func TestDeleteWorkspace(t *testing.T) {
	assert := assert.New(t)

	fakeProv, mongo := getFakeProvisioner()
	fakeProv.DeleteDatabaseReturns(nil)

	workspaceid := "8b490a70-a892-4eff-a495-81e905f3960f"

	response, err := mongo.DeleteWorkspace(workspaceid)
	assert.NotNil(response)
	assert.Equal("successful", response.Status)
	assert.NoError(err)
}

func TestCreateConnection(t *testing.T) {
	assert := assert.New(t)

	fakeProv, mongo := getFakeProvisioner()
	fakeProv.CreateUserReturns(nil)

	workspaceid := "8b490a70-a892-4eff-a495-81e905f3960f"
	connectionid := "8b490a70-a892-4eff-a495-81e905f3961d"

	response, err := mongo.CreateConnection(workspaceid, connectionid)
	assert.NotNil(response)
	assert.Equal("successful", response.Status)
	assert.NoError(err)
}

func TestGetConnection(t *testing.T) {
	assert := assert.New(t)

	fakeProv, mongo := getFakeProvisioner()
	fakeProv.IsUserCreatedReturns(true, nil)

	workspaceid := "8b490a70-a892-4eff-a495-81e905f3960f"
	connectionid := "8b490a70-a892-4eff-a495-81e905f3961d"

	response, err := mongo.GetConnection(workspaceid, connectionid)
	assert.NotNil(response)
	assert.Equal("successful", response.Status)
	assert.NoError(err)
}

func TestDeleteConnection(t *testing.T) {
	assert := assert.New(t)

	fakeProv, mongo := getFakeProvisioner()
	fakeProv.DeleteUserReturns(nil)

	workspaceid := "8b490a70-a892-4eff-a495-81e905f3960f"
	connectionid := "8b490a70-a892-4eff-a495-81e905f3961d"

	response, err := mongo.DeleteConnection(workspaceid, connectionid)
	assert.NotNil(response)
	assert.Equal("successful", response.Status)
	assert.NoError(err)
}
