package status

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/hpcloud/catalog-service-manager/src/common"
	"github.com/hpcloud/catalog-service-manager/src/common/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	DEFAULT_GET_STATUS_EXTENSION = "/catalog-service-manager/status"
	FAKE_GET_STATUS_EXTENSION    = "/tmp/fake/status/status"
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
	var args mock.Arguments
	if len(params) > 0 {
		args = l.Called(extensionPath, params)
	} else {
		args = l.Called(extensionPath)
	}
	if args.Get(1) == nil {
		if len(args) >= 3 && args.Get(2) == nil {
			return args.Bool(0), nil, nil
		}
		if len(args) >= 3 {
			retString := args.String(2)
			return args.Bool(0), nil, &retString
		}
		return args.Bool(0), nil, nil
	} else if args.Get(1) != nil {
		arg2 := args.String(1)
		tmpfile, _ := ioutil.TempFile("", "example1")
		if arg2 == "DeletedOutputFile" {
			os.Remove(tmpfile.Name())
		} else if arg2 == "UnaccessibleOuputFile" {
			if len(args) >= 3 && args.Get(2) != nil {
				tmpfile.Write([]byte(args.String(2)))
			}
			os.Chown(tmpfile.Name(), 0, 0)
			os.Chmod(tmpfile.Name(), 0000)
			return true, tmpfile, nil
		} else if arg2 != "" {
			tmpfile.Write([]byte(arg2))
		}
		if len(args) >= 3 && args.Get(2) == nil {
			return args.Bool(0), tmpfile, nil
		}
		if len(args) >= 3 && args.Get(2) != nil {
			retString := args.String(2)
			return args.Bool(0), tmpfile, &retString
		} else {
			return args.Bool(0), tmpfile, nil
		}

	}
	return false, nil, nil
}

func setup(cmsFileHelper utils.CSMFileHelperInterface) (*common.ServiceManagerConfiguration, *CSMStatus) {
	config := common.NewServiceManagerConfiguration()
	logger := common.NewLogger(strings.ToLower(*config.LOG_LEVEL))
	if cmsFileHelper == nil {
		cmsFileHelper = utils.CSMFileHelper{
			Logger: logger,
		}
	}

	CSMStatus := NewCSMStatus(logger, config, cmsFileHelper)
	return config, CSMStatus
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

func Test_GetStatus_NoExtension(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_STATUS_EXTENSION).Return(false, nil, nil)
	_, csmStatus := setup(csmMockedFileExtension)
	response, modelserr := csmStatus.GetStatus()
	assert.Nil(t, modelserr)
	assert.Equal(t, response.ProcessingType, "none")
	assert.Equal(t, response.Status, "none")
}

func Test_GetStatus_NullExtension(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_STATUS_EXTENSION).Return(false, "", "extension not found")
	_, csmStatus := setup(csmMockedFileExtension)
	response, modelserr := csmStatus.GetStatus()
	assert.Nil(t, modelserr)
	assert.Equal(t, response.ProcessingType, "none")
	assert.Equal(t, response.Status, "none")
}

func Test_GetStatus_RunExtensionSuccessful(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_STATUS_EXTENSION).Return(true, FAKE_GET_STATUS_EXTENSION)
	status := "successful"
	processingType := "Extension"
	statusString := getStatusString(&status, &processingType, nil)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_STATUS_EXTENSION).Return(true, statusString)
	_, csmStatus := setup(csmMockedFileExtension)
	response, modelserr := csmStatus.GetStatus()
	assert.Equal(t, response.ProcessingType, "extension")
	assert.Equal(t, response.Status, "successful")
	assert.Nil(t, modelserr)
}

func Test_GetStatus_RunExtensionFailed(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_STATUS_EXTENSION).Return(true, FAKE_GET_STATUS_EXTENSION)
	status := "failed"
	processingType := "Extension"
	statusString := getStatusString(&status, &processingType, nil)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_STATUS_EXTENSION).Return(false, statusString, "An Error")
	_, csmStatus := setup(csmMockedFileExtension)
	response, modelserr := csmStatus.GetStatus()
	assert.Nil(t, response)
	assert.Equal(t, &utils.HTTP_500, modelserr.Code)
	assert.Equal(t, "An Error", modelserr.Message)
}

func Test_GetStatus_RunExtensionIncorrectJsonOutput(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_STATUS_EXTENSION).Return(true, FAKE_GET_STATUS_EXTENSION)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_STATUS_EXTENSION).Return(true, "Incorrect", nil)
	_, csmStatus := setup(csmMockedFileExtension)
	response, modelserr := csmStatus.GetStatus()
	assert.Nil(t, response)
	assert.Equal(t, modelserr.Code, &utils.HTTP_500)
	assert.Contains(t, modelserr.Message, "Invalid json response from extension")
}

func Test_GetStatus_RunExtensionEmptyOuput(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_STATUS_EXTENSION).Return(true, FAKE_GET_STATUS_EXTENSION)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_STATUS_EXTENSION).Return(true, " ", nil)
	_, csmStatus := setup(csmMockedFileExtension)
	response, modelserr := csmStatus.GetStatus()
	assert.Nil(t, response)
	assert.Equal(t, modelserr.Code, &utils.HTTP_500)
	assert.Contains(t, modelserr.Message, "Invalid json response from extension")
}

func Test_GetStatus_RunExtensionDeletedOuputFile(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_STATUS_EXTENSION).Return(true, FAKE_GET_STATUS_EXTENSION)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_STATUS_EXTENSION).Return(true, "DeletedOutputFile", nil)
	_, csmStatus := setup(csmMockedFileExtension)
	response, modelserr := csmStatus.GetStatus()
	assert.Nil(t, response)
	assert.Equal(t, modelserr.Code, &utils.HTTP_500)
	assert.Contains(t, modelserr.Message, "Error reading response from extension:")
}

func Test_GetStatus_RunExtensionUnAccessibleFile(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_STATUS_EXTENSION).Return(true, FAKE_GET_STATUS_EXTENSION)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_STATUS_EXTENSION).Return(true, "UnaccessibleOuputFile", nil)
	_, csmStatus := setup(csmMockedFileExtension)
	response, modelserr := csmStatus.GetStatus()
	assert.Nil(t, response)
	assert.Equal(t, modelserr.Code, &utils.HTTP_500)
}

func TestCheck_GetStatus(t *testing.T) {
	_, csmStatus := setup(nil)
	response, modelserr := csmStatus.GetStatus()
	assert.Nil(t, modelserr)
	assert.Equal(t, response.ProcessingType, "none")
	assert.Equal(t, response.Status, "none")
}
