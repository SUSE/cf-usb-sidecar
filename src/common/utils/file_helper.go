package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/SUSE/cf-usb-sidecar/src/common"
	"github.com/Sirupsen/logrus"
)

type CSMFileHelperInterface interface {
	GetExtension(extPath string) (bool, string)
	RunExtension(extensionPath string, params ...string) (bool, string)
	RunExtensionFileGen(extensionPath string, params ...string) (bool, *os.File, string)
}

type CSMFileHelper struct {
	CSMFileHelperInterface
	Logger *logrus.Logger
	Config *common.ServiceManagerConfiguration
}

//GetExtension verify if the extension file exists
func (c CSMFileHelper) GetExtension(extPath string) (bool, string) {
	c.Logger.Debug("GetExtension", "Path: "+extPath)
	if _, err := os.Stat(extPath); os.IsNotExist(err) {
		return false, ""
	}

	filename := filepath.Join(extPath, filepath.Base(extPath))
	_, err := ioutil.ReadFile(filename)
	if err != nil {
		return false, ""
	}

	return true, filename
}

//RunExtension executes the extension
func (c CSMFileHelper) RunExtension(extensionPath string, params ...string) (bool, string) {
	var (
		err    error
		cmdOut bytes.Buffer
		cmdErr bytes.Buffer
	)

	//Remove empty params
	paramsString := strings.Join(params, " ")
	params = strings.Split(strings.TrimSpace(paramsString), " ")

	cmd := fmt.Sprintf("%s %s", extensionPath, params)
	c.Logger.Debug("RunExtension", "Running command : "+cmd)

	cmdExec := exec.Command(extensionPath, params...)
	cmdExec.Stderr = &cmdErr
	cmdExec.Stdout = &cmdOut

	err = cmdExec.Start()

	if err != nil {
		return false, ""
	}

	bCommandExitOk := RunCmd(cmdExec, c)

	//the command was forced stopped
	if !bCommandExitOk {
		c.Logger.Info("RunExtension", fmt.Sprintf("Extension timeout forced stop PID : %d ", cmdExec.Process.Pid)+"Successful")
		return false, ""
	}

	exitStatus := ReadExitStatus(err, cmdExec, c)

	if exitStatus != 0 { //the command returned in error state
		resp := cmdErr.String()
		return false, resp
	}
	resp := cmdErr.String()
	return true, resp

}

//returns false if the process is forced stopped and true otherwise
func RunCmd(cmdExec *exec.Cmd, c CSMFileHelper) bool {
	var (
		timeout        time.Duration
		timeoutOnErr   time.Duration
		timerErr       *time.Timer
		bCommandExitOk bool
	)
	bCommandExitOk = true
	sTimeout := *c.Config.EXT_TIMEOUT
	sTimeoutErr := *c.Config.EXT_TIMEOUT_ERROR

	itimeout, err := strconv.Atoi(sTimeout)
	if err != nil {
		timeout = time.Duration(30) //default 30 secs
	} else {
		timeout = time.Duration(itimeout)
	}

	itimeout, err = strconv.Atoi(sTimeoutErr)
	if err != nil {
		timeoutOnErr = time.Duration(2) //default 2 secs
	} else {
		timeoutOnErr = time.Duration(itimeout)
	}

	//if the extension does not finish in the timeout seconds
	//we send it an interrupt request.
	//If that does not help either, we just kill the extension
	timer := time.AfterFunc(timeout*time.Second, func() {
		c.Logger.Info("RunExtension", fmt.Sprintf("Extension timeout PID : %d ", cmdExec.Process.Pid)+"stopping")

		err := cmdExec.Process.Signal(os.Interrupt)

		c.Logger.Info("RunExtension", fmt.Sprintf("Extension timeout PID : %d", cmdExec.Process.Pid)+" stopping gracefully")
		bCommandExitOk = false
		timerErr = time.AfterFunc(timeoutOnErr*time.Second, func() {
			c.Logger.Info("RunExtension", fmt.Sprintf("Extension timeout PID : %d", cmdExec.Process.Pid)+" stopping forcefully")

			cmdExec.Process.Kill()

		})

		if err != nil {
			//the extension could not be gracefully kiled.
			//Now we will try to kill it forcefully
			c.Logger.Info("RunExtension", fmt.Sprintf("Extension timeout stop PID : ", cmdExec.Process.Pid)+"Failed")
		}

	})

	cmdExec.Wait()

	if timerErr != nil {
		timerErr.Stop()
	}
	if timer != nil {
		bCommandExitOk = timer.Stop()
	}

	return bCommandExitOk

}

func ReadExitStatus(err error, cmdExec *exec.Cmd, c CSMFileHelper) int {
	var waitStatus syscall.WaitStatus
	exitStatus := 0

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			exitStatus = waitStatus.ExitStatus()
		}
		c.Logger.Debug("RunExtension", "command execution failed"+fmt.Sprintf("%d", exitStatus))
		return exitStatus
	} else {
		waitStatus = cmdExec.ProcessState.Sys().(syscall.WaitStatus)
		exitStatus = waitStatus.ExitStatus()

		return exitStatus
	}

	return 0
}

//RunExtensionFileGen executes the extension
func (c CSMFileHelper) RunExtensionFileGen(extensionPath string, params ...string) (bool, *os.File, string) {
	var (
		cmdOut  bytes.Buffer
		cmdErr  bytes.Buffer
		err     error
		tmpfile *os.File

		bCommandExitOk bool
	)

	tmpfile, err = ioutil.TempFile("", "example")
	outputFilePath := tmpfile.Name()

	newParams := append([]string{outputFilePath}, params...)

	//Remove empty params
	paramsString := strings.Join(newParams, " ")
	newParams = strings.Split(strings.TrimSpace(paramsString), " ")

	cmd := fmt.Sprintf("%s %s", extensionPath, newParams)
	c.Logger.Debug("RunExtension", "Running command : "+cmd)

	cmdExec := exec.Command(extensionPath, newParams...)

	cmdExec.Stdout = &cmdOut
	cmdExec.Stderr = &cmdErr

	err = cmdExec.Start()

	//if we could not start the process we consider the command returned an error code
	if err != nil {
		sExitErrorWithStatus := fmt.Sprintf("The extension process could not be started: %s", err.Error())
		return false, tmpfile, sExitErrorWithStatus
	}

	bCommandExitOk = RunCmd(cmdExec, c)

	//the command was forced stopped
	if !bCommandExitOk {
		c.Logger.Info("RunExtension", fmt.Sprintf("Extension timeout forced stop PID : %d", cmdExec.Process.Pid)+"Successful")
		return false, tmpfile, ""
	}

	exitStatus := ReadExitStatus(err, cmdExec, c)

	//if the command exited with an error status
	if exitStatus != 0 {
		c.Logger.Info("RunExtension", "Extension_Executed "+"Error state")
		sCmdErr := fmt.Sprintf("%s, %s", cmdErr.String(), cmdOut.String())
		c.Logger.Debug("RunExtension", "Output : "+sCmdErr)

		return false, tmpfile, sCmdErr

	}

	c.Logger.Info("RunExtension", "Extension_Executed"+"Successfuly")
	c.Logger.Debug("RunExtension", "Output : "+cmdOut.String())
	sCmdOut := cmdOut.String()
	return true, tmpfile, sCmdOut

}
