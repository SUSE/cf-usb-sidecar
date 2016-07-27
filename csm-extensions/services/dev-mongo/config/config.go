package config

type MongoDriverConfig struct {
	User string `env:"SERVICE_MONGO_USER"`
	Pass string `env:"SERVICE_MONGO_PASS"`
	Host string `env:"SERVICE_MONGO_HOST"`
	Port string `env:"SERVICE_MONGO_PORT" envDefault:"27017"`
}
