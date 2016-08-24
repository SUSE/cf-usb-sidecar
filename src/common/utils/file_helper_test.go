package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/hpcloud/catalog-service-manager/src/common"
)

func GetFileToRunPath(filename string) *string {
	env, _ := os.Getwd()
	if env == "" {
		return nil
	} else {
		fp := filepath.Join(env, "../../../tests/integration-tests/test_assets", filename)
		return &fp
	}
}

/////////////////////
///////WARNING///////
////this can take quite a while
////if you do not want to wait for (SIDECAR_EXT_TIMEOUT + SIDECAR_EXT_TIMEOUT_ERROR) secs, disable it
////////////////////
func TestRunExtensionShouldKillAfterTimeout(t *testing.T) {
	//t.Skip()

	t.Log("Should die after (SIDECAR_EXT_TIMEOUT + SIDECAR_EXT_TIMEOUT_ERROR) secs with false and no string output")
	os.Setenv("SIDECAR_API_KEY", "NotARealKey")
	defer os.Unsetenv("SIDECAR_API_KEY")
	config := common.NewServiceManagerConfiguration()
	logger := common.NewLogger(*config.LOG_LEVEL, "test-long-running-file")
	*config.EXT_TIMEOUT = "2"
	*config.EXT_TIMEOUT_ERROR = "2" //this will lower the time to wait to about 4 secs

	csm := CSMFileHelper{Logger: logger, Config: config}

	fileToRun := GetFileToRunPath("long_running_task.sh")

	if _, err := os.Stat(*fileToRun); os.IsNotExist(err) {
		t.Skipf("The file %s needed to run this test, does not exist. Skipping test for now.", *fileToRun)
	}

	os.Chmod(*fileToRun, 0777)

	info, _ := os.Stat(*fileToRun)

	if !strings.Contains(info.Mode().Perm().String(), "-rwx") {
		t.Skipf("The file %s is not runnable: %s", *fileToRun, info.Mode().Perm().String())
	}

	t.Logf("The file %s needed to run this test, was found. and has permissions: %s", *fileToRun, info.Mode())

	if fileToRun == nil {
		t.Error("$TEST_ASSETS not set?")
		t.Fail()
		return
	}

	param := "workspace_id"

	bOk, file, strout := csm.RunExtensionFileGen(*fileToRun, param)
	if file != nil {
		defer os.Remove(file.Name())
	}
	if bOk || strout != "" {
		t.Error("For ", *fileToRun, param, "expected", false, nil, "got", bOk, strout)
	}
}

func TestRunExtensionShouldFalse(t *testing.T) {

	param := 10
	os.Setenv("SIDECAR_API_KEY", "NotARealKey")
	defer os.Unsetenv("SIDECAR_API_KEY")
	t.Log(fmt.Sprintf("Should  return an exitStatus %d ", param))

	config := common.NewServiceManagerConfiguration()
	logger := common.NewLogger(*config.LOG_LEVEL, "file-helper-test")
	csm := CSMFileHelper{Logger: logger, Config: config}

	sParam := strconv.Itoa(param)

	fileToRun := GetFileToRunPath("non_clean_exit_task.sh")

	if fileToRun == nil {
		t.Error("$TEST_ASSETS not set?")
		t.Fail()
		return
	}

	bOk, file, strout := csm.RunExtensionFileGen(*fileToRun, sParam)
	if file != nil {
		defer os.Remove(file.Name())
	}
	if strout == "" {
		t.Error("For ", *fileToRun, param, "expected", false, "not nil", "got", bOk, strout)
		t.Fail()
		return
	}
	if bOk || strout == "" {
		t.Error("For ", *fileToRun, param, "expected", false, param, "got", bOk, strout)
	}
}

func TestRunExtensionShouldOk(t *testing.T) {
	t.Log("Everything should be ok")
	os.Setenv("SIDECAR_API_KEY", "NotARealKey")
	defer os.Unsetenv("SIDECAR_API_KEY")
	config := common.NewServiceManagerConfiguration()
	logger := common.NewLogger(*config.LOG_LEVEL, "ok-file")
	csm := CSMFileHelper{Logger: logger, Config: config}

	fileToRun := GetFileToRunPath("normal_response_task.sh")

	if fileToRun == nil {
		t.Error("$TEST_ASSETS not set?")
		t.Fail()
		return
	}

	param := "workspace_id"

	t.Log("testing for", *fileToRun, param)

	bOk, file, _ := csm.RunExtensionFileGen(*fileToRun, param, "2")

	if file == nil {
		t.Error("For ", fileToRun, param, "expected file not nul", "got", file)
		return
	}

	fileContent, _ := ioutil.ReadFile(file.Name())
	sFileContent := strings.Trim(string(fileContent), "\r\n ")

	defer os.Remove(file.Name())

	if !bOk || sFileContent != "{\"response\":\"OK\"}" {
		t.Error("For ", *fileToRun, param, "expected", true, "|{\"response\":\"OK\"}|", "got", bOk, "|"+sFileContent+"|")
	}
}
