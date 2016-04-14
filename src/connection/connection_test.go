package connection

import (
	"fmt"
	"io/ioutil"
	"os"
	"encoding/json"
	"strings"
	"testing"

	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager/models"
	"github.com/hpcloud/catalog-service-manager/src/common"
	"github.com/hpcloud/catalog-service-manager/src/common/utils"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/mock"
)


var (
    DEFAULT_GET_CONNECTION_EXTENSION="/catalog-service-manager/connection/get"
    DEFAULT_DELETE_CONNECTION_EXTENSION="/catalog-service-manager/connection/delete"
    DEFAULT_CREATE_CONNECTION_EXTENSION="/catalog-service-manager/connection/create"
    FAKE_GET_CONNECTION_EXTENSION="/tmp/fake/connection/get/getConnection.sh"
)

type MockedFileExtension struct {
	mock.Mock
	utils.CSMFileHelperInterface
}

func (l MockedFileExtension) GetExtension(extPath string) (bool, *string) {
	args := l.Called(extPath)
	if args.Get(1) == nil {
		return args.Bool(0), nil
	} else if args.Get(1) != nil {
		arg2 := args.String(1)
		return args.Bool(0), &arg2
	}
	return false, nil
}

func (l MockedFileExtension) RunExtension(extensionPath string, params ...string) (bool, *string) {
	args := l.Called(extensionPath, params)
	if args.Get(1) == nil {
		return args.Bool(0), nil
	} else {
		arg2 := args.String(1)
		return args.Bool(0), &arg2
	}
}

func (l MockedFileExtension) RunExtensionFileGen(extensionPath string, params ...string) (bool, *os.File, *string) {

	args := l.Called(extensionPath, params)
	if args.Get(1) == nil {
		return args.Bool(0), nil, nil
	} else if args.Get(1) != nil {
		arg2 := args.String(1)
		tmpfile, _ := ioutil.TempFile("", "example")
		if arg2 == "DeletedOutputFile" {
			os.Remove(tmpfile.Name())
		} else if arg2 == "UnaccessibleOuputFile" {
			tmpfile.Write([]byte(args.String(2)))
			os.Chown(tmpfile.Name(), 0, 0)
			os.Chmod(tmpfile.Name(), 0000)
			fmt.Println(tmpfile.Name())
		} else if arg2 != "" {
			tmpfile.Write([]byte(arg2))
		}
		retString := &arg2
		return args.Bool(0), tmpfile, retString
	}
	return false, nil, nil
}

func setup(cmsFileHelper utils.CSMFileHelperInterface) (*common.ServiceManagerConfiguration, *CSMConnection) {
	config := common.NewServiceManagerConfiguration()
	logger := common.NewLogger(strings.ToLower(*config.LOG_LEVEL))
	if cmsFileHelper == nil {
		cmsFileHelper = utils.CSMFileHelper{
			Logger: logger,
		}
	}

	CSMConnection := NewCSMConnection(logger, config, cmsFileHelper)
	return config, CSMConnection
}

func getStatusString(status *string, processingType *string, details map[string]interface{}) string {
    test := models.ServiceManagerWorkspaceResponse{
        Status: *status,
    }
    if processingType != nil {
        test.ProcessingType = *processingType
    }
    if details != nil {
        test.Details = details
    }
    out, _ := json.Marshal(test)
    return string(out)
}

func Test_GetWorkspace_NoExtension(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_CONNECTION_EXTENSION).Return(false, nil)
	_, csmConnection := setup(csmMockedFileExtension)
	connection := csmConnection.GetConnection("123", "123")
	assert.Equal(t, connection.ProcessingType, "None")
}
func Test_GetWorkspace_NullExtension(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_CONNECTION_EXTENSION).Return(true, nil)
	_, csmConnection := setup(csmMockedFileExtension)
	connection := csmConnection.GetConnection("123", "123")
	assert.Equal(t, connection.ProcessingType, "None")
}

func Test_GetWorkspace_FailedToRunExtension(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_CONNECTION_EXTENSION).Return(true, FAKE_GET_CONNECTION_EXTENSION)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_CONNECTION_EXTENSION, []string{"123", "123"}).Return(false, nil)
	_, csmConnection := setup(csmMockedFileExtension)
	connection := csmConnection.GetConnection("123", "123")
	assert.Equal(t, connection.ProcessingType, "Extension")
}

func Test_GetWorkspace_RunExtensionSuccessful(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_CONNECTION_EXTENSION).Return(true, FAKE_GET_CONNECTION_EXTENSION)
	status := "successful"
	statusString := getStatusString(&status, nil, nil)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_CONNECTION_EXTENSION, []string{"123", "123"}).Return(true, statusString)
	_, csmConnection := setup(csmMockedFileExtension)
	connection := csmConnection.GetConnection("123", "123")
	assert.Equal(t, connection.ProcessingType, "Extension")
	assert.Equal(t, connection.Status, "successful")
}

func Test_GetWorkspace_RunExtensionFailed(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_CONNECTION_EXTENSION).Return(true, FAKE_GET_CONNECTION_EXTENSION)
	status := "failed"
	statusString := getStatusString(&status, nil, nil)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_CONNECTION_EXTENSION, []string{"123", "123"}).Return(true, statusString)
	_, csmConnection := setup(csmMockedFileExtension)
	connection := csmConnection.GetConnection("123", "123")
	assert.Equal(t, connection.ProcessingType, "Extension")
	assert.Equal(t, connection.Status, "failed")
}

func Test_GetWorkspace_RunExtensionIncorrectJsonOutput(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_CONNECTION_EXTENSION).Return(true, FAKE_GET_CONNECTION_EXTENSION)

	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_CONNECTION_EXTENSION, []string{"123", "123"}).Return(true, "Incorrect")
	_, csmConnection := setup(csmMockedFileExtension)
	connection := csmConnection.GetConnection("123", "123")
	assert.Equal(t, connection.ProcessingType, "Extension")
}

func Test_GetWorkspace_RunExtensionEmptyOuput(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_CONNECTION_EXTENSION).Return(true, FAKE_GET_CONNECTION_EXTENSION)

	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_CONNECTION_EXTENSION, []string{"123", "123"}).Return(true, "")
	_, csmConnection := setup(csmMockedFileExtension)
	connection := csmConnection.GetConnection("123", "123")
	assert.Equal(t, connection.ProcessingType, "Extension")
}

func Test_GetWorkspace_RunExtensionDeletedOuputFile(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_CONNECTION_EXTENSION).Return(true, FAKE_GET_CONNECTION_EXTENSION)

	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_CONNECTION_EXTENSION, []string{"123", "123"}).Return(true, "")
	_, csmConnection := setup(csmMockedFileExtension)
	connection := csmConnection.GetConnection("123", "123")
	assert.Equal(t, connection.ProcessingType, "Extension")
}

func Test_GetWorkspace_RunExtensionUnAccessibleFile(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_CONNECTION_EXTENSION).Return(true, FAKE_GET_CONNECTION_EXTENSION)
	status := "successful"
	statusString := getStatusString(&status, nil, nil)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_CONNECTION_EXTENSION, []string{"123", "123"}).Return(true, "UnaccessibleOuputFile", statusString)
	_, csmConnection := setup(csmMockedFileExtension)
	connection := csmConnection.GetConnection("123", "123")
	assert.Equal(t, connection.ProcessingType, "Extension")
}
func TestCheck_GetConnection(t *testing.T) {
	_, csmConnection := setup(nil)
	connection := csmConnection.GetConnection("123", "123")
	assert.Equal(t, connection.ProcessingType, "None")
}

func TestCheck_CreateWorkspace(t *testing.T) {
	_, csmConnection := setup(nil)

	connectionID := "123"
	connectionDetails := models.ServiceManagerConnectionCreateRequest{
		ConnectionID: connectionID,
	}
	connection := csmConnection.CreateConnection("123", &connectionDetails)
	assert.Equal(t, connection.ProcessingType, "Default")
}

func TestCheck_DeleteWorkspace(t *testing.T) {
	_, csmConnection := setup(nil)
	connection := csmConnection.DeleteConnection("123", "123")
	assert.Equal(t, connection.ProcessingType, "None")
}

func TestCheck_CheckExtensions(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_CONNECTION_EXTENSION).Return(true, nil).Once()
	csmMockedFileExtension.On("GetExtension", DEFAULT_CREATE_CONNECTION_EXTENSION).Return(true, nil).Once()
	csmMockedFileExtension.On("GetExtension", DEFAULT_DELETE_CONNECTION_EXTENSION).Return(true, nil).Once()

	_, csmSetup := setup(csmMockedFileExtension)
	csmSetup.CheckExtensions()
}
