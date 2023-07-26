package interfaces

import (
	"github.com/skoona/ggapcmon/internal/entities"
)

type Configuration interface {
	VerifyHostConnection(h *entities.ApcHost) error
	Hosts() []*entities.ApcHost
	HubHosts() []*entities.HubHost
	HostKeys() []string
	HostByName(hostName string) *entities.ApcHost
	AddHost(host *entities.ApcHost)
	Apply(h *entities.ApcHost) Configuration
	Save()
	Remove(hostName string)
	ResetConfig()
	Shutdown()
}
