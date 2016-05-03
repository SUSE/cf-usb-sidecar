package config

type MySQLBinding struct {
	Host     string `json:"host"`
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
	Username string `json:"username"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}
