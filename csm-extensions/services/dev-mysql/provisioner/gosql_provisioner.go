package provisioner

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hpcloud/catalog-service-manager/csm-extensions/services/dev-mysql/config"
	"github.com/pivotal-golang/lager"
)

type GoSqlProvisioner struct {
	Conf        config.MySQLConfig
	Connection  *sql.DB
	logger      lager.Logger
	isConnected bool
}

func NewGoSQL(logger lager.Logger, config config.MySQLConfig) MySQLProvisioner {
	return &GoSqlProvisioner{logger: logger, Conf: config}
}

func (e *GoSqlProvisioner) connect() error {
	var err error

	if e.isConnected {
		return nil
	}

	e.Connection, err = e.openSqlConnection()
	if err != nil {
		return err
	}
	err = e.Connection.Ping()
	if err == nil {
		e.isConnected = true
	}
	return err
}

func (e *GoSqlProvisioner) IsDatabaseCreated(databaseName string) (bool, error) {
	err := e.connect()
	if err != nil {
		return false, err
	}
	rows, err := e.Query("SHOW DATABASES WHERE `database` = ?", databaseName)
	if err != nil {
		return false, err
	}

	var (
		result    [][]string
		container []string
		pointers  []interface{}
	)

	cols, err := rows.Columns()
	if err != nil {
		return false, err
	}

	length := len(cols)

	for rows.Next() {
		pointers = make([]interface{}, length)
		container = make([]string, length)

		for i := range pointers {
			pointers[i] = &container[i]
		}

		err = rows.Scan(pointers...)
		if err != nil {
			return false, err
		}

		result = append(result, container)
	}
	for _, cont := range result {
		if cont[0] == databaseName {
			return true, nil
		}
	}

	return false, nil
}

func (e *GoSqlProvisioner) IsUserCreated(userName string) (bool, error) {
	err := e.connect()
	if err != nil {
		return false, err
	}
	rows, err := e.Query("SELECT user from mysql.user WHERE user = ?", userName)

	if err != nil {
		return false, err
	}

	var (
		result    [][]string
		container []string
		pointers  []interface{}
	)

	cols, err := rows.Columns()
	if err != nil {
		return false, err
	}

	length := len(cols)

	for rows.Next() {
		pointers = make([]interface{}, length)
		container = make([]string, length)

		for i := range pointers {
			pointers[i] = &container[i]
		}

		err = rows.Scan(pointers...)
		if err != nil {
			return false, err
		}

		result = append(result, container)
	}
	for _, cont := range result {
		if cont[0] == userName {
			return true, nil
		}
	}

	return false, nil
}

func (e *GoSqlProvisioner) CreateDatabase(databaseName string) error {
	err := e.connect()
	if err != nil {
		return err
	}

	err = e.executeTransaction(e.Connection, fmt.Sprintf("CREATE DATABASE %s", databaseName))
	if err != nil {
		e.logger.Error("create database", err)
		return err
	}

	return nil
}

func (e *GoSqlProvisioner) DeleteDatabase(databaseName string) error {
	err := e.connect()
	if err != nil {
		return err
	}

	err = e.executeTransaction(e.Connection, fmt.Sprintf("DROP DATABASE %s", databaseName))
	if err != nil {
		e.logger.Error("delete database", err)
		return err
	}
	return nil
}

func (e *GoSqlProvisioner) Query(query string, args ...interface{}) (*sql.Rows, error) {

	var err error
	var result *sql.Rows
	err = e.connect()
	if err != nil {
		return nil, err
	}

	if len(args) > 0 {
		result, err = e.Connection.Query(query, args...)
	} else {
		result, err = e.Connection.Query(query)
	}
	if err != nil {
		e.logger.Error("query", err)
		return nil, err
	}
	return result, nil
}

func (e *GoSqlProvisioner) CreateUser(databaseName string, username string, password string) error {

	e.logger.Info("Connection open - executing transaction")
	err := e.connect()
	if err != nil {
		return err
	}

	err = e.executeTransaction(e.Connection,
		fmt.Sprintf("CREATE USER '%s' IDENTIFIED BY '%s';", username, password),
		fmt.Sprintf("GRANT ALL PRIVILEGES ON %s.* TO '%s'@'%%';", databaseName, username),
		"FLUSH PRIVILEGES;")
	e.logger.Info("Transaction done")
	if err != nil {
		e.logger.Error("create user", err)
		return err
	}

	return nil

}

func (e *GoSqlProvisioner) DeleteUser(username string) error {
	err := e.connect()
	if err != nil {
		return err
	}

	err = e.executeTransaction(e.Connection, fmt.Sprintf("DROP USER '%s'", username))
	if err != nil {
		e.logger.Error("delete user", err)
		return err
	}

	return nil
}

func (e *GoSqlProvisioner) openSqlConnection() (*sql.DB, error) {

	con, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/mysql?interpolateParams=true", e.Conf.User, e.Conf.Pass, e.Conf.Host, e.Conf.Port))
	if err != nil {
		return nil, err
	}
	return con, nil
}

func (e *GoSqlProvisioner) executeTransaction(con *sql.DB, querys ...string) error {
	tx, err := con.Begin()
	if err != nil {
		e.logger.Error("execute transaction", err)
		return err
	} else {
		for _, query := range querys {
			e.logger.Info(query)
			_, err = tx.Exec(query)
			if err != nil {
				e.logger.Error("execute transaction query", err)
				tx.Rollback()
				break
			}
		}
		tx.Commit()
	}
	if err != nil {
		return err
	}
	return nil
}
