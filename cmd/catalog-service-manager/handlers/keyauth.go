package handlers

import (
	errors "github.com/go-openapi/errors"

	"github.com/SUSE/cf-usb-sidecar/src/csm_manager"
)

// APIKeyAuth authorizatizes an API call
// If the environment variable SIDECAR_API_KEY is set then this will validate
// that the token given matches the SIDECAR_API_KEY configured, otherwise if
// the SIDECAR_API_KEY was not set then all calls are valid.
func OldApiKeyAuth(token string) (interface{}, error) {
	return ApiKeyAuth(token)
}
func ApiKeyAuth(token string) (interface{}, error) {
	config := csm_manager.GetConfig()
	if *config.API_KEY == token {
		return true, nil // valid request
	}
	// invalid request
	return nil, errors.Unauthenticated("Authentication Failed for Sidecar")
}
