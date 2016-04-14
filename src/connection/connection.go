package connection

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"errors"
	"strings"

	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager/models"
	"github.com/hpcloud/catalog-service-manager/src/common"
	"github.com/hpcloud/catalog-service-manager/src/common/utils"
	"github.com/pivotal-golang/lager"
)

type CSMConnection struct {
	common.CSMSetupInterface
	Logger     lager.Logger
	Config     *common.ServiceManagerConfiguration
	FileHelper utils.CSMFileHelperInterface
}

// NewCSMConnection creates CSMConnection
func NewCSMConnection(logger lager.Logger, config *common.ServiceManagerConfiguration, fileHelper utils.CSMFileHelperInterface) *CSMConnection {
	return &CSMConnection{Logger: logger.Session("CSM-Workspace"), Config: config, FileHelper: fileHelper}
}

func (c *CSMConnection) getConnectionsGetExtension(homePath string) (bool, *string) {
	return c.FileHelper.GetExtension(homePath + "connection/get")
}

func (c *CSMConnection) getConnectionsCreateExtension(homePath string) (bool, *string) {
	return c.FileHelper.GetExtension(homePath + "connection/create")
}

func (c *CSMConnection) getConnectionsDeleteExtension(homePath string) (bool, *string) {
	return c.FileHelper.GetExtension(homePath + "connection/delete")
}

func (c *CSMConnection) executeExtension(workspaceID *string, connectionID *string, extensionPath *string, connection *models.ServiceManagerConnectionResponse) {
    if workspaceID == nil {
        c.Logger.Error("executeExtension", errors.New("workspaceID is nil"))
        return
    }
    if connectionID == nil {
        c.Logger.Error("executeExtension", errors.New("connectionID is nil"))
        return
    }
    if extensionPath == nil {
        c.Logger.Error("executeExtension", errors.New("extensionPath is nil"))
        return
    }
	c.Logger.Info("executeExtension", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID, "extension Path ": extensionPath})
	if success, outputFile, output := c.FileHelper.RunExtensionFileGen(*extensionPath, *workspaceID, *connectionID); success {
		c.Logger.Info("executeExtension", lager.Data{"extension execution status ": success})
		c.Logger.Debug("executeExtension", lager.Data{"extension execution Result: ": output})
		if outputFile != nil {
			if *c.Config.DEV_MODE != "true" {
				// clean up if not running in dev mode
				defer os.Remove(outputFile.Name())
			}
			// checking the file size of the extension output
			if fileStat, err := os.Stat(outputFile.Name()); err == nil && fileStat != nil {
				if fileStat.Size() > 0 {
					file, e := ioutil.ReadFile(outputFile.Name())
					if e != nil {
						c.Logger.Info("executeExtension", lager.Data{"File error while reading extension output file": e})
						c.Logger.Error("executeExtension", e)
						connection.Status = common.PROCESSING_STATUS_FAILED
					}
					err := json.Unmarshal(file, &connection)
					if err != nil {
						c.Logger.Info("executeExtension", lager.Data{"Failed to parse the extension output": ""})
						c.Logger.Error("executeExtension", err)
					}
					c.Logger.Info("executeExtension", lager.Data{"extension processing status ": connection.Status})
				} else {
					// file size of extension output file is not greater than 0
					connection.Status = common.PROCESSING_STATUS_FAILED
					c.Logger.Debug("executeExtension", lager.Data{"extension output file is empty": success})
				}
			} else {
				c.Logger.Info("executeExtension", lager.Data{"File error while reading extension output file ": err})
				c.Logger.Error("executeExtension", err)
			}
		}
	} else {
		// extension couldn't be executed
		connection.Status = common.PROCESSING_STATUS_FAILED
		c.Logger.Debug("executeExtension", lager.Data{"extension execution failed ": success})
	}

	connection.ProcessingType = common.PROCESSING_TYPE_EXTENSION
	c.Logger.Debug("executeExtension", lager.Data{"Updated workspace processing type to ": common.PROCESSING_TYPE_EXTENSION})
}

// CheckExtensions checks for workspace extensions
func (c *CSMConnection) CheckExtensions() {
	_, file := c.getConnectionsGetExtension(*c.Config.MANAGER_HOME)
	c.Logger.Info("CheckExtensions", lager.Data{"Connections Get extension ": file})

	_, file = c.getConnectionsCreateExtension(*c.Config.MANAGER_HOME)
	c.Logger.Info("CheckExtensions", lager.Data{"Connections Create extension ": file})

	_, file = c.getConnectionsDeleteExtension(*c.Config.MANAGER_HOME)
	c.Logger.Info("CheckExtensions", lager.Data{"Connections Delete extension ": file})
}

// GetConnection get connections
func (c *CSMConnection) GetConnection(workspaceID string, connectionID string) *models.ServiceManagerConnectionResponse {
	c.Logger.Info("GetConnection", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})
	connection := models.ServiceManagerConnectionResponse{
		ProcessingType: common.PROCESSING_TYPE_NONE,
	}

	serviceManagerConfig := common.NewServiceManagerConfiguration()
	exists, filename := c.getConnectionsGetExtension(*serviceManagerConfig.MANAGER_HOME)
	if exists && filename != nil {
		c.executeExtension(&workspaceID, &connectionID, filename, &connection)
	} else {
		c.Logger.Info("GetConnection", lager.Data{"extension not found ": exists})
	}
	return &connection
}

// CreateConnection create connections
func (c *CSMConnection) CreateConnection(workspaceID string, createConnection *models.ServiceManagerConnectionCreateRequest) *models.ServiceManagerConnectionResponse {
	c.Logger.Info("CreateConnection", lager.Data{"workspaceID": workspaceID, "connectionID": createConnection.ConnectionID})
	connection := models.ServiceManagerConnectionResponse{
		ProcessingType: common.PROCESSING_TYPE_NONE,
	}
	serviceManagerConfig := common.NewServiceManagerConfiguration()
	exists, filename := c.getConnectionsCreateExtension(*serviceManagerConfig.MANAGER_HOME)
	if exists && filename != nil {
		c.executeExtension(&workspaceID, &createConnection.ConnectionID, filename, &connection)
	} else if (!exists) && (serviceManagerConfig.PARAMETERS != nil) {
		c.Logger.Info("GetConnection", lager.Data{"extension not found ": exists})
		parametersNameList := strings.Split(*serviceManagerConfig.PARAMETERS, " ")
		c.Logger.Info("GetConnection", lager.Data{"Parameter List ": parametersNameList})
		connection.Details = make(map[string]interface{})
		for _, parameterName := range parametersNameList {
			parameterValue, ok := os.LookupEnv(parameterName)
			if ok {
				connection.Details[parameterName] = parameterValue
			}
		}
		connection.ProcessingType = common.PROCESSING_TYPE_DEFAULT
	}
	return &connection
}

// DeleteConnection delete connections
func (c *CSMConnection) DeleteConnection(workspaceID string, connectionID string) *models.ServiceManagerConnectionResponse {
	c.Logger.Info("DeleteConnection", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})
	connection := models.ServiceManagerConnectionResponse{
		ProcessingType: common.PROCESSING_TYPE_NONE,
	}
	serviceManagerConfig := common.NewServiceManagerConfiguration()
	exists, filename := c.getConnectionsDeleteExtension(*serviceManagerConfig.MANAGER_HOME)
	if exists && filename != nil {
		c.executeExtension(&workspaceID, &connectionID, filename, &connection)
	} else {
		c.Logger.Info("DeleteConnection", lager.Data{"extension not found ": exists})
	}
	return &connection
}
