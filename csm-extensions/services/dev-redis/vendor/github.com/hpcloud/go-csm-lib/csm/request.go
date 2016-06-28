package csm

import (
	"errors"
)

type CSMRequest struct {
	WorkspaceID  string
	ConnectionID string
	OutputPath   string
}

//This assumes that that the args are passed in a specific order as follows:
//1. the filepath of the output file
//2. the workspace ID
//3. the connection ID if present
func GetCSMRequest(args []string) (*CSMRequest, error) {
	if len(args) > 4 {
		return nil, errors.New("Invalid number of arguments")
	}

	request := CSMRequest{}
	request.OutputPath = args[1]
	if len(args) >= 3 {
		request.WorkspaceID = args[2]
	}
	if len(args) == 4 {
		request.ConnectionID = args[3]
	}

	return &request, nil
}
