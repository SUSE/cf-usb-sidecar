package config

type MySQLConfig struct {
	User            string `env:"MYSQL_USER"`
	Pass            string `env:"MYSQL_PASS"`
	Host            string `env:"MYSQL_HOST"`
	Port            string `env:"MYSQL_PORT" envDefault:"3306"`
	UcpDomainSuffix string `env:"UCP_SERVICE_DOMAIN_SUFFIX"`
}
