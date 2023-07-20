package services

import (
	"context"
	"errors"
	"github.com/skoona/ggapcmon/internal/entities"
	"github.com/skoona/ggapcmon/internal/interfaces"
	"github.com/skoona/ggapcmon/internal/providers"
	"log"
)

type service struct {
	ctx        context.Context
	cfg        interfaces.Configuration
	providers  map[string]interfaces.Provider
	publishers map[string]entities.ChannelTuple
	log        *log.Logger
}

var _ interfaces.Service = (*service)(nil)

func NewService(ctx context.Context, cfg interfaces.Configuration, log *log.Logger) (interfaces.Service, error) {
	if len(cfg.HostKeys()) == 0 {
		return nil, errors.New("hosts param cannot be empty.")
	}
	s := &service{
		ctx:        ctx,
		providers:  map[string]interfaces.Provider{},
		publishers: map[string]entities.ChannelTuple{},
		cfg:        cfg,
		log:        log,
	}
	err := s.begin()

	return s, err
}

func (s *service) begin() error {
	var err error
failure:
	for _, host := range s.cfg.Hosts() {
		if host.Enabled {

			s.publishers[host.Name] = *entities.NewChannelTuple(16)
			apc, err := providers.NewAPCProvider(s.ctx, host, s.publishers[host.Name], s.log)
			if err != nil {
				s.log.Println("Service::begin(", host.Name, ") failed: ", err.Error())
				break failure
			}
			s.providers[host.Name] = apc
		}
	}

	return err
}
func (s *service) Shutdown() {
	s.log.Println("Service::Shutdown() called.")
	for key, v := range s.publishers {
		v.Close()
		s.providers[key].Shutdown()
	}
}
func (s *service) HostMessageChannel(hostName string) entities.ChannelTuple {
	return s.publishers[hostName]
}
