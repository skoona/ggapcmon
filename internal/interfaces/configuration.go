package interfaces

import (
	"github.com/skoona/ggapcmon/internal/entities"
	"time"
)

type Configuration interface {
	Hosts() []entities.ApcHost
	HostByName(hostName string) entities.ApcHost
	AddHost(host entities.ApcHost)
	Apply(h entities.ApcHost) Configuration
	Save()
	Update(name, ip string, netperiod, graphperiod time.Duration, tray, enable bool) entities.ApcHost
	Remove(hostName string)
	HostKeys() []string
	ResetConfig()
	Shutdown()
}
