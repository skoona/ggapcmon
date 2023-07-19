package services

import (
	"context"
	"fyne.io/fyne/v2"
	"github.com/skoona/ggapcmon/internal/interfaces"
)

type service struct {
	ctx       context.Context
	settings  fyne.Preferences
	providers []interfaces.Provider
}

var _ interfaces.Service = (*service)(nil)

func NewService(ctx context.Context, settings fyne.Preferences) interfaces.Service {
	return &service{
		ctx:       ctx,
		settings:  settings,
		providers: []interfaces.Provider{},
	}
}
