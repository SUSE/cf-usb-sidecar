package workspace

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

//create ServiceManagerWorkspaceResponse from the json we received in file
func marshalResponseFromMessage(message []byte, ok_resp int) (*models.ServiceManagerWorkspaceResponse, *models.Error, error) {
	workspace := utils.NewWorkspace()
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
	workspace.Details = make(map[string]interface{})
	workspace.Details["data"] = jsonresp.Details
	workspace.Status = jsonresp.Status
	workspace.ProcessingType = jsonresp.ProcessingType

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

func (w *CSMWorkspace) executeExtension(workspaceID *string, extensionPath *string, ok_resp int) (*models.ServiceManagerWorkspaceResponse, *models.Error, error) {
	if err := checkParamsOk(workspaceID, extensionPath); err != nil {
		return nil, nil, err
	}
	w.Logger.Info("executeExtension", lager.Data{"workspaceID": workspaceID, "extension Path ": extensionPath})

	if success, outputFile, output := w.FileHelper.RunExtensionFileGen(*extensionPath, *workspaceID); success {
		w.Logger.Info("executeExtension", lager.Data{"extension execution status ": success})
		w.Logger.Debug("executeExtension", lager.Data{"extension execution Result: ": output})

		fileContent, err := utils.ReadOutputFile(outputFile, *w.Config.DEV_MODE != "true")
		return marshalResponseFromMessage(fileContent, ok_resp)

		if err != nil {
			return nil, nil, err
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
func (w *CSMWorkspace) CheckExtensions() {
	_, file := w.getWorkspaceGetExtension(*w.Config.MANAGER_HOME)
	w.Logger.Info("CheckExtensions", lager.Data{"Workspaces Get extension ": file})

	_, file = w.getWorkspaceCreateExtension(*w.Config.MANAGER_HOME)
	w.Logger.Info("CheckExtensions", lager.Data{"Workspaces Create extension ": file})

	_, file = w.getWorkspaceDeleteExtension(*w.Config.MANAGER_HOME)
	w.Logger.Info("CheckExtensions", lager.Data{"Workspaces Delete extension ": file})
}

func (w *CSMWorkspace) executeRequest(workspaceID string, requestType string, ok_resp int, filename *string) (*models.ServiceManagerWorkspaceResponse, *models.Error) {
	var modelserr *models.Error
	var workspace *models.ServiceManagerWorkspaceResponse
	var err error

	workspace, modelserr, err = w.executeExtension(&workspaceID, filename, ok_resp)

	if err != nil {
		w.Logger.Error(requestType, err)
		modelserr = utils.GenerateErrorResponse(&utils.HTTP_500, err.Error())
	}
	return workspace, modelserr
}

// GetWorkspace get workspaces
func (w *CSMWorkspace) GetWorkspace(workspaceID string) (*models.ServiceManagerWorkspaceResponse, *models.Error) {

	w.Logger.Info("GetWorkspace", lager.Data{"workspaceID": workspaceID})
	exists, filename := w.getWorkspaceGetExtension(*w.Config.MANAGER_HOME)

	if !exists || filename == nil {
		w.Logger.Info("GetWorkspace", lager.Data{"extension not found ": exists})
		return nil, utils.GenerateErrorResponse(&utils.HTTP_500, "extension not found")
	}
	return w.executeRequest(workspaceID, "GetWorkspace", common.GET_WORKSPACE_OK_RESPONSE, filename)

}

// CreateWorkspace create workspaces
func (w *CSMWorkspace) CreateWorkspace(workspaceID string) (*models.ServiceManagerWorkspaceResponse, *models.Error) {
	w.Logger.Info("CreateWorkspace", lager.Data{"workspaceID": workspaceID})
	exists, filename := w.getWorkspaceCreateExtension(*w.Config.MANAGER_HOME)

	if !exists || filename == nil {
		w.Logger.Info("CreateWorkspace", lager.Data{"extension not found ": exists})
		return nil, utils.GenerateErrorResponse(&utils.HTTP_500, "extension not found")
	}
	return w.executeRequest(workspaceID, "CreateWorkspace", common.CREATE_WORKSPACE_OK_RESPONSE, filename)
}

// DeleteWorkspace delete workspaces
func (w *CSMWorkspace) DeleteWorkspace(workspaceID string) (*models.ServiceManagerWorkspaceResponse, *models.Error) {

	w.Logger.Info("DeleteWorkspace", lager.Data{"workspaceID": workspaceID})
	exists, filename := w.getWorkspaceDeleteExtension(*w.Config.MANAGER_HOME)

	if !exists || filename == nil {
		w.Logger.Info("DeleteWorkspace", lager.Data{"extension not found ": exists})
		return nil, utils.GenerateErrorResponse(&utils.HTTP_500, "extension not found")

	}
	return w.executeRequest(workspaceID, "DeleteWorkspace", common.DELETE_WORKSPACE_OK_RESPONSE, filename)
}
