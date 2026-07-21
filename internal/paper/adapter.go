package paper

import (
	"context"
	"time"

	"htf-trader/internal/domain"
)

type StaticExchange struct{ now func() time.Time }

func NewStaticExchange(now func() time.Time) *StaticExchange { return &StaticExchange{now: now} }
func (e *StaticExchange) Name() string                       { return "paper-static" }
func (e *StaticExchange) DiscoverMarkets(ctx context.Context) ([]domain.Market, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	n := e.now()
	return []domain.Market{{ID: "BTC-1H-UPDOWN", Symbol: "BTC/USD", Exchange: e.Name(), Expiry: n.Add(time.Hour), Price: 65000, PayoutMultiplier: 1.82}, {ID: "ETH-1H-UPDOWN", Symbol: "ETH/USD", Exchange: e.Name(), Expiry: n.Add(time.Hour), Price: 3200, PayoutMultiplier: 1.76}}, nil
}

type MomentumStrategy struct {
	name          string
	minConfidence float64
}

func NewMomentumStrategy(minConfidence float64) *MomentumStrategy {
	return &MomentumStrategy{name: "paper-momentum", minConfidence: minConfidence}
}
func (s *MomentumStrategy) Evaluate(m domain.Market, now time.Time) (domain.Signal, bool) {
	confidence := 0.62
	if m.Symbol == "BTC/USD" {
		confidence = 0.71
	}
	if confidence < s.minConfidence {
		return domain.Signal{}, false
	}
	return domain.Signal{MarketID: m.ID, Strategy: s.name, Direction: domain.DirectionUp, Confidence: confidence, CreatedAt: now}, true
}
