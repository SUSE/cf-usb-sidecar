package setup

import (
	"encoding/json"
	"errors"
	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager/models"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hpcloud/catalog-service-manager/src/common"
	"github.com/hpcloud/catalog-service-manager/src/common/utils"
	"github.com/pivotal-golang/lager"
)

type CSMSetup struct {
	common.CSMSetupInterface
	Logger     lager.Logger
	Config     *common.ServiceManagerConfiguration
	FileHelper utils.CSMFileHelperInterface
}

func NewCSMSetup(logger lager.Logger,
	config *common.ServiceManagerConfiguration,
	fileHelper utils.CSMFileHelperInterface) *CSMSetup {
	return &CSMSetup{Logger: logger.Session("CSM-Setup"), Config: config, FileHelper: fileHelper}
}

func (s *CSMSetup) getSetupStartupExtension(homePath string) (bool, *string) {
	return s.FileHelper.GetExtension(filepath.Join(homePath, "setup", "startup"))
}

func (s *CSMSetup) getSetupShutdownExtension(homePath string) (bool, *string) {
	return s.FileHelper.GetExtension(filepath.Join(homePath, "setup", "shutdown"))
}

func (s *CSMSetup) executeExtension(extensionPath *string, setup *models.ServiceManagerWorkspaceResponse) {
	if extensionPath == nil {
		s.Logger.Error("executeExtension", errors.New("extensionPath is nil"))
		return
	}
	s.Logger.Info("executeExtension", lager.Data{"extension Path ": extensionPath})
	if success, outputFile, output := s.FileHelper.RunExtensionFileGen(*extensionPath, ""); success {
		s.Logger.Info("executeExtension", lager.Data{"extension execution status ": success})
		s.Logger.Debug("executeExtension", lager.Data{"extension execution Result: ": output})
		if outputFile != nil {
			if *s.Config.DEV_MODE != "true" {
				// clean up if not running in dev mode
				defer os.Remove(outputFile.Name())
			}
			if fileStat, err := os.Stat(outputFile.Name()); err == nil {
				// checking the file size of the extension output
				if fileStat.Size() > 0 {
					file, e := ioutil.ReadFile(outputFile.Name())
					if e != nil {
						s.Logger.Info("executeExtension", lager.Data{"File error while reading extension output file ": e})
						s.Logger.Error("executeExtension", e)
						setup.Status = common.PROCESSING_STATUS_FAILED
					}
					err := json.Unmarshal(file, &setup)
					if err != nil {
						s.Logger.Info("executeExtension", lager.Data{"Failed to parse the extension output": ""})
						s.Logger.Error("executeExtension", err)
					}
					s.Logger.Info("executeExtension", lager.Data{"extension processing status ": setup.Status})
				} else {
					// file size of extension output file is not greater than 0
					setup.Status = common.PROCESSING_STATUS_FAILED
					s.Logger.Info("executeExtension", lager.Data{"extension output file is empty": success})
				}
			} else {
				s.Logger.Info("executeExtension", lager.Data{"File error while reading extension output file ": err})
				s.Logger.Error("executeExtension", err)
			}
		}
	} else {
		// extension couldn't be executed
		setup.Status = common.PROCESSING_STATUS_FAILED
		s.Logger.Info("executeExtension", lager.Data{"extension execution failed": success})
	}
	s.Logger.Debug("executeExtension", lager.Data{"Updated workspace processing type to ": common.PROCESSING_TYPE_EXTENSION})
}

// CheckExtensions checks for workspace extensions
func (s *CSMSetup) CheckExtensions() {
	_, file := s.getSetupStartupExtension(*s.Config.MANAGER_HOME)
	s.Logger.Info("CheckExtensions", lager.Data{"Setup Startup extension": file})

	_, file = s.getSetupShutdownExtension(*s.Config.MANAGER_HOME)
	s.Logger.Info("CheckExtensions", lager.Data{"Setup Shutdown extension": file})
}

// Startup executes startup extension if any
func (s *CSMSetup) Startup() bool {
	s.Logger.Info("Startup", lager.Data{"Run Startup Extension": ""})
	setup := models.ServiceManagerWorkspaceResponse{
		ProcessingType: common.PROCESSING_TYPE_NONE,
	}
	exists, filename := s.getSetupStartupExtension(*s.Config.MANAGER_HOME)
	if exists && filename != nil {
		s.executeExtension(filename, &setup)
	} else {
		s.Logger.Info("Startup", lager.Data{"Extension not found ": exists})
	}
	return setup.Status == common.PROCESSING_STATUS_SUCCESSFUL
}

// Shutdown executes shutdown extension if any
func (s *CSMSetup) Shutdown() bool {
	s.Logger.Info("Shutdown", lager.Data{"Run Shutdown Extension": ""})
	setup := models.ServiceManagerWorkspaceResponse{
		ProcessingType: common.PROCESSING_TYPE_NONE,
	}
	exists, filename := s.getSetupShutdownExtension(*s.Config.MANAGER_HOME)
	if exists && filename != nil {
		s.executeExtension(filename, &setup)
	} else {
		s.Logger.Info("Shutdown", lager.Data{"Extension not found ": exists})
	}
	return setup.Status == common.PROCESSING_STATUS_SUCCESSFUL
}
