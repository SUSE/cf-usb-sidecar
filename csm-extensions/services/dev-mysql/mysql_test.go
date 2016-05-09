package mysql

import (
	"errors"
	"testing"

	"github.com/hpcloud/catalog-service-manager/csm-extensions/services/dev-mysql/config"
	"github.com/hpcloud/catalog-service-manager/csm-extensions/services/dev-mysql/provisioner/provisionerfakes"
	"github.com/hpcloud/go-csm-lib/extension"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("mysql-provisioner-test")

func getMySQLExtension() (extension.Extension, *provisionerfakes.FakeMySQLProvisioner) {
	logger = lagertest.NewTestLogger("process-controller")

	conf := config.MySQLConfig{
		User: "testuser",
		Pass: "testpass",
		Host: "testhost",
		Port: "3306",
	}

	var fakeProvisioner = new(provisionerfakes.FakeMySQLProvisioner)

	extension := NewMySQLExtension(fakeProvisioner, conf, logger)
	return extension, fakeProvisioner

}

func TestCreateConnection(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMySQLExtension()
	fakeProv.CreateUserReturns(nil)

	response, err := ext.CreateConnection("workspace", "connection")

	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(0, response.ErrorCode)
	assert.Equal("successful", response.Status)

	creds := response.Details.(config.MySQLBinding)

	assert.NotEmpty(creds.Database)
	assert.NotEmpty(creds.Host)
	assert.NotEmpty(creds.Hostname)
	assert.NotEmpty(creds.Password)
	assert.NotEmpty(creds.Port)
	assert.NotEmpty(creds.User)
	assert.NotEmpty(creds.Password)
	assert.Equal(creds.Host, creds.Hostname)
	assert.Equal("3306", creds.Port)
	assert.Equal("testhost", creds.Host)
}

func TestCreateConnectionError(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMySQLExtension()
	fakeProv.CreateUserReturns(errors.New("db error"))

	_, err := ext.CreateConnection("workspace", "connection")

	assert.NotNil(err)
}

func TestCreateWorkspace(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMySQLExtension()
	fakeProv.CreateDatabaseReturns(nil)

	response, err := ext.CreateWorkspace("workspaceID")

	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(0, response.ErrorCode)
	assert.Equal("successful", response.Status)
}

func TestCreateWorkspaceError(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMySQLExtension()
	fakeProv.CreateDatabaseReturns(errors.New("this is an error"))

	response, err := ext.CreateWorkspace("workspaceID")

	assert.NotNil(err)
	assert.Nil(response)
}

func TestDeleteConnection(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMySQLExtension()
	fakeProv.DeleteDatabaseReturns(nil)

	response, err := ext.DeleteConnection("workspaceID", "connectionID")

	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(0, response.ErrorCode)
	assert.Equal("successful", response.Status)
}

func TestDeleteConnectionError(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMySQLExtension()
	fakeProv.DeleteUserReturns(errors.New("db creation error"))

	response, err := ext.DeleteConnection("workspaceID", "connectionID")

	assert.NotNil(err)
	assert.Nil(response)
}

func TestDeleteWorkspace(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMySQLExtension()
	fakeProv.DeleteDatabaseReturns(nil)

	response, err := ext.DeleteWorkspace("workspaceID")

	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(0, response.ErrorCode)
	assert.Equal("successful", response.Status)
}

func TestDeleteWorkspaceError(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMySQLExtension()
	fakeProv.DeleteDatabaseReturns(errors.New("delete workspace error"))

	response, err := ext.DeleteWorkspace("workspaceID")

	assert.NotNil(err)
	assert.Nil(response)
}

func TestGetConnection(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMySQLExtension()
	fakeProv.IsUserCreatedReturns(true, nil)

	response, err := ext.GetConnection("workspaceID", "connectionID")

	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(0, response.ErrorCode)
	assert.Equal("successful", response.Status)
}

func TestGetConnectionUserDoesNotExist(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMySQLExtension()
	fakeProv.IsUserCreatedReturns(false, nil)

	response, err := ext.GetConnection("workspaceID", "connectionID")

	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(404, response.ErrorCode)
	assert.Equal("failed", response.Status)
}

func TestGetConnectionError(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMySQLExtension()
	fakeProv.IsUserCreatedReturns(true, errors.New("getconnectionError"))

	response, err := ext.GetConnection("workspaceID", "connectionID")

	assert.NotNil(err)
	assert.Nil(response)
}

func TestGetWorkspace(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMySQLExtension()
	fakeProv.IsDatabaseCreatedReturns(true, nil)

	response, err := ext.GetWorkspace("workspaceID")

	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(0, response.ErrorCode)
	assert.Equal("successful", response.Status)
}

func TestGetWorkspaceDoesNotExist(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMySQLExtension()
	fakeProv.IsDatabaseCreatedReturns(false, nil)

	response, err := ext.GetWorkspace("workspaceID")

	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(404, response.ErrorCode)
	assert.Equal("failed", response.Status)
}

func TestGetWorkspaceDoesError(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMySQLExtension()
	fakeProv.IsDatabaseCreatedReturns(false, errors.New("getWorkspace error"))

	response, err := ext.GetWorkspace("workspaceID")

	assert.NotNil(err)
	assert.Nil(response)

}
