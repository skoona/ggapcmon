package providers

import (
	"context"
	"github.com/skoona/ggapcmon/internal/interfaces"
)

type viewProvider struct {
	ctx context.Context
}

var (
	_ interfaces.ViewProvider = (*viewProvider)(nil)
	_ interfaces.Provider     = (*viewProvider)(nil)
)

func NewViewProvider(ctx context.Context) interfaces.ViewProvider {
	return &viewProvider{ctx: ctx}
}
