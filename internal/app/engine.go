package app

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"htf-trader/internal/config"
	"htf-trader/internal/domain"
	"htf-trader/internal/exchange"
	"htf-trader/internal/metrics"
	"htf-trader/internal/paper"
)

type Strategy interface {
	Evaluate(domain.Market, time.Time) (domain.Signal, bool)
}

type Engine struct {
	cfg      config.Config
	exchange exchange.Adapter
	strategy Strategy
	metrics  *metrics.Collector
	logger   *slog.Logger
	now      func() time.Time
}

func NewEngine(cfg config.Config, ex exchange.Adapter, strategy Strategy, collector *metrics.Collector, logger *slog.Logger, now func() time.Time) (*Engine, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if ex == nil || strategy == nil || collector == nil || logger == nil || now == nil {
		return nil, fmt.Errorf("engine dependencies are required")
	}
	return &Engine{cfg: cfg, exchange: ex, strategy: strategy, metrics: collector, logger: logger, now: now}, nil
}

func NewPaperEngine(cfg config.Config, logger *slog.Logger) (*Engine, error) {
	now := time.Now
	return NewEngine(cfg, paper.NewStaticExchange(now), paper.NewMomentumStrategy(cfg.MinConfidence), metrics.NewCollector(cfg.StartingBalance), logger, now)
}

func (e *Engine) RunOnce(ctx context.Context) (metrics.Snapshot, error) {
	markets, err := e.exchange.DiscoverMarkets(ctx)
	if err != nil {
		return e.metrics.Snapshot(), fmt.Errorf("discover markets: %w", err)
	}
	e.metrics.AddMarkets(len(markets))
	for _, market := range markets {
		if err := market.Validate(e.now()); err != nil {
			e.metrics.AddRejected()
			e.logger.Warn("market rejected", "component", "market", "market_id", market.ID, "error", err)
			continue
		}
		signal, ok := e.strategy.Evaluate(market, e.now())
		if !ok {
			continue
		}
		e.metrics.AddSignal()
		if err := signal.Validate(); err != nil {
			return e.metrics.Snapshot(), fmt.Errorf("validate signal: %w", err)
		}
		decision := e.approve(signal)
		if !decision.Approved {
			e.metrics.AddRejected()
			e.logger.Info("risk rejected signal", "component", "risk", "market_id", signal.MarketID, "reason", decision.Reason)
			continue
		}
		e.metrics.AddApproved()
		e.metrics.AddTrade(decision.Stake)
		e.logger.Info("paper trade opened", "component", "paper", "exchange", market.Exchange, "market_id", market.ID, "stake", decision.Stake, "direction", signal.Direction)
	}
	return e.metrics.Snapshot(), nil
}

func (e *Engine) approve(signal domain.Signal) domain.RiskDecision {
	snapshot := e.metrics.Snapshot()
	if signal.Confidence < e.cfg.MinConfidence {
		return domain.RiskDecision{Approved: false, Reason: "confidence below minimum"}
	}
	if snapshot.Exposure+e.cfg.MaxStake > e.cfg.MaxExposure {
		return domain.RiskDecision{Approved: false, Reason: "max exposure exceeded"}
	}
	if snapshot.Balance < e.cfg.MaxStake {
		return domain.RiskDecision{Approved: false, Reason: "insufficient paper balance"}
	}
	return domain.RiskDecision{Approved: true, Stake: e.cfg.MaxStake}
}
