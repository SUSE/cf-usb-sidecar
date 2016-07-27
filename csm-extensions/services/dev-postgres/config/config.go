package config

type PostgresConfig struct {
	User     string `env:"SERVICE_POSTGRES_USER"`
	Password string `env:"SERVICE_POSTGRES_PASSWORD"`
	Host     string `env:"SERVICE_POSTGRES_HOST"`
	Port     string `env:"SERVICE_POSTGRES_PORT" envDefault:"5432"`
	Dbname   string `env:"SERVICE_POSTGRES_DBNAME"`
	Sslmode  string `env:"SERVICE_POSTGRES_SSLMODE"`
}
