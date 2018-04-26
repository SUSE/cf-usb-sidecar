package provisioner

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/SUSE/cf-usb-sidecar/csm-extensions/services/dev-rabbitmq/config"
	"github.com/SUSE/go-csm-lib/util"
	"github.com/michaelklishin/rabbit-hole"
	"github.com/pivotal-golang/lager"

	dockerclient "github.com/fsouza/go-dockerclient"
)

const CONTAINER_START_TIMEOUT int = 30

type RabbitHoleProvisioner struct {
	rabbitmqConfig config.RabbitmqConfig
	client         *dockerclient.Client
	logger         lager.Logger
	connected      bool
}

const DockerImage = "rabbitmq"
const ImageTag = "hsm"

func NewRabbitHoleProvisioner(logger lager.Logger, conf config.RabbitmqConfig) RabbitmqProvisionerInterface {
	return &RabbitHoleProvisioner{logger: logger, rabbitmqConfig: conf}
}

func (provisioner *RabbitHoleProvisioner) CreateContainer(containerName string) error {
	err := provisioner.connect()
	if err != nil {
		return err
	}

	admin_user, err := util.SecureRandomString(32)
	if err != nil {
		return err
	}
	admin_pass, err := util.SecureRandomString(32)
	if err != nil {
		return err
	}

	svcPort, mgmtPort, err := provisioner.findNextOpenPorts()
	if err != nil {
		return err
	}

	hostConfig := dockerclient.HostConfig{
		PortBindings: map[dockerclient.Port][]dockerclient.PortBinding{
			"5672/tcp":  {{HostIP: "", HostPort: strconv.Itoa(svcPort)}},
			"15672/tcp": {{HostIP: "", HostPort: strconv.Itoa(mgmtPort)}},
		},
		Binds:         []string{"/var/lib/rabbitmq/" + containerName + ":/var/lib/rabbitmq/mnesia/:rw"},
		RestartPolicy: dockerclient.RestartPolicy{Name: "always"},
	}

	createOpts := dockerclient.CreateContainerOptions{
		Config: &dockerclient.Config{
			Image: DockerImage + ":" + ImageTag,
			Env: []string{"RABBITMQ_DEFAULT_USER=" + admin_user,
				"RABBITMQ_DEFAULT_PASS=" + admin_pass},
		},
		HostConfig: &hostConfig,
		Name:       containerName,
	}

	container, err := provisioner.client.CreateContainer(createOpts)
	if err != nil {
		return err
	}

	provisioner.client.StartContainer(container.ID, nil)
	if err != nil {
		return err
	}

	for i := 0; i < CONTAINER_START_TIMEOUT; i++ {
		state, err := provisioner.getContainerState(containerName)
		if err != nil {
			provisioner.logger.Debug("create-container", lager.Data{"err": err.Error()})
			continue
		}
		if state.Running {
			break
		}
	}

	return nil
}

func (provisioner *RabbitHoleProvisioner) DeleteContainer(containerName string) error {
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

func (provisioner *RabbitHoleProvisioner) ContainerExists(containerName string) (bool, error) {
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

func (provisioner *RabbitHoleProvisioner) DeleteUser(containerName, user string) error {
	err := provisioner.connect()
	if err != nil {
		return err
	}

	host, err := provisioner.getHost()
	if err != nil {
		return err
	}

	admin, err := provisioner.getAdminCredentials(containerName)
	if err != nil {
		return err
	}

	rmqc, err := rabbithole.NewClient(fmt.Sprintf("http://%s:%s", host, admin["mgmt_port"]), admin["user"], admin["password"])
	if err != nil {
		return err
	}

	_, err = rmqc.DeleteUser(user)
	if err != nil {
		return err
	}
	return nil
}

func (provisioner *RabbitHoleProvisioner) UserExists(containerName, user string) (bool, error) {
	err := provisioner.connect()
	if err != nil {
		return false, err
	}

	host, err := provisioner.getHost()
	if err != nil {
		return false, err
	}

	admin, err := provisioner.getAdminCredentials(containerName)
	if err != nil {
		return false, err
	}

	rmqc, err := rabbithole.NewClient(fmt.Sprintf("http://%s:%s", host, admin["mgmt_port"]), admin["user"], admin["password"])
	if err != nil {
		return false, err
	}

	users, err := rmqc.ListUsers()
	if err != nil {
		return false, err
	}
	if users == nil {
		return false, err
	}

	for _, u := range users {
		if u.Name == user {
			return true, nil
		}
	}

	return false, nil
}

func (provisioner *RabbitHoleProvisioner) getClient() (*dockerclient.Client, error) {
	client, err := dockerclient.NewClient(provisioner.rabbitmqConfig.DockerEndpoint)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (provisioner *RabbitHoleProvisioner) FindImage(imageName string) (*dockerclient.Image, error) {
	var err error
	if provisioner.client == nil {
		provisioner.client, err = provisioner.getClient()
		if err != nil {
			return nil, err
		}
	}
	image, err := provisioner.client.InspectImage(imageName)
	if err != nil {
		return nil, fmt.Errorf("Could not find base image %s: %s", imageName, err.Error())
	}

	return image, nil
}

func (provisioner *RabbitHoleProvisioner) getContainerId(containerName string) (string, error) {
	container, err := provisioner.getContainer(containerName)
	if err != nil {
		return "", err
	}

	if container == nil {
		return "", fmt.Errorf("Could not find container %s", containerName)
	}
	return container.ID, nil
}

func (provisioner *RabbitHoleProvisioner) getContainers() ([]dockerclient.APIContainers, error) {
	opts := dockerclient.ListContainersOptions{
		All: true,
	}
	return provisioner.client.ListContainers(opts)
}

func (provisioner *RabbitHoleProvisioner) getContainer(containerName string) (*dockerclient.APIContainers, error) {
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

func (provisioner *RabbitHoleProvisioner) inspectContainer(containerId string) (*dockerclient.Container, error) {
	return provisioner.client.InspectContainer(containerId)
}

func (provisioner *RabbitHoleProvisioner) getAdminCredentials(containerName string) (map[string]string, error) {

	m := make(map[string]string)
	containerId, err := provisioner.getContainerId(containerName)
	if err != nil {
		provisioner.logger.Debug(err.Error())
		return nil, err
	}

	container, err := provisioner.inspectContainer(containerId)
	if err != nil {
		provisioner.logger.Debug(err.Error())
		return nil, err
	}

	var env dockerclient.Env
	env = make([]string, len(container.Config.Env)) // container.Config.Env.(dockerclient.Env)  // dockerclient.Env( []string{ container.Config.Env })
	copy(env, container.Config.Env)
	m["user"] = env.Get("RABBITMQ_DEFAULT_USER")
	m["password"] = env.Get("RABBITMQ_DEFAULT_PASS")
	for k, v := range container.NetworkSettings.Ports {
		if k == "15672/tcp" {
			m["mgmt_port"] = v[0].HostPort
		}
		if k == "5672/tcp" {
			m["port"] = v[0].HostPort
		}
	}
	return m, nil
}

func (provisioner *RabbitHoleProvisioner) getContainerState(containerName string) (dockerclient.State, error) {
	container, err := provisioner.getContainer(containerName)
	if err != nil {
		return dockerclient.State{}, err
	}

	if container == nil {
		return dockerclient.State{}, fmt.Errorf("Container %s does not exist", containerName)
	}

	c, err := provisioner.inspectContainer(container.ID)
	if err != nil {
		return dockerclient.State{}, err
	}
	return c.State, nil
}

func (provisioner *RabbitHoleProvisioner) CreateUser(containerName, newUser, userPass string) (map[string]string, error) {
	err := provisioner.connect()
	if err != nil {
		return nil, err
	}

	host, err := provisioner.getHost()
	if err != nil {
		return nil, err
	}

	admin, err := provisioner.getAdminCredentials(containerName)
	if err != nil {
		return nil, err
	}

	rmqc, err := rabbithole.NewClient(fmt.Sprintf("http://%s:%s", host, admin["mgmt_port"]), admin["user"], admin["password"])
	if err != nil {
		return nil, err
	}
	m := make(map[string]string)

	_, err = rmqc.PutUser(newUser, rabbithole.UserSettings{Password: userPass, Tags: "management,policymaker"})
	if err != nil {
		return nil, err
	}

	_, err = rmqc.UpdatePermissionsIn("/", newUser, rabbithole.Permissions{Configure: ".*", Write: ".*", Read: ".*"})
	if err != nil {
		return nil, err
	}
	m["host"] = host
	m["user"] = newUser
	m["password"] = userPass
	m["mgmt_port"] = admin["mgmt_port"]
	m["port"] = admin["port"]
	x, err := rmqc.GetVhost("/")
	if err != nil {
		return nil, err
	}
	m["vhost"] = x.Name

	return m, nil
}

func (provisioner *RabbitHoleProvisioner) getHost() (string, error) {
	host := ""
	dockerUrl, err := url.Parse(provisioner.rabbitmqConfig.DockerEndpoint)
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

func (provisioner *RabbitHoleProvisioner) connect() error {
	if provisioner.connected {
		return nil
	}

	var err error

	dockerUrl, err := url.Parse(provisioner.rabbitmqConfig.DockerEndpoint)
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

func (provisioner *RabbitHoleProvisioner) findNextOpenPorts() (int, int, error) {

	startPort, err := strconv.Atoi(provisioner.rabbitmqConfig.RabbitServicesPortsStart)
	if err != nil {
		return 0, 0, err
	}
	endPort, err := strconv.Atoi(provisioner.rabbitmqConfig.RabbitServicesPortsEnd)
	if err != nil {
		return 0, 0, err
	}

	containers, err := provisioner.getContainers()
	if err != nil {
		provisioner.logger.Debug(err.Error())
		return 0, 0, err
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
	mgmtPort := 0

	freePort := startPort
	usedPortsIndex := 0

	for {
		if freePort > endPort {
			break
		}

		if len(ports) <= usedPortsIndex || ports[usedPortsIndex] > freePort {
			if svcPort == 0 {
				svcPort = freePort
				freePort++
			} else if mgmtPort == 0 {
				mgmtPort = freePort
				break
			}
		} else if ports[usedPortsIndex] < freePort {
			usedPortsIndex++
		} else if ports[usedPortsIndex] == freePort {
			usedPortsIndex++
			freePort++
		}
	}

	if svcPort == 0 || mgmtPort == 0 {
		return 0, 0, fmt.Errorf("Could not find any available ports")
	}

	return svcPort, mgmtPort, nil
}
