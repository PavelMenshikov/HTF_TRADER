package config

import "testing"

func TestValidateRequiresPaperMode(t *testing.T) {
	cfg := Config{Mode: "live", StartingBalance: 100, MaxStake: 10, MaxExposure: 20, MinConfidence: .5, TickInterval: 1}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected live mode to be rejected")
	}
}

func TestValidateAcceptsPaperMode(t *testing.T) {
	cfg := Config{Mode: "paper", StartingBalance: 100, MaxStake: 10, MaxExposure: 20, MinConfidence: .5, TickInterval: 1}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid config: %v", err)
	}
}
