package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	spec "github.com/go-swagger/go-swagger/spec"
	flags "github.com/jessevdk/go-flags"
	graceful "github.com/tylerb/graceful"

	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager/restapi"
	"github.com/hpcloud/catalog-service-manager/generated/CatalogServiceManager/restapi/operations"
	srvManagerAPI "github.com/hpcloud/catalog-service-manager/src/api"
	"github.com/hpcloud/catalog-service-manager/src/common"
	"github.com/hpcloud/catalog-service-manager/src/csm_manager"
	"github.com/hpcloud/catalog-service-manager/src/tls"
)

// This file was generated by the swagger tool.
// Make sure not to overwrite this file after you generated it because all your edits would be lost!

func main() {
	csm_manager.InitServiceCatalogManager()
	logger := csm_manager.GetLogger()

	swaggerSpec, err := spec.New(restapi.SwaggerJSON, "")
	if err != nil {
		logger.Fatalln(err)
	}

	api := operations.NewCatlogServiceManagerAPI(swaggerSpec)
	server := restapi.NewServer(api)
	handler := srvManagerAPI.ConfigureAPI(api)

	defer server.Shutdown()

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = `Catalog Service Manager API`
	parser.LongDescription = `This API will be available on the Catalog Service
Manager container which runs along side your service and serves some of the
service management capabilities.
`

	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
	}

	if _, err := parser.Parse(); err != nil {
		os.Exit(1)
	}

	exposedPort := os.Getenv("PORT")
	if exposedPort == "" {
		exposedPort = "8081"
	}

	httpServer := &graceful.Server{Server: new(http.Server)}
	httpServer.Addr = "0.0.0.0:" + exposedPort
	httpServer.Handler = handler
	configuration := common.NewServiceManagerConfiguration()

	// transform logrus to a writer for the http server errors
	// httpServer.ErrorLog = log.New(common.Logger.Writer(), "", 0)

	if !isTLSCertsProvided(configuration) {
		fmt.Printf("ssl certs not given so generate self-signed cert")
		cert, err := tls.GenCert()
		if err != nil {
			shutdown(api, err)
		}
		os.Setenv("TLS_CERT_FILE", cert.Public)
		os.Setenv("TLS_PRIVATE_KEY_FILE", cert.Private)
		configuration = common.NewServiceManagerConfiguration()
	}

	fmt.Printf("TLS_CERT_FILE: %s\n", *configuration.TLS_CERT_FILE)
	fmt.Printf("TLS_PRIVATE_KEY_FILE: %s\n", *configuration.TLS_PRIVATE_KEY_FILE)

	// serving https
	logger.Printf("Helion Service Manager listening at https://%s", httpServer.Addr)
	if err := httpServer.ListenAndServeTLS(*configuration.TLS_CERT_FILE, *configuration.TLS_PRIVATE_KEY_FILE); err != nil {
		shutdown(api, err)
	}

	go func() {
		<-httpServer.StopChan()
		api.ServerShutdown()
	}()

}

// shutdown closes down the api server
func shutdown(api *operations.CatlogServiceManagerAPI, err error) {
	api.ServerShutdown()
	// common.Logger.Fatalln(err)
	log.Fatalln(err)
}

// isSSLEnabled returns true if all the SSL config variables are set
// false otherwise
func isTLSCertsProvided(configuration *common.ServiceManagerConfiguration) bool {
	if configuration.TLS_CERT_FILE == nil {
		return false
	}
	if _, err := os.Stat(*configuration.TLS_CERT_FILE); err != nil {
		return false
	}
	if configuration.TLS_PRIVATE_KEY_FILE == nil {
		return false
	}
	if _, err := os.Stat(*configuration.TLS_PRIVATE_KEY_FILE); err != nil {
		return false
	}
	return true
}

// tcpKeepAliveListener is copied from the stdlib net/http package

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}
