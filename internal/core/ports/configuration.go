package ports

import "github.com/skoona/ggapcmon/internal/core/domain"

type Configuration interface {
	VerifyHostConnection(h *domain.ApcHost) error
	Hosts() []*domain.ApcHost
	HostKeys() []string
	HostByName(hostName string) *domain.ApcHost
	AddHost(host *domain.ApcHost)
	Apply(h *domain.ApcHost) Configuration
	Save()
	Remove(hostName string)
	ResetConfig()
	Provider
}
