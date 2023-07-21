package interfaces

import "github.com/skoona/ggapcmon/internal/entities"

type Service interface {
	MessageChannelByName(hostName string) entities.ChannelTuple
	ParseStatus(status []string) map[string]string
	Shutdown()
}
