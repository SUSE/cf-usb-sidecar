package mongo

import (
	"github.com/hpcloud/catalog-service-manager/services/dev-mongo/config"
	"github.com/hpcloud/catalog-service-manager/services/dev-mongo/provisioner"
	"github.com/hpcloud/go-csm-lib/csm"
	"github.com/hpcloud/go-csm-lib/csm/status"
	"github.com/hpcloud/go-csm-lib/extension"
	"github.com/hpcloud/go-csm-lib/util"
	"github.com/pivotal-golang/lager"
)

const userSize = 16

type mongoExtension struct {
	prov   provisioner.MongoProvisionerInterface
	conf   config.MongoDriverConfig
	logger lager.Logger
}

func NewMongoExtension(prov provisioner.MongoProvisionerInterface,
	conf config.MongoDriverConfig, logger lager.Logger) extension.Extension {
	return &mongoExtension{prov: prov, conf: conf, logger: logger}
}

func (e *mongoExtension) CreateConnection(workspaceID, connectionID string) (*csm.CSMResponse, error) {
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

	binding := config.MongoBindingCredentials{
		Hostname: e.conf.Host,
		Host:     e.conf.Host,
		Port:     e.conf.Port,
		Username: username,
		Password: password,
		Uri:      config.GenerateConnectionString(e.conf.Host, e.conf.Port, dbName, username, password),
		Name:     dbName,
		Db:       dbName,
	}

	response := csm.NewCSMResponse(200, binding, status.Successful)
	return &response, err
}
func (e *mongoExtension) CreateWorkspace(workspaceID string) (*csm.CSMResponse, error) {
	e.logger.Info("create-workspace", lager.Data{"workspaceID": workspaceID})
	dbName := util.NormalizeGuid(workspaceID)
	err := e.prov.CreateDatabase(dbName)
	if err != nil {
		return nil, err
	}
	response := csm.NewCSMResponse(200, "", status.Successful)

	return &response, nil
}
func (e *mongoExtension) DeleteConnection(workspaceID, connectionID string) (*csm.CSMResponse, error) {
	e.logger.Info("delete-connection", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})

	username, err := util.GetMD5Hash(connectionID, userSize)
	if err != nil {
		return nil, err
	}

	dbName := util.NormalizeGuid(workspaceID)

	err = e.prov.DeleteUser(dbName, username)
	if err != nil {
		return nil, err
	}

	response := csm.NewCSMResponse(200, "", status.Successful)

	return &response, nil
}
func (e *mongoExtension) DeleteWorkspace(workspaceID string) (*csm.CSMResponse, error) {
	e.logger.Info("delete-workspace", lager.Data{"workspaceID": workspaceID})

	database := util.NormalizeGuid(workspaceID)

	err := e.prov.DeleteDatabase(database)
	if err != nil {
		return nil, err
	}

	response := csm.NewCSMResponse(200, "", status.Successful)

	return &response, nil
}
func (e *mongoExtension) GetConnection(workspaceID, connectionID string) (*csm.CSMResponse, error) {
	username, err := util.GetMD5Hash(connectionID, userSize)
	if err != nil {
		return nil, err
	}
	dbName := util.NormalizeGuid(workspaceID)

	exists, err := e.prov.IsUserCreated(dbName, username)
	if err != nil {
		return nil, err
	}

	response := csm.CSMResponse{}

	if exists {
		response = csm.NewCSMResponse(200, "", status.Successful)
	} else {
		response = csm.NewCSMResponse(404, "", status.Successful)
	}

	return &response, nil
}
func (e *mongoExtension) GetWorkspace(workspaceID string) (*csm.CSMResponse, error) {
	database := util.NormalizeGuid(workspaceID)

	exists, err := e.prov.IsDatabaseCreated(database)
	if err != nil {
		return nil, err
	}

	response := csm.CSMResponse{}

	if exists {
		response = csm.NewCSMResponse(200, "", status.Successful)
	} else {
		response = csm.NewCSMResponse(404, "", status.Successful)
	}

	return &response, nil
}
