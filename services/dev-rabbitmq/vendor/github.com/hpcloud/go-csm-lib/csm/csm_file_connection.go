package csm

import (
	"encoding/json"
	"os"

	"github.com/hpcloud/go-csm-lib/csm/status"
	"github.com/pivotal-golang/lager"
)

type csmFileConnection struct {
	filePath string
	logger   lager.Logger
}

func NewCSMFileConnection(filePath string, logger lager.Logger) CSMConnection {
	return &csmFileConnection{filePath: filePath, logger: logger}
}

func (c *csmFileConnection) Write(response CSMResponse) error {

	f, err := os.OpenFile(c.filePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	defer f.Close()
	if err != nil {
		return err
	}

	c.logger.Debug("csm-connection-write", lager.Data{"filename": f.Name()})

	b, err := json.Marshal(response)
	if err != nil {
		return err
	}

	_, err = f.Write(b)

	return err
}

func (c *csmFileConnection) WriteError(input error) error {
	response := NewCSMResponse(500, input.Error(), status.Failed)

	err := c.Write(response)
	return err
}
