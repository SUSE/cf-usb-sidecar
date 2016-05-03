package main

import (
	"os"

	"github.com/hpcloud/go-csm-lib/csm"
	"github.com/hpcloud/catalog-service-manager/services/dev-postgres"
	"github.com/hpcloud/catalog-service-manager/services/dev-postgres/config"
	"github.com/hpcloud/catalog-service-manager/services/dev-postgres/provisioner"
	"github.com/pivotal-golang/lager"
	"gopkg.in/caarlos0/env.v2"
)

func main() {

	var logger = lager.NewLogger("postgres-extension")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	conf := config.PostgresConfig{}
	err := env.Parse(&conf)
	if err != nil {
		logger.Fatal("main", err)
	}

	request, err := csm.GetCSMRequest(os.Args)
	if err != nil {
		logger.Fatal("main", err)
	}

	csmConnection := csm.NewCSMFileConnection(request.OutputPath, logger)
	prov := provisioner.NewPqProvisioner(logger, conf)

	extension := postgres.NewPostgresExtension(prov, conf, logger)

	response, err := extension.CreateWorkspace(request.WorkspaceID)
	if err != nil {
		err := csmConnection.WriteError(err)
		if err != nil {
			logger.Fatal("main", err)
		}
	}

	err = csmConnection.Write(*response)
	if err != nil {
		logger.Fatal("main", err)
	}
}
