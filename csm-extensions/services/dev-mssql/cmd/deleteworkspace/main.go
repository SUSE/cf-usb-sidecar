package main

import (
	"os"

	"github.com/hpcloud/catalog-service-manager/csm-extensions/services/dev-mssql"
	"github.com/hpcloud/catalog-service-manager/csm-extensions/services/dev-mssql/config"
	"github.com/hpcloud/catalog-service-manager/csm-extensions/services/dev-mssql/provisioner"
	"github.com/hpcloud/go-csm-lib/csm"
	"github.com/pivotal-golang/lager"
	"gopkg.in/caarlos0/env.v2"
)

func main() {

	var logger = lager.NewLogger("mssql-extension")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	conf := config.MssqlConfig{}
	err := env.Parse(&conf)
	if err != nil {
		logger.Fatal("main", err)
	}

	request, err := csm.GetCSMRequest(os.Args)
	if err != nil {
		logger.Fatal("main", err)
	}

	csmConnection := csm.NewCSMFileConnection(request.OutputPath, logger)
	prov := provisioner.NewGoMssqlProvisioner(logger, conf)

	extension := mssql.NewMSSQLExtension(prov, conf, logger)

	response, err := extension.DeleteWorkspace(request.WorkspaceID)
	if err != nil {
		err := csmConnection.WriteError(err)
		if err != nil {
			logger.Fatal("main", err)
		}
		os.Exit(0)
	}

	err = csmConnection.Write(*response)
	if err != nil {
		logger.Fatal("main", err)
	}
}
