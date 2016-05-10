package common

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func unsetEnvVariable() {

	os.Unsetenv("CSM_HOME")
	os.Unsetenv("CSM_PARAMETERS")
	os.Unsetenv("CSM_LOG_LEVEL")
	os.Unsetenv("CSM_DEV_MODE")
	os.Unsetenv("CSM_API_KEY")
	os.Unsetenv("CSM_EXT_TIMEOUT")
	os.Unsetenv("CSM_EXT_TIMEOUT_ERROR")
}

func TestCheck_DefaultConfig(t *testing.T) {
	unsetEnvVariable()
	os.Setenv("CSM_API_KEY", "foobar")
	config := NewServiceManagerConfiguration()
	assert.Equal(t, *config.PARAMETERS, "")
	assert.Equal(t, *config.MANAGER_HOME, "/catalog-service-manager/")
	assert.Equal(t, *config.LOG_LEVEL, "info")
	assert.Equal(t, *config.DEV_MODE, "false")
	assert.Equal(t, *config.API_KEY, "foobar")
	assert.Equal(t, *config.EXT_TIMEOUT, "30")
	assert.Equal(t, *config.EXT_TIMEOUT_ERROR, "2")
}

func TestCheck_CSM_HOME(t *testing.T) {
	unsetEnvVariable()
	os.Setenv("CSM_API_KEY", "foobar")
	os.Setenv("CSM_HOME", "/tmp/csm_home")
	config := NewServiceManagerConfiguration()
	assert.Equal(t, *config.MANAGER_HOME, "/tmp/csm_home")
	assert.Equal(t, *config.PARAMETERS, "")
	assert.Equal(t, *config.LOG_LEVEL, "info")
	assert.Equal(t, *config.DEV_MODE, "false")
	assert.Equal(t, *config.API_KEY, "foobar")
	assert.Equal(t, *config.EXT_TIMEOUT, "30")
	assert.Equal(t, *config.EXT_TIMEOUT_ERROR, "2")
}

func TestCheck_PARAMETERS(t *testing.T) {
	unsetEnvVariable()
	os.Setenv("CSM_API_KEY", "foobar")
	os.Setenv("CSM_PARAMETERS", "Username")
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
	os.Setenv("CSM_API_KEY", "foobar")
	os.Setenv("CSM_LOG_LEVEL", "debug")
	config := NewServiceManagerConfiguration()
	assert.Equal(t, *config.LOG_LEVEL, "debug")
	assert.Equal(t, *config.PARAMETERS, "")
	assert.Equal(t, *config.MANAGER_HOME, "/catalog-service-manager/")
	assert.Equal(t, *config.DEV_MODE, "false")
	assert.Equal(t, *config.API_KEY, "foobar")
	assert.Equal(t, *config.EXT_TIMEOUT, "30")
	assert.Equal(t, *config.EXT_TIMEOUT_ERROR, "2")
}

func TestCheck_DEV_MODE(t *testing.T) {
	unsetEnvVariable()
	os.Setenv("CSM_API_KEY", "foobar")
	os.Setenv("CSM_DEV_MODE", "true")
	config := NewServiceManagerConfiguration()
	assert.Equal(t, *config.DEV_MODE, "true")
	assert.Equal(t, *config.PARAMETERS, "")
	assert.Equal(t, *config.MANAGER_HOME, "/catalog-service-manager/")
	assert.Equal(t, *config.LOG_LEVEL, "info")
	assert.Equal(t, *config.API_KEY, "foobar")
	assert.Equal(t, *config.EXT_TIMEOUT, "30")
	assert.Equal(t, *config.EXT_TIMEOUT_ERROR, "2")
}

func TestCheck_API_KEY(t *testing.T) {
	unsetEnvVariable()
	os.Setenv("CSM_API_KEY", "foobar")
	config := NewServiceManagerConfiguration()
	assert.Equal(t, *config.DEV_MODE, "false")
	assert.Equal(t, *config.PARAMETERS, "")
	assert.Equal(t, *config.MANAGER_HOME, "/catalog-service-manager/")
	assert.Equal(t, *config.LOG_LEVEL, "info")
	assert.Equal(t, *config.API_KEY, "foobar")
	assert.Equal(t, *config.EXT_TIMEOUT, "30")
	assert.Equal(t, *config.EXT_TIMEOUT_ERROR, "2")
}

func TestCheck_CSM_EXT_TIMEOUTs(t *testing.T) {
	unsetEnvVariable()
	os.Setenv("CSM_EXT_TIMEOUT", "29")
	os.Setenv("CSM_EXT_TIMEOUT_ERROR", "3")
	os.Setenv("CSM_API_KEY", "foobar")
	config := NewServiceManagerConfiguration()
	assert.Equal(t, *config.DEV_MODE, "false")
	assert.Equal(t, *config.PARAMETERS, "")
	assert.Equal(t, *config.MANAGER_HOME, "/catalog-service-manager/")
	assert.Equal(t, *config.LOG_LEVEL, "info")
	assert.Equal(t, *config.API_KEY, "foobar")
	assert.Equal(t, *config.EXT_TIMEOUT, "29")
	assert.Equal(t, *config.EXT_TIMEOUT_ERROR, "3")
}

func TestCheck_All(t *testing.T) {
	unsetEnvVariable()
	os.Setenv("CSM_HOME", "/tmp/csm_home")
	os.Setenv("CSM_PARAMETERS", "Username")
	os.Setenv("CSM_LOG_LEVEL", "debug")
	os.Setenv("CSM_DEV_MODE", "true")
	os.Setenv("CSM_API_KEY", "foobar")
	os.Setenv("CSM_EXT_TIMEOUT", "29")
	os.Setenv("CSM_EXT_TIMEOUT_ERROR", "3")
	config := NewServiceManagerConfiguration()
	assert.Equal(t, *config.MANAGER_HOME, "/tmp/csm_home")
	assert.Equal(t, *config.PARAMETERS, "Username")
	assert.Equal(t, *config.LOG_LEVEL, "debug")
	assert.Equal(t, *config.DEV_MODE, "true")
	assert.Equal(t, *config.API_KEY, "foobar")
	assert.Equal(t, *config.EXT_TIMEOUT, "29")
	assert.Equal(t, *config.EXT_TIMEOUT_ERROR, "3")
}
