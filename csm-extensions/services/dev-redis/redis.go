package redis

import (
	"fmt"

	"github.com/hpcloud/catalog-service-manager/csm-extensions/services/dev-redis/config"
	"github.com/hpcloud/catalog-service-manager/csm-extensions/services/dev-redis/provisioner"
	"github.com/hpcloud/go-csm-lib/csm"
	"github.com/hpcloud/go-csm-lib/extension"
	"github.com/hpcloud/go-csm-lib/util"
	"github.com/pivotal-golang/lager"
)

type redisExtension struct {
	conf   config.RedisConfig
	prov   provisioner.RedisProvisionerInterface
	logger lager.Logger
}

func NewRedisExtension(prov provisioner.RedisProvisionerInterface, conf config.RedisConfig, logger lager.Logger) extension.Extension {
	return &redisExtension{prov: prov, conf: conf, logger: logger}
}

func (e *redisExtension) CreateConnection(workspaceID, connectionID string) (*csm.CSMResponse, error) {
	e.logger.Info("create-connection", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})

	dbName := util.NormalizeGuid(workspaceID)

	credentials, err := e.prov.GetCredentials(dbName)
	if err != nil {
		return nil, err
	}

	binding := config.RedisBinding{
		Password: credentials["password"],
		Port:     credentials["port"],
		Host:     credentials["host"],
		Hostname: credentials["host"],
		Uri:      fmt.Sprintf("redis://:%s@%s:%s/", credentials["password"], credentials["host"], credentials["port"]),
	}

	response := csm.CreateCSMResponse(binding)
	return &response, err
}
func (e *redisExtension) CreateWorkspace(workspaceID string) (*csm.CSMResponse, error) {
	e.logger.Info("create-workspace", lager.Data{"workspaceID": workspaceID})
	dbName := util.NormalizeGuid(workspaceID)
	err := e.prov.CreateContainer(dbName)
	if err != nil {
		return nil, err
	}

	response := csm.CreateCSMResponse("")

	return &response, nil
}
func (e *redisExtension) DeleteConnection(workspaceID, connectionID string) (*csm.CSMResponse, error) {
	e.logger.Info("delete-connection", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})

	response := csm.CreateCSMResponse("")

	return &response, nil
}
func (e *redisExtension) DeleteWorkspace(workspaceID string) (*csm.CSMResponse, error) {
	e.logger.Info("delete-workspace", lager.Data{"workspaceID": workspaceID})

	database := util.NormalizeGuid(workspaceID)
	err := e.prov.DeleteContainer(database)
	if err != nil {
		return nil, err
	}

	response := csm.CreateCSMResponse("")

	return &response, nil
}
func (e *redisExtension) GetConnection(workspaceID, connectionID string) (*csm.CSMResponse, error) {

	response := csm.CreateCSMErrorResponse(404, "Connection does not exist")

	return &response, nil
}
func (e *redisExtension) GetWorkspace(workspaceID string) (*csm.CSMResponse, error) {
	database := util.NormalizeGuid(workspaceID)

	exists, err := e.prov.ContainerExists(database)
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

func (e *redisExtension) GetStatus() (*csm.CSMResponse, error) {
	response := csm.CSMResponse{}
	_, err := e.prov.ContainerExists("test")

	if err != nil {
		response.Status = "failed"
		response.ErrorMessage = "Could not connect to redis docker host"
		response.Diagnostics = append(response.Diagnostics, &csm.StatusDiagnostic{Name: "Database", Message: err.Error(), Description: "Server reply", Status: "failed"})

		return &response, err
	}
	response.Status = "successful"
	return &response, nil
}
