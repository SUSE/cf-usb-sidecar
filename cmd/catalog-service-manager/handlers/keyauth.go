package handlers

import (
	errors "github.com/go-swagger/go-swagger/errors"

	"github.com/hpcloud/catalog-service-manager/src/csm_manager"
)

// APIKeyAuth authorizatizes an API call
// If the environment variable CSM_API_KEY is set then this will validate
// that the token given matches the CSM_API_KEY configured, otherwise if
// the CSM_API_KEY was not set then all calls are valid.
func ApiKeyAuth(token string) (interface{}, error) {
	config := csm_manager.GetConfig()
	if *config.API_KEY == token {
		return true, nil // valid request
	}
	// invalid request
	return nil, errors.Unauthenticated("Authentication Failed for CSM")
}
