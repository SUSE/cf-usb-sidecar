package utils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager/models"
	"github.com/hpcloud/catalog-service-manager/src/common"
)

var HTTP_500 int64 = 500
var HTTP_408 int64 = 408

type JsonResponse struct {
	HttpCode       int         `json:"http_code"`
	Details        interface{} `json:"details"`
	Status         string      `json:"status"`
	ProcessingType string      `json:"processing_type"`
}

func IsValidJSON(s string) bool {
	var vjson map[string]interface{}
	err := json.Unmarshal([]byte(s), &vjson)
	return err == nil
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
				return fileContentB, nil
			} else {
				return nil, errors.New("File size is 0")
			}
		}
	} else {
		return nil, errors.New("No output file specified")
	}
}
