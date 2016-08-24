package workspace

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
	DEFAULT_GET_WORKSPACE_EXTENSION    = "/catalog-service-manager/workspace/get"
	DEFAULT_DELETE_WORKSPACE_EXTENSION = "/catalog-service-manager/workspace/delete"
	DEFAULT_CREATE_WORKSPACE_EXTENSION = "/catalog-service-manager/workspace/create"
	FAKE_GET_WORKSPACE_EXTENSION       = "/tmp/fake/workspace/get/getWorkspace.sh"
	FAKE_CREATE_WORKSPACE_EXTENSION    = "/tmp/fake/workspace/create/createWorkspace.sh"
	FAKE_DELETE_WORKSPACE_EXTENSION    = "/tmp/fake/workspace/delete/deleteWorkspace.sh"
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
		if len(args) >= 3 {
			retString := args.String(2)
			return args.Bool(0), nil, retString
		}
		return args.Bool(0), nil, ""
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

func setup(cmsFileHelper utils.CSMFileHelperInterface) (*common.ServiceManagerConfiguration, *CSMWorkspace) {
	os.Setenv("SIDECAR_API_KEY", "NotARealKey")
	config := common.NewServiceManagerConfiguration()
	logger := common.NewLogger(strings.ToLower(*config.LOG_LEVEL), "test-workspace")
	if cmsFileHelper == nil {
		cmsFileHelper = utils.CSMFileHelper{
			Logger: logger,
		}
	}

	CSMWorkspace := NewCSMWorkspace(logger, config, cmsFileHelper)
	return config, CSMWorkspace
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

func Test_GetWorkspace_NoExtension(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_WORKSPACE_EXTENSION).Return(false, nil, nil)
	_, csmWorkspace := setup(csmMockedFileExtension)
	workspace, modelserr := csmWorkspace.GetWorkspace("123")
	assert.Nil(t, modelserr)
	assert.Equal(t, workspace.ProcessingType, "none")
	assert.Equal(t, workspace.Status, "none")
}

func Test_GetWorkspace_NullExtension(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_WORKSPACE_EXTENSION).Return(false, "", "extension not found")
	_, csmWorkspace := setup(csmMockedFileExtension)
	workspace, modelserr := csmWorkspace.GetWorkspace("123")
	assert.Nil(t, modelserr)
	assert.Equal(t, workspace.ProcessingType, "none")
	assert.Equal(t, workspace.Status, "none")
}

func Test_GetWorkspace_FailedToRunExtension(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_WORKSPACE_EXTENSION).Return(true, FAKE_GET_WORKSPACE_EXTENSION)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_WORKSPACE_EXTENSION, []string{"123", "{}"}).Return(false, nil, nil)
	_, csmWorkspace := setup(csmMockedFileExtension)
	workspace, modelserr := csmWorkspace.GetWorkspace("123")
	assert.Nil(t, workspace)
	assert.Equal(t, modelserr.Message, utils.ERR_TIMEOUT)
	assert.Equal(t, modelserr.Code, &utils.HTTP_408)
}

func Test_GetWorkspace_RunExtensionSuccessful(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_WORKSPACE_EXTENSION).Return(true, FAKE_GET_WORKSPACE_EXTENSION)
	status := "successful"
	processingType := "Extension"
	statusString := getStatusString(&status, &processingType, nil)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_WORKSPACE_EXTENSION, []string{"123", "{}"}).Return(true, statusString)
	_, csmWorkspace := setup(csmMockedFileExtension)
	workspace, modelserr := csmWorkspace.GetWorkspace("123")
	assert.Equal(t, workspace.ProcessingType, "Extension")
	assert.Equal(t, workspace.Status, "successful")
	assert.Nil(t, modelserr)
}

func Test_GetWorkspace_RunExtensionFailed(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_WORKSPACE_EXTENSION).Return(true, FAKE_GET_WORKSPACE_EXTENSION)
	status := "failed"
	processingType := "Extension"
	statusString := getStatusString(&status, &processingType, nil)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_WORKSPACE_EXTENSION, []string{"123", "{}"}).Return(false, statusString, "An Error")
	_, csmWorkspace := setup(csmMockedFileExtension)
	workspace, modelserr := csmWorkspace.GetWorkspace("123")
	assert.Nil(t, workspace)
	assert.Equal(t, &utils.HTTP_500, modelserr.Code)
	assert.Equal(t, "An Error", modelserr.Message)
}

func Test_GetWorkspace_RunExtensionIncorrectJsonOutput(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_WORKSPACE_EXTENSION).Return(true, FAKE_GET_WORKSPACE_EXTENSION)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_WORKSPACE_EXTENSION, []string{"123", "{}"}).Return(true, "Incorrect", nil)
	_, csmWorkspace := setup(csmMockedFileExtension)
	workspace, modelserr := csmWorkspace.GetWorkspace("123")
	assert.Nil(t, workspace)
	assert.Equal(t, modelserr.Code, &utils.HTTP_500)
	assert.Contains(t, modelserr.Message, "Invalid json response from extension")
}

func Test_GetWorkspace_RunExtensionEmptyOuput(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_WORKSPACE_EXTENSION).Return(true, FAKE_GET_WORKSPACE_EXTENSION)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_WORKSPACE_EXTENSION, []string{"123", "{}"}).Return(true, " ", nil)
	_, csmWorkspace := setup(csmMockedFileExtension)
	workspace, modelserr := csmWorkspace.GetWorkspace("123")
	assert.Nil(t, workspace)
	assert.Equal(t, modelserr.Code, &utils.HTTP_500)
	assert.Contains(t, modelserr.Message, "Invalid json response from extension")
}

func Test_GetWorkspace_RunExtensionDeletedOuputFile(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_WORKSPACE_EXTENSION).Return(true, FAKE_GET_WORKSPACE_EXTENSION)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_WORKSPACE_EXTENSION, []string{"123", "{}"}).Return(true, "DeletedOutputFile", nil)
	_, csmWorkspace := setup(csmMockedFileExtension)
	workspace, modelserr := csmWorkspace.GetWorkspace("123")
	assert.Nil(t, workspace)
	assert.Equal(t, modelserr.Code, &utils.HTTP_500)
	assert.Contains(t, modelserr.Message, "Error reading response from extension:")
}

func Test_GetWorkspace_RunExtensionUnAccessibleFile(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_WORKSPACE_EXTENSION).Return(true, FAKE_GET_WORKSPACE_EXTENSION)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_WORKSPACE_EXTENSION, []string{"123", "{}"}).Return(true, "UnaccessibleOuputFile", nil)
	_, csmWorkspace := setup(csmMockedFileExtension)
	workspace, modelserr := csmWorkspace.GetWorkspace("123")
	assert.Nil(t, workspace)
	assert.Equal(t, modelserr.Code, &utils.HTTP_500)
}

func TestCheck_GetWorkspace(t *testing.T) {
	_, csmWorkspace := setup(nil)
	workspace, modelserr := csmWorkspace.GetWorkspace("123")
	assert.Nil(t, modelserr)
	assert.Equal(t, workspace.ProcessingType, "none")
	assert.Equal(t, workspace.Status, "none")
}

func TestCheck_CreateWorkspace(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_CREATE_WORKSPACE_EXTENSION).Return(true, FAKE_CREATE_WORKSPACE_EXTENSION)
	status := "successful"
	processingType := "Extension"
	statusString := getStatusString(&status, &processingType, nil)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_CREATE_WORKSPACE_EXTENSION, []string{"123", "{}"}).Return(true, statusString)
	_, csmWorkspace := setup(csmMockedFileExtension)
	workspaceID := "123"
	details := make(map[string]interface{})
	workspace, modelserr := csmWorkspace.CreateWorkspace(workspaceID, details)
	assert.Equal(t, workspace.ProcessingType, "Extension")
	assert.Nil(t, modelserr)
}

func TestCheck_CreateWorkspace_WithDetails(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_CREATE_WORKSPACE_EXTENSION).Return(true, FAKE_CREATE_WORKSPACE_EXTENSION)
	status := "successful"
	processingType := "Extension"
	statusString := getStatusString(&status, &processingType, nil)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_CREATE_WORKSPACE_EXTENSION, []string{"123", "{\"custom-api-request-token\":\"Custom-key\"}"}).Return(true, statusString)
	_, csmWorkspace := setup(csmMockedFileExtension)
	workspaceID := "123"
	details := make(map[string]interface{})
	details["custom-api-request-token"] = "Custom-key"
	workspace, modelserr := csmWorkspace.CreateWorkspace(workspaceID, details)
	assert.Equal(t, workspace.ProcessingType, "Extension")
	assert.NotNil(t, workspace.Details, "Details should not be nil")
	assert.Equal(t, "Custom-key", workspace.Details["custom-api-request-token"])
	assert.Nil(t, modelserr)
}

func Test_CreateWorkspaceUnAccessibleFile(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_CREATE_WORKSPACE_EXTENSION).Return(true, FAKE_CREATE_WORKSPACE_EXTENSION)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_CREATE_WORKSPACE_EXTENSION, []string{"123", "{}"}).Return(true, "UnaccessibleOuputFile", nil)
	_, csmWorkspace := setup(csmMockedFileExtension)
	workspaceID := "123"
	details := make(map[string]interface{})
	workspace, modelserr := csmWorkspace.CreateWorkspace(workspaceID, details)
	assert.Nil(t, workspace)
	assert.Equal(t, modelserr.Code, &utils.HTTP_500)
	assert.Contains(t, modelserr.Message, "Error reading response from extension:")
}

func TestCheck_DeleteWorkspace(t *testing.T) {
	_, csmWorkspace := setup(nil)
	workspace, modelserr := csmWorkspace.DeleteWorkspace("123")
	assert.Nil(t, modelserr)
	assert.Equal(t, workspace.ProcessingType, "none")
	assert.Equal(t, workspace.Status, "none")
}

func Test_DeleteWorkspaceUnAccessibleFile(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_DELETE_WORKSPACE_EXTENSION).Return(true, FAKE_GET_WORKSPACE_EXTENSION)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_WORKSPACE_EXTENSION, []string{"123", "{}"}).Return(true, "UnaccessibleOuputFile", nil)
	_, csmWorkspace := setup(csmMockedFileExtension)
	workspace, modelserr := csmWorkspace.DeleteWorkspace("123")
	assert.Nil(t, workspace)
	assert.Equal(t, modelserr.Code, &utils.HTTP_500)
	assert.Contains(t, modelserr.Message, "Error reading response from extension:")
}

func TestCheck_CheckExtensions(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_WORKSPACE_EXTENSION).Return(true, nil).Once()
	csmMockedFileExtension.On("GetExtension", DEFAULT_CREATE_WORKSPACE_EXTENSION).Return(true, nil).Once()
	csmMockedFileExtension.On("GetExtension", DEFAULT_DELETE_WORKSPACE_EXTENSION).Return(true, nil).Once()
	_, csmSetup := setup(csmMockedFileExtension)
	csmSetup.CheckExtensions()
}
