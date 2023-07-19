package interfaces

import "github.com/skoona/ggapcmon/internal/entities"

type Configuration interface {
	Hosts() []entities.ApcHost
	HostByName(hostName string) entities.ApcHost
}
