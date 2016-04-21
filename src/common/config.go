package common

import (
	"fmt"
	"os"
)

type ServiceManagerConfiguration struct {
	PARAMETERS   *string
	MANAGER_HOME *string
	LOG_LEVEL    *string
	DEV_MODE     *string
	API_KEY      *string
}

var paramDefaultList = map[string]string{
	"CSM_HOME":       "/catalog-service-manager/",
	"CSM_PARAMETERS": "",
	"CSM_LOG_LEVEL":  "info",
	"CSM_DEV_MODE":   "false",
}

// NewServiceManagerConfiguration : Creates object of ServiceManagerConfiguration
func NewServiceManagerConfiguration() *ServiceManagerConfiguration {
	config := ServiceManagerConfiguration{}
	config.PARAMETERS = config.getConfigFromEnv("CSM_PARAMETERS")
	config.MANAGER_HOME = config.getConfigFromEnv("CSM_HOME")
	config.LOG_LEVEL = config.getConfigFromEnv("CSM_LOG_LEVEL")
	config.DEV_MODE = config.getConfigFromEnv("CSM_DEV_MODE")
	config.API_KEY = config.getConfigFromEnv("CSM_API_KEY")

	return &config
}

// getConfigFromEnv reads value of the provided environment variable
func (*ServiceManagerConfiguration) getConfigFromEnv(key string) *string {
	value, ok := os.LookupEnv(key)

	if ok {
		return &value
	}
	defValue, found := paramDefaultList[key]
	if found {
		return &defValue
	}
	fmt.Fprintf(os.Stderr, "error: Please configure "+key+".")
	os.Exit(1)
	return nil
}
