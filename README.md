<div align="center">

# ⚡ HTF Trader

**Production-oriented Go execution-engine foundation for binary-options and EVM prediction-market strategies.**

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Architecture](https://img.shields.io/badge/Architecture-Clean%20Architecture-6C5CE7?style=for-the-badge)](ARCHITECTURE.md)
[![Mode](https://img.shields.io/badge/Mode-Paper%20Trading-2ECC71?style=for-the-badge)](#-current-product-status)
[![Target](https://img.shields.io/badge/Target-Buffer%20Finance-F39C12?style=for-the-badge)](#-exchange-roadmap)

</div>

---

## ✅ Current product status

HTF Trader **does run today** as a deterministic paper-trading engine.

The current runnable product intentionally supports **paper trading only**. Live-money execution is rejected by configuration validation until live execution, wallet handling, exchange adapters, signing, and transaction recovery are implemented and reviewed.

What works now:

- 🔎 deterministic paper-market discovery;
- 🧠 strategy evaluation through an application-level strategy interface;
- 🛡️ risk checks for confidence, exposure, and paper balance;
- 📈 paper trade accounting and metrics snapshots;
- 🪵 structured JSON logs for operator visibility;
- 🧪 unit tests for configuration validation and engine behavior.

---

## 🧰 Technology stack

| Area | Technology |
| --- | --- |
| Language | ![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&logoColor=white) |
| Runtime | ![CLI](https://img.shields.io/badge/CLI-Engine-111827?logo=gnometerminal&logoColor=white) |
| Architecture | ![Clean Architecture](https://img.shields.io/badge/Clean%20Architecture-SOLID-6C5CE7) |
| Logging | ![slog](https://img.shields.io/badge/Go-slog-00ADD8) |
| Config | ![JSON](https://img.shields.io/badge/JSON-Config-000000?logo=json&logoColor=white) |
| Trading mode | ![Paper](https://img.shields.io/badge/Paper-Trading-2ECC71) |
| Future chain target | ![Arbitrum](https://img.shields.io/badge/Arbitrum-EVM-28A0F0) |

---

## 🚀 Quick start

### 1. Verify the repository

```bash
go test ./...
```

### 2. Run the paper engine

```bash
go run ./cmd/engine -config configs/paper.json
```

The command discovers deterministic paper markets, evaluates configured strategy confidence, applies risk limits, opens paper trades when approved, and prints structured JSON logs with the final run summary.

Example output shape:

```json
{
  "msg": "paper run complete",
  "component": "engine",
  "markets": 2,
  "signals": 1,
  "approved": 1,
  "rejected": 0,
  "trades": 1,
  "balance": 990,
  "exposure": 10
}
```

---

## ⚙️ Configuration

Default paper configuration lives in [`configs/paper.json`](configs/paper.json).

Key controls:

| Field | Purpose |
| --- | --- |
| `mode` | Must be `paper` for the current runnable product. |
| `exchange` | Exchange label used by adapters and logs. |
| `starting_balance` | Initial paper balance. |
| `max_stake` | Maximum stake per approved paper trade. |
| `max_exposure` | Maximum total open paper exposure. |
| `min_confidence` | Minimum strategy confidence required for risk approval. |

---

## 🏗️ Architecture

HTF Trader follows Clean Architecture and dependency inversion so strategy logic stays independent from exchange, wallet, RPC, ABI, gas, nonce, transaction, and persistence concerns.

```text
cmd/engine            ── CLI composition root
internal/app          ── use-case orchestration and risk approval
internal/domain       ── markets, signals, risk decisions, invariants
internal/exchange     ── exchange adapter contract
internal/paper        ── deterministic paper exchange and strategy
internal/config       ── typed config loading and validation
internal/metrics      ── paper balance, exposure, counters
internal/logging      ── structured JSON logging
```

For the full engineering blueprint, see [`ARCHITECTURE.md`](ARCHITECTURE.md).

---

## 🧭 Exchange roadmap

| Exchange | Status |
| --- | --- |
| Buffer Finance | 🎯 first live adapter target |
| OptionBlitz | 🧩 planned future adapter |
| Limitless | 🧩 planned future adapter |
| Other EVM prediction markets | 🧩 supported by adapter-oriented design |

---

## 🛡️ Safety model

Live execution is deliberately unavailable until the production-critical pieces exist:

- wallet and private-key handling;
- RPC and chain-id validation;
- ABI and contract binding strategy;
- gas, nonce, signing, broadcasting, and confirmation flows;
- durable persistence and idempotency;
- transaction recovery and operator alerting;
- exchange-specific review and test coverage.

This keeps the repository safe to run while the engine evolves from paper trading toward production live execution.

---

## 🧪 Development checks

```bash
go test ./...
```

Recommended before every commit:

```bash
go test ./...
go run ./cmd/engine -config configs/paper.json
```

---

## 📌 Project intent

HTF Trader is not a one-off bot. It is intended to grow into a maintainable quantitative trading execution platform where:

- strategies produce signals only;
- risk owns position sizing and exposure limits;
- executors own transaction creation and delivery;
- exchanges expose markets, prices, liquidity, and contract interactions;
- persistence owns durable trading state;
- observability is built in from the start.

