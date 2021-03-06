package main

import (
	"os"

	"github.com/SUSE/cf-usb-sidecar/csm-extensions/services/dev-postgres"
	"github.com/SUSE/cf-usb-sidecar/csm-extensions/services/dev-postgres/config"
	"github.com/SUSE/cf-usb-sidecar/csm-extensions/services/dev-postgres/provisioner"
	"github.com/SUSE/go-csm-lib/csm"
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

	if conf.Host == "" {
		logger.Fatal("SERVICE_POSTGRES_HOST environment variable is not set", nil)
	}

	request, err := csm.GetCSMRequest(os.Args)
	if err != nil {
		logger.Fatal("main", err)
	}

	csmConnection := csm.NewCSMFileConnection(request.OutputPath, logger)
	prov := provisioner.NewPqProvisioner(logger, conf)

	extension := postgres.NewPostgresExtension(prov, conf, logger)
	response, err := extension.GetStatus()
	if err != nil {
		err := csmConnection.Write(*response)
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
