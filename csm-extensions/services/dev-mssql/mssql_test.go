package mssql

import (
	"errors"
	"testing"

	"github.com/hpcloud/catalog-service-manager/csm-extensions/services/dev-mssql/config"
	"github.com/hpcloud/catalog-service-manager/csm-extensions/services/dev-mssql/provisioner/provisionerfakes"
	"github.com/hpcloud/go-csm-lib/extension"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("mssql-provisioner-test")

func getMssqlExtension() (extension.Extension, *provisionerfakes.FakeMssqlProvisioner) {
	logger = lagertest.NewTestLogger("process-controller")

	conf := config.MssqlConfig{
		User: "testuser",
		Pass: "testpass",
		Host: "testhost",
		Port: "1433",
	}

	var fakeProvisioner = new(provisionerfakes.FakeMssqlProvisioner)

	extension := NewMSSQLExtension(fakeProvisioner, conf, logger)
	return extension, fakeProvisioner

}

func TestCreateConnection(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMssqlExtension()
	fakeProv.CreateUserReturns(nil)

	response, err := ext.CreateConnection("workspace", "connection")

	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(0, response.ErrorCode)
	assert.Equal("successful", response.Status)

	creds := response.Details.(config.MssqlBindingCredentials)

	assert.NotEmpty(creds.Host)
	assert.NotEmpty(creds.Hostname)
	assert.NotEmpty(creds.Password)
	assert.NotEmpty(creds.Port)
	assert.NotEmpty(creds.User)
	assert.NotEmpty(creds.Password)
	assert.Equal(creds.Host, creds.Hostname)
	assert.Equal("1433", creds.Port)
	assert.Equal("testhost", creds.Host)
}

func TestCreateConnectionError(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMssqlExtension()
	fakeProv.CreateUserReturns(errors.New("db error"))

	_, err := ext.CreateConnection("workspace", "connection")

	assert.NotNil(err)
}

func TestCreateWorkspace(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMssqlExtension()
	fakeProv.CreateDatabaseReturns(nil)

	response, err := ext.CreateWorkspace("workspaceID")

	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(0, response.ErrorCode)
	assert.Equal("successful", response.Status)
}

func TestCreateWorkspaceError(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMssqlExtension()
	fakeProv.CreateDatabaseReturns(errors.New("this is an error"))

	response, err := ext.CreateWorkspace("workspaceID")

	assert.NotNil(err)
	assert.Nil(response)
}

func TestDeleteConnection(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMssqlExtension()
	fakeProv.DeleteDatabaseReturns(nil)

	response, err := ext.DeleteConnection("workspaceID", "connectionID")

	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(0, response.ErrorCode)
	assert.Equal("successful", response.Status)
}

func TestDeleteConnectionError(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMssqlExtension()
	fakeProv.DeleteUserReturns(errors.New("db creation error"))

	response, err := ext.DeleteConnection("workspaceID", "connectionID")

	assert.NotNil(err)
	assert.Nil(response)
}

func TestDeleteWorkspace(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMssqlExtension()
	fakeProv.DeleteDatabaseReturns(nil)

	response, err := ext.DeleteWorkspace("workspaceID")

	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(0, response.ErrorCode)
	assert.Equal("successful", response.Status)
}

func TestDeleteWorkspaceError(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMssqlExtension()
	fakeProv.DeleteDatabaseReturns(errors.New("delete workspace error"))

	response, err := ext.DeleteWorkspace("workspaceID")

	assert.NotNil(err)
	assert.Nil(response)
}

func TestGetConnection(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMssqlExtension()
	fakeProv.IsUserCreatedReturns(true, nil)

	response, err := ext.GetConnection("workspaceID", "connectionID")

	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(0, response.ErrorCode)
	assert.Equal("successful", response.Status)
}

func TestGetConnectionUserDoesNotExist(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMssqlExtension()
	fakeProv.IsUserCreatedReturns(false, nil)

	response, err := ext.GetConnection("workspaceID", "connectionID")

	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(404, response.ErrorCode)
	assert.Equal("failed", response.Status)
	assert.Equal("Connection does not exist", response.ErrorMessage)
}

func TestGetConnectionError(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMssqlExtension()
	fakeProv.IsUserCreatedReturns(true, errors.New("getconnectionError"))

	response, err := ext.GetConnection("workspaceID", "connectionID")

	assert.NotNil(err)
	assert.Nil(response)
}

func TestGetWorkspace(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMssqlExtension()
	fakeProv.IsDatabaseCreatedReturns(true, nil)

	response, err := ext.GetWorkspace("workspaceID")

	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(0, response.ErrorCode)
	assert.Equal("successful", response.Status)
}

func TestGetWorkspaceDoesNotExist(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMssqlExtension()
	fakeProv.IsDatabaseCreatedReturns(false, nil)

	response, err := ext.GetWorkspace("workspaceID")

	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(404, response.ErrorCode)
	assert.Equal("failed", response.Status)
	assert.Equal("Workspace does not exist", response.ErrorMessage)
}

func TestGetWorkspaceDoesError(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getMssqlExtension()
	fakeProv.IsDatabaseCreatedReturns(false, errors.New("getWorkspace error"))

	response, err := ext.GetWorkspace("workspaceID")

	assert.NotNil(err)
	assert.Nil(response)

}
