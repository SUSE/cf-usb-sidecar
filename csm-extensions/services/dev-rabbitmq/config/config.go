package config

type RabbitmqConfig struct {
	DockerEndpoint           string `env:"DOCKER_ENDPOINT"`
	DockerHost               string `env:"DOCKER_HOST"`
	DockerPort               string `env:"DOCKER_PORT"`
	RabbitServicesPortsStart string `env:"RABBIT_SERVICE_PORTS_POOL_START"`
	RabbitServicesPortsEnd   string `env:"RABBIT_SERVICE_PORTS_POOL_END"`
}
