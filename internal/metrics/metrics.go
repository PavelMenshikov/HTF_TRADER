package metrics

import "sync"

type Snapshot struct {
	Markets, Signals, Approved, Rejected, Trades int
	Balance, Exposure                            float64
}

type Collector struct {
	mu sync.Mutex
	s  Snapshot
}

func NewCollector(startingBalance float64) *Collector {
	return &Collector{s: Snapshot{Balance: startingBalance}}
}
func (c *Collector) AddMarkets(n int) { c.mu.Lock(); defer c.mu.Unlock(); c.s.Markets += n }
func (c *Collector) AddSignal()       { c.mu.Lock(); defer c.mu.Unlock(); c.s.Signals++ }
func (c *Collector) AddApproved()     { c.mu.Lock(); defer c.mu.Unlock(); c.s.Approved++ }
func (c *Collector) AddRejected()     { c.mu.Lock(); defer c.mu.Unlock(); c.s.Rejected++ }
func (c *Collector) AddTrade(stake float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.s.Trades++
	c.s.Exposure += stake
	c.s.Balance -= stake
}
func (c *Collector) Snapshot() Snapshot { c.mu.Lock(); defer c.mu.Unlock(); return c.s }
