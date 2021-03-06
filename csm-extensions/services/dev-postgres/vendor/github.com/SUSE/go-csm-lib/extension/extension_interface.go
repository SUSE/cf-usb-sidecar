package extension

import (
	"github.com/SUSE/go-csm-lib/csm"
)

type Extension interface {
	CreateConnection(workspaceID, connectionID string) (*csm.CSMResponse, error)
	CreateWorkspace(workspaceID string) (*csm.CSMResponse, error)
	DeleteConnection(workspaceID, connectionID string) (*csm.CSMResponse, error)
	DeleteWorkspace(wokspaceID string) (*csm.CSMResponse, error)
	GetConnection(workspaceID, connectionID string) (*csm.CSMResponse, error)
	GetWorkspace(workspaceID string) (*csm.CSMResponse, error)
	GetStatus() (*csm.CSMResponse, error)
}
