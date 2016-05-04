package csm

import "github.com/hpcloud/go-csm-lib/csm/status"

type CSMResponse struct {
	HttpCode       int         `json:"http_code"`
	Details        interface{} `json:"details"`
	Status         string      `json:"status"`
	ProcessingType string      `json:"processing_type"`
}

func NewCSMResponse(httpCode int, details interface{}, stat status.Status) CSMResponse {
	response := CSMResponse{}

	switch stat {
	case status.None:
		response.Status = "none"
	case status.Failed:
		response.Status = "failed"
	case status.Successful:
		response.Status = "successful"
	case status.Unknown:
		response.Status = "unknown"

	}
	response.HttpCode = httpCode
	response.Details = details
	response.ProcessingType = "Extension"
	return response
}
