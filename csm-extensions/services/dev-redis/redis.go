package redis

import (
	"fmt"

	"github.com/hpcloud/go-csm-lib/csm"
	"github.com/hpcloud/go-csm-lib/csm/status"
	"github.com/hpcloud/go-csm-lib/extension"
	"github.com/hpcloud/catalog-service-manager/csm-extensions/services/dev-redis/config"
	"github.com/hpcloud/catalog-service-manager/csm-extensions/services/dev-redis/provisioner"
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

	response := csm.NewCSMResponse(200, binding, status.Successful)
	return &response, err
}
func (e *redisExtension) CreateWorkspace(workspaceID string) (*csm.CSMResponse, error) {
	e.logger.Info("create-workspace", lager.Data{"workspaceID": workspaceID})
	dbName := util.NormalizeGuid(workspaceID)
	err := e.prov.CreateContainer(dbName)
	if err != nil {
		return nil, err
	}

	response := csm.NewCSMResponse(200, "", status.Successful)

	return &response, nil
}
func (e *redisExtension) DeleteConnection(workspaceID, connectionID string) (*csm.CSMResponse, error) {
	e.logger.Info("delete-connection", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})

	response := csm.NewCSMResponse(200, "", status.Successful)

	return &response, nil
}
func (e *redisExtension) DeleteWorkspace(workspaceID string) (*csm.CSMResponse, error) {
	e.logger.Info("delete-workspace", lager.Data{"workspaceID": workspaceID})

	database := util.NormalizeGuid(workspaceID)
	err := e.prov.DeleteContainer(database)
	if err != nil {
		return nil, err
	}

	response := csm.NewCSMResponse(200, "", status.Successful)

	return &response, nil
}
func (e *redisExtension) GetConnection(workspaceID, connectionID string) (*csm.CSMResponse, error) {

	response := csm.NewCSMResponse(200, "", status.Successful)

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
		response = csm.NewCSMResponse(200, "", status.Successful)
	} else {
		response = csm.NewCSMResponse(404, "", status.Successful)
	}

	return &response, nil
}
