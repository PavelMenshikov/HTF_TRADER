package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Mode             string        `json:"mode"`
	StartingBalance  float64       `json:"starting_balance"`
	MaxStake         float64       `json:"max_stake"`
	MaxExposure      float64       `json:"max_exposure"`
	MinConfidence    float64       `json:"min_confidence"`
	TickInterval     time.Duration `json:"-"`
	TickIntervalText string        `json:"tick_interval"`
}

func Load(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config: %w", err)
	}
	var c Config
	if err := json.Unmarshal(b, &c); err != nil {
		return Config{}, fmt.Errorf("parse config: %w", err)
	}
	d, err := time.ParseDuration(c.TickIntervalText)
	if err != nil {
		return Config{}, fmt.Errorf("parse tick_interval: %w", err)
	}
	c.TickInterval = d
	if err := c.Validate(); err != nil {
		return Config{}, err
	}
	return c, nil
}

func (c Config) Validate() error {
	if c.Mode != "paper" {
		return errors.New("mode must be paper until live trading is explicitly enabled")
	}
	if c.StartingBalance <= 0 {
		return errors.New("starting_balance must be positive")
	}
	if c.MaxStake <= 0 || c.MaxStake > c.StartingBalance {
		return errors.New("max_stake must be positive and not exceed starting_balance")
	}
	if c.MaxExposure <= 0 || c.MaxExposure > c.StartingBalance {
		return errors.New("max_exposure must be positive and not exceed starting_balance")
	}
	if c.MinConfidence < 0 || c.MinConfidence > 1 {
		return errors.New("min_confidence must be between 0 and 1")
	}
	if c.TickInterval <= 0 {
		return errors.New("tick_interval must be positive")
	}
	return nil
}
