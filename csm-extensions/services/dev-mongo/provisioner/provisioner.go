package provisioner

import (
	"fmt"

	"github.com/hpcloud/catalog-service-manager/csm-extensions/services/dev-mongo/config"
	"github.com/pivotal-golang/lager"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoProvisioner struct {
	Config     config.MongoDriverConfig
	Connection *mgo.Session
	logger     lager.Logger
	connected  bool
}

func New(config config.MongoDriverConfig, logger lager.Logger) MongoProvisionerInterface {
	return &MongoProvisioner{Config: config, logger: logger}
}

func (e *MongoProvisioner) connect(mongoConfig config.MongoDriverConfig) error {
	var err error
	e.Config = mongoConfig

	var connString string
	if e.Config.User != "" && e.Config.Pass != "" {
		connString = fmt.Sprintf("mongodb://%s:%s@%s:%s", mongoConfig.User, mongoConfig.Pass, mongoConfig.Host, mongoConfig.Port)
	} else {
		connString = fmt.Sprintf("mongodb://%s:%s", mongoConfig.Host, mongoConfig.Port)
	}

	e.Connection, err = mgo.Dial(connString)
	if err != nil {
		e.logger.Info("connect", lager.Data{"connection": connString})
		e.logger.Error("Error loging into the mongo db service", err)
	}
	return err
}

func (e *MongoProvisioner) Close() {
	e.Connection.Close()
}

func (e *MongoProvisioner) IsDatabaseCreated(databaseName string) (bool, error) {
	if e.connected == false {
		err := e.connect(e.Config)
		if err != nil {
			return false, err
		}
	}
	databases, err := e.Connection.DatabaseNames()
	if err != nil {
		return false, err
	}

	for _, db := range databases {
		if db == databaseName {
			return true, nil
		}
	}

	return false, nil
}

func (e *MongoProvisioner) IsUserCreated(databaseName string, userName string) (bool, error) {
	if e.connected == false {
		err := e.connect(e.Config)
		if err != nil {
			return false, err
		}
	}
	userDB := e.Connection.DB(databaseName)
	result := bson.M{}
	err := userDB.Run(bson.M{"usersInfo": userName}, &result)
	if err != nil {
		return false, err
	}
	userInfo := result["users"].([]interface{})
	if len(userInfo) > 0 {
		return true, nil
	}
	return false, nil
}

func (e *MongoProvisioner) CreateDatabase(databaseName string) error {
	if e.connected == false {
		err := e.connect(e.Config)
		if err != nil {
			return err
		}
	}
	//this should create the db with empty users collection
	coll := e.Connection.DB(databaseName).C("sample")
	coll.Insert(bson.M{"a": 1, "b": 2})

	result := struct{ A, B int }{}

	err := coll.Find(bson.M{"a": 1}).One(&result)
	err = coll.DropCollection()
	if err != nil {
		return err
	}

	return nil
}

func (e *MongoProvisioner) DeleteDatabase(databaseName string) error {
	if e.connected == false {
		err := e.connect(e.Config)
		if err != nil {
			return err
		}
	}
	err := e.Connection.DB(databaseName).DropDatabase()
	if err != nil {
		return err
	}
	return nil
}

func (e *MongoProvisioner) CreateUser(databaseName string, username string, password string) error {
	if e.connected == false {
		err := e.connect(e.Config)
		if err != nil {
			return err
		}
	}
	userDB := e.Connection.DB(databaseName)

	err := userDB.AddUser(username, password, false)
	if err != nil {
		return err
	}
	return nil
}

func (e *MongoProvisioner) DeleteUser(databaseName string, username string) error {
	if e.connected == false {
		err := e.connect(e.Config)
		if err != nil {
			return err
		}
	}
	err := e.Connection.DB(databaseName).RemoveUser(username)
	if err != nil {
		return err
	}
	return nil
}
