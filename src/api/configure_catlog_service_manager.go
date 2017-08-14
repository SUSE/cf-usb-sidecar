package api

import (
	"net/http"

	errors "github.com/go-swagger/go-swagger/errors"
	httpkit "github.com/go-swagger/go-swagger/httpkit"
	middleware "github.com/go-swagger/go-swagger/httpkit/middleware"

	"github.com/SUSE/cf-usb-sidecar/cmd/catalog-service-manager/handlers"

	"github.com/SUSE/cf-usb-sidecar/generated/CatalogServiceManager/restapi/operations"

	"github.com/SUSE/cf-usb-sidecar/generated/CatalogServiceManager/restapi/operations/connection"
	"github.com/SUSE/cf-usb-sidecar/generated/CatalogServiceManager/restapi/operations/status"
	"github.com/SUSE/cf-usb-sidecar/generated/CatalogServiceManager/restapi/operations/workspace"
)

// This file is safe to edit. Once it exists it will not be overwritten

func configureFlags(api *operations.CatlogServiceManagerAPI) {

}

func ConfigureAPI(api *operations.CatlogServiceManagerAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	api.JSONConsumer = httpkit.JSONConsumer()

	api.JSONProducer = httpkit.JSONProducer()

	api.APIKeyAuth = handlers.ApiKeyAuth

	api.OldAPIKeyAuth = handlers.OldApiKeyAuth

	api.ConnectionCreateConnectionHandler = connection.CreateConnectionHandlerFunc(func(params connection.CreateConnectionParams, principal interface{}) middleware.Responder {
		return handlers.CreateConnection(params.WorkspaceID, params.ConnectionCreateRequest)
	})

	api.WorkspaceCreateWorkspaceHandler = workspace.CreateWorkspaceHandlerFunc(func(params workspace.CreateWorkspaceParams, principal interface{}) middleware.Responder {
		return handlers.CreateWorkspace(params.CreateWorkspaceRequest)
	})

	api.WorkspaceDeleteWorkspaceHandler = workspace.DeleteWorkspaceHandlerFunc(func(params workspace.DeleteWorkspaceParams, principal interface{}) middleware.Responder {
		return handlers.DeleteWorkspace(params.WorkspaceID)
	})

	api.ConnectionDeleteConnectionHandler = connection.DeleteConnectionHandlerFunc(func(params connection.DeleteConnectionParams, principal interface{}) middleware.Responder {
		return handlers.DeleteConnection(params.WorkspaceID, params.ConnectionID)
	})

	api.WorkspaceGetWorkspaceHandler = workspace.GetWorkspaceHandlerFunc(func(params workspace.GetWorkspaceParams, principal interface{}) middleware.Responder {
		return handlers.GetWorkspace(params.WorkspaceID)
	})

	api.ConnectionGetConnectionHandler = connection.GetConnectionHandlerFunc(func(params connection.GetConnectionParams, principal interface{}) middleware.Responder {
		return handlers.GetConnection(params.WorkspaceID, params.ConnectionID)
	})

	api.StatusStatusHandler = status.StatusHandlerFunc(func(principal interface{}) middleware.Responder {
		return handlers.GetStatus()
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
