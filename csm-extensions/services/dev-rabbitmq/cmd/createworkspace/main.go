package main

import (
	"fmt"
	"os"

	"github.com/SUSE/cf-usb-sidecar/csm-extensions/services/dev-rabbitmq"
	"github.com/SUSE/cf-usb-sidecar/csm-extensions/services/dev-rabbitmq/config"
	"github.com/SUSE/cf-usb-sidecar/csm-extensions/services/dev-rabbitmq/provisioner"
	"github.com/hpcloud/go-csm-lib/csm"
	"github.com/pivotal-golang/lager"
	"gopkg.in/caarlos0/env.v2"
)

func main() {

	var logger = lager.NewLogger("rabbitmq-extension")
	logger.RegisterSink(lager.NewWriterSink(os.Stdout, lager.DEBUG))

	conf := config.RabbitmqConfig{}
	err := env.Parse(&conf)
	if err != nil {
		logger.Fatal("main", err)
	}

	if conf.DockerEndpoint == "" {
		if conf.DockerHost == "" {
			logger.Fatal("DOCKER_ENDPOINT or DOCKER_HOST environment variables not set", nil)
		}

		conf.DockerEndpoint = fmt.Sprintf("http://%s:%s", conf.DockerHost, conf.DockerPort)
	}

	request, err := csm.GetCSMRequest(os.Args)
	if err != nil {
		logger.Fatal("main", err)
	}

	csmConnection := csm.NewCSMFileConnection(request.OutputPath, logger)
	prov := provisioner.NewRabbitHoleProvisioner(logger, conf)

	extension := rabbitmq.NewRabbitmqExtension(prov, conf, logger)

	response, err := extension.CreateWorkspace(request.WorkspaceID)
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
