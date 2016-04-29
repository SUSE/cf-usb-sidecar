package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/rds"
)

//SUCCESS represents successful
const SUCCESS = "successful"

//FAIL represents fail
const FAIL = "failed"

//ServiceManagerConnectionResponse with Detailed msg and status
type ServiceManagerConnectionResponse struct {
	Details map[string]interface{} `json:"details,omitempty"`
	Status  string                 `json:"status,omitempty"`
}

//create vpc SecurityGroup
func createSecurityGroup(region string) string {
	svc := ec2.New(session.New(), aws.NewConfig().WithRegion(region))

	randonID := randStringBytesRmndr(8)
	vpcsecgroupname := "sgname-" + randonID

	params := &ec2.CreateSecurityGroupInput{
		Description: aws.String("VPC security group"), // Required
		GroupName:   aws.String(vpcsecgroupname),      // Required
		DryRun:      aws.Bool(false),
	}
	resp, err := svc.CreateSecurityGroup(params)
	
	if err != nil {
		return ""
	}
	createSecurityGroupInboumdRule(region, *resp.GroupId)
	return *resp.GroupId
}

//create inbound vpc SecurityGroup rule tcp 0.0.0.0/0
func createSecurityGroupInboumdRule(region string, secgroupID string) {

	svc := ec2.New(session.New(), aws.NewConfig().WithRegion(region))

	params := &ec2.AuthorizeSecurityGroupIngressInput{
		CidrIp:     aws.String("0.0.0.0/0"),
		DryRun:     aws.Bool(false),
		FromPort:   aws.Int64(3306),
		GroupId:    aws.String(secgroupID),
		IpProtocol: aws.String("tcp"),
		ToPort:     aws.Int64(3306),
	}
	_, err := svc.AuthorizeSecurityGroupIngress(params)

	if err != nil {
		return
	}
}

// Create an rds database instance in the specified region
func createRdsDbInstance(region string, dbidentifier string, dbclass string, dbtype string, size int64,
	multiAZ bool, masteruser string, masterpwd string, outputfile string) error {

	connectionresponse := make(map[string]interface{})
	var smcr ServiceManagerConnectionResponse

	vpcSecGroupID := createSecurityGroup(region)

	svc := rds.New(session.New(), aws.NewConfig().WithRegion(region))

	params := &rds.CreateDBInstanceInput{
		DBInstanceClass:      aws.String(dbclass),      // Required
		DBInstanceIdentifier: aws.String(dbidentifier), // Required
		Engine:               aws.String(dbtype),       // Required
		AllocatedStorage:     aws.Int64(size),
		MasterUserPassword:   aws.String(masterpwd),
		MasterUsername:       aws.String(masteruser),
		MultiAZ:              aws.Bool(multiAZ),
		VpcSecurityGroupIds: []*string{
			aws.String(vpcSecGroupID), // Required
		},
	}
	_, err := svc.CreateDBInstance(params)

	if err != nil {
		marshalAndWriteToOutputFile(outputfile, ServiceManagerConnectionResponse{Status: FAIL})
		return err
	}
	connectionresponse["region"] = region
	connectionresponse["dbidentifier"] = dbidentifier
	connectionresponse["masteruser"] = masteruser
	connectionresponse["masterpwd"] = masterpwd
	smcr = ServiceManagerConnectionResponse{Details: connectionresponse, Status: SUCCESS}
	marshalAndWriteToOutputFile(outputfile, smcr)
	return nil
}

// Get the rds database instance given the database identifier
func getRdsDbSecGroupID(region string, dbidentifier string) (string, error) {

	svc := rds.New(session.New(), aws.NewConfig().WithRegion(region))

	params := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(dbidentifier),
	}
	resp, err := svc.DescribeDBInstances(params)

	if err != nil {
		return "", err
	}
	//return security group if hostflag is false
	secgroupID := *resp.DBInstances[0].VpcSecurityGroups[0].VpcSecurityGroupId
	return secgroupID, nil

}

// Get the rds database instance host address given the database identifier
func getRdsDbHostAddress(region string, dbidentifier string) (string, error) {

	svc := rds.New(session.New(), aws.NewConfig().WithRegion(region))

	params := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(dbidentifier),
	}
	resp, err := svc.DescribeDBInstances(params)

	if err != nil {
		return "", err
	}
	//amazon-rds host if hostflag is true
	rdshost := *resp.DBInstances[0].Endpoint.Address
	return rdshost, nil
}

//create workspace in the database amazon-rds-mysql
func createWorkspace(region string, dbidentifier string, masteruser string, masterpwd string, workspace string, outputfile string) error {

	connectionresponse := make(map[string]interface{})
	var smcr ServiceManagerConnectionResponse

	rdshost, err := getRdsDbHostAddress(region, dbidentifier)
	if err != nil {
		return err
	}
	connectionstring := masteruser + ":" + masterpwd + "@tcp(" + rdshost + ":3306)/"
	db, err := sql.Open("mysql", connectionstring)
	logErr(err)
	
	// create database
	_, dberr := db.Exec("CREATE DATABASE " + workspace + ";")
	logErr(dberr)
	if dberr == nil {
		_, uberr := db.Exec("USE " + workspace + ";")
		if uberr == nil {
			connectionresponse["workspace"] = workspace
			connectionresponse["port"] = 3306
			connectionresponse["host"] = rdshost
			smcr = ServiceManagerConnectionResponse{Details: connectionresponse, Status: SUCCESS}
			marshalAndWriteToOutputFile(outputfile, smcr)
		} else {
			marshalAndWriteToOutputFile(outputfile, ServiceManagerConnectionResponse{Status: FAIL})
		}
	} else {
		marshalAndWriteToOutputFile(outputfile, ServiceManagerConnectionResponse{Status: FAIL})
	}
	db.Close()
	return nil
}

//get workspace in the database amazon-rds-mysql
func getWorkspace(region string, dbidentifier string, masteruser string, masterpwd string, workspace string, outputfile string) error {

	connectionresponse := make(map[string]interface{})
	var smcr ServiceManagerConnectionResponse

	rdshost, err := getRdsDbHostAddress(region, dbidentifier)
	if err != nil {
		return err
	}
	connectionstring := masteruser + ":" + masterpwd + "@tcp(" + rdshost + ":3306)/"
	db, err := sql.Open("mysql", connectionstring)
	logErr(err)

	// returns error if teh workspace doesnot exist
	_, dberr := db.Exec("USE " + workspace + ";")
	logErr(dberr)
	if dberr == nil {
		connectionresponse["workspace"] = workspace
		connectionresponse["port"] = 3306
		connectionresponse["host"] = rdshost
		smcr = ServiceManagerConnectionResponse{Details: connectionresponse, Status: SUCCESS}
		marshalAndWriteToOutputFile(outputfile, smcr)
	} else {
		marshalAndWriteToOutputFile(outputfile, ServiceManagerConnectionResponse{Status: FAIL})
	}
	db.Close()
	return nil
}

//delete workspace in the database amazon-rds-mysql
func deleteWorkspace(region string, dbidentifier string, masteruser string, masterpwd string, workspace string, outputfile string) error {

	rdshost, err := getRdsDbHostAddress(region, dbidentifier)
	if err != nil {
		return err
	}
	connectionstring := masteruser + ":" + masterpwd + "@tcp(" + rdshost + ":3306)/"
	db, err := sql.Open("mysql", connectionstring)
	logErr(err)

	// delete workspace
	_, dberr := db.Exec("DROP DATABASE " + workspace + ";")
	logErr(dberr)
	if dberr != nil {
		marshalAndWriteToOutputFile(outputfile, ServiceManagerConnectionResponse{Status: FAIL})
		return dberr
	}
	marshalAndWriteToOutputFile(outputfile, ServiceManagerConnectionResponse{Status: SUCCESS})
	db.Close()
	return nil
}

//create connection in the database amazon-rds-mysql
func createConnection(region string, dbidentifier string, masteruser string, masterpwd string, workspace string, user string, outputfile string) error {

	connectionresponse := make(map[string]interface{})
	var smcr ServiceManagerConnectionResponse

	rdshost, err := getRdsDbHostAddress(region, dbidentifier)
	if err != nil {
		return err
	}
	connectionstring := masteruser + ":" + masterpwd + "@tcp(" + rdshost + ":3306)/"
	db, err := sql.Open("mysql", connectionstring)
	logErr(err)

	// check if the database exists
	_, dberr := db.Exec("USE " + workspace + ";")
	logErr(err)
	if dberr == nil {
		// create random password
		randonID := randStringBytesRmndr(10)
		pwd := "pwd_" + randonID
		_, uerr := db.Exec("CREATE USER " + user + " IDENTIFIED BY '" + pwd + "';")
		logErr(uerr)
		if uerr == nil {
			_, perr := db.Exec("GRANT ALL PRIVILEGES ON " + workspace + ".* TO " + user)
			if perr == nil {
				connectionresponse["username"] = user
				connectionresponse["userpassword"] = pwd
				connectionresponse["port"] = 3306
				connectionresponse["host"] = rdshost
				smcr = ServiceManagerConnectionResponse{Details: connectionresponse, Status: SUCCESS}
				marshalAndWriteToOutputFile(outputfile, smcr)
			} else {
				marshalAndWriteToOutputFile(outputfile, ServiceManagerConnectionResponse{Status: FAIL})
				return perr
			}
		} else {
			marshalAndWriteToOutputFile(outputfile, ServiceManagerConnectionResponse{Status: FAIL})
			return uerr
		}
	} else {
		marshalAndWriteToOutputFile(outputfile, ServiceManagerConnectionResponse{Status: FAIL})
		return dberr
	}
	db.Close()
	return nil
}

//get connection in the database amazon-rds-mysql
func getConnection(region string, dbidentifier string, masteruser string, masterpwd string, workspace string, user string, outputfile string) error {

	connectionresponse := make(map[string]interface{})
	var smcr ServiceManagerConnectionResponse

	rdshost, err := getRdsDbHostAddress(region, dbidentifier)
	if err != nil {
		return err
	}
	connectionstring := masteruser + ":" + masterpwd + "@tcp(" + rdshost + ":3306)/"
	db, err := sql.Open("mysql", connectionstring)
	logErr(err)

	// check if the database exists
	_, dberr := db.Exec("USE " + workspace + ";")
	logErr(dberr)
	if dberr == nil {
		// check user
		_, uerr := db.Exec("SHOW GRANTS FOR " + user + ";")
		logErr(uerr)
		if uerr == nil {
			connectionresponse["username"] = user
			connectionresponse["port"] = 3306
			connectionresponse["host"] = rdshost
			smcr = ServiceManagerConnectionResponse{Details: connectionresponse, Status: SUCCESS}
			marshalAndWriteToOutputFile(outputfile, smcr)
		} else {
			marshalAndWriteToOutputFile(outputfile, ServiceManagerConnectionResponse{Status: FAIL})
			return uerr
		}
	} else {
		marshalAndWriteToOutputFile(outputfile, ServiceManagerConnectionResponse{Status: FAIL})
		return dberr
	}
	db.Close()
	return nil
}

//delete connection in the database amazon-rds-mysql
func deleteConnection(region string, dbidentifier string, masteruser string, masterpwd string, workspace string, user string, outputfile string) error {

	rdshost, err := getRdsDbHostAddress(region, dbidentifier)
	if err != nil {
		return err
	}
	connectionstring := masteruser + ":" + masterpwd + "@tcp(" + rdshost + ":3306)/"
	db, err := sql.Open("mysql", connectionstring)
	logErr(err)

	// delete the user
	_, dberr := db.Exec("DROP USER " + user + ";")
	logErr(dberr)
	if dberr == nil {
		// check if user is deleted
		_, uerr := db.Exec("SHOW GRANTS FOR " + user + ";")
		logErr(uerr)
		if uerr != nil {
			marshalAndWriteToOutputFile(outputfile, ServiceManagerConnectionResponse{Status: SUCCESS})
		} else {
			marshalAndWriteToOutputFile(outputfile, ServiceManagerConnectionResponse{Status: FAIL})
			return uerr
		}
	} else {
		marshalAndWriteToOutputFile(outputfile, ServiceManagerConnectionResponse{Status: FAIL})
		return dberr
	}
	db.Close()
	return nil
}

// Delete the rds database instance given the database identifier
func deleteRdsDbInstance(region string, dbidentifier string, outputfile string) error {

	secgroupname, sgerr := getRdsDbSecGroupID(region, dbidentifier)
	logErr(sgerr)

	svc := rds.New(session.New(), aws.NewConfig().WithRegion(region))
	param := &rds.DeleteDBInstanceInput{
		DBInstanceIdentifier: aws.String(dbidentifier), // Required
		SkipFinalSnapshot:    aws.Bool(true),
	}
	_, err := svc.DeleteDBInstance(param)

	if err != nil {
		marshalAndWriteToOutputFile(outputfile, ServiceManagerConnectionResponse{Status: FAIL})
		return err
	}
	marshalAndWriteToOutputFile(outputfile, ServiceManagerConnectionResponse{Status: SUCCESS})

	//wait until the db is deleted
	dbparams := &rds.DescribeDBInstancesInput{
		DBInstanceIdentifier: aws.String(dbidentifier),
	}
	dberr := svc.WaitUntilDBInstanceDeleted(dbparams)
	logErr(dberr)

	dsgerr := deleteSecurityGroup(region, secgroupname, 5, errors.New(""))
	if dsgerr != nil {
		return dsgerr
	}
	return nil
}

//delete security group
func deleteSecurityGroup(region string, secgroupname string, retry int, err error) error {

	if err == nil {
		return nil
	}
	if retry == 0 {
		errmsg := "Cannot delete SecurityGroup: " + secgroupname
		return errors.New(errmsg)
	}
	param := &ec2.DeleteSecurityGroupInput{
		GroupId: aws.String(secgroupname),
	}
	svc := ec2.New(session.New(), aws.NewConfig().WithRegion(region))
	_, updatedErr := svc.DeleteSecurityGroup(param)
	return deleteSecurityGroup(region, secgroupname, retry-1, updatedErr)
}

//helper method to generating a random string
func randStringBytesRmndr(n int) string {
	uniqueInt := time.Now().Unix()
	uniqueStr := strconv.FormatInt(uniqueInt, 10)
	strWithoutdash := strings.Replace(uniqueStr, "-", "", -1)
	validStr := strings.Replace(strWithoutdash, " ", "", -1)
	b := make([]byte, n)
	for i := range b {
		b[i] = validStr[rand.Intn(len(validStr))]
	}
	return string(b)
}

//helper method to check error
func logErr(err error) {
	if err != nil {
		fmt.Println("oops something went wrong!!")
		fmt.Println(err)
	}
}

//helper method to marshal and write to the output file
func marshalAndWriteToOutputFile(outputfile string, smcr ServiceManagerConnectionResponse) {
	msgBytes, err := json.Marshal(smcr)
	logErr(err)
	if err == nil {
		writeerr := ioutil.WriteFile(outputfile, msgBytes, 0666)
		logErr(writeerr)
	}
}

func main() {

	if len(os.Args) < 5 {
		fmt.Println("Please enter valid command.")
	} else {
		command := os.Args[1]
		region := os.Args[2]
		dbidentifier := os.Args[3]
		if command == "createdb" {
			if len(os.Args) != 10 {
				fmt.Println("Please enter valid options for createdb: region, dbclass, dbidentifier, size, multiAZ, user, password and outputfile")
			} else {
				dbclass := os.Args[4]
				size := os.Args[5]
				multiAZ := os.Args[6]
				masteruser := os.Args[7]
				masterpwd := os.Args[8]
				outputfile := os.Args[9]
				size64, serr := strconv.ParseInt(size, 10, 64)
				if serr != nil {
					fmt.Println("Size should be of type int64")
				}
				multiAZbool, merr := strconv.ParseBool(multiAZ)
				if merr != nil {
					fmt.Println("multiAZ accepts 1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False")
				}
				createRdsDbInstance(region, dbidentifier, dbclass, "mysql", size64, multiAZbool, masteruser, masterpwd, outputfile)
			}
		} else if command == "createworkspace" || command == "getworkspace" || command == "deleteworkspace" {
			if len(os.Args) != 8 {
				fmt.Println("Please enter valid options for createworkspace: region, dbidentifier, user, password, workspaceID and outputfile")
			} else {
				masteruser := os.Args[4]
				masterpwd := os.Args[5]
				workspace := os.Args[6]
				outputfile := os.Args[7]
				if command == "createworkspace" {
					createWorkspace(region, dbidentifier, masteruser, masterpwd, workspace, outputfile)
				} else if command == "deleteworkspace" {
					deleteWorkspace(region, dbidentifier, masteruser, masterpwd, workspace, outputfile)
				} else {
					getWorkspace(region, dbidentifier, masteruser, masterpwd, workspace, outputfile)
				}

			}
		} else if command == "createconnection" || command == "getconnection" || command == "deleteconnection" {
			if len(os.Args) != 9 {
				fmt.Println("Please enter valid options for createconnection: region, workspaceID, user, password, dbidentifier and outputfile")
			} else {
				masteruser := os.Args[4]
				masterpwd := os.Args[5]
				workspace := os.Args[6]
				user := os.Args[7]
				outputfile := os.Args[8]
				if command == "createconnection" {
					createConnection(region, dbidentifier, masteruser, masterpwd, workspace, user, outputfile)
				} else if command == "deleteconnection" {
					deleteConnection(region, dbidentifier, masteruser, masterpwd, workspace, user, outputfile)
				} else {
					getConnection(region, dbidentifier, masteruser, masterpwd, workspace, user, outputfile)
				}

			}
		} else if command == "deletedb" {
			if len(os.Args) != 5 {
				fmt.Println("Please enter valid options for deletedb: region, dbID and outputfile")
			} else {
				outputfile := os.Args[4]
				deleteRdsDbInstance(region, dbidentifier, outputfile)
			}
		} else {
			fmt.Println("Please enter valid command")
		}
	}

}
