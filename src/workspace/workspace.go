package workspace

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"

	"strings"

	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager/models"
	"github.com/hpcloud/catalog-service-manager/src/common"
	"github.com/hpcloud/catalog-service-manager/src/common/utils"
	"github.com/pivotal-golang/lager"
)

type CSMWorkspace struct {
	common.CSMWorkspaceInterface
	Logger     lager.Logger
	Config     *common.ServiceManagerConfiguration
	FileHelper utils.CSMFileHelperInterface
}

func NewCSMWorkspace(logger lager.Logger,
	config *common.ServiceManagerConfiguration,
	fileHelper utils.CSMFileHelperInterface) *CSMWorkspace {
	return &CSMWorkspace{Logger: logger.Session("CSM-Workspace"), Config: config, FileHelper: fileHelper}
}

func (w *CSMWorkspace) getWorkspaceGetExtension(homePath string) (bool, *string) {
	return w.FileHelper.GetExtension(filepath.Join(homePath, "workspace", "get"))
}

func (w *CSMWorkspace) getWorkspaceCreateExtension(homePath string) (bool, *string) {
	return w.FileHelper.GetExtension(filepath.Join(homePath, "workspace", "create"))
}

func (w *CSMWorkspace) getWorkspaceDeleteExtension(homePath string) (bool, *string) {
	return w.FileHelper.GetExtension(filepath.Join(homePath, "workspace", "delete"))
}

func generateNoopResponse() *models.ServiceManagerWorkspaceResponse {
	resp := models.ServiceManagerWorkspaceResponse{
		ProcessingType: common.PROCESSING_TYPE_NONE,
		Status:         common.PROCESSING_STATUS_NONE,
	}
	return &resp
}

//create ServiceManagerWorkspaceResponse from the json we received in file
func marshalResponseFromMessage(message []byte) (*models.ServiceManagerWorkspaceResponse, *models.Error, error) {
	workspace := utils.NewWorkspace()
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

	workspace.Details = make(map[string]interface{})
	switch t := jsonresp.Details.(type) {
	default:
		workspace.Details["data"] = t
	case map[string]interface{}:
		workspace.Details = jsonresp.Details.(map[string]interface{})
	case map[string]string:
		workspace.Details = jsonresp.Details.(map[string]interface{})
	}

	workspace.Status = jsonresp.Status
	workspace.ProcessingType = "Extension"

	return &workspace, nil, nil
}

func checkParamsOk(workspaceID *string, extensionPath *string) error {
	if workspaceID == nil {
		err := errors.New("workspaceID is nil")
		return err
	}
	if extensionPath == nil {
		err := errors.New("extensionPath is nil")
		return err
	}
	return nil
}

func (w *CSMWorkspace) executeExtension(workspaceID *string, extensionPath *string, details map[string]interface{}) (*models.ServiceManagerWorkspaceResponse, *models.Error, error) {
	if err := checkParamsOk(workspaceID, extensionPath); err != nil {
		return nil, nil, err
	}

	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return nil, nil, err
	}

	detailsStr := string(detailsJSON)
	w.Logger.Info("executeExtension", lager.Data{"workspaceID": workspaceID, "extension Path ": extensionPath, "details": details})

	if success, outputFile, output := w.FileHelper.RunExtensionFileGen(*extensionPath, *workspaceID, detailsStr); success {
		w.Logger.Info("executeExtension", lager.Data{"extension execution status ": success})
		w.Logger.Debug("executeExtension", lager.Data{"extension execution Result: ": output})

		fileContent, err := utils.ReadOutputFile(outputFile, *w.Config.DEV_MODE != "true")
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
func (w *CSMWorkspace) CheckExtensions() {
	_, file := w.getWorkspaceGetExtension(*w.Config.MANAGER_HOME)
	w.Logger.Info("CheckExtensions", lager.Data{"Workspaces Get extension ": file})

	_, file = w.getWorkspaceCreateExtension(*w.Config.MANAGER_HOME)
	w.Logger.Info("CheckExtensions", lager.Data{"Workspaces Create extension ": file})

	_, file = w.getWorkspaceDeleteExtension(*w.Config.MANAGER_HOME)
	w.Logger.Info("CheckExtensions", lager.Data{"Workspaces Delete extension ": file})
}

func (w *CSMWorkspace) executeRequest(workspaceID string, requestType string, filename *string, details map[string]interface{}) (*models.ServiceManagerWorkspaceResponse, *models.Error) {
	var modelserr *models.Error
	var workspace *models.ServiceManagerWorkspaceResponse
	var err error

	workspace, modelserr, err = w.executeExtension(&workspaceID, filename, details)

	if err != nil {
		w.Logger.Error(requestType, err)
		modelserr = utils.GenerateErrorResponse(&utils.HTTP_500, err.Error())
	}

	if workspace != nil {
		workspace.Details = details
	}

	return workspace, modelserr
}

// GetWorkspace get workspaces
func (w *CSMWorkspace) GetWorkspace(workspaceID string) (*models.ServiceManagerWorkspaceResponse, *models.Error) {

	w.Logger.Info("GetWorkspace", lager.Data{"workspaceID": workspaceID})
	exists, filename := w.getWorkspaceGetExtension(*w.Config.MANAGER_HOME)

	if !exists || filename == nil {
		w.Logger.Info("GetWorkspace", lager.Data{utils.ERR_EXTENSION_NOT_FOUND: exists})
		return generateNoopResponse(), nil
	}
	return w.executeRequest(workspaceID, "GetWorkspace", filename, make(map[string]interface{}))

}

// CreateWorkspace create workspaces
func (w *CSMWorkspace) CreateWorkspace(workspaceID string, details map[string]interface{}) (*models.ServiceManagerWorkspaceResponse, *models.Error) {
	w.Logger.Info("CreateWorkspace", lager.Data{"workspaceID": workspaceID, "details": details})
	exists, filename := w.getWorkspaceCreateExtension(*w.Config.MANAGER_HOME)

	if !exists || filename == nil {
		w.Logger.Info("CreateWorkspace", lager.Data{utils.ERR_EXTENSION_NOT_FOUND: exists})
		return generateNoopResponse(), nil
	}
	return w.executeRequest(workspaceID, "CreateWorkspace", filename, details)
}

// DeleteWorkspace delete workspaces
func (w *CSMWorkspace) DeleteWorkspace(workspaceID string) (*models.ServiceManagerWorkspaceResponse, *models.Error) {

	w.Logger.Info("DeleteWorkspace", lager.Data{"workspaceID": workspaceID})
	exists, filename := w.getWorkspaceDeleteExtension(*w.Config.MANAGER_HOME)

	if !exists || filename == nil {
		w.Logger.Info("DeleteWorkspace", lager.Data{utils.ERR_EXTENSION_NOT_FOUND: exists})
		return generateNoopResponse(), nil
	}
	return w.executeRequest(workspaceID, "DeleteWorkspace", filename, make(map[string]interface{}))
}
