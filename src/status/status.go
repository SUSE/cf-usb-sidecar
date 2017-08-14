package status

import (
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"strconv"
	"time"

	"github.com/SUSE/cf-usb-sidecar/generated/CatalogServiceManager/models"
	"github.com/SUSE/cf-usb-sidecar/src/common"
	"github.com/SUSE/cf-usb-sidecar/src/common/utils"
	"github.com/Sirupsen/logrus"
)

type CSMStatus struct {
	common.CSMStatusInterface
	Logger     *logrus.Logger
	Config     *common.ServiceManagerConfiguration
	FileHelper utils.CSMFileHelperInterface
}

func NewCSMStatus(logger *logrus.Logger,
	config *common.ServiceManagerConfiguration,
	fileHelper utils.CSMFileHelperInterface) *CSMStatus {
	return &CSMStatus{Logger: logger, Config: config, FileHelper: fileHelper}
}

func (w *CSMStatus) GetStatus() (*models.StatusResponse, *models.Error) {
	exists, filename := w.getStatusExtension(*w.Config.MANAGER_HOME)

	if !exists || filename == "" {
		w.Logger.WithFields(logrus.Fields{utils.ERR_EXTENSION_NOT_FOUND: exists}).Info("GetStatus")

		return w.statusExtentionNotFound(utils.ERR_EXTENSION_NOT_FOUND)
	}
	return w.executeRequest("GetStatus", filename)
}

func (w *CSMStatus) statusExtentionNotFound(message string) (*models.StatusResponse, *models.Error) {
	status := utils.NewStatus()

	if *w.Config.HEALTHCHECK_HOST == "" || *w.Config.HEALTHCHECK_PORT == "" {
		status.Status = common.PROCESSING_STATUS_NONE
		status.Message = ""
	} else {
		sTimeout := *w.Config.EXT_TIMEOUT
		status.ProcessingType = common.PROCESSING_TYPE_DEFAULT

		var timeout time.Duration
		itimeout, err := strconv.Atoi(sTimeout)
		if err != nil {
			timeout = time.Duration(30) * time.Second //default 30 secs
		} else {
			timeout = time.Duration(itimeout) * time.Second
		}

		_, err = net.DialTimeout("tcp", fmt.Sprintf("%s:%s", *w.Config.HEALTHCHECK_HOST, *w.Config.HEALTHCHECK_PORT), timeout)
		if err != nil {
			status.Message = err.Error()
			status.Status = common.PROCESSING_STATUS_FAILED
		} else {
			status.Status = common.PROCESSING_STATUS_SUCCESSFUL
			status.Message = ""
		}
	}
	return &status, nil
}
func (w *CSMStatus) getStatusExtension(homePath string) (bool, string) {
	return w.FileHelper.GetExtension(filepath.Join(homePath, "status"))
}

func (w *CSMStatus) checkParamsOk(extensionPath string) error {
	if extensionPath == "" {
		err := errors.New("extensionPath is nil")
		return err
	}
	return nil
}

func (w *CSMStatus) executeRequest(requestType string, filename string) (*models.StatusResponse, *models.Error) {
	status, err := w.executeExtension(filename)

	if err != nil {
		w.Logger.Error(requestType, err)
		return nil, utils.GenerateErrorResponse(&utils.HTTP_500, err.Error())
	}
	return status, nil
}

func (w *CSMStatus) executeExtension(extensionPath string) (*models.StatusResponse, error) {
	if err := w.checkParamsOk(extensionPath); err != nil {
		return nil, err
	}
	w.Logger.WithFields(logrus.Fields{"extension Path": extensionPath}).Info("executeExtension")

	if success, outputFile, output := w.FileHelper.RunExtensionFileGen(extensionPath); success {
		w.Logger.WithFields(logrus.Fields{"extension execution status": strconv.FormatBool(success)}).Info("executeExtension")
		var stringOutput string
		if output == "" {
			stringOutput = ""
		} else {
			stringOutput = output
		}
		w.Logger.WithFields(logrus.Fields{"extension execution Result": stringOutput}).Debug("executeExtension")

		fileContent, err := utils.ReadOutputFile(outputFile, *w.Config.DEV_MODE != "true")
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error reading response from extension: %s", err.Error()))
		}
		return marshalResponseFromMessage(fileContent)

	} else {
		// extension couldn't be executed, returned an error or timedout
		//first we check for timeout (success=false,  output==nil)
		if output == "" {
			return nil, errors.New(utils.ERR_TIMEOUT)
		}
		//else it means that the extension did not return a zero code	 ("success = false, output != nil)
		err := errors.New(output)
		return nil, err
	}
}

//create PingResponse from the json we received in file
func marshalResponseFromMessage(message []byte) (*models.StatusResponse, error) {
	status := utils.NewStatus()
	jsonresp := utils.JsonResponse{}
	if len(message) == 0 {
		return nil, errors.New("Empty response")
	}
	err := jsonresp.Unmarshal(message)
	if err != nil {
		return nil, err
	}

	status.Status = jsonresp.Status
	status.Message = jsonresp.ErrorMessage
	status.ProcessingType = common.PROCESSING_TYPE_EXTENSION
	status.ServiceType = &jsonresp.ServiceType
	status.Diagnostics = jsonresp.Diagnostics

	return &status, nil
}
