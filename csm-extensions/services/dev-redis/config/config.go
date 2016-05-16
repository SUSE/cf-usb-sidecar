package config

type RedisConfig struct {
	DockerEndpoint string `env:"DOCKER_ENDPOINT"`
	DockerImage    string `env:"DOCKER_IMAGE"`
	ImageTag       string `env:"DOCKER_IMAGE_TAG"`
}
