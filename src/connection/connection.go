package connection

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager/models"
	"github.com/hpcloud/catalog-service-manager/src/common"
	"github.com/hpcloud/catalog-service-manager/src/common/utils"
	"github.com/pivotal-golang/lager"
)

// CSMConnection object for managing the connection.
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
	return c.FileHelper.GetExtension(filepath.Join(homePath, "connection", "get"))
}

func (c *CSMConnection) getConnectionsCreateExtension(homePath string) (bool, *string) {
	return c.FileHelper.GetExtension(filepath.Join(homePath, "connection", "create"))
}

func (c *CSMConnection) getConnectionsDeleteExtension(homePath string) (bool, *string) {
	return c.FileHelper.GetExtension(filepath.Join(homePath, "connection", "delete"))
}

//create ServiceManagerConnectionResponse from the json we received in file
func marshalResponseFromMessage(message []byte) (*models.ServiceManagerConnectionResponse, *models.Error, error) {
	connection := utils.NewConnection()
	jsonresp := utils.JsonResponse{}
	if len(message) == 0 {
		return nil, nil, errors.New("Empty response")
	}
	err := jsonresp.Unmarshal(message)

	if err != nil {
		return nil, nil, err
	}
	if strings.ToLower(jsonresp.Status) != "successful" { //the extension is giving us an error responses
		var code int64
		var message string
		if jsonresp.ErrorCode == 0 {
			code = utils.HTTP_500
		} else {
			code = int64(jsonresp.ErrorCode)
		}

		message = jsonresp.ErrorMessage

		return nil, utils.GenerateErrorResponse(&code, message), nil

	}

	connection.Details = make(map[string]interface{})
	switch t := jsonresp.Details.(type) {
	default:
		connection.Details["data"] = t
	case map[string]interface{}:
		connection.Details = jsonresp.Details.(map[string]interface{})
	}
	//connection.Details = jsonresp.Details.(map[string]interface{})
	connection.Status = "successful"
	connection.ProcessingType = "Extension"

	return &connection, nil, nil
}

func checkParamsOk(workspaceID *string, connectionID *string, extensionPath *string) error {
	if workspaceID == nil {
		err := errors.New("workspaceID is nil")
		return err
	}
	if connectionID == nil {
		err := errors.New("connectionID is nil")
		return err
	}
	if extensionPath == nil {
		err := errors.New("extensionPath is nil")
		return err
	}
	return nil
}

func (c *CSMConnection) executeExtension(workspaceID *string, connectionID *string, details map[string]interface{}, extensionPath *string) (*models.ServiceManagerConnectionResponse, *models.Error, error) {
	if err := checkParamsOk(workspaceID, connectionID, extensionPath); err != nil {
		return nil, nil, err
	}

	detailsStr := ""

	if details != nil {
		detailsJSON, err := json.Marshal(details)
		if err != nil {
			return nil, nil, err
		}

		detailsStr = string(detailsJSON)
	}

	c.Logger.Info("executeExtension", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID, "extension Path ": extensionPath, "details": details})
	if success, outputFile, output := c.FileHelper.RunExtensionFileGen(*extensionPath, *workspaceID, *connectionID, detailsStr); success {
		c.Logger.Info("executeExtension", lager.Data{"extension execution status ": success})
		c.Logger.Debug("executeExtension", lager.Data{"extension execution Result: ": output})

		fileContent, err := utils.ReadOutputFile(outputFile, *c.Config.DEV_MODE != "true")
		if err != nil {
			return nil, nil, errors.New(fmt.Sprintf("Error reading response from extension: %s", err.Error()))
		}
		return marshalResponseFromMessage(fileContent)
	} else {
		// extension couldn't be executed, returned an error or timedout
		//first we check for timeout (success=false,  output==nil)
		if output == nil {
			return nil, utils.GenerateErrorResponse(&utils.HTTP_408, utils.ERR_TIMEOUT), nil
		}
		//else it means that the extension did not return a zero code	 ("success = false, output != nil)
		err := errors.New(*output)
		return nil, nil, err
	}
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

func (c *CSMConnection) executeRequest(workspaceID string, connectionID string, details map[string]interface{}, requestType string, filename *string) (*models.ServiceManagerConnectionResponse, *models.Error) {
	var modelserr *models.Error
	var connection *models.ServiceManagerConnectionResponse
	var err error

	connection, modelserr, err = c.executeExtension(&workspaceID, &connectionID, details, filename)
	if err != nil {
		c.Logger.Error(requestType, err)
		modelserr = utils.GenerateErrorResponse(&utils.HTTP_500, err.Error())
	}

	if connection != nil {
		if connection.Details == nil {
			connection.Details = details
		} else {
			//the "data" item is added if the response of the extension is a string or nil
			if _, ok := connection.Details["data"]; ok {
				//if the response is nil, set the details
				if connection.Details["data"] == nil {
					connection.Details = details
				}
			}
		}
	}

	return connection, modelserr
}

func generateNoopResponse() *models.ServiceManagerConnectionResponse {
	resp := models.ServiceManagerConnectionResponse{
		ProcessingType: common.PROCESSING_TYPE_NONE,
		Status:         common.PROCESSING_STATUS_NONE,
	}
	return &resp
}

// GetConnection get connections
func (c *CSMConnection) GetConnection(workspaceID string, connectionID string) (*models.ServiceManagerConnectionResponse, *models.Error) {
	c.Logger.Info("GetConnection", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})

	serviceManagerConfig := common.NewServiceManagerConfiguration()
	exists, filename := c.getConnectionsGetExtension(*serviceManagerConfig.MANAGER_HOME)
	if !exists || filename == nil {
		c.Logger.Info("GetConnection", lager.Data{utils.ERR_EXTENSION_NOT_FOUND: exists})
		return generateNoopResponse(), nil
	}
	return c.executeRequest(workspaceID, connectionID, make(map[string]interface{}), "GetConnection", filename)
}

// CreateConnection create connections
func (c *CSMConnection) CreateConnection(workspaceID string, connectionID string, details map[string]interface{}) (*models.ServiceManagerConnectionResponse, *models.Error) {
	c.Logger.Info("CreateConnection", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID, "details": details})

	serviceManagerConfig := common.NewServiceManagerConfiguration()
	exists, filename := c.getConnectionsCreateExtension(*serviceManagerConfig.MANAGER_HOME)
	if (!exists) && (serviceManagerConfig.PARAMETERS != nil) {
		connection := utils.NewConnection()
		c.Logger.Info("GetConnection", lager.Data{"Extension not found ": exists})
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
		connection.Status = common.PROCESSING_STATUS_SUCCESSFUL
		return &connection, nil

	} else if !exists || filename == nil {
		c.Logger.Info("CreateConnection", lager.Data{utils.ERR_EXTENSION_NOT_FOUND: exists})
		return generateNoopResponse(), nil
	}
	return c.executeRequest(workspaceID, connectionID, details, "CreateConnection", filename)
}

// DeleteConnection delete connections
func (c *CSMConnection) DeleteConnection(workspaceID string, connectionID string) (*models.ServiceManagerConnectionResponse, *models.Error) {
	c.Logger.Info("DeleteConnection", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})

	serviceManagerConfig := common.NewServiceManagerConfiguration()
	exists, filename := c.getConnectionsDeleteExtension(*serviceManagerConfig.MANAGER_HOME)
	if !exists || filename == nil {
		c.Logger.Info("DeleteConnection", lager.Data{utils.ERR_EXTENSION_NOT_FOUND: exists})
		return generateNoopResponse(), nil
	}
	return c.executeRequest(workspaceID, connectionID, make(map[string]interface{}), "DeleteConnection", filename)
}
