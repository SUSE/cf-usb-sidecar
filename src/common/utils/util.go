package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager/models"
	"github.com/hpcloud/catalog-service-manager/src/common"
)

var HTTP_500 int64 = 500
var HTTP_408 int64 = 408

const ERR_EXTENSION_NOT_FOUND string = "extension not found"
const ERR_TIMEOUT string = "Timeout while executing the extension. The extension did not respond in a reasonable amount of time."

type JsonResponse struct {
	ErrorCode    int                        `json:"error_code,omitempty"`
	ErrorMessage string                     `json:"error_message,omitempty"`
	Details      interface{}                `json:"details,omitempty"`
	Status       string                     `json:"status"`
	Diagnostics  []*models.StatusDiagnostic `json:"diagnostics,omitempty"`
}

func (j *JsonResponse) Unmarshal(value []byte) error {
	err := json.Unmarshal(value, j)
	if err != nil {
		return errors.New(fmt.Sprintf("Invalid json response from extension: %s", err.Error()))
	}
	if strings.ToLower((*j).Status) != "successful" &&
		strings.ToLower((*j).Status) != "failed" {
		return errors.New(fmt.Sprintf("Invalid status received from extension: %s", (*j).Status))
	}

	return nil
}

func NewWorkspace() models.ServiceManagerWorkspaceResponse {
	workspace := models.ServiceManagerWorkspaceResponse{
		ProcessingType: common.PROCESSING_TYPE_NONE,
	}

	return workspace
}

func NewConnection() models.ServiceManagerConnectionResponse {
	workspace := models.ServiceManagerConnectionResponse{
		ProcessingType: common.PROCESSING_TYPE_NONE,
	}

	return workspace
}

func NewStatus() models.StatusResponse {
	status := models.StatusResponse{
		ProcessingType: common.PROCESSING_TYPE_NONE,
	}

	return status
}

func GenerateErrorResponse(code *int64, message string) *models.Error {
	errReturn := models.Error{Code: code, Message: message}
	return &errReturn
}

func ReadOutputFile(outputFile *os.File, removeAfter bool) ([]byte, error) {
	if outputFile != nil {
		if removeAfter {
			defer os.Remove(outputFile.Name())
		}

		if fileStat, err := os.Stat(outputFile.Name()); err != nil {
			return nil, errors.New("No output file found")
		} else {
			if fileStat.Size() > 0 {
				fileContentB, err := ioutil.ReadFile(outputFile.Name())
				if err != nil {
					return nil, err
				}
				if len(fileContentB) == 0 {
					return nil, errors.New("The generated file is empty")
				}
				return fileContentB, nil
			} else {
				return nil, errors.New("File size is 0")
			}
		}
	} else {
		return nil, errors.New("No output file specified")
	}
}
