package interfaces

import "github.com/skoona/ggapcmon/internal/entities"

type Service interface {
	HostMessageChannel(hostName string) entities.ChannelTuple
	Shutdown()
}
