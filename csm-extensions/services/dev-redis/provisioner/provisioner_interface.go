package provisioner

import dockerclient "github.com/fsouza/go-dockerclient"

type RedisProvisionerInterface interface {
	CreateContainer(string) error
	DeleteContainer(string) error
	ContainerExists(string) (bool, error)
	GetCredentials(string) (map[string]string, error)
	FindImage(string) (*dockerclient.Image, error)
}
