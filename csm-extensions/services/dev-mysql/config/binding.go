package config

import "fmt"

type MySQLBinding struct {
	Host     string `json:"host"`
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	Username string `json:"username"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
	JdbcUrl  string `json:"jdbcUrl"`
}

var JdbcUrilTemplate = "jdbc:mysql://%[1]v:%[2]v/%[3]v?user=%[4]v&password=%[5]v"

func GenerateConnectionString(input string, hostname string, port string, databaseName string, username string, password string) string {
	return fmt.Sprintf(input, hostname, port, databaseName, username, password)
}
