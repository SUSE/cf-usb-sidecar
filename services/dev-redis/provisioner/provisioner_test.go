package provisioner

import (
	"os"
	"testing"

	"github.com/hpcloud/catalog-service-manager/services/dev-redis/config"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("redis-provisioner")

var testRedisProv = struct {
	redisProvisioner RedisProvisionerInterface
	driverConfig     config.RedisConfig
}{}

func init() {
	testRedisProv.driverConfig = config.RedisConfig{
		DockerEndpoint: os.Getenv("DOCKER_ENDPOINT"),
		DockerImage:    os.Getenv("REDIS_DOCKER_IMAGE"),
		ImageVersion:   os.Getenv("REDIS_DOCKER_IMAGE_VERSION"),
	}

	testRedisProv.redisProvisioner = NewRedisProvisioner(logger, testRedisProv.driverConfig)
}

func TestRedisProvisioner(t *testing.T) {
	if !envVarsOk() {
		t.SkipNow()
	}

	assert := assert.New(t)

	name := "testContainer"

	// Create container

	err := testRedisProv.redisProvisioner.CreateContainer(name)
	assert.NoError(err)

	// Check container exists

	exists, err := testRedisProv.redisProvisioner.ContainerExists(name)
	assert.NoError(err)
	assert.True(exists)

	// Get Credentials

	credentials, err := testRedisProv.redisProvisioner.GetCredentials(name)
	assert.NoError(err)
	assert.NotNil(credentials["password"])
	assert.NotNil(credentials["port"])

	// Delete Container

	err = testRedisProv.redisProvisioner.DeleteContainer(name)
	assert.NoError(err)

	// Check container does not exist

	exists, err = testRedisProv.redisProvisioner.ContainerExists(name)
	assert.NoError(err)
	assert.False(exists)
}

func envVarsOk() bool {
	return testRedisProv.driverConfig.DockerEndpoint != "" && testRedisProv.driverConfig.DockerImage != "" && testRedisProv.driverConfig.ImageVersion != ""
}
