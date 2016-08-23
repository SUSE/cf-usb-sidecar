package common

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func unsetEnvVariable() {

	os.Unsetenv("SIDECAR_HOME")
	os.Unsetenv("SIDECAR_PARAMETERS")
	os.Unsetenv("SIDECAR_LOG_LEVEL")
	os.Unsetenv("SIDECAR_DEV_MODE")
	os.Unsetenv("SIDECAR_API_KEY")
	os.Unsetenv("SIDECAR_EXT_TIMEOUT")
	os.Unsetenv("SIDECAR_EXT_TIMEOUT_ERROR")
}

func TestCheck_DefaultConfig(t *testing.T) {
	unsetEnvVariable()
	os.Setenv("SIDECAR_PARAMETERS", "username")
	os.Setenv("SIDECAR_API_KEY", "foobar")
	config := NewServiceManagerConfiguration()
	assert.Equal(t, *config.PARAMETERS, "username")
	assert.Equal(t, *config.MANAGER_HOME, "/catalog-service-manager/")
	assert.Equal(t, *config.LOG_LEVEL, "info")
	assert.Equal(t, *config.DEV_MODE, "false")
	assert.Equal(t, *config.API_KEY, "foobar")
	assert.Equal(t, *config.EXT_TIMEOUT, "30")
	assert.Equal(t, *config.EXT_TIMEOUT_ERROR, "2")
}

func TestCheck_SIDECAR_HOME(t *testing.T) {
	unsetEnvVariable()
	os.Setenv("SIDECAR_PARAMETERS", "username")
	os.Setenv("SIDECAR_API_KEY", "foobar")
	os.Setenv("SIDECAR_HOME", "/tmp/SIDECAR_home")
	config := NewServiceManagerConfiguration()
	assert.Equal(t, *config.MANAGER_HOME, "/tmp/SIDECAR_home")
	assert.Equal(t, *config.PARAMETERS, "username")
	assert.Equal(t, *config.LOG_LEVEL, "info")
	assert.Equal(t, *config.DEV_MODE, "false")
	assert.Equal(t, *config.API_KEY, "foobar")
	assert.Equal(t, *config.EXT_TIMEOUT, "30")
	assert.Equal(t, *config.EXT_TIMEOUT_ERROR, "2")
}

func TestCheck_PARAMETERS(t *testing.T) {
	unsetEnvVariable()
	os.Setenv("SIDECAR_API_KEY", "foobar")
	os.Setenv("SIDECAR_PARAMETERS", "Username")
	config := NewServiceManagerConfiguration()
	assert.Equal(t, *config.PARAMETERS, "Username")
	assert.Equal(t, *config.MANAGER_HOME, "/catalog-service-manager/")
	assert.Equal(t, *config.LOG_LEVEL, "info")
	assert.Equal(t, *config.DEV_MODE, "false")
	assert.Equal(t, *config.API_KEY, "foobar")
	assert.Equal(t, *config.EXT_TIMEOUT, "30")
	assert.Equal(t, *config.EXT_TIMEOUT_ERROR, "2")
}

func TestCheck_LOGLEVEL(t *testing.T) {
	unsetEnvVariable()
	os.Setenv("SIDECAR_PARAMETERS", "username")
	os.Setenv("SIDECAR_API_KEY", "foobar")
	os.Setenv("SIDECAR_LOG_LEVEL", "debug")
	config := NewServiceManagerConfiguration()
	assert.Equal(t, *config.LOG_LEVEL, "debug")
	assert.Equal(t, *config.PARAMETERS, "username")
	assert.Equal(t, *config.MANAGER_HOME, "/catalog-service-manager/")
	assert.Equal(t, *config.DEV_MODE, "false")
	assert.Equal(t, *config.API_KEY, "foobar")
	assert.Equal(t, *config.EXT_TIMEOUT, "30")
	assert.Equal(t, *config.EXT_TIMEOUT_ERROR, "2")
}

func TestCheck_DEV_MODE(t *testing.T) {
	unsetEnvVariable()
	os.Setenv("SIDECAR_PARAMETERS", "username")
	os.Setenv("SIDECAR_API_KEY", "foobar")
	os.Setenv("SIDECAR_DEV_MODE", "true")
	config := NewServiceManagerConfiguration()
	assert.Equal(t, *config.DEV_MODE, "true")
	assert.Equal(t, *config.PARAMETERS, "username")
	assert.Equal(t, *config.MANAGER_HOME, "/catalog-service-manager/")
	assert.Equal(t, *config.LOG_LEVEL, "info")
	assert.Equal(t, *config.API_KEY, "foobar")
	assert.Equal(t, *config.EXT_TIMEOUT, "30")
	assert.Equal(t, *config.EXT_TIMEOUT_ERROR, "2")
}

func TestCheck_API_KEY(t *testing.T) {
	unsetEnvVariable()
	os.Setenv("SIDECAR_PARAMETERS", "username")
	os.Setenv("SIDECAR_API_KEY", "foobar")
	config := NewServiceManagerConfiguration()
	assert.Equal(t, *config.DEV_MODE, "false")
	assert.Equal(t, *config.PARAMETERS, "username")
	assert.Equal(t, *config.MANAGER_HOME, "/catalog-service-manager/")
	assert.Equal(t, *config.LOG_LEVEL, "info")
	assert.Equal(t, *config.API_KEY, "foobar")
	assert.Equal(t, *config.EXT_TIMEOUT, "30")
	assert.Equal(t, *config.EXT_TIMEOUT_ERROR, "2")
}

func TestCheck_SIDECAR_EXT_TIMEOUTs(t *testing.T) {
	unsetEnvVariable()
	os.Setenv("SIDECAR_PARAMETERS", "username")
	os.Setenv("SIDECAR_EXT_TIMEOUT", "29")
	os.Setenv("SIDECAR_EXT_TIMEOUT_ERROR", "3")
	os.Setenv("SIDECAR_API_KEY", "foobar")
	config := NewServiceManagerConfiguration()
	assert.Equal(t, *config.DEV_MODE, "false")
	assert.Equal(t, *config.PARAMETERS, "username")
	assert.Equal(t, *config.MANAGER_HOME, "/catalog-service-manager/")
	assert.Equal(t, *config.LOG_LEVEL, "info")
	assert.Equal(t, *config.API_KEY, "foobar")
	assert.Equal(t, *config.EXT_TIMEOUT, "29")
	assert.Equal(t, *config.EXT_TIMEOUT_ERROR, "3")
}

func TestCheck_All(t *testing.T) {
	unsetEnvVariable()
	os.Setenv("SIDECAR_HOME", "/tmp/SIDECAR_home")
	os.Setenv("SIDECAR_PARAMETERS", "Username")
	os.Setenv("SIDECAR_LOG_LEVEL", "debug")
	os.Setenv("SIDECAR_DEV_MODE", "true")
	os.Setenv("SIDECAR_API_KEY", "foobar")
	os.Setenv("SIDECAR_EXT_TIMEOUT", "29")
	os.Setenv("SIDECAR_EXT_TIMEOUT_ERROR", "3")
	config := NewServiceManagerConfiguration()
	assert.Equal(t, *config.MANAGER_HOME, "/tmp/SIDECAR_home")
	assert.Equal(t, *config.PARAMETERS, "Username")
	assert.Equal(t, *config.LOG_LEVEL, "debug")
	assert.Equal(t, *config.DEV_MODE, "true")
	assert.Equal(t, *config.API_KEY, "foobar")
	assert.Equal(t, *config.EXT_TIMEOUT, "29")
	assert.Equal(t, *config.EXT_TIMEOUT_ERROR, "3")
}
