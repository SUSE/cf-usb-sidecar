package workspace

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager/models"
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

func setup(cmsFileHelper utils.CSMFileHelperInterface) (*common.ServiceManagerConfiguration, *CSMWorkspace) {
	config := common.NewServiceManagerConfiguration()
	logger := common.NewLogger(strings.ToLower(*config.LOG_LEVEL))
	if cmsFileHelper == nil {
		cmsFileHelper = utils.CSMFileHelper{
			Logger: logger,
		}
	}

	CSMWorkspace := NewCSMWorkspace(logger, config, cmsFileHelper)
	return config, CSMWorkspace
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
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_WORKSPACE_EXTENSION).Return(false, nil)
	_, csmWorkspace := setup(csmMockedFileExtension)
	workspace := csmWorkspace.GetWorkspace("123")
	assert.Equal(t, workspace.ProcessingType, "None")
}

func Test_GetWorkspace_NullExtension(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_WORKSPACE_EXTENSION).Return(true, nil)
	_, csmWorkspace := setup(csmMockedFileExtension)
	workspace := csmWorkspace.GetWorkspace("123")
	assert.Equal(t, workspace.ProcessingType, "None")
}

func Test_GetWorkspace_FailedToRunExtension(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_WORKSPACE_EXTENSION).Return(true, FAKE_GET_WORKSPACE_EXTENSION)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_WORKSPACE_EXTENSION, []string{"123"}).Return(false, nil)
	_, csmWorkspace := setup(csmMockedFileExtension)
	workspace := csmWorkspace.GetWorkspace("123")
	assert.Equal(t, workspace.ProcessingType, "Extension")
}

func Test_GetWorkspace_RunExtensionSuccessful(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_WORKSPACE_EXTENSION).Return(true, FAKE_GET_WORKSPACE_EXTENSION)
	status := "successful"
	statusString := getStatusString(&status, nil, nil)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_WORKSPACE_EXTENSION, []string{"123"}).Return(true, statusString)
	_, csmWorkspace := setup(csmMockedFileExtension)
	workspace := csmWorkspace.GetWorkspace("123")
	assert.Equal(t, workspace.ProcessingType, "Extension")
	assert.Equal(t, workspace.Status, "successful")
}

func Test_GetWorkspace_RunExtensionFailed(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_WORKSPACE_EXTENSION).Return(true, FAKE_GET_WORKSPACE_EXTENSION)
	status := "failed"
	statusString := getStatusString(&status, nil, nil)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_WORKSPACE_EXTENSION, []string{"123"}).Return(true, statusString)
	_, csmWorkspace := setup(csmMockedFileExtension)
	workspace := csmWorkspace.GetWorkspace("123")
	assert.Equal(t, workspace.ProcessingType, "Extension")
	assert.Equal(t, workspace.Status, "failed")
}

func Test_GetWorkspace_RunExtensionIncorrectJsonOutput(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_WORKSPACE_EXTENSION).Return(true, FAKE_GET_WORKSPACE_EXTENSION)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_WORKSPACE_EXTENSION, []string{"123"}).Return(true, "Incorrect")
	_, csmWorkspace := setup(csmMockedFileExtension)
	workspace := csmWorkspace.GetWorkspace("123")
	assert.Equal(t, workspace.ProcessingType, "Extension")
	assert.Equal(t, workspace.Status, "")
	assert.Nil(t, workspace.Details)
}

func Test_GetWorkspace_RunExtensionEmptyOuput(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_WORKSPACE_EXTENSION).Return(true, FAKE_GET_WORKSPACE_EXTENSION)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_WORKSPACE_EXTENSION, []string{"123"}).Return(true, "")
	_, csmWorkspace := setup(csmMockedFileExtension)
	workspace := csmWorkspace.GetWorkspace("123")
	assert.Equal(t, workspace.ProcessingType, "Extension")
}

func Test_GetWorkspace_RunExtensionDeletedOuputFile(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_WORKSPACE_EXTENSION).Return(true, FAKE_GET_WORKSPACE_EXTENSION)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_WORKSPACE_EXTENSION, []string{"123"}).Return(true, "DeletedOutputFile")
	_, csmWorkspace := setup(csmMockedFileExtension)
	workspace := csmWorkspace.GetWorkspace("123")
	assert.Equal(t, workspace.ProcessingType, "Extension")
}

func Test_GetWorkspace_RunExtensionUnAccessibleFile(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_WORKSPACE_EXTENSION).Return(true, FAKE_GET_WORKSPACE_EXTENSION)
	status := "successful"
	statusString := getStatusString(&status, nil, nil)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_WORKSPACE_EXTENSION, []string{"123"}).Return(true, "UnaccessibleOuputFile", statusString)
	_, csmWorkspace := setup(csmMockedFileExtension)
	workspace := csmWorkspace.GetWorkspace("123")
	assert.Equal(t, workspace.ProcessingType, "Extension")
}

func TestCheck_GetWorkspace(t *testing.T) {
	_, csmWorkspace := setup(nil)
	workspace := csmWorkspace.GetWorkspace("123")
	assert.Equal(t, workspace.ProcessingType, "None")
}

func TestCheck_CreateWorkspace(t *testing.T) {
	_, csmWorkspace := setup(nil)
	workspaceID := "123"
	workspaceDetails := models.ServiceManagerWorkspaceCreateRequest{
		WorkspaceID: workspaceID,
	}
	workspace := csmWorkspace.CreateWorkspace(&workspaceDetails)
	assert.Equal(t, workspace.ProcessingType, "None")
}

func Test_CreateWorkspaceUnAccessibleFile(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_CREATE_WORKSPACE_EXTENSION).Return(true, FAKE_GET_WORKSPACE_EXTENSION)
	status := "successful"
	statusString := getStatusString(&status, nil, nil)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_WORKSPACE_EXTENSION, []string{"123"}).Return(true, "UnaccessibleOuputFile", statusString)
	_, csmWorkspace := setup(csmMockedFileExtension)
	workspaceID := "123"
	workspaceDetails := models.ServiceManagerWorkspaceCreateRequest{
		WorkspaceID: workspaceID,
	}
	workspace := csmWorkspace.CreateWorkspace(&workspaceDetails)
	assert.Equal(t, workspace.ProcessingType, "Extension")
}

func TestCheck_DeleteWorkspace(t *testing.T) {
	_, csmWorkspace := setup(nil)
	workspace := csmWorkspace.DeleteWorkspace("123")
	assert.Equal(t, workspace.ProcessingType, "None")
}

func Test_DeleteWorkspaceUnAccessibleFile(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_DELETE_WORKSPACE_EXTENSION).Return(true, FAKE_GET_WORKSPACE_EXTENSION)
	status := "successful"
	statusString := getStatusString(&status, nil, nil)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_GET_WORKSPACE_EXTENSION, []string{"123"}).Return(true, "UnaccessibleOuputFile", statusString)
	_, csmWorkspace := setup(csmMockedFileExtension)
	workspace := csmWorkspace.DeleteWorkspace("123")
	assert.Equal(t, workspace.ProcessingType, "Extension")
}

func TestCheck_CheckExtensions(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_GET_WORKSPACE_EXTENSION).Return(true, nil).Once()
	csmMockedFileExtension.On("GetExtension", DEFAULT_CREATE_WORKSPACE_EXTENSION).Return(true, nil).Once()
	csmMockedFileExtension.On("GetExtension", DEFAULT_DELETE_WORKSPACE_EXTENSION).Return(true, nil).Once()
	_, csmSetup := setup(csmMockedFileExtension)
	csmSetup.CheckExtensions()
}
