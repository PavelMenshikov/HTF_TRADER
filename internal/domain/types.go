package domain

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrInvalidMarket = errors.New("invalid market")
	ErrInvalidSignal = errors.New("invalid signal")
	ErrRiskRejected  = errors.New("risk rejected")
)

type Direction string

const (
	DirectionUp   Direction = "UP"
	DirectionDown Direction = "DOWN"
)

type Market struct {
	ID               string
	Symbol           string
	Exchange         string
	Expiry           time.Time
	Price            float64
	PayoutMultiplier float64
}

func (m Market) Validate(now time.Time) error {
	if m.ID == "" || m.Symbol == "" || m.Exchange == "" {
		return fmt.Errorf("%w: identity fields are required", ErrInvalidMarket)
	}
	if !m.Expiry.After(now) {
		return fmt.Errorf("%w: expiry must be in the future", ErrInvalidMarket)
	}
	if m.Price <= 0 {
		return fmt.Errorf("%w: price must be positive", ErrInvalidMarket)
	}
	if m.PayoutMultiplier <= 1 {
		return fmt.Errorf("%w: payout multiplier must exceed 1", ErrInvalidMarket)
	}
	return nil
}

type Signal struct {
	MarketID   string
	Strategy   string
	Direction  Direction
	Confidence float64
	CreatedAt  time.Time
}

func (s Signal) Validate() error {
	if s.MarketID == "" || s.Strategy == "" {
		return fmt.Errorf("%w: market and strategy are required", ErrInvalidSignal)
	}
	if s.Direction != DirectionUp && s.Direction != DirectionDown {
		return fmt.Errorf("%w: unsupported direction %q", ErrInvalidSignal, s.Direction)
	}
	if s.Confidence < 0 || s.Confidence > 1 {
		return fmt.Errorf("%w: confidence must be between 0 and 1", ErrInvalidSignal)
	}
	if s.CreatedAt.IsZero() {
		return fmt.Errorf("%w: created_at is required", ErrInvalidSignal)
	}
	return nil
}

type RiskDecision struct {
	Approved bool
	Stake    float64
	Reason   string
}

type PaperTrade struct {
	ID         string
	Market     Market
	Signal     Signal
	Stake      float64
	EntryPrice float64
	OpenedAt   time.Time
}
