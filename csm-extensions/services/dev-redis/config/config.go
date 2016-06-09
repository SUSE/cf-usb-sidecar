package config

type RedisConfig struct {
	DockerEndpoint          string `env:"DOCKER_ENDPOINT"`
	DockerHost              string `env:"DOCKER_HOST"`
	DockerPort              string `env:"DOCKER_PORT"`
	DockerImage             string `env:"DOCKER_IMAGE"`
	ImageTag                string `env:"DOCKER_IMAGE_TAG"`
	UcpDomainSuffix         string `env:"UCP_SERVICE_DOMAIN_SUFFIX"`
	RedisServicesPortsStart string `env:"REDIS_SERVICE_PORTS_POOL_START"`
	RedisServicesPortsEnd   string `env:"REDIS_SERVICE_PORTS_POOL_END"`
}
