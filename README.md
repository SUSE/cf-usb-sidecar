# Catalog Service Manager

This is the repository for the catalog service manager(aka Side car container).
This repository holds the rest API for Catalog service manager.

<!-- TOC depthFrom:1 depthTo:6 withLinks:1 updateOnSave:1 orderedList:0 -->

- [Setting up project environment/dependencies](#setting-up-dependencies)
- [Environment Variables](#environment-variables)
- [Run the service](#run-the-service)
  - [In go](#In-go)
  - [On docker](#With-container-in-docker)
- [Testing CSM Service] (#testing-csm-service)
- [Docker Containers for Catalog Service Manager](#docker-containers-for-catalog-service-manager)

## Setting up Dependencies

Please make sure you have git and mercurial installed on your machine.

### install mercurial on mac

Either you can use brew to install mercurial

```
brew install mercurial
```

or you can also download latest mercurial package from their site.

### install go tools/dependencies for catalog service manager

```
make tools
```

## Environment Variables
Catalog Service Manager reads it configuration/settings from environment variables. 

### CSM_PARAMETERS 
This is the list of environment variables that Catalog Service Manager should read, this needs to be set if you are running catalog service manager in the default mode (specifically without create connections extensions). When extensions are not catalog service manager reads this environment and returns all the variables listes here.
present 
```
export CSM_PARAMETERS="HOSTNAME USERNAME"
```
 

### CSM_HOME
This is the path of the CSM home directory, default value for this is /catalog-service-manager which is the path on the Catalog Service Manager's base docker image where all the extensions are available. There is no need to change the value of this variable in the SDL, if you are copying your extensions in the default path (which is strongly recommended). If you are running catalog service manager on your machine and trying to test the extensions scripts, you need to set this variable appropriately to point to the directory where your workspaces and connections extensions are located.   

```
export CSM_HOME=${GOPATH}/src/github.com/hpcloud/catalog-service-manager/examples/mysql 
``` 

### CSM_DEBUG
Set this environment variable if you want catalog service manager to generate logs in Debug mode. By default catalog service manager runs with log level Info, but for debugging/testing you can set this environment variable and get more details from catalog service manager logs. It is strongly recommended to not set this environment variable in your SDL (that you plan on adding to service catalog), as debug logs may give our sensitive information which is a security risk.
```
export CSM_DEBUG=true
``` 

### DEV_MODE
Set this environment variable if you want catalog service manager to keep the output output files written by the extensions. By default catalog service manager runs with DEV_MOD set to off, and it alwasy deletes the output file where extensions write their output. 
It is strongly recommended to not set this environment variable in your SDL (that you plan on adding to service catalog), as leaving these output files on the disk may pose a is a security risk, as these files may contain credentials to the service.
```
export DEV_MODE=true
``` 


## Run the service
### In go
you can use make run to run the service on command line
```
make run
```
This will start service listening on port http://0.0.0.0:8081 on the machine where you run this command.
You can use above mentioned environment variables (CSM_HOME and CSM_DEBUG) to alter the behavior of the service. (specifically if you are debugging an issue or while writing new extensions )


### Run service with exampl/mysql extensions
1. To start with you need mysql instance to which you can connect to, if you have one already you can skipt this first step, if you don't then run following commands
```
cd  example/mysql
make tools
```
This should start a mysql container and bind port 3306 on the docker host to mysql port on the container
After this is done, you should be able to connect to the mysql with 

```
MY_DOCKER_HOST_IP=`env | grep DOCKER_HOST | cut -d "/" -f 3 | cut -d ":" -f 1`
mysql -h ${MY_DOCKER_HOST_IP} -P 3306 -u root -proot123

```

2. Set following environment variables

```
export CSM_HOME=`pwd`/examples/mysql/CSM_HOME
export MYSQL_SERVICE_HOST=`env | grep DOCKER_HOST | cut -d "/" -f 3 | cut -d ":" -f 1`
export MYSQL_SERVICE_PORT_MYSQL=3306
export MYSQL_ROOT_PASSWORD=root123
```

3. Run the make command to start the service
```
make run
```

## Testing CSM Service

### create workspace

curl -X POST http://localhost:8081/workspaces -H "content-type: application/json" -d '{"workspace_id":"test_workspace"}'

### get workspace

curl -X GET http://localhost:8081/workspaces/test_workspace

### create credentials

curl -X POST http://localhost:8081/workspaces/test_workspace/connections -H "content-type: application/json" -d '{"connection_id":"test_user"}'

### get credentials

curl -X GET http://localhost:8081/workspaces/test_workspace/connections/test_user

### delete credentials

curl -X DELETE http://localhost:8081/workspaces/test_workspace/connections/test_user

### delete workspace

curl -X DELETE http://localhost:8081/workspaces/test_workspace