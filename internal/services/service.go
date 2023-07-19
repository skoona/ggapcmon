package services

import (
	"context"
	"github.com/skoona/ggapcmon/internal/entities"
	"github.com/skoona/ggapcmon/internal/interfaces"
	"github.com/skoona/ggapcmon/internal/providers"
	"log"
)

type service struct {
	ctx        context.Context
	apchosts   []entities.ApcHost
	providers  map[string]interfaces.Provider
	publishers map[string]chan []string
	log        *log.Logger
}

var _ interfaces.Service = (*service)(nil)

func NewService(ctx context.Context, hosts []entities.ApcHost, log *log.Logger) (interfaces.Service, error) {
	s := &service{
		ctx:        ctx,
		providers:  map[string]interfaces.Provider{},
		publishers: map[string]chan []string{},
		apchosts:   hosts,
		log:        log,
	}
	err := s.begin()

	return s, err
}

func (s *service) begin() error {
	var err error
failure:
	for _, host := range s.apchosts {
		c := make(chan []string, 16)
		s.publishers[host.Name] = c
		apc, err := providers.NewAPCProvider(s.ctx, host, c, s.log)
		if err != nil {
			s.log.Println("Service::begin(", host.Name, ") failed: ", err.Error())
			break failure
		}
		s.providers[host.Name] = apc
	}

	return err
}
func (s *service) Shutdown() {
	s.log.Println("Service::Shutdown() called.")
	for key, v := range s.publishers {
		close(v)
		s.providers[key].Shutdown()
	}
}
func (s *service) HostMessageChannel(hostName string) chan []string {
	return s.publishers[hostName]
}
