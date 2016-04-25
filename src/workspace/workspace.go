package workspace

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"errors"
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
	return w.FileHelper.GetExtension(filepath.Join(homePath, "workspace","get"))
}

func (w *CSMWorkspace) getWorkspaceCreateExtension(homePath string) (bool, *string) {
	return w.FileHelper.GetExtension(filepath.Join(homePath, "workspace","create"))
}

func (w *CSMWorkspace) getWorkspaceDeleteExtension(homePath string) (bool, *string) {
	return w.FileHelper.GetExtension(filepath.Join(homePath , "workspace","delete"))
}

func (w *CSMWorkspace) executeExtension(workspaceID *string, extensionPath *string, workspace *models.ServiceManagerWorkspaceResponse) {
    if workspaceID == nil {
        w.Logger.Error("executeExtension", errors.New("workspaceID is nil"))
        return
    }
    if extensionPath == nil {
        w.Logger.Error("executeExtension", errors.New("extensionPath is nil"))
        return
    }
	w.Logger.Info("executeExtension", lager.Data{"workspaceID": workspaceID, "extension Path ": extensionPath})
	if success, outputFile, output := w.FileHelper.RunExtensionFileGen(*extensionPath, *workspaceID); success {
		w.Logger.Info("executeExtension", lager.Data{"extension execution status ": success})
		w.Logger.Debug("executeExtension", lager.Data{"extension execution Result: ": output})

		if outputFile != nil {
			if *w.Config.DEV_MODE != "true" {
				// clean up if not running in dev mode
				defer os.Remove(outputFile.Name())
			}
			if fileStat, err := os.Stat(outputFile.Name()); err == nil { // checking the file size of the extension output
				if fileStat.Size() > 0 {
					file, e := ioutil.ReadFile(outputFile.Name())
					if e != nil {
						w.Logger.Info("executeExtension", lager.Data{"File error while reading extension output file ": e})
						w.Logger.Error("executeExtension", e)
						workspace.Status = common.PROCESSING_STATUS_FAILED
					}
					err := json.Unmarshal(file, &workspace)
					if err != nil {
						w.Logger.Info("executeExtension", lager.Data{"Failed to parse the extension output": ""})
						w.Logger.Error("executeExtension", err)
					}
					w.Logger.Info("executeExtension", lager.Data{"extension processing status ": workspace.Status})
				} else {
					// file size of extension output file is not greater than 0
					workspace.Status = common.PROCESSING_STATUS_FAILED
					w.Logger.Debug("executeExtension", lager.Data{"extension execution failed": success})
				}
			} else {
				w.Logger.Info("executeExtension", lager.Data{"File error while reading extension output file ": err})
				w.Logger.Error("executeExtension", err)
			}
		}
	} else {
		// extension couldn't be executed
		workspace.Status = common.PROCESSING_STATUS_FAILED
		w.Logger.Debug("executeExtension", lager.Data{"extension execution failed": success})
	}
	workspace.ProcessingType = common.PROCESSING_TYPE_EXTENSION
	w.Logger.Debug("executeExtension", lager.Data{"Updated workspace processing type to ": common.PROCESSING_TYPE_EXTENSION})
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

// GetWorkspace get workspaces
func (w *CSMWorkspace) GetWorkspace(workspaceID string) *models.ServiceManagerWorkspaceResponse {
	w.Logger.Info("GetWorkspace", lager.Data{"workspaceID": workspaceID})
	workspace := models.ServiceManagerWorkspaceResponse{
		ProcessingType: common.PROCESSING_TYPE_NONE,
	}
	exists, filename := w.getWorkspaceGetExtension(*w.Config.MANAGER_HOME)
	w.Logger.Info("GetWorkspace", lager.Data{"exists": exists, "filename": filename})
	if exists && filename != nil {
		w.executeExtension(&workspaceID, filename, &workspace)
	} else {
		w.Logger.Info("GetWorkspace", lager.Data{"extension not found ": exists})
	}
	return &workspace
}

// CreateWorkspace create workspaces
func (w *CSMWorkspace) CreateWorkspace(workspaceCreate *models.ServiceManagerWorkspaceCreateRequest) *models.ServiceManagerWorkspaceResponse {
	w.Logger.Info("CreateWorkspace", lager.Data{"workspaceID": workspaceCreate.WorkspaceID})
	workspace := models.ServiceManagerWorkspaceResponse{
		ProcessingType: common.PROCESSING_TYPE_NONE,
	}
	exists, filename := w.getWorkspaceCreateExtension(*w.Config.MANAGER_HOME)
	if exists && filename != nil {
		w.executeExtension(&workspaceCreate.WorkspaceID, filename, &workspace)
	} else {
		w.Logger.Info("CreateWorkspace", lager.Data{"extension not found ": exists})
	}

	return &workspace
}

// DeleteWorkspace delete workspaces
func (w *CSMWorkspace) DeleteWorkspace(workspaceID string) *models.ServiceManagerWorkspaceResponse {
	w.Logger.Info("CreateWorkspace", lager.Data{"workspaceID": workspaceID})
	workspace := models.ServiceManagerWorkspaceResponse{
		ProcessingType: common.PROCESSING_TYPE_NONE,
	}
	exists, filename := w.getWorkspaceDeleteExtension(*w.Config.MANAGER_HOME)
	if exists && filename != nil {
		w.executeExtension(&workspaceID, filename, &workspace)
	} else {
		w.Logger.Info("CreateWorkspace", lager.Data{"extension not found": exists})
	}

	return &workspace
}
