package provisioner

type RedisProvisionerInterface interface {
	CreateContainer(string) error
	DeleteContainer(string) error
	ContainerExists(string) (bool, error)
	GetCredentials(string) (map[string]string, error)
}
