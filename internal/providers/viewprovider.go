package providers

import (
	"context"
	"github.com/skoona/ggapcmon/internal/interfaces"
	"log"
)

type viewProvider struct {
	ctx context.Context
	log *log.Logger
}

var (
	_ interfaces.ViewProvider = (*viewProvider)(nil)
	_ interfaces.Provider     = (*viewProvider)(nil)
)

func NewViewProvider(ctx context.Context, log *log.Logger) interfaces.ViewProvider {
	return &viewProvider{
		ctx: ctx,
		log: log,
	}
}

// End closes the apc connection and stops go routines
func (a *viewProvider) Shutdown() {
	a.log.Println("ViewProvider::Shutdown() called.")
}
