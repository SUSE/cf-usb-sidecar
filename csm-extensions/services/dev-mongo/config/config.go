package config

type MongoDriverConfig struct {
	User string `env:"MONGO_USER"`
	Pass string `env:"MONGO_PASS"`
	Host string `env:"MONGO_HOST"`
	Port string `env:"MONGO_PORT" envDefault:"27017"`
}
