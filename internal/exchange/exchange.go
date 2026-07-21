package exchange

import (
	"context"
	"htf-trader/internal/domain"
)

type Adapter interface {
	Name() string
	DiscoverMarkets(context.Context) ([]domain.Market, error)
}
