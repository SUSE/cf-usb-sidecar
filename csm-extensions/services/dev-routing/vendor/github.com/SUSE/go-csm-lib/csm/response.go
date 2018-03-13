package csm

type CSMResponse struct {
	ErrorCode    int                 `json:"error_code,omitempty"`
	ErrorMessage string              `json:"error_message,omitempty"`
	Details      interface{}         `json:"details,omitempty"`
	Status       string              `json:"status"`
	ServiceType  string              `json:"service_type"`
	Diagnostics  []*StatusDiagnostic `json:"diagnostics,omitempty"`
}

type StatusDiagnostic struct {
	Description string `json:"description"`
	Message     string `json:"message"`
	Name        string `json:"name"`
	Status      string `json:"status"`
}

func CreateCSMResponse(details interface{}) CSMResponse {
	response := CSMResponse{
		Status:  "successful",
		Details: details,
	}
	response.Status = "successful"
	response.Details = details
	return response
}

func CreateCSMErrorResponse(errorCode int, errorMessage string) CSMResponse {
	response := CSMResponse{
		Status:       "failed",
		ErrorCode:    errorCode,
		ErrorMessage: errorMessage,
	}
	return response
}
