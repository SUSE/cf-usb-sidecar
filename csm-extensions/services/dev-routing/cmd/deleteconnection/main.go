package main

import (
	"fmt"
	"os"

	"github.com/hpcloud/catalog-service-manager/csm-extensions/services/dev-routing"
	"github.com/hpcloud/catalog-service-manager/csm-extensions/services/dev-routing/config"
	"github.com/hpcloud/go-csm-lib/csm"
	"github.com/pivotal-golang/lager"
	"gopkg.in/caarlos0/env.v2"
)

func main() {

	var logger = lager.NewLogger("routing-extension")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	conf := config.RoutingConfig{}
	err := env.Parse(&conf)
	if err != nil {
		logger.Fatal("main", err)
	}

	request, err := csm.GetCSMRequest(os.Args[1:])
	if err != nil {
		logger.Fatal("main", err)
	}

	csmConnection := csm.NewCSMFileConnection(request.OutputPath, logger)

	extension := routing.NewRoutingExtension(conf, logger)

	response, err := extension.DeleteConnection(request.WorkspaceID, request.ConnectionID)
	if err != nil {
		err := csmConnection.WriteError(err)
		if err != nil {
			logger.Fatal("main", err)
		}
		os.Exit(0)
	}

	err = csmConnection.Write(*response)
	if err != nil {
		fmt.Println(*response)
		logger.Fatal("main", err)
	}
}
