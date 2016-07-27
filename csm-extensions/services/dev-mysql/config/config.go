package config

type MySQLConfig struct {
	User string `env:"SERVICE_MYSQL_USER"`
	Pass string `env:"SERVICE_MYSQL_PASS"`
	Host string `env:"SERVICE_MYSQL_HOST"`
	Port string `env:"SERVICE_MYSQL_PORT" envDefault:"3306"`
}
