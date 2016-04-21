# Extending Catalog Service Manager

catalog-service-manager's docker image named catalog-service-manager:base can be used to extend CSM for service specific requirements. Following document describes various CSM actions which can be controlled by providing necessary extensions. All these actions are optional.

By default all extensions are optional.

Only Create connection action has a default implementation where CSM tries to read environment variable CSM_PARAMETERS. This environment variable is defined as list of ENV variables that CSM should expose in the create connections API. If this variable is not defined then even the create connection action is treated as NOP.

<!-- TOC depthFrom:1 depthTo:6 withLinks:1 updateOnSave:1 orderedList:0 -->

- [CSM Extension contract] (#csm-extension-contract)
- [CSM Extensions](#csm-extensions)

## CSM Extension contract:
When CSM executes any extension, it supplies set of command line parameter. These parameters change as per extension and are described in detail under corresponding extension section.

### Output file
The first parameter to the extension being executed is always the path of the output file. CSM generates a random file before executing an extension, supplies that path to the extension as the first argument and extension is suppose to write the necessary details about the execution to that file. After extension returns, CSM tries to interpret the output file (if the output conforms to the expected schema), there after the file is removed from the file system.

### Output format
Extension is exepcted to return the results of the operation via writing JSON object to the provided output file. The format of the file should conform to the json schema as described in [extension_output_schema.json](extension_output_schema.json).

* **status** field should indicate the status of the operation
* **details** field is a map of key value pair, where key is the name of the parameter and its corresponding value for the running instance. 

After reading the output file CSM identifies the status of the operation by looking at status field, if the status is successful, the details section is returned in the API call back to the user, so details section should list all the parameters that the user of the service will need in order to interact with the workspaces or connections that are being referenced in the API.

#### Example output
```
{
	"status": "successful",
	"details": {
		"host" : "csm-mysql-kjd8pqo",
		"port" : "3306",
		"username" : "payroll-service-test",
		"password" : "lpq902kjjd01pd",
		"database" : "payroll-service"
	}
}
```
## CSM Extensions
CSM exposes 2 primary set of entities Workspaces and Credentials. Both of these entities support GET,CREATE,DELETE actions, which can be extended by the service providers.

CSM also allows extending setup (lifecycle) phases of the CSM container, allows service authors to extended startup and shutdown of the container. (current implementation of the CSM doesn't support the shutdown workflow completely).

### setup (lifecycle)
These action correspond to the lifecycle of the CSM container, these include startups and shutdown
Setup actions can be extended to change the behavior of startup and shutdown phases of the CSM container.

#### startup
startup extension will run every time CSM service starts on the container. If the extension is present, CSM service run the extension and reads the response. If extension is present then CSM checks for the status of the execution, **if extension fails for any reason then CSM service will not start.**

#### shutdown
shutdown extension will run every time  USB driver instance is deleted or particular service instance is disabled from a service.


### workspaces
Workspace in the CSM represents logical space created for a particular application/user inside a multi-tenant system.

For example, workspace on a mysql instance would represent a logical database that can be created on that instance. Since we can create multiple databases on a single mysql instance, we can provisioning multiple workspaces and associate each application with one workspace/database. So this makes mysql a multi-workspace/multi-tenant service.

#### get
Extension is executed with following command line parameters
```
/catalog-service-manager/workspace/get/getWorkspace <path of the output file> WorkspaceID
```
The workspaceID is usually the uniq ID that is supplied by the user of the API so that the workspace can be uniquely identified.

This extension is expected to return the status and details of the workspace, via writing appropriate JSON to the output file.

***Output***

1. Status (required)
	It should be `successful`, if workspace already exist, active and ready to use. Otherwise status should be `failed`
2. Details (optional)
	It can be list of parameters, users of the service will need to interact with the workspace

This extension is called for GET /workspaces/{workspace_id} API.

#### delete
Extension is executed with following command line parameters
```
/catalog-service-manager/worksapces/get/deleteWorkspace <path of the output file> workspaceID
```
The workspaceID is usually the uniq ID that is supplied by the user of the API so that the workspace can be uniquely identified.

This extension is expected to delete/cleanup the workspace and return the status of the operation via writing appropriate JSON to the output file.
This extension is called for DELETE /workspaces/{workspace_id} API.

***Output***

1. Status (required)
	It should be returned `successful`, if workspace is deleted successfully. Otherwise status should be `failed`

This extension is called for GET /workspaces/{workspace_id} API.

#### create
Extension is executed with following command line parameters
```
/catalog-service-manager/worksapces/get/createWorkspace <path of the output file> workspaceID
```
The workspaceID is usually the uniq ID that is supplied by the user of the API so that the workspace can be uniquely identified.

This extension is expected to create new workspace and return the status of the operation and details of the workspace via writing appropriate JSON to the output file.

***Output***

1. Status (required)
	It should be `successful`, if workspace is created, active and ready to use. Otherwise status should be `failed`
2. Details (optional)
	It can be list of parameters, users of the service will need to interact with the workspace

This extension is called for POST /workspaces API.

### connections
Connection in the CSM context represent logical connection/connections created for a particular application/users inside a multi-tenant system.

For example, connection on a mysql instance would represent a separate mysql user whose access is restricted to specific database schema and/or is controlled via specific grants. The database which represents a workspace can be accessed by multiple users/connections so Workspaces to connections is a one to many

#### get
Extension is executed with following command line parameters
```
/catalog-service-manager/connections/get/getConnection <path of the output file> workspaceID connectionID
```
The workspaceID and connectionID are usually the uniq IDs that are supplied by the user of the API so that the workspace and connection can be uniquely identified.

***Output***

1. Status (required)
	It should be `successful`, if connection already exist, active and ready to use. Otherwise status should be `failed`
2. Details (optional)
	It can be list of parameters, users of the service will need to interact with the connection (it is recommended that workspace details be included here as well)


This extension is called for GET /workspaces/{workspace_id}/connections/{connection_id} API.

#### delete
Extension is executed with following command line parameters
```
/catalog-service-manager/connections/get/deleteConnection <path of the output file> workspaceID connectionID
```
The workspaceID and connectionID are usually the uniq IDs that are supplied by the user of the API so that the workspace and connection can be uniquely identified.

This extension is expected to delete/cleanup the connection, via writing appropriate JSON to the output file.

***Output***

1. Status (required)
	It should be returned `successful`, if connection is deleted successfully. Otherwise status should be `failed`

This extension is called for DELETE /workspaces/{workspace_id}/connections/{connection_id} API.


#### create
Extension is executed with following command line parameters
```
/catalog-service-manager/connections/get/createConnection <path of the output file> workspaceID connectionID
```
The workspaceID and connectionID are usually the uniq IDs that are supplied by the user of the API so that the workspace and connection can be uniquely identified.

This extension is expected to create new connection and return the status of the operation and details of the connection via writing appropriate JSON to the output file.

***Output***

1. Status (required)
	It should be `successful`, if connection is created, active and ready to use. Otherwise status should be `failed`
2. Details (optional)
	It can be list of parameters, users of the service will need to interact with the workspace(it is recommended that workspace details be included here as well)


This extension is called for POST /workspaces/{workspace_id}/connections API.
