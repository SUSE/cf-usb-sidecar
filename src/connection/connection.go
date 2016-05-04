package connection

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"

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
func marshalResponseFromMessage(message []byte, ok_resp int) (*models.ServiceManagerConnectionResponse, *models.Error, error) {
	connection := utils.NewConnection()
	jsonresp := utils.JsonResponse{}
	err := json.Unmarshal(message, &jsonresp)
	if err != nil {
		return nil, nil, errors.New(fmt.Sprintf("Invalid json response from extension: %s", err.Error()))
	}
	if jsonresp.HttpCode != ok_resp { //the extension is giving us an error response
		if jsonresp.HttpCode == 0 {
			err = errors.New("invalid response received from extension")
			return nil, nil, err
		}
		code := int64(jsonresp.HttpCode)
		//if the client wants us to use a specific error code we must create a specific error here
		//all the other errors will have code 500
		return nil, utils.GenerateErrorResponse(&code, fmt.Sprintf("%v", jsonresp.Details)), nil
	}
	connection.Details = make(map[string]interface{})
	connection.Details["data"] = jsonresp.Details
	connection.Status = jsonresp.Status
	connection.ProcessingType = jsonresp.ProcessingType

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

func (c *CSMConnection) executeExtension(workspaceID *string, connectionID *string, extensionPath *string, ok_resp int) (*models.ServiceManagerConnectionResponse, *models.Error, error) {
	if err := checkParamsOk(workspaceID, connectionID, extensionPath); err != nil {
		return nil, nil, err
	}
	c.Logger.Info("executeExtension", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID, "extension Path ": extensionPath})
	if success, outputFile, output := c.FileHelper.RunExtensionFileGen(*extensionPath, *workspaceID, *connectionID); success {
		c.Logger.Info("executeExtension", lager.Data{"extension execution status ": success})
		c.Logger.Debug("executeExtension", lager.Data{"extension execution Result: ": output})

		fileContent, err := utils.ReadOutputFile(outputFile, *c.Config.DEV_MODE != "true")

		if err != nil {
			return nil, nil, errors.New(fmt.Sprintf("Invalid json response from extension: %", err.Error()))
		}

		return marshalResponseFromMessage(fileContent, ok_resp)

	} else {
		// extension couldn't be executed, returned an error or timedout
		//first we check for timeout (success=false,  output==nil)
		if output == nil {
			return nil, utils.GenerateErrorResponse(&utils.HTTP_408, "Timeout while executing the extension. The extension did not respond in a reasonable ammount of time."), nil
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

func (w *CSMConnection) executeRequest(workspaceID string, connectionID string, requestType string, ok_resp int, filename *string) (*models.ServiceManagerConnectionResponse, *models.Error) {
	var modelserr *models.Error
	var connection *models.ServiceManagerConnectionResponse
	var err error
	connection, modelserr, err = w.executeExtension(&workspaceID, &connectionID, filename, common.GET_WORKSPACE_OK_RESPONSE)
	if err != nil {
		w.Logger.Error(requestType, err)
		modelserr = utils.GenerateErrorResponse(&utils.HTTP_500, err.Error())
	}
	return connection, modelserr
}

// GetConnection get connections
func (c *CSMConnection) GetConnection(workspaceID string, connectionID string) (*models.ServiceManagerConnectionResponse, *models.Error) {
	c.Logger.Info("GetConnection", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})

	serviceManagerConfig := common.NewServiceManagerConfiguration()
	exists, filename := c.getConnectionsGetExtension(*serviceManagerConfig.MANAGER_HOME)
	if !exists || filename == nil {
		c.Logger.Info("GetConnection", lager.Data{"extension not found ": exists})
		return nil, utils.GenerateErrorResponse(&utils.HTTP_500, "extension not found")
	}
	return c.executeRequest(workspaceID, connectionID, "GetConnection", common.GET_CONNECTION_OK_RESPONSE, filename)
}

// CreateConnection create connections
func (c *CSMConnection) CreateConnection(workspaceID string, connectionID string) (*models.ServiceManagerConnectionResponse, *models.Error) {
	c.Logger.Info("CreateConnection", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})

	serviceManagerConfig := common.NewServiceManagerConfiguration()
	exists, filename := c.getConnectionsCreateExtension(*serviceManagerConfig.MANAGER_HOME)
	if !exists || filename == nil {
		c.Logger.Info("CreateConnection", lager.Data{"extension not found ": exists})
		return nil, utils.GenerateErrorResponse(&utils.HTTP_500, "extension not found")
	}
	return c.executeRequest(workspaceID, connectionID, "CreateConnection", common.CREATE_CONNECTION_OK_RESPONSE, filename)
}

// DeleteConnection delete connections
func (c *CSMConnection) DeleteConnection(workspaceID string, connectionID string) (*models.ServiceManagerConnectionResponse, *models.Error) {
	c.Logger.Info("DeleteConnection", lager.Data{"workspaceID": workspaceID, "connectionID": connectionID})

	serviceManagerConfig := common.NewServiceManagerConfiguration()
	exists, filename := c.getConnectionsDeleteExtension(*serviceManagerConfig.MANAGER_HOME)
	if !exists || filename == nil {
		c.Logger.Info("DeleteConnection", lager.Data{"extension not found ": exists})
		return nil, utils.GenerateErrorResponse(&utils.HTTP_500, "extension not found")
	}
	return c.executeRequest(workspaceID, connectionID, "DeleteConnection", common.DELETE_CONNECTION_OK_RESPONSE, filename)
}
