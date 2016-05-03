package provisioner

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/hpcloud/catalog-service-manager/services/dev-mssql/config"
	"github.com/pivotal-golang/lager"
)

// fmt template paramters: 1.databaseId
var createDatabaseTemplate = []string{
	"create database [%[1]v] containment = partial",
}

// fmt template parameters: 1.databaseId
var deleteDatabaseTemplate = []string{
	"alter database [%[1]v] set single_user with rollback immediate",
	"drop database [%[1]v]",
}

// fmt template parameters: 1.databaseId, 2.userId, 3.password
var createUserTemplate = []string{
	"use [%[1]v]",
	"create user [%[2]v] with password='%[3]v'",
	"alter role [db_owner] add member [%[2]v]",
	"use master",
}

// fmt template parameters: 1.databaseId, 2.userId
var deleteUserTemplate = []string{
	"use [%[1]v]",
	"drop user [%[2]v]",
	"use master",
}

// fmt template paramters: 1.databaseId
var isDatabaseCreatedTemplate = "select count(*)  from [master].sys.databases  where name = '%[1]v'"

// fmt template parameters: 1.databaseId, 2.userId
var isUserCreatedTemplate = "select count(*)  from [%[1]v].sys.database_principals  where name = '%[2]v'"

type GoMssqlProvisioner struct {
	dbClient    *sql.DB
	logger      lager.Logger
	isConnected bool
	conf        config.MssqlConfig
}

func NewGoMssqlProvisioner(logger lager.Logger, conf config.MssqlConfig) MssqlProvisioner {
	return &GoMssqlProvisioner{
		dbClient:    nil,
		logger:      logger,
		isConnected: false,
		conf:        conf,
	}
}

func buildConnectionString(conf config.MssqlConfig) string {
	return fmt.Sprintf("server=%[1]v;port=%[2]v;user id=%[3]v;password=%[4]v;",
		conf.Host, conf.Port, conf.User, conf.Pass)
}

func (provisioner *GoMssqlProvisioner) connect() error {

	if provisioner.isConnected {
		return nil
	}

	var err error = nil
	connString := buildConnectionString(provisioner.conf)

	provisioner.dbClient, err = sql.Open("mssql", connString)
	if err != nil {
		return err
	}

	// Set idle connections to 0 to prevent keeping open databases
	provisioner.dbClient.SetMaxIdleConns(0)

	err = provisioner.dbClient.Ping()
	if err != nil {
		return err
	}
	provisioner.isConnected = true

	return nil
}

func (provisioner *GoMssqlProvisioner) CreateDatabase(databaseId string) error {
	err := provisioner.connect()
	if err != nil {
		return err
	}

	err = provisioner.executeTemplateWithoutTx(createDatabaseTemplate, databaseId)
	return err
}

func (provisioner *GoMssqlProvisioner) DeleteDatabase(databaseId string) error {
	err := provisioner.connect()
	if err != nil {
		return err
	}
	return provisioner.executeTemplateWithoutTx(deleteDatabaseTemplate, databaseId)
}

func (provisioner *GoMssqlProvisioner) CreateUser(databaseId, userId, password string) error {
	err := provisioner.connect()
	if err != nil {
		return err
	}
	return provisioner.executeTemplateWithTx(createUserTemplate, databaseId, userId, password)
}

func (provisioner *GoMssqlProvisioner) DeleteUser(databaseId, userId string) error {
	err := provisioner.connect()
	if err != nil {
		return err
	}
	return provisioner.executeTemplateWithTx(deleteUserTemplate, databaseId, userId)
}

func (provisioner *GoMssqlProvisioner) IsDatabaseCreated(databaseId string) (bool, error) {
	err := provisioner.connect()
	if err != nil {
		return false, err
	}
	res := 0

	err = provisioner.queryScalarTemplate(isDatabaseCreatedTemplate, &res, databaseId)
	if err != nil {
		return false, err
	}
	if res == 1 {
		return true, nil
	}
	return false, nil
}

func (provisioner *GoMssqlProvisioner) IsUserCreated(databaseId, userId string) (bool, error) {
	err := provisioner.connect()
	if err != nil {
		return false, err
	}
	res := 0

	err = provisioner.queryScalarTemplate(isUserCreatedTemplate, &res, databaseId, userId)
	if err != nil {
		return false, err
	}
	if res == 1 {
		return true, nil
	}
	return false, nil
}

func (provisioner *GoMssqlProvisioner) queryScalarTemplate(template string, output interface{}, targs ...interface{}) error {
	sqlLine := compileTemplate(template, targs...)

	provisioner.logger.Debug("mssql-exec", lager.Data{"query": sqlLine})
	rowRes := provisioner.dbClient.QueryRow(sqlLine)

	err := rowRes.Scan(output)
	if err != nil {
		provisioner.logger.Error("mssql-exec", err, lager.Data{"query": sqlLine})
		return err
	}

	return nil
}

func (provisioner *GoMssqlProvisioner) executeTemplateWithTx(template []string, targs ...interface{}) error {
	tx, err := provisioner.dbClient.Begin()
	if err != nil {
		return err
	}

	for _, templateLine := range template {
		sqlLine := compileTemplate(templateLine, targs...)

		provisioner.logger.Debug("mssql-exec", lager.Data{"query": sqlLine})
		_, err = tx.Exec(sqlLine)
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				panic(rollbackErr.Error())
			}
			provisioner.logger.Error("mssql-exec", err, lager.Data{"query": sqlLine})
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (provisioner *GoMssqlProvisioner) executeTemplateWithoutTx(template []string, targs ...interface{}) error {
	for _, templateLine := range template {
		sqlLine := compileTemplate(templateLine, targs...)

		provisioner.logger.Debug("mssql-exec", lager.Data{"query": sqlLine})
		_, err := provisioner.dbClient.Exec(sqlLine)
		if err != nil {
			provisioner.logger.Error("mssql-exec", err, lager.Data{"query": sqlLine})
			return err
		}
	}

	return nil
}

func compileTemplate(template string, targs ...interface{}) string {
	compiled := fmt.Sprintf(template, targs...)
	extraErrorStart := strings.LastIndex(compiled, "%!(EXTRA")
	if extraErrorStart != -1 {
		// trim the extra args errs from sprintf
		compiled = compiled[0:extraErrorStart]
	}
	return compiled
}
