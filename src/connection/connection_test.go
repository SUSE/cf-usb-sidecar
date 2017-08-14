package connection

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/SUSE/cf-usb-sidecar/generated/CatalogServiceManager/models"
	"github.com/SUSE/cf-usb-sidecar/src/common"
	"github.com/SUSE/cf-usb-sidecar/src/common/utils"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
)

var (
	DEFAULT_GET_CONNECTION_EXTENSION    = "/catalog-service-manager/connection/get"
	DEFAULT_DELETE_CONNECTION_EXTENSION = "/catalog-service-manager/connection/delete"
	DEFAULT_CREATE_CONNECTION_EXTENSION = "/catalog-service-manager/connection/create"
	FAKE_GET_CONNECTION_EXTENSION       = "/tmp/fake/connection/get/getConnection.sh"
	FAKE_CREATE_CONNECTION_EXTENSION    = "/tmp/fake/connection/create/createConnection.sh"
	FAKE_DELETE_CONNECTION_EXTENSION    = "/tmp/fake/connection/delete/deleteConnection.sh"
)

type MockedFileExtension struct {
	mock.Mock
	utils.CSMFileHelperInterface
}

func (l MockedFileExtension) GetExtension(extPath string) (bool, string) {
	args := l.Called(extPath)
	if args.Get(1) == nil {
		return args.Bool(0), ""
	} else if args.Get(1) != nil {
		arg2 := args.String(1)
		return args.Bool(0), arg2
	}
	return false, ""
}

func (l MockedFileExtension) RunExtension(extensionPath string, params ...string) (bool, string) {
	args := l.Called(extensionPath, params)
	if args.Get(1) == nil {
		return args.Bool(0), ""
	} else {
		arg2 := args.String(1)
		return args.Bool(0), arg2
	}
}

func (l MockedFileExtension) RunExtensionFileGen(extensionPath string, params ...string) (bool, *os.File, string) {

	args := l.Called(extensionPath, params)

	if args.Get(1) == nil {

		if len(args) >= 3 && args.Get(2) == nil {
			return args.Bool(0), nil, ""
		}
		if len(args) >= 3 && args.Get(2) != nil {
			retString := args.String(2)
			return args.Bool(0), nil, retString
		}
		return args.Bool(0), nil, ""
	} else if args.Get(1) != nil {
		arg2 := args.String(1)
		tmpfile, _ := ioutil.TempFile("", "example")
		if arg2 == "DeletedOutputFile" {
			os.Remove(tmpfile.Name())
		} else if arg2 == "UnaccessibleOuputFile" {
			if len(args) >= 3 && args.Get(2) != nil {
				tmpfile.Write([]byte(args.String(2)))
			}
			os.Chown(tmpfile.Name(), 0, 0)
			os.Chmod(tmpfile.Name(), 0000)
			return true, tmpfile, ""
		} else if arg2 != "" {
			tmpfile.Write([]byte(arg2))
		}
		if len(args) >= 3 && args.Get(2) == nil {
			return args.Bool(0), tmpfile, ""
		}
		if len(args) >= 3 && args.Get(2) != nil {
			retString := args.String(2)
			return args.Bool(0), tmpfile, retString
		} else {
			return args.Bool(0), tmpfile, ""
		}

	}
	return false, nil, ""
}

func setup(cmsFileHelper utils.CSMFileHelperInterface) (*common.ServiceManagerConfiguration, *CSMConnection) {
	os.Setenv("SIDECAR_API_KEY", "NotARealKey")
	config := common.NewServiceManagerConfiguration()
	logger := common.NewLogger(strings.ToLower(*config.LOG_LEVEL), "test-connection")
	if cmsFileHelper == nil {
		cmsFileHelper = utils.CSMFileHelper{
			Logger: logger,
		}
	}

	CSMConnection := NewCSMConnection(logger, config, cmsFileHelper)
	return config, CSMConnection
}

func getStatusString(status *string, processingType *string, details map[string]interface{}) string {
	test := utils.JsonResponse{
		Status: *status,
	}
	if details != nil {
		test.Details = details
	}
	out, _ := json.Marshal(test)
	return string(out)
}

func Test_GetConnection_NoExtension(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_CONNECTION_EXTENSION).Return(false, nil, nil)
	_, csmConnection := setup(csmMockedFileExtension)
	connection, modelserr := csmConnection.GetConnection("123", "123")
	assert.Nil(t, modelserr)
	assert.Equal(t, connection.ProcessingType, "none")
	assert.Equal(t, connection.Status, "none")

}

func Test_GetConnection_NullExtension(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_CONNECTION_EXTENSION).Return(false, "", "extension not found")
	_, csmConnection := setup(csmMockedFileExtension)
	connection, modelserr := csmConnection.GetConnection("123", "123")
	assert.Nil(t, modelserr)
	assert.Equal(t, connection.ProcessingType, "none")
	assert.Equal(t, connection.Status, "none")
}

func Test_GetConnection_FailedToRunExtension(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_CONNECTION_EXTENSION).Return(true, FAKE_GET_CONNECTION_EXTENSION)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_CONNECTION_EXTENSION, []string{"123", "123", "{}"}).Return(false, nil, nil)
	_, csmConnection := setup(csmMockedFileExtension)
	connection, modelserr := csmConnection.GetConnection("123", "123")
	assert.Nil(t, connection)
	assert.Equal(t, modelserr.Message, utils.ERR_TIMEOUT)
	assert.Equal(t, modelserr.Code, &utils.HTTP_408)
}

func Test_GetConnection_RunExtensionSuccessful(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_CONNECTION_EXTENSION).Return(true, FAKE_GET_CONNECTION_EXTENSION)
	status := "successful"
	processingType := "Extension"
	statusString := getStatusString(&status, &processingType, nil)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_CONNECTION_EXTENSION, []string{"123", "123", "{}"}).Return(true, statusString, nil)
	_, csmConnection := setup(csmMockedFileExtension)
	connection, modelserr := csmConnection.GetConnection("123", "123")
	assert.Nil(t, modelserr)
	assert.Equal(t, "Extension", connection.ProcessingType)
	assert.Equal(t, "successful", connection.Status)
}

func Test_GetConnection_RunExtensionFailed(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_CONNECTION_EXTENSION).Return(true, FAKE_GET_CONNECTION_EXTENSION)
	status := "failed"
	processingType := "Extension"
	statusString := getStatusString(&status, &processingType, nil)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_CONNECTION_EXTENSION, []string{"123", "123", "{}"}).Return(false, statusString, "An error")
	_, csmConnection := setup(csmMockedFileExtension)
	connection, modelserr := csmConnection.GetConnection("123", "123")
	assert.Equal(t, &utils.HTTP_500, modelserr.Code)
	assert.Equal(t, "An error", modelserr.Message)
	assert.Nil(t, connection)
}

func Test_GetConnection_RunExtensionIncorrectJsonOutput(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_CONNECTION_EXTENSION).Return(true, FAKE_GET_CONNECTION_EXTENSION)

	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_CONNECTION_EXTENSION, []string{"123", "123", "{}"}).Return(true, "Incorrect", nil)
	_, csmConnection := setup(csmMockedFileExtension)
	connection, modelserr := csmConnection.GetConnection("123", "123")
	assert.Nil(t, connection)
	assert.Equal(t, modelserr.Code, &utils.HTTP_500)
	assert.Contains(t, modelserr.Message, "Invalid json response from extension")
}

func Test_GetConnection_RunExtensionEmptyOuput(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_CONNECTION_EXTENSION).Return(true, FAKE_GET_CONNECTION_EXTENSION)

	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_CONNECTION_EXTENSION, []string{"123", "123", "{}"}).Return(true, " ", nil)
	_, csmConnection := setup(csmMockedFileExtension)
	connection, modelserr := csmConnection.GetConnection("123", "123")
	assert.Nil(t, connection)
	assert.Equal(t, modelserr.Code, &utils.HTTP_500)
	assert.Contains(t, modelserr.Message, "Invalid json response from extension")
}

func Test_GetConnection_RunExtensionDeletedOuputFile(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_CONNECTION_EXTENSION).Return(true, FAKE_GET_CONNECTION_EXTENSION)

	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_CONNECTION_EXTENSION, []string{"123", "123", "{}"}).Return(true, "DeletedOutputFile", nil)
	_, csmConnection := setup(csmMockedFileExtension)
	connection, modelserr := csmConnection.GetConnection("123", "123")
	assert.Nil(t, connection)
	assert.Equal(t, modelserr.Code, &utils.HTTP_500)
	assert.Contains(t, modelserr.Message, "Error reading response from extension:")
}

func Test_GetConnection_RunExtensionUnAccessibleFile(t *testing.T) {
	t.Skip()
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_CONNECTION_EXTENSION).Return(true, FAKE_GET_CONNECTION_EXTENSION)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_CONNECTION_EXTENSION, []string{"123", "123", "{}"}).Return(true, "UnaccessibleOuputFile", nil)
	_, csmConnection := setup(csmMockedFileExtension)
	connection, modelserr := csmConnection.GetConnection("123", "123")
	assert.Nil(t, connection)
	assert.Equal(t, modelserr.Code, &utils.HTTP_500)

}

func TestCheck_GetConnection(t *testing.T) {
	_, csmConnection := setup(nil)
	connection, modelserr := csmConnection.GetConnection("123", "123")
	assert.Nil(t, modelserr)
	assert.Equal(t, connection.ProcessingType, "none")
	assert.Equal(t, connection.Status, "none")
}

func TestCheck_CreateConnection(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_CREATE_CONNECTION_EXTENSION).Return(true, FAKE_CREATE_CONNECTION_EXTENSION)
	status := "successful"
	processingType := "Extension"
	statusString := getStatusString(&status, &processingType, nil)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_CREATE_CONNECTION_EXTENSION, []string{"123", "123", "{\"custom-api-request-token\":\"Custom-key\"}"}).Return(true, statusString, nil)
	_, csmConnection := setup(csmMockedFileExtension)

	connectionID := "123"
	details := make(map[string]interface{})
	details["custom-api-request-token"] = "Custom-key"
	connection, modelserr := csmConnection.CreateConnection("123", connectionID, details)
	assert.Equal(t, connection.ProcessingType, "Extension")
	assert.NotNil(t, connection.Details, "Details should not be nil")
	assert.Equal(t, "Custom-key", connection.Details["custom-api-request-token"])
	assert.Nil(t, modelserr)
}

func TestCheck_CreateConnectionDefault(t *testing.T) {
	os.Setenv("test-param", "test-value")
	os.Setenv("SIDECAR_PARAMETERS", "test-param")
	_, csmConnection := setup(nil)

	connectionID := "123"
	details := make(map[string]interface{})
	connectionCreate := models.ServiceManagerConnectionCreateRequest{
		ConnectionID: connectionID,
		Details:      details}
	connection, modelserr := csmConnection.CreateConnection("123", connectionCreate.ConnectionID, connectionCreate.Details)
	assert.Nil(t, modelserr)
	assert.Equal(t, connection.ProcessingType, "default")
	assert.Equal(t, connection.Status, "successful")
	assert.Equal(t, connection.Details["test-param"], "test-value")
	os.Unsetenv("SIDECAR_PARAMETERS")
	os.Unsetenv("test-param")
}

func TestCheck_CreateConnectionDefaultMultipleParameters(t *testing.T) {
	os.Setenv("test-param1", "test-value1")
	os.Setenv("test-param2", "test-value2")
	os.Setenv("SIDECAR_PARAMETERS", "test-param1 test-param2")
	_, csmConnection := setup(nil)

	connectionID := "123"
	details := make(map[string]interface{})
	connectionCreate := models.ServiceManagerConnectionCreateRequest{
		ConnectionID: connectionID,
		Details:      details}
	connection, modelserr := csmConnection.CreateConnection("123", connectionCreate.ConnectionID, connectionCreate.Details)
	assert.Nil(t, modelserr)
	assert.Equal(t, connection.ProcessingType, "default")
	assert.Equal(t, connection.Status, "successful")
	assert.Equal(t, connection.Details["test-param1"], "test-value1")
	assert.Equal(t, connection.Details["test-param2"], "test-value2")
	os.Unsetenv("SIDECAR_PARAMETERS")
	os.Unsetenv("test-param1")
	os.Unsetenv("test-param2")
}

func TestCheck_DeleteConnectionWithNone(t *testing.T) {
	_, csmConnection := setup(nil)
	connection, modelserr := csmConnection.DeleteConnection("123", "123")
	assert.Nil(t, modelserr)
	assert.Equal(t, connection.ProcessingType, "none")
	assert.Equal(t, connection.Status, "none")
}

func TestCheck_DeleteConnection(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_DELETE_CONNECTION_EXTENSION).Return(true, FAKE_DELETE_CONNECTION_EXTENSION)
	status := "successful"
	processingType := "Extension"
	statusString := getStatusString(&status, &processingType, nil)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_DELETE_CONNECTION_EXTENSION, []string{"123", "123", "{}"}).Return(true, statusString)
	_, csmConnection := setup(csmMockedFileExtension)

	// _, csmConnection := setup(nil)
	connection, modelserr := csmConnection.DeleteConnection("123", "123")
	assert.Equal(t, connection.ProcessingType, "Extension")
	assert.Nil(t, modelserr)
}

func TestCheck_CheckExtensions(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_CONNECTION_EXTENSION).Return(true, nil).Once()
	csmMockedFileExtension.On("GetExtension", DEFAULT_CREATE_CONNECTION_EXTENSION).Return(true, nil).Once()
	csmMockedFileExtension.On("GetExtension", DEFAULT_DELETE_CONNECTION_EXTENSION).Return(true, nil).Once()

	_, csmSetup := setup(csmMockedFileExtension)
	csmSetup.CheckExtensions()
}
