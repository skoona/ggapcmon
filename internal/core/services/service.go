package services

import (
	"context"
	"errors"
	"github.com/skoona/ggapcmon/internal/adapters/repository"
	"github.com/skoona/ggapcmon/internal/commons"
	"github.com/skoona/ggapcmon/internal/core/domain"
	"github.com/skoona/ggapcmon/internal/core/ports"
	"strings"
)

type service struct {
	ctx        context.Context
	cfg        ports.Configuration
	providers  map[string]ports.Provider
	publishers map[string]domain.ChannelTuple
}

var _ ports.Service = (*service)(nil)

func NewService(ctx context.Context, cfg ports.Configuration) (ports.Service, error) {
	var err error

	if len(cfg.HostKeys()) == 0 {
		return nil, errors.New("hosts param cannot be empty.")
	}
	s := &service{
		ctx:        ctx,
		providers:  map[string]ports.Provider{},
		publishers: map[string]domain.ChannelTuple{},
		cfg:        cfg,
	}

failure:
	for _, host := range s.cfg.Hosts() {
		if host.Enabled {
			commons.DebugLog("Service::begin(", host.Name, "::", host.Id, ") Init ")
			s.publishers[host.Id] = *domain.NewChannelTuple(16)
			apc, err := repository.NewAPCProvider(s.ctx, host, s.publishers[host.Id])
			if err != nil {
				commons.DebugLog("Service::begin(", host.Name, ") failed: ", err.Error())
				break failure
			}
			s.providers[host.Id] = apc
		}
	}

	return s, err
}
func (s *service) MessageChannelById(id string) domain.ChannelTuple {
	return s.publishers[id]
}
func (s *service) ParseStatus(status []string) map[string]string {
	params := map[string]string{}
	var key, value string

	//DATE     : Fri, 21 Jul 2023 00:16:52 EDT
	//0123456789012345678901234567890123456789
	//         1         2         3         4
	for _, line := range status {
		key = strings.TrimSpace(line[0:9])
		value = strings.TrimSpace(line[11:])
		params[key] = value
	}
	return params
}
func (s *service) Close() {
	commons.DebugLog("Service::Close() called.")
	for key, v := range s.publishers {
		v.Close()
		if z, ok := s.providers[key]; ok {
			z.Close()
		}
	}
}
