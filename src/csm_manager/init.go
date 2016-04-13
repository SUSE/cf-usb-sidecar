package csm_manager

import (
	"strings"

	"github.com/hpcloud/catalog-service-manager/src/common"
	"github.com/hpcloud/catalog-service-manager/src/common/utils"
	internalConnection "github.com/hpcloud/catalog-service-manager/src/connection"
	internalSetup "github.com/hpcloud/catalog-service-manager/src/setup"
	internalWorkspaces "github.com/hpcloud/catalog-service-manager/src/workspace"
	"github.com/pivotal-golang/lager"
)

var (
    csmSetup *internalSetup.CSMSetup
    csmWorkspace *internalWorkspaces.CSMWorkspace
    csmConnection *internalConnection.CSMConnection

    config *common.ServiceManagerConfiguration
    logger lager.Logger
)

//InitServiceCatalogManager initilizes service catalog manager
func InitServiceCatalogManager() {
	config = common.NewServiceManagerConfiguration()
	logger = common.NewLogger(strings.ToLower(*config.LOG_LEVEL))
	fileHelper := utils.CSMFileHelper{
		Logger: logger,
	}
	csmSetup = internalSetup.NewCSMSetup(logger, config, fileHelper)
	csmWorkspace = internalWorkspaces.NewCSMWorkspace(logger, config, fileHelper)
	csmConnection = internalConnection.NewCSMConnection(logger, config, fileHelper)

	logInitDetails(logger, config)
	checkExtensions()

	csmSetup.Startup()
}

// GetSetup returns the CSMSetup object
func GetSetup() *internalSetup.CSMSetup {
	return csmSetup
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
func GetLogger() lager.Logger {
	return logger
}

// logs the initialization details
func logInitDetails(logger lager.Logger, config *common.ServiceManagerConfiguration) {
	logger.Info("InitServiceCatalogManager", lager.Data{"Initialisizing ": "Catalog Service Manager"})
	logger.Info("InitServiceCatalogManager", lager.Data{"CSM_LOG_LEVEL ": *config.LOG_LEVEL})
	if *config.DEV_MODE == "true" {
		// log this only if dev mode is enabled
		// this is developer flag, this is not inteded to be available for production
		logger.Info("InitServiceCatalogManager", lager.Data{"CSM_DEV_MODE ": *config.DEV_MODE})
	}
	logger.Info("InitServiceCatalogManager", lager.Data{"CSM_HOME ": *config.MANAGER_HOME})
	logger.Info("InitServiceCatalogManager", lager.Data{"CSM_PARAMETERS ": *config.PARAMETERS})

}

// runs CheckExtension on setup, workspace and connection objects
// We log this on startup so that we can let user know what all
// extensions are visible/known to the CSM service
func checkExtensions() {
	csmSetup.CheckExtensions()
	csmWorkspace.CheckExtensions()
	csmConnection.CheckExtensions()
}
