package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/pivotal-golang/lager"
)

type CSMFileHelperInterface interface {
	GetExtension(extPath string) (bool, *string)
	RunExtension(extensionPath string, params ...string) (bool, *string)
	RunExtensionFileGen(extensionPath string, params ...string) (bool, *os.File, *string)
}

type CSMFileHelper struct {
	CSMFileHelperInterface
	Logger lager.Logger
}

//GetExtension verify if the extension file exists
func (c CSMFileHelper) GetExtension(extPath string) (bool, *string) {
	c.Logger.Debug("GetExtension", lager.Data{"Path": extPath})
	if _, err := os.Stat(extPath); os.IsNotExist(err) {
		return false, nil
	}
	files, err := ioutil.ReadDir(extPath)
	if err != nil {
		return false, nil
	}
	for _, file := range files {
		filename := extPath + "/" + file.Name()
		return true, &filename
	}
	return false, nil
}

//RunExtension executes the extension
func (c CSMFileHelper) RunExtension(extensionPath string, params ...string) (bool, *string) {
	var (
		cmdOut []byte
		err    error
	)
	cmd := fmt.Sprintf("%s %s %s", extensionPath, params)
	c.Logger.Debug("RunExtension", lager.Data{"Running command : ": cmd})

	if cmdOut, err = exec.Command(extensionPath, params...).Output(); err != nil {
		c.Logger.Debug("RunExtension", lager.Data{"command execution failed : ": err})
		c.Logger.Error("RunExtension", err)
		return false, nil
	}
	cmdOutString := string(cmdOut)
	c.Logger.Info("RunExtension", lager.Data{"Extension_Executed": "Successfuly"})
	c.Logger.Debug("RunExtension", lager.Data{"Output : ": cmdOutString})
	return true, &cmdOutString
}

//RunExtensionFileGen executes the extension
func (c CSMFileHelper) RunExtensionFileGen(extensionPath string, params ...string) (bool, *os.File, *string) {
	var (
		cmdOut  []byte
		err     error
		tmpfile *os.File
	)
	tmpfile, err = ioutil.TempFile("", "example")
	outputFilePath := tmpfile.Name()

	newParams := append([]string{outputFilePath}, params...)

	cmd := fmt.Sprintf("%s %s %s", extensionPath, newParams)
	c.Logger.Debug("RunExtension", lager.Data{"Running command : ": cmd})

	if cmdOut, err = exec.Command(extensionPath, newParams...).Output(); err != nil {
		c.Logger.Debug("RunExtension", lager.Data{"command execution failed : ": err})
		c.Logger.Error("RunExtension", err)
		return false, nil, nil
	}
	cmdOutString := string(cmdOut)
	c.Logger.Info("RunExtension", lager.Data{"Extension_Executed": "Successfuly"})
	c.Logger.Debug("RunExtension", lager.Data{"Output : ": cmdOutString})
	return true, tmpfile, &cmdOutString
}
