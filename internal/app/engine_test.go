package app

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"htf-trader/internal/config"
	"htf-trader/internal/metrics"
	"htf-trader/internal/paper"
)

func TestRunOnceOpensPaperTradesWithinRisk(t *testing.T) {
	cfg := config.Config{Mode: "paper", StartingBalance: 100, MaxStake: 25, MaxExposure: 50, MinConfidence: .6, TickInterval: time.Second}
	now := func() time.Time { return time.Date(2026, 7, 21, 0, 0, 0, 0, time.UTC) }
	engine, err := NewEngine(cfg, paper.NewStaticExchange(now), paper.NewMomentumStrategy(cfg.MinConfidence), metrics.NewCollector(cfg.StartingBalance), slog.Default(), now)
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	got, err := engine.RunOnce(context.Background())
	if err != nil {
		t.Fatalf("run once: %v", err)
	}
	if got.Trades != 2 {
		t.Fatalf("trades = %d, want 2", got.Trades)
	}
	if got.Exposure != 50 {
		t.Fatalf("exposure = %f, want 50", got.Exposure)
	}
	if got.Balance != 50 {
		t.Fatalf("balance = %f, want 50", got.Balance)
	}
}

func TestRunOnceRejectsAboveExposure(t *testing.T) {
	cfg := config.Config{Mode: "paper", StartingBalance: 100, MaxStake: 25, MaxExposure: 25, MinConfidence: .6, TickInterval: time.Second}
	now := func() time.Time { return time.Date(2026, 7, 21, 0, 0, 0, 0, time.UTC) }
	engine, err := NewEngine(cfg, paper.NewStaticExchange(now), paper.NewMomentumStrategy(cfg.MinConfidence), metrics.NewCollector(cfg.StartingBalance), slog.Default(), now)
	if err != nil {
		t.Fatalf("new engine: %v", err)
	}
	got, err := engine.RunOnce(context.Background())
	if err != nil {
		t.Fatalf("run once: %v", err)
	}
	if got.Trades != 1 {
		t.Fatalf("trades = %d, want 1", got.Trades)
	}
	if got.Rejected != 1 {
		t.Fatalf("rejected = %d, want 1", got.Rejected)
	}
}
