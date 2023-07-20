package interfaces

import (
	"github.com/skoona/ggapcmon/internal/entities"
	"time"
)

type Configuration interface {
	Hosts() []entities.ApcHost
	HostByName(hostName string) entities.ApcHost
	Apply(h entities.ApcHost) entities.ApcHost
	Save(hosts []entities.ApcHost) error
	Update(name, ip string, netperiod, graphperiod time.Duration, tray, enable bool) entities.ApcHost
	HostKeys() []string
	ResetConfig()
}
