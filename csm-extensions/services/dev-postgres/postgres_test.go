package postgres

import (
	"errors"
	"testing"

	"github.com/SUSE/cf-usb-sidecar/csm-extensions/services/dev-postgres/config"
	"github.com/SUSE/cf-usb-sidecar/csm-extensions/services/dev-postgres/provisioner/provisionerfakes"
	"github.com/SUSE/go-csm-lib/extension"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("postgres-provisioner-test")

func getPostgresExtension() (extension.Extension, *provisionerfakes.FakePostgresProvisionerInterface) {
	logger = lagertest.NewTestLogger("process-controller")

	conf := config.PostgresConfig{
		User:     "testuser",
		Password: "testpass",
		Host:     "testhost",
		Port:     "5432",
		Dbname:   "db",
		Sslmode:  "ssl",
	}

	fakeProvisioner := new(provisionerfakes.FakePostgresProvisionerInterface)

	extension := NewPostgresExtension(fakeProvisioner, conf, logger)
	return extension, fakeProvisioner

}

func TestCreateConnection(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getPostgresExtension()
	fakeProv.CreateUserReturns(nil)

	response, err := ext.CreateConnection("workspace", "connection")

	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(0, response.ErrorCode)
	assert.Equal("successful", response.Status)

	creds := response.Details.(config.PostgresBindingCredentials)

	assert.NotEmpty(creds.Database)
	assert.NotEmpty(creds.Host)
	assert.NotEmpty(creds.Hostname)
	assert.NotEmpty(creds.Password)
	assert.NotEmpty(creds.Port)
	assert.NotEmpty(creds.User)
	assert.NotEmpty(creds.Password)
	assert.Equal(creds.Host, creds.Hostname)
	assert.Equal("5432", creds.Port)
	assert.Equal("testhost", creds.Host)
}

func TestCreateConnectionError(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getPostgresExtension()
	fakeProv.CreateUserReturns(errors.New("db error"))

	_, err := ext.CreateConnection("workspace", "connection")

	assert.NotNil(err)
}

func TestCreateWorkspace(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getPostgresExtension()
	fakeProv.CreateDatabaseReturns(nil)

	response, err := ext.CreateWorkspace("workspaceID")

	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(0, response.ErrorCode)
	assert.Equal("successful", response.Status)
}

func TestCreateWorkspaceError(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getPostgresExtension()
	fakeProv.CreateDatabaseReturns(errors.New("this is an error"))

	response, err := ext.CreateWorkspace("workspaceID")

	assert.NotNil(err)
	assert.Nil(response)
}

func TestDeleteConnection(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getPostgresExtension()
	fakeProv.DeleteDatabaseReturns(nil)

	response, err := ext.DeleteConnection("workspaceID", "connectionID")

	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(0, response.ErrorCode)
	assert.Equal("successful", response.Status)
}

func TestDeleteConnectionError(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getPostgresExtension()
	fakeProv.DeleteUserReturns(errors.New("db creation error"))

	response, err := ext.DeleteConnection("workspaceID", "connectionID")

	assert.NotNil(err)
	assert.Nil(response)
}

func TestDeleteWorkspace(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getPostgresExtension()
	fakeProv.DeleteDatabaseReturns(nil)

	response, err := ext.DeleteWorkspace("workspaceID")

	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(0, response.ErrorCode)
	assert.Equal("successful", response.Status)
}

func TestDeleteWorkspaceError(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getPostgresExtension()
	fakeProv.DeleteDatabaseReturns(errors.New("delete workspace error"))

	response, err := ext.DeleteWorkspace("workspaceID")

	assert.NotNil(err)
	assert.Nil(response)
}

func TestGetConnection(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getPostgresExtension()
	fakeProv.UserExistsReturns(true, nil)

	response, err := ext.GetConnection("workspaceID", "connectionID")

	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(0, response.ErrorCode)
	assert.Equal("successful", response.Status)
}

func TestGetConnectionUserDoesNotExist(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getPostgresExtension()
	fakeProv.UserExistsReturns(false, nil)

	response, err := ext.GetConnection("workspaceID", "connectionID")

	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(404, response.ErrorCode)
	assert.Equal("failed", response.Status)
}

func TestGetConnectionError(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getPostgresExtension()
	fakeProv.UserExistsReturns(true, errors.New("getconnectionError"))

	response, err := ext.GetConnection("workspaceID", "connectionID")

	assert.NotNil(err)
	assert.Nil(response)
}

func TestGetWorkspace(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getPostgresExtension()
	fakeProv.DatabaseExistsReturns(true, nil)

	response, err := ext.GetWorkspace("workspaceID")

	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(0, response.ErrorCode)
	assert.Equal("successful", response.Status)
}

func TestGetWorkspaceDoesNotExist(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getPostgresExtension()
	fakeProv.DatabaseExistsReturns(false, nil)

	response, err := ext.GetWorkspace("workspaceID")

	assert.Nil(err)
	assert.NotNil(response)
	assert.Equal(404, response.ErrorCode)
	assert.Equal("failed", response.Status)
}

func TestGetWorkspaceDoesError(t *testing.T) {
	assert := assert.New(t)

	ext, fakeProv := getPostgresExtension()
	fakeProv.DatabaseExistsReturns(false, errors.New("getWorkspace error"))

	response, err := ext.GetWorkspace("workspaceID")

	assert.NotNil(err)
	assert.Nil(response)

}
