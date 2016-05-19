package config

type PostgresConfig struct {
	User            string `env:"POSTGRES_USER"`
	Password        string `env:"POSTGRES_PASSWORD"`
	Host            string `env:"POSTGRES_HOST"`
	Port            string `env:"POSTGRES_PORT" envDefault:"5432"`
	Dbname          string `env:"POSTGRES_DBNAME"`
	Sslmode         string `env:"POSTGRES_SSLMODE"`
	UcpDomainSuffix string `env:"UCP_SERVICE_DOMAIN_SUFFIX"`
}
