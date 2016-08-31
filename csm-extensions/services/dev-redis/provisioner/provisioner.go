package provisioner

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"sort"
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

	pass, err := util.SecureRandomString(12)
	if err != nil {
		return err
	}

	svcPort, err := provisioner.findNextOpenPort()
	if err != nil {
		return err
	}

	hostConfig := dockerclient.HostConfig{
		PortBindings: map[dockerclient.Port][]dockerclient.PortBinding{
			"6379/tcp": {{HostIP: "", HostPort: strconv.Itoa(svcPort)}},
		},
		RestartPolicy: dockerclient.RestartPolicy{Name: "always"},
	}

	createOpts := dockerclient.CreateContainerOptions{
		Config: &dockerclient.Config{
			Image: provisioner.redisConfig.DockerImage + ":" + provisioner.redisConfig.ImageTag,
			Cmd:   []string{"redis-server", fmt.Sprintf("--requirepass %s", pass), "--appendonly yes"},
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

	container, err := provisioner.getContainer(containerName)
	if err != nil {
		return false, err
	}

	if container == nil {
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

	if container == nil {
		return "", fmt.Errorf("Could not find container %s", containerName)
	}

	return container.ID, nil
}

func (provisioner *RedisProvisioner) getContainers() ([]dockerclient.APIContainers, error) {
	opts := dockerclient.ListContainersOptions{
		All: true,
	}
	return provisioner.client.ListContainers(opts)
}

func (provisioner *RedisProvisioner) getContainer(containerName string) (*dockerclient.APIContainers, error) {
	containers, err := provisioner.getContainers()
	if err != nil {
		return nil, err
	}

	for _, c := range containers {
		for _, n := range c.Names {
			if strings.TrimPrefix(n, "/") == containerName {
				return &c, nil
			}
		}
	}

	return nil, nil
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

func (provisioner *RedisProvisioner) findNextOpenPort() (int, error) {

	startPort, err := strconv.Atoi(provisioner.redisConfig.RedisServicesPortsStart)
	if err != nil {
		return 0, err
	}
	endPort, err := strconv.Atoi(provisioner.redisConfig.RedisServicesPortsEnd)
	if err != nil {
		return 0, err
	}

	containers, err := provisioner.getContainers()
	if err != nil {
		provisioner.logger.Debug(err.Error())
		return 0, err
	}

	ports := []int{}

	for _, container := range containers {
		for _, p := range container.Ports {
			if p.PublicPort != 0 {
				port := int(p.PublicPort)
				if port > startPort || port < endPort {
					ports = append(ports, port)
				}
			}
		}
	}

	sort.Ints(ports)
	svcPort := 0

	for p := startPort; p <= endPort; p++ {
		used := false
		for _, port := range ports {
			if port == p {
				used = true
				break
			}
		}

		if used == false {
			svcPort = p
			break
		}
	}

	if svcPort == 0 {
		return 0, fmt.Errorf("Could not find any available ports")
	}

	return svcPort, nil
}
