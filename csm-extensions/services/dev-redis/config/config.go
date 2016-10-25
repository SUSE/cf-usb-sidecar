package config

type RedisConfig struct {
	DockerEndpoint          string `env:"DOCKER_ENDPOINT"`
	DockerHost              string `env:"DOCKER_HOST"`
	DockerPort              string `env:"DOCKER_PORT"`
	RedisServicesPortsStart string `env:"REDIS_SERVICE_PORTS_POOL_START"`
	RedisServicesPortsEnd   string `env:"REDIS_SERVICE_PORTS_POOL_END"`
}
