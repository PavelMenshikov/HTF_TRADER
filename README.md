# HTF Trader

HTF Trader is a Go execution-engine foundation for binary-options and prediction-market strategies.
The current runnable product is intentionally limited to paper trading. Live money mode is rejected by configuration validation until explicit live execution, wallet, and exchange adapters are implemented and reviewed.

## Run the paper engine

```bash
go run ./cmd/engine -config configs/paper.json
```

The command discovers deterministic paper markets, evaluates a strategy, applies risk limits, opens paper trades, and prints structured JSON logs with the run summary.
