package common

import (
	"fmt"
	"os"
)

type ServiceManagerConfiguration struct {
	PARAMETERS           *string
	MANAGER_HOME         *string
	LOG_LEVEL            *string
	DEV_MODE             *string
	API_KEY              *string
	EXT_TIMEOUT          *string
	EXT_TIMEOUT_ERROR    *string
	HEALTHCHECK_HOST     *string
	HEALTHCHECK_PORT     *string
	TLS_CERT_FILE        *string
	TLS_PRIVATE_KEY_FILE *string
	PORT                 *string `env:"PORT"`
	KUBE_INSTANCE_ID      *string `env:"KUBE_INSTANCE_ID"`
}

var paramDefaultList = map[string]string{
	"SIDECAR_HOME":              "/catalog-service-manager/",
	"SIDECAR_PARAMETERS":        "",
	"SIDECAR_LOG_LEVEL":         "info",
	"SIDECAR_DEV_MODE":          "false",
	"SIDECAR_EXT_TIMEOUT":       "30",
	"SIDECAR_EXT_TIMEOUT_ERROR": "2",
	"HEALTHCHECK_HOST":          "",
	"HEALTHCHECK_PORT":          "",
	"TLS_CERT_FILE":             "/etc/secrets/tls-cert-file",
	"TLS_PRIVATE_KEY_FILE":      "/etc/secrets/tls-private-key-file",
	"PORT":                      "8081",
	"KUBE_INSTANCE_ID":           "Sidecar",
}

// NewServiceManagerConfiguration : Creates object of ServiceManagerConfiguration
func NewServiceManagerConfiguration() *ServiceManagerConfiguration {
	config := ServiceManagerConfiguration{}
	config.PARAMETERS = config.getConfigFromEnv("SIDECAR_PARAMETERS")
	config.MANAGER_HOME = config.getConfigFromEnv("SIDECAR_HOME")
	config.LOG_LEVEL = config.getConfigFromEnv("SIDECAR_LOG_LEVEL")
	config.DEV_MODE = config.getConfigFromEnv("SIDECAR_DEV_MODE")
	config.API_KEY = config.getConfigFromEnv("SIDECAR_API_KEY")
	config.EXT_TIMEOUT = config.getConfigFromEnv("SIDECAR_EXT_TIMEOUT")
	config.EXT_TIMEOUT_ERROR = config.getConfigFromEnv("SIDECAR_EXT_TIMEOUT_ERROR")
	if *config.LOG_LEVEL == "" {
		config.PARAMETERS = config.getConfigFromEnv("CSM_PARAMETERS")
		config.MANAGER_HOME = config.getConfigFromEnv("CSM_HOME")
		config.LOG_LEVEL = config.getConfigFromEnv("CSM_LOG_LEVEL")
		config.DEV_MODE = config.getConfigFromEnv("CSM_DEV_MODE")
		config.API_KEY = config.getConfigFromEnv("CSM_API_KEY")
		config.EXT_TIMEOUT = config.getConfigFromEnv("CSM_EXT_TIMEOUT")
		config.EXT_TIMEOUT_ERROR = config.getConfigFromEnv("CSM_EXT_TIMEOUT_ERROR")
	}

	config.HEALTHCHECK_HOST = config.getConfigFromEnv("HEALTHCHECK_HOST")
	config.HEALTHCHECK_PORT = config.getConfigFromEnv("HEALTHCHECK_PORT")
	config.TLS_CERT_FILE = config.getConfigFromEnv("TLS_CERT_FILE")
	config.TLS_PRIVATE_KEY_FILE = config.getConfigFromEnv("TLS_PRIVATE_KEY_FILE")
	config.PORT = config.getConfigFromEnv("PORT")
	config.KUBE_INSTANCE_ID = config.getConfigFromEnv("KUBE_INSTANCE_ID")
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
