package provisioner

import (
	"bytes"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/hpcloud/catalog-service-manager/csm-extensions/services/dev-redis/config"
	"github.com/hpcloud/go-csm-lib/util"
	"github.com/pivotal-golang/lager"

	dockerclient "github.com/fsouza/go-dockerclient"
)

type RedisProvisioner struct {
	redisConfig config.RedisConfig
	client      *dockerclient.Client
	logger      lager.Logger
	connected   bool
}

func NewRedisProvisioner(logger lager.Logger, conf config.RedisConfig) RedisProvisionerInterface {
	return &RedisProvisioner{logger: logger, redisConfig: conf}
}

func (provisioner *RedisProvisioner) CreateContainer(containerName string) error {
	err := provisioner.connect()
	if err != nil {
		return err
	}
	err = provisioner.pullImage(provisioner.redisConfig.DockerImage, provisioner.redisConfig.ImageTag)
	if err != nil {
		return err
	}

	pass, err := util.SecureRandomString(12)
	if err != nil {
		return err
	}

	hostConfig := dockerclient.HostConfig{PublishAllPorts: true}
	createOpts := dockerclient.CreateContainerOptions{
		Config: &dockerclient.Config{
			Image: provisioner.redisConfig.DockerImage + ":" + provisioner.redisConfig.ImageTag,
			Cmd:   []string{"redis-server", fmt.Sprintf("--requirepass %s", pass)},
		},
		HostConfig: &hostConfig,
		Name:       containerName,
	}

	container, err := provisioner.client.CreateContainer(createOpts)
	if err != nil {
		return err
	}

	provisioner.client.StartContainer(container.ID, &hostConfig)
	if err != nil {
		return err
	}

	return nil
}

func (provisioner *RedisProvisioner) DeleteContainer(containerName string) error {
	err := provisioner.connect()
	if err != nil {
		return err
	}
	containerID, err := provisioner.getContainerId(containerName)
	if err != nil {
		return err
	}

	err = provisioner.client.StopContainer(containerID, 5)
	if err != nil {
		return err
	}

	return provisioner.client.RemoveContainer(dockerclient.RemoveContainerOptions{
		ID:    containerID,
		Force: true,
	})
}

func (provisioner *RedisProvisioner) GetCredentials(containerName string) (map[string]string, error) {
	err := provisioner.connect()
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)

	container, err := provisioner.getContainer(containerName)
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`'--requirepass\s(\S+)'`)
	submatch := re.FindStringSubmatch(container.Command)
	if submatch == nil {
		return nil, fmt.Errorf("Could not get password")
	}

	host, err := provisioner.getHost()
	if err != nil {
		return nil, err
	}

	m["host"] = host
	m["password"] = submatch[1]
	m["port"] = strconv.FormatInt(container.Ports[0].PublicPort, 10)

	return m, nil
}

func (provisioner *RedisProvisioner) ContainerExists(containerName string) (bool, error) {
	err := provisioner.connect()
	if err != nil {
		return false, err
	}

	_, err = provisioner.getContainer(containerName)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (provisioner *RedisProvisioner) getClient() (*dockerclient.Client, error) {
	client, err := dockerclient.NewClient(provisioner.redisConfig.DockerEndpoint)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (provisioner *RedisProvisioner) pullImage(imageName, version string) error {
	var buf bytes.Buffer
	pullOpts := dockerclient.PullImageOptions{
		Repository:   imageName,
		Tag:          version,
		OutputStream: &buf,
	}

	err := provisioner.client.PullImage(pullOpts, dockerclient.AuthConfiguration{})
	if err != nil {
		return err
	}
	return nil
}

func (provisioner *RedisProvisioner) findImage(imageName string) (*dockerclient.Image, error) {
	image, err := provisioner.client.InspectImage(imageName)
	if err != nil {
		return nil, fmt.Errorf("Could not find base image %s: %s", imageName, err.Error())
	}

	return image, nil
}

func (provisioner *RedisProvisioner) getContainerId(containerName string) (string, error) {
	container, err := provisioner.getContainer(containerName)
	if err != nil {
		return "", err
	}
	return container.ID, nil
}

func (provisioner *RedisProvisioner) getContainer(containerName string) (dockerclient.APIContainers, error) {
	opts := dockerclient.ListContainersOptions{
		All: true,
	}
	containers, err := provisioner.client.ListContainers(opts)
	if err != nil {
		return dockerclient.APIContainers{}, err
	}

	for _, c := range containers {
		for _, n := range c.Names {
			if strings.TrimPrefix(n, "/") == containerName {
				return c, nil
			}
		}
	}

	return dockerclient.APIContainers{}, fmt.Errorf("Container %s not found", containerName)
}

func (provisioner *RedisProvisioner) connect() error {
	if provisioner.connected {
		return nil
	}

	var err error

	dockerUrl, err := url.Parse(provisioner.redisConfig.DockerEndpoint)
	if err != nil {
		return err
	}

	if dockerUrl.Scheme == "" {
		return errors.New("Invalid URL format")
	}

	provisioner.client, err = provisioner.getClient()

	if err != nil {
		return err
	}

	provisioner.connected = true
	return nil
}

func (provisioner *RedisProvisioner) getHost() (string, error) {

	host := ""
	dockerUrl, err := url.Parse(provisioner.redisConfig.DockerEndpoint)
	if err != nil {
		return "", err
	}

	if dockerUrl.Scheme == "unix" {
		host, err = util.GetLocalIP()
		if err != nil {
			return "", err
		}
	} else {
		host = strings.Split(dockerUrl.Host, ":")[0]
	}

	return host, nil
}
