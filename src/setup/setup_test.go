package setup

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
	DEFAULT_STARTUP_EXTENSION  = "/catalog-service-manager/setup/startup"
	DEFAULT_SHUTDOWN_EXTENSION = "/catalog-service-manager/setup/shutdown"
	FAKE_STARTUP_EXTENSION     = "/tmp/fake/setup/startup/startup.sh"
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

func setup(cmsFileHelper utils.CSMFileHelperInterface) (*common.ServiceManagerConfiguration, *CSMSetup) {
	config := common.NewServiceManagerConfiguration()
	logger := common.NewLogger(strings.ToLower(*config.LOG_LEVEL))
	if cmsFileHelper == nil {
		cmsFileHelper = utils.CSMFileHelper{
			Logger: logger,
		}
	}

	CSMSetup := NewCSMSetup(logger, config, cmsFileHelper)
	return config, CSMSetup
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

func Test_Startup_NoExtension(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_STARTUP_EXTENSION).Return(false, nil)
	_, csmSetup := setup(csmMockedFileExtension)
	setupStatus := csmSetup.Startup()
	assert.Equal(t, setupStatus, false)

}

func Test_Startup_NullExtension(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_STARTUP_EXTENSION).Return(true, nil)
	_, csmSetup := setup(csmMockedFileExtension)
	setupStatus := csmSetup.Startup()
	assert.Equal(t, setupStatus, false)

}

func Test_Startup_FailedToRunExtension(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_STARTUP_EXTENSION).Return(true, FAKE_STARTUP_EXTENSION)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_STARTUP_EXTENSION, []string{""}).Return(false, nil)
	_, csmSetup := setup(csmMockedFileExtension)
	setupStatus := csmSetup.Startup()
	assert.Equal(t, setupStatus, false)
}

func Test_Startup_RunExtensionSuccessful(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_STARTUP_EXTENSION).Return(true, FAKE_STARTUP_EXTENSION)
	status := "successful"
	statusString := getStatusString(&status, nil, nil)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_STARTUP_EXTENSION, []string{""}).Return(true, statusString)
	_, csmSetup := setup(csmMockedFileExtension)
	setupStatus := csmSetup.Startup()
	assert.Equal(t, setupStatus, false)
}

func Test_Startup_RunExtensionFailed(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_STARTUP_EXTENSION).Return(true, FAKE_STARTUP_EXTENSION)
	status := "failed"
	statusString := getStatusString(&status, nil, nil)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_STARTUP_EXTENSION, []string{""}).Return(true, statusString)
	_, csmSetup := setup(csmMockedFileExtension)
	setupStatus := csmSetup.Startup()
	assert.Equal(t, setupStatus, false)
}

func Test_Startup_RunExtensionIncorrectJsonOutput(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_STARTUP_EXTENSION).Return(true, FAKE_STARTUP_EXTENSION)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_STARTUP_EXTENSION, []string{""}).Return(true, "Incorrect")
	_, csmSetup := setup(csmMockedFileExtension)
	setupStatus := csmSetup.Startup()
	assert.Equal(t, setupStatus, false)
}

func Test_Startup_RunExtensionEmptyOuput(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_STARTUP_EXTENSION).Return(true, FAKE_STARTUP_EXTENSION)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_STARTUP_EXTENSION, []string{""}).Return(true, "")
	_, csmSetup := setup(csmMockedFileExtension)
	setupStatus := csmSetup.Startup()
	assert.Equal(t, setupStatus, false)
}

func Test_Startup_RunExtensionDeletedOuputFile(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_STARTUP_EXTENSION).Return(true, FAKE_STARTUP_EXTENSION)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_STARTUP_EXTENSION, []string{""}).Return(true, "DeletedOutputFile")
	_, csmSetup := setup(csmMockedFileExtension)
	setupStatus := csmSetup.Startup()
	assert.Equal(t, setupStatus, false)
}

func Test_Startup_RunExtensionUnAccessibleFile(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_STARTUP_EXTENSION).Return(true, FAKE_STARTUP_EXTENSION)
	status := "successful"
	statusString := getStatusString(&status, nil, nil)
	csmMockedFileExtension.On("RunExtensionFileGen", FAKE_STARTUP_EXTENSION, []string{""}).Return(true, "UnaccessibleOuputFile", statusString)
	_, csmSetup := setup(csmMockedFileExtension)
	setupStatus := csmSetup.Startup()
	assert.Equal(t, setupStatus, false)
}

func TestCheck_Startup(t *testing.T) {
	_, csmSetup := setup(nil)
	setupStatus := csmSetup.Startup()
	assert.Equal(t, setupStatus, false)
}

func TestCheck_Shutdown(t *testing.T) {
	_, csmSetup := setup(nil)
	setupStatus := csmSetup.Shutdown()
	assert.Equal(t, setupStatus, false)
}
func TestCheck_CheckExtensions(t *testing.T) {
	csmMockedFileExtension := MockedFileExtension{}
	csmMockedFileExtension.On("GetExtension", DEFAULT_STARTUP_EXTENSION).Return(true, nil).Once()
	csmMockedFileExtension.On("GetExtension", DEFAULT_SHUTDOWN_EXTENSION).Return(true, nil).Once()
	_, csmSetup := setup(csmMockedFileExtension)
	csmSetup.CheckExtensions()
}
