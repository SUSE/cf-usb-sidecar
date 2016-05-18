package provisioner

import "database/sql"

type MySQLProvisioner interface {
	IsDatabaseCreated(string) (bool, error)
	IsUserCreated(string) (bool, error)
	CreateDatabase(string) error
	DeleteDatabase(string) error
	Query(string, ...interface{}) (*sql.Rows, error)
	CreateUser(string, string, string) error
	DeleteUser(string) error
}
