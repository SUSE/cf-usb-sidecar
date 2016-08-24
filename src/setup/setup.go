package setup

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager/models"

	"github.com/Sirupsen/logrus"
	"github.com/hpcloud/catalog-service-manager/src/common"
	"github.com/hpcloud/catalog-service-manager/src/common/utils"
)

type CSMSetup struct {
	common.CSMSetupInterface
	Logger     *logrus.Logger
	Config     *common.ServiceManagerConfiguration
	FileHelper utils.CSMFileHelperInterface
}

func NewCSMSetup(logger *logrus.Logger,
	config *common.ServiceManagerConfiguration,
	fileHelper utils.CSMFileHelperInterface) *CSMSetup {
	return &CSMSetup{Logger: logger, Config: config, FileHelper: fileHelper}
}

func (s *CSMSetup) getSetupStartupExtension(homePath string) (bool, string) {
	return s.FileHelper.GetExtension(filepath.Join(homePath, "setup", "startup"))
}

func (s *CSMSetup) getSetupShutdownExtension(homePath string) (bool, string) {
	return s.FileHelper.GetExtension(filepath.Join(homePath, "setup", "shutdown"))
}

func (s *CSMSetup) executeExtension(extensionPath string, setup *models.ServiceManagerWorkspaceResponse) {
	if extensionPath == "" {
		s.Logger.Error("executeExtension", errors.New("extensionPath is not set"))
		return
	}
	s.Logger.WithFields(logrus.Fields{"extension Path ": extensionPath}).Info("executeExtension")
	if success, outputFile, output := s.FileHelper.RunExtensionFileGen(extensionPath, ""); success {
		s.Logger.WithFields(logrus.Fields{"extension execution status ": success}).Info("executeExtension")
		s.Logger.WithFields(logrus.Fields{"extension execution Result: ": output}).Debug("executeExtension")
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
						s.Logger.WithFields(logrus.Fields{"File error while reading extension output file ": e}).Info("executeExtension")
						s.Logger.Error("executeExtension", e)
						setup.Status = common.PROCESSING_STATUS_FAILED
					}
					err := json.Unmarshal(file, &setup)
					if err != nil {
						s.Logger.WithFields(logrus.Fields{"Failed to parse the extension output": ""}).Info("executeExtension")
						s.Logger.Error("executeExtension", err)
					}
					s.Logger.WithFields(logrus.Fields{"extension processing status ": setup.Status}).Info("executeExtension")
				} else {
					// file size of extension output file is not greater than 0
					setup.Status = common.PROCESSING_STATUS_FAILED
					s.Logger.WithFields(logrus.Fields{"extension output file is empty": success}).Info("executeExtension")
				}
			} else {
				s.Logger.WithFields(logrus.Fields{"File error while reading extension output file ": err}).Info("executeExtension")
				s.Logger.Error("executeExtension", err)
			}
		}
	} else {
		// extension couldn't be executed
		setup.Status = common.PROCESSING_STATUS_FAILED
		s.Logger.WithFields(logrus.Fields{"extension execution failed": success}).Info("executeExtension")
	}
	s.Logger.WithFields(logrus.Fields{"Updated workspace processing type to ": common.PROCESSING_TYPE_EXTENSION}).Debug("executeExtension")
}

// CheckExtensions checks for workspace extensions
func (s *CSMSetup) CheckExtensions() {
	_, file := s.getSetupStartupExtension(*s.Config.MANAGER_HOME)
	s.Logger.WithFields(logrus.Fields{"Setup Startup extension": file}).Info("CheckExtensions")

	_, file = s.getSetupShutdownExtension(*s.Config.MANAGER_HOME)
	s.Logger.WithFields(logrus.Fields{"Setup Shutdown extension": file}).Info("CheckExtensions")
}

// Startup executes startup extension if any
func (s *CSMSetup) Startup() bool {
	s.Logger.WithFields(logrus.Fields{"Run Startup Extension": ""}).Info("Startup")
	setup := models.ServiceManagerWorkspaceResponse{
		ProcessingType: common.PROCESSING_TYPE_NONE,
	}
	exists, filename := s.getSetupStartupExtension(*s.Config.MANAGER_HOME)
	if exists && filename != "" {
		s.executeExtension(filename, &setup)
	} else {
		s.Logger.WithFields(logrus.Fields{"Extension not found ": exists}).Info("Startup")
	}
	return setup.Status == common.PROCESSING_STATUS_SUCCESSFUL
}

// Shutdown executes shutdown extension if any
func (s *CSMSetup) Shutdown() bool {
	s.Logger.WithFields(logrus.Fields{"Run Shutdown Extension": ""}).Info("Shutdown")
	setup := models.ServiceManagerWorkspaceResponse{
		ProcessingType: common.PROCESSING_TYPE_NONE,
	}
	exists, filename := s.getSetupShutdownExtension(*s.Config.MANAGER_HOME)
	if exists && filename != "" {
		s.executeExtension(filename, &setup)
	} else {
		s.Logger.WithFields(logrus.Fields{"Extension not found ": exists}).Info("Shutdown")
	}
	return setup.Status == common.PROCESSING_STATUS_SUCCESSFUL
}
