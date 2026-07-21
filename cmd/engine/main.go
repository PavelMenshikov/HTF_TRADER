package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"htf-trader/internal/app"
	"htf-trader/internal/config"
	"htf-trader/internal/logging"
)

func main() {
	configPath := flag.String("config", "configs/paper.json", "path to engine configuration")
	flag.Parse()
	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "configuration error: %v\n", err)
		os.Exit(1)
	}
	logger := logging.New()
	engine, err := app.NewPaperEngine(cfg, logger)
	if err != nil {
		fmt.Fprintf(os.Stderr, "startup error: %v\n", err)
		os.Exit(1)
	}
	snapshot, err := engine.RunOnce(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "run error: %v\n", err)
		os.Exit(1)
	}
	logger.Info("paper run complete", "component", "engine", "markets", snapshot.Markets, "signals", snapshot.Signals, "approved", snapshot.Approved, "rejected", snapshot.Rejected, "trades", snapshot.Trades, "balance", snapshot.Balance, "exposure", snapshot.Exposure)
}
