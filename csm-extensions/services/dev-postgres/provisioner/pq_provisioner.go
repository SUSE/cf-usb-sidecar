package provisioner

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	"github.com/SUSE/cf-usb-sidecar/csm-extensions/services/dev-postgres/config"
	"github.com/hpcloud/go-csm-lib/util"
	_ "github.com/lib/pq"
	"github.com/pivotal-golang/lager"
)

var createDatabaseQuery = "CREATE DATABASE {{.Database}}"
var createDatabaseAdminQuery = "CREATE ROLE \"{{.DbAdmin}}\" LOGIN PASSWORD '{{.DbAdminPassword}}'"
var dbCountQuery = "SELECT COUNT(*) FROM pg_database WHERE datname = '{{.Database}}'"
var createRoleQuery = "CREATE ROLE \"{{.User}}\" LOGIN PASSWORD '{{.Password}}' INHERIT"
var grantAllPrivToRoleQuery = "GRANT ALL PRIVILEGES ON DATABASE {{.Database}} TO \"{{.User}}\""
var grantDbAdminPrivToRoleQuery = "GRANT {{.DbAdmin}} TO \"{{.User}}\""
var userCountQuery = "SELECT COUNT(*) FROM pg_roles WHERE rolname = '{{.User}}'"
var revokeAllPrivFromRoleQuery = "REVOKE ALL PRIVILEGES ON DATABASE {{.Database}} FROM \"{{.User}}\""
var reassignRoleOwned = "REASSIGN OWNED BY \"{{.User}}\" TO {{.DbAdmin}}"
var deleteRoleQuery = "DROP ROLE \"{{.User}}\""
var terminateDatabaseConnQuery = "SELECT pg_terminate_backend(pg_stat_activity.{{ .PidColumn }}) FROM pg_stat_activity WHERE pg_stat_activity.datname = '{{.Database}}' AND {{ .PidColumn }} <> pg_backend_pid()"
var deleteDatabaseQuery = "DROP DATABASE {{.Database}}"

type PqProvisioner struct {
	pgClient  *sql.DB
	conf      config.PostgresConfig
	logger    lager.Logger
	connected bool
}

func NewPqProvisioner(logger lager.Logger, config config.PostgresConfig) PostgresProvisionerInterface {
	return &PqProvisioner{logger: logger, conf: config}
}

func (provisioner *PqProvisioner) connect() error {
	if provisioner.connected {
		return nil
	}
	var err error = nil
	connString := buildConnectionString(provisioner.conf)
	provisioner.pgClient, err = sql.Open("postgres", connString)

	if err != nil {
		return err
	}

	err = provisioner.pgClient.Ping()
	if err != nil {
		return err
	}

	provisioner.connected = true
	return nil
}

func (provisioner *PqProvisioner) CreateDatabase(dbname string) error {
	err := provisioner.connect()
	if err != nil {
		return err
	}

	// for pg driver, create database can not be executed in transaction
	err = provisioner.executeQueryNoTx([]string{createDatabaseQuery}, map[string]string{"Database": dbname})
	if err != nil {
		return err
	}

	password, err := util.SecureRandomString(15)
	if err != nil {
		return err
	}

	err = provisioner.executeQueryTx([]string{createDatabaseAdminQuery, grantAllPrivToRoleQuery}, map[string]string{"DbAdmin": dbname, "DbAdminPassword": password, "Database": dbname, "User": dbname})
	if err != nil {
		return err
	}

	return nil
}

func (provisioner *PqProvisioner) DeleteDatabase(dbname string) error {
	err := provisioner.connect()
	if err != nil {
		return err
	}

	version, err := provisioner.getServerVersion()
	if err != nil {
		return err
	}

	var pidColumn string
	if version > 90200 {
		pidColumn = "pid"
	} else {
		pidColumn = "procpid"
	}

	err = provisioner.executeQueryTx([]string{terminateDatabaseConnQuery}, map[string]string{
		"Database":  dbname,
		"PidColumn": pidColumn,
	})
	if err != nil {
		return err
	}

	// for pg driver, drop database can not be executed in transaction
	err = provisioner.executeQueryNoTx([]string{deleteDatabaseQuery}, map[string]string{"Database": dbname})
	if err != nil {
		return err
	}

	return nil
}

func (provisioner *PqProvisioner) DatabaseExists(dbname string) (bool, error) {
	err := provisioner.connect()
	if err != nil {
		return false, err
	}

	res, err := provisioner.executeQueryRow(dbCountQuery, map[string]string{"Database": dbname})
	if err != nil {
		return false, err
	}

	if res.(int64) == 1 {
		return true, nil
	}

	return false, nil
}

func (provisioner *PqProvisioner) CreateUser(dbname string, username string, password string) error {
	err := provisioner.connect()
	if err != nil {
		return err
	}

	err = provisioner.executeQueryTx([]string{createRoleQuery}, map[string]string{"User": username, "Password": password})
	if err != nil {
		return err
	}

	err = provisioner.executeQueryTx([]string{grantDbAdminPrivToRoleQuery}, map[string]string{"User": username, "DbAdmin": dbname})
	if err != nil {
		return err
	}
	return nil
}

func (provisioner *PqProvisioner) DeleteUser(dbname string, username string) error {

	originalDb := provisioner.conf.Dbname
	provisioner.conf.Dbname = dbname
	err := provisioner.connect()
	if err != nil {
		return err
	}

	err = provisioner.executeQueryTx([]string{reassignRoleOwned}, map[string]string{"User": username, "DbAdmin": dbname})
	if err != nil {
		return err
	}

	provisioner.conf.Dbname = originalDb
	err = provisioner.connect()
	if err != nil {
		return err
	}

	err = provisioner.executeQueryTx([]string{revokeAllPrivFromRoleQuery, deleteRoleQuery}, map[string]string{"User": username, "Database": dbname})
	if err != nil {
		return err
	}

	return nil
}

func (provisioner *PqProvisioner) UserExists(username string) (bool, error) {
	err := provisioner.connect()
	if err != nil {
		return false, err
	}

	res, err := provisioner.executeQueryRow(userCountQuery, map[string]string{"User": username})
	if err != nil {
		return false, err
	}

	if res.(int64) == 1 {
		return true, nil
	}

	return false, nil
}

func buildConnectionString(connectionParams config.PostgresConfig) string {
	var res string = fmt.Sprintf("user=%v password=%v host=%v port=%v dbname=%v sslmode=%v",
		connectionParams.User, connectionParams.Password, connectionParams.Host, connectionParams.Port, connectionParams.Dbname, connectionParams.Sslmode)
	return res
}

func parametrizeQuery(query string, params map[string]string) (string, error) {
	queryTemplate := template.Must(template.New("query").Parse(query))
	output := bytes.Buffer{}
	queryTemplate.Execute(&output, params)

	queryString := output.String()

	if strings.Contains(queryString, "<no value>") {
		return queryString, errors.New("Invalid parameter passed to query")
	}

	return queryString, nil
}

func (provisioner *PqProvisioner) executeQueryNoTx(queries []string, params map[string]string) error {
	for _, query := range queries {
		pQuery, err := parametrizeQuery(query, params)

		provisioner.logger.Debug("postgres-exec", lager.Data{"query": pQuery})
		if err != nil {
			provisioner.logger.Error("postgres-exec", err, lager.Data{"query": pQuery})
			return err
		}

		_, err = provisioner.pgClient.Exec(pQuery)
		if err != nil {
			provisioner.logger.Error("postgres-exec", err, lager.Data{"query": pQuery})
			return err
		}
	}

	return nil
}

func (provisioner *PqProvisioner) executeQueryTx(queries []string, params map[string]string) error {
	tx, err := provisioner.pgClient.Begin()
	if err != nil {
		return err
	}

	for _, query := range queries {
		pQuery, err := parametrizeQuery(query, params)
		provisioner.logger.Debug("postgres-exec", lager.Data{"query": pQuery})
		if err != nil {
			provisioner.logger.Error("postgres-exec", err, lager.Data{"query": pQuery})
			return err
		}

		_, err = tx.Exec(pQuery)
		if err != nil {
			tx.Rollback()
			provisioner.logger.Error("postgres-exec", err, lager.Data{"query": pQuery})
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (provisioner *PqProvisioner) executeQueryRow(query string, params map[string]string) (interface{}, error) {
	pQuery, err := parametrizeQuery(query, params)
	provisioner.logger.Debug("postgres-exec", lager.Data{"query": pQuery})
	if err != nil {
		provisioner.logger.Error("postgres-exec", err, lager.Data{"query": pQuery})
		return nil, err
	}

	var res interface{}
	err = provisioner.pgClient.QueryRow(pQuery).Scan(&res)
	if err != nil && err == sql.ErrNoRows {
		provisioner.logger.Error("postgres-exec", err, lager.Data{"query": pQuery})
		return nil, err
	}

	return res, nil
}

func (provisioner *PqProvisioner) getServerVersion() (int, error) {
	res, err := provisioner.executeQueryRow("SHOW server_version_num", map[string]string{})
	if err != nil {
		return 0, err
	}

	i := res.([]uint8)
	b := make([]byte, len(i))
	for i, v := range i {
		if v < 0 {
			b[i] = byte(256 + int(v))
		} else {
			b[i] = byte(v)
		}
	}

	return strconv.Atoi(string(b))
}
