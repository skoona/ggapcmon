package ports

import (
	"github.com/skoona/ggapcmon/internal/core/domain"
)

type Service interface {
	MessageChannelByName(hostName string) domain.ChannelTuple
	ParseStatus(status []string) map[string]string
	Provider
}
