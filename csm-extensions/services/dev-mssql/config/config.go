package config

type MssqlConfig struct {
	User               string `env:"MSSQL_USER"`
	Pass               string `env:"MSSQL_PASS"`
	Host               string `env:"MSSQL_HOST"`
	Port               string `env:"MSSQL_PORT" envDefault:"1433"`
	DbIdentifierPrefix string `env:"MSSQL_DBPREFIX" envDefault:"d"`
	UcpDomainSuffix    string `env:"UCP_SERVICE_DOMAIN_SUFFIX"`
}
