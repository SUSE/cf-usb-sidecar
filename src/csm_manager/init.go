package csm_manager

import (
	"strings"

	"github.com/SUSE/cf-usb-sidecar/src/common"
	"github.com/SUSE/cf-usb-sidecar/src/common/utils"
	internalConnection "github.com/SUSE/cf-usb-sidecar/src/connection"
	internalStatus "github.com/SUSE/cf-usb-sidecar/src/status"
	internalWorkspaces "github.com/SUSE/cf-usb-sidecar/src/workspace"
	"github.com/Sirupsen/logrus"
)

var (
	csmWorkspace  *internalWorkspaces.CSMWorkspace
	csmConnection *internalConnection.CSMConnection
	csmStatus     *internalStatus.CSMStatus

	config *common.ServiceManagerConfiguration
	logger *logrus.Logger
)

//InitServiceCatalogManager initilizes service catalog manager
func InitServiceCatalogManager() {
	config = common.NewServiceManagerConfiguration()
	logger = common.NewLogger(strings.ToLower(*config.LOG_LEVEL), *config.HCP_INSTANCE_ID)
	fileHelper := utils.CSMFileHelper{
		Logger: logger,
		Config: config,
	}
	csmWorkspace = internalWorkspaces.NewCSMWorkspace(logger, config, fileHelper)
	csmStatus = internalStatus.NewCSMStatus(logger, config, fileHelper)
	csmConnection = internalConnection.NewCSMConnection(logger, config, fileHelper)

	logInitDetails(logger, config)
	checkExtensions()
}

//GetStatus returns the CSMStatus object
func GetStatus() *internalStatus.CSMStatus {
	return csmStatus
}

// GetWorkspace returns the CSMWorkspace object
func GetWorkspace() *internalWorkspaces.CSMWorkspace {
	return csmWorkspace
}

// GetConnection returns the CSMConnection object
func GetConnection() *internalConnection.CSMConnection {
	return csmConnection
}

// GetConfig returns the ServiceManagerConfiguration object
func GetConfig() *common.ServiceManagerConfiguration {
	return config
}

// GetLogger returns the Logger object
func GetLogger() *logrus.Logger {
	return logger
}

// logInitDetails logs the initialization details
func logInitDetails(logger *logrus.Logger, config *common.ServiceManagerConfiguration) {
	logger.Info("InitServiceCatalogManager ", "Initializing: "+"Catalog Service Manager")
	logger.Info("InitServiceCatalogManager ", "SIDECAR_LOG_LEVEL: "+*config.LOG_LEVEL)
	if *config.DEV_MODE == "true" {
		// log this only if dev mode is enabled
		// this is developer flag, this is not inteded to be available for production
		logger.Info("InitServiceCatalogManager ", "SIDECAR_DEV_MODE: "+*config.DEV_MODE)
	}
	logger.Info("InitServiceCatalogManager ", "SIDECAR_HOME: "+*config.MANAGER_HOME)
	logger.Info("InitServiceCatalogManager ", "SIDECAR_PARAMETERS: "+*config.PARAMETERS)
	logger.Info("InitServiceCatalogManager ", "SIDECAR_API_KEY: "+*config.API_KEY)
}

// CheckExtension runs on workspace and connection objects
// We log this on startup so that we can let user know what all
// extensions are visible/known to the CSM service
func checkExtensions() {
	csmWorkspace.CheckExtensions()
	csmConnection.CheckExtensions()
}
