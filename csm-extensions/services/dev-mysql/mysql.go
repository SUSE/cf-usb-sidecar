package mysql

import (
	"github.com/hpcloud/catalog-service-manager/csm-extensions/services/dev-mysql/config"
	"github.com/hpcloud/catalog-service-manager/csm-extensions/services/dev-mysql/provisioner"
	"github.com/hpcloud/go-csm-lib/csm"
	"github.com/hpcloud/go-csm-lib/extension"
	"github.com/hpcloud/go-csm-lib/util"
	"github.com/pivotal-golang/lager"
)

const userSize = 16

type mysqlExtension struct {
	prov   provisioner.MySQLProvisioner
	conf   config.MySQLConfig
	logger lager.Logger
}

func NewMySQLExtension(prov provisioner.MySQLProvisioner,
	conf config.MySQLConfig, logger lager.Logger) extension.Extension {
	return &mysqlExtension{prov: prov, conf: conf, logger: logger}
}

func (e *mysqlExtension) CreateConnection(workspaceID, connectionID string) (*csm.CSMResponse, error) {
	e.logger.Info("create-connection", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})

	dbName := util.NormalizeGuid(workspaceID)

	username, err := util.GetMD5Hash(connectionID, userSize)
	if err != nil {
		return nil, err
	}

	password, err := util.SecureRandomString(32)
	if err != nil {
		return nil, err
	}

	err = e.prov.CreateUser(dbName, username, password)

	if err != nil {
		return nil, err
	}

	binding := config.MySQLBinding{
		Hostname: e.conf.Host,
		Host:     e.conf.Host,
		Port:     e.conf.Port,
		Username: username,
		User:     username,
		Password: password,
		Database: dbName,
	}

	response := csm.CreateCSMResponse(binding)
	return &response, err
}
func (e *mysqlExtension) CreateWorkspace(workspaceID string) (*csm.CSMResponse, error) {
	e.logger.Info("create-workspace", lager.Data{"workspaceID": workspaceID})
	dbName := util.NormalizeGuid(workspaceID)
	err := e.prov.CreateDatabase(dbName)
	if err != nil {
		return nil, err
	}

	response := csm.CreateCSMResponse("")

	return &response, nil
}
func (e *mysqlExtension) DeleteConnection(workspaceID, connectionID string) (*csm.CSMResponse, error) {
	e.logger.Info("delete-connection", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})

	username, err := util.GetMD5Hash(connectionID, userSize)
	if err != nil {
		return nil, err
	}

	err = e.prov.DeleteUser(username)
	if err != nil {
		return nil, err
	}

	response := csm.CreateCSMResponse("")

	return &response, nil
}
func (e *mysqlExtension) DeleteWorkspace(workspaceID string) (*csm.CSMResponse, error) {
	e.logger.Info("delete-workspace", lager.Data{"workspaceID": workspaceID})

	database := util.NormalizeGuid(workspaceID)
	err := e.prov.DeleteDatabase(database)
	if err != nil {
		return nil, err
	}

	response := csm.CreateCSMResponse("")

	return &response, nil
}
func (e *mysqlExtension) GetConnection(workspaceID, connectionID string) (*csm.CSMResponse, error) {
	username, err := util.GetMD5Hash(connectionID, userSize)
	if err != nil {
		return nil, err
	}

	exists, err := e.prov.IsUserCreated(username)
	if err != nil {
		return nil, err
	}

	var response csm.CSMResponse

	if exists {
		response = csm.CreateCSMResponse("")
	} else {
		response = csm.CreateCSMErrorResponse(404, "Connection does not exist")
	}

	return &response, nil
}
func (e *mysqlExtension) GetWorkspace(workspaceID string) (*csm.CSMResponse, error) {
	database := util.NormalizeGuid(workspaceID)

	exists, err := e.prov.IsDatabaseCreated(database)
	if err != nil {
		return nil, err
	}

	response := csm.CSMResponse{}

	if exists {
		response = csm.CreateCSMResponse("")
	} else {
		response = csm.CreateCSMErrorResponse(404, "Workspace does not exist")
	}

	return &response, nil
}

func (e *mysqlExtension) GetStatus() (*csm.CSMResponse, error) {
	response := csm.CSMResponse{}

	_, err := e.prov.Query("SHOW DATABASES")
	if err != nil {
		response.Status = "failed"
		response.ErrorMessage = "Could not connect to database"
		response.Diagnostics = append(response.Diagnostics, &csm.StatusDiagnostic{Name: "Database", Message: err.Error(), Description: "Server reply", Status: "failed"})

		return &response, err
	}
	response.Status = "successful"
	return &response, nil
}
