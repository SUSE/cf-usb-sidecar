package provisioner

import (
	"os"
	"testing"
	"time"

	"github.com/hpcloud/catalog-service-manager/csm-extensions/services/dev-rabbitmq/config"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/stretchr/testify/assert"
)

var logger *lagertest.TestLogger = lagertest.NewTestLogger("rabbitmq-provisioner")

var testRabbitmqProv = struct {
	rabbitmqProvisioner RabbitmqProvisionerInterface
	driverConfig        config.RabbitmqConfig
}{}

func init() {
	testRabbitmqProv.driverConfig = config.RabbitmqConfig{
		DockerEndpoint: os.Getenv("DOCKER_ENDPOINT"),
		DockerImage:    os.Getenv("RABBIT_DOCKER_IMAGE"),
		ImageVersion:   os.Getenv("RABBIT_DOCKER_IMAGE_VERSION"),
	}

	testRabbitmqProv.rabbitmqProvisioner = NewRabbitHoleProvisioner(logger, testRabbitmqProv.driverConfig)
}

func TestRabbitholeProvisioner(t *testing.T) {
	if !envVarsOk() {
		t.SkipNow()
	}

	assert := assert.New(t)

	name := "rabbitContainer"

	// Create Container

	err := testRabbitmqProv.rabbitmqProvisioner.CreateContainer(name)
	assert.NoError(err)

	// Check container exists

	exists, err := testRabbitmqProv.rabbitmqProvisioner.ContainerExists(name)
	assert.NoError(err)
	assert.True(exists)

	// Create User

	// Wait for container to initialize
	time.Sleep(10 * time.Second)

	user := "user"
	password := "password"

	credentials, err := testRabbitmqProv.rabbitmqProvisioner.CreateUser(name, user, password)
	assert.NoError(err)
	assert.NotNil(credentials["password"])
	assert.NotNil(credentials["port"])
	assert.NotNil(credentials["mgmt_port"])
	assert.NotNil(credentials["host"])
	assert.NotNil(credentials["user"])
	assert.NotNil(credentials["vhost"])

	// Check user exists

	exists, err = testRabbitmqProv.rabbitmqProvisioner.UserExists(name, user)
	assert.NoError(err)
	assert.True(exists)

	// Delete User

	err = testRabbitmqProv.rabbitmqProvisioner.DeleteUser(name, user)
	assert.NoError(err)

	// Make sure user exists fails

	exists, err = testRabbitmqProv.rabbitmqProvisioner.UserExists(name, user)
	assert.NoError(err)
	assert.False(exists)

	// Delete container

	err = testRabbitmqProv.rabbitmqProvisioner.DeleteContainer(name)
	assert.NoError(err)

	// Make sure container does not exist

	exists, err = testRabbitmqProv.rabbitmqProvisioner.ContainerExists(name)
	assert.False(exists)
}

func envVarsOk() bool {
	return testRabbitmqProv.driverConfig.DockerEndpoint != "" && testRabbitmqProv.driverConfig.DockerImage != "" && testRabbitmqProv.driverConfig.ImageVersion != ""
}
