package status

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager/models"
	"github.com/hpcloud/catalog-service-manager/src/common"
	"github.com/hpcloud/catalog-service-manager/src/common/utils"
	"github.com/pivotal-golang/lager"
)

type CSMStatus struct {
	common.CSMStatusInterface
	Logger     lager.Logger
	Config     *common.ServiceManagerConfiguration
	FileHelper utils.CSMFileHelperInterface
}

func NewCSMStatus(logger lager.Logger,
	config *common.ServiceManagerConfiguration,
	fileHelper utils.CSMFileHelperInterface) *CSMStatus {
	return &CSMStatus{Logger: logger.Session("CSM-Status"), Config: config, FileHelper: fileHelper}
}

func (w *CSMStatus) GetStatus() *models.StatusResponse {
	exists, filename := w.getStatusExtension(*w.Config.MANAGER_HOME)

	if !exists || filename == nil {
		w.Logger.Info("GetStatus", lager.Data{utils.ERR_EXTENSION_NOT_FOUND: exists})

		return w.statusExtentionNotFound(utils.ERR_EXTENSION_NOT_FOUND)
	}
	return w.executeRequest("GetStatus", filename)
}

func (w *CSMStatus) statusExtentionNotFound(message string) *models.StatusResponse {
	var status models.StatusResponse

	status.ProcessingType = &common.PROCESSING_STATUS_NONE

	host := os.Getenv("HEALTHCHECK_HOST")
	port := os.Getenv("HEALTHCHECK_PORT")

	if host == "" || port == "" {
		status.Status = "successful"
		status.Message = ""
	} else {
		_, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
		if err != nil {
			status.Message = err.Error()
			status.Status = "failed"
		}
	}
	return &status
}
func (w *CSMStatus) getStatusExtension(homePath string) (bool, *string) {
	return w.FileHelper.GetExtension(filepath.Join(homePath, "status"))
}

func (w *CSMStatus) checkParamsOk(extensionPath *string) error {
	if extensionPath == nil {
		err := errors.New("extensionPath is nil")
		return err
	}
	return nil
}

func (w *CSMStatus) executeRequest(requestType string, filename *string) *models.StatusResponse {
	var status *models.StatusResponse
	var err error

	status, err = w.executeExtension(filename)

	if err != nil {
		w.Logger.Error(requestType, err)
		status.Status = "false"
		status.Message = err.Error()
	}
	return status
}

func (w *CSMStatus) executeExtension(extensionPath *string) (*models.StatusResponse, error) {
	if err := w.checkParamsOk(extensionPath); err != nil {
		return nil, err
	}
	w.Logger.Info("executeExtension", lager.Data{"extension Path ": extensionPath})

	if success, outputFile, output := w.FileHelper.RunExtensionFileGen(*extensionPath); success {
		w.Logger.Info("executeExtension", lager.Data{"extension execution status ": success})
		w.Logger.Debug("executeExtension", lager.Data{"extension execution Result: ": output})

		fileContent, err := utils.ReadOutputFile(outputFile, *w.Config.DEV_MODE != "true")
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error reading response from extension: %s", err.Error()))
		}
		return marshalResponseFromMessage(fileContent)

	} else {
		// extension couldn't be executed, returned an error or timedout
		//first we check for timeout (success=false,  output==nil)
		if output == nil {
			return nil, errors.New(utils.ERR_TIMEOUT)
		}
		//else it means that the extension did not return a zero code	 ("success = false, output != nil)
		err := errors.New(*output)
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
	status.Message = fmt.Sprintf("%d - %s", jsonresp.ErrorCode, jsonresp.ErrorMessage)
	status.ProcessingType = &common.PROCESSING_TYPE_EXTENSION

	return &status, nil
}
