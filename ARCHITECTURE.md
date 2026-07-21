# ARCHITECTURE.md

You are the Principal Software Architect of this project.

You are NOT writing code.

You are designing and maintaining the engineering blueprint for a production-grade quantitative trading execution engine.

This document is the single source of truth for every future developer and AI agent working in this repository.

The initial target exchange is Buffer Finance (Supurr) running on Arbitrum.

Future exchanges include:

- OptionBlitz
- Limitless
- Any EVM prediction market

The architecture must support adding exchanges without changing the trading engine.

## Vision

This project solves the problem of reliably executing quantitative binary-options and prediction-market strategies across EVM-based exchanges while keeping strategy logic independent from blockchain, exchange, wallet, gas, transaction, and persistence concerns.

The project is an execution engine, not a Buffer-specific bot. A bot couples decisions, exchange APIs, wallet operations, contract calls, and operational recovery into one program. That structure is fast to prototype and expensive to maintain. It fails when a second exchange is added, when execution rules change, or when production reliability requirements appear.

An execution engine provides stable internal contracts for markets, signals, risk decisions, orders, positions, transactions, statistics, and observability. Buffer Finance is the first adapter, not the center of the system. OptionBlitz, Limitless, and future EVM prediction markets must integrate by implementing the same exchange adapter contract without changing strategy, risk, execution, portfolio, metrics, or persistence code.

Strategies are separated from infrastructure because strategies express trading intent, not operational mechanics. A strategy evaluates market state and produces a signal. It must not know RPC endpoints, wallets, private keys, ABI encoding, gas prices, nonces, transaction hashes, contract addresses, retry behavior, database schemas, or notification providers. This separation makes strategies testable, portable, deterministic, and safe to run in live, paper, and backtest modes.

## Design Principles

### Clean Architecture

Business rules sit at the center. External systems sit at the edge. Dependencies point inward.

- Domain defines the vocabulary and invariants.
- Application orchestrates use cases.
- Infrastructure implements external integrations.
- Presentation starts processes and exposes operator interfaces.

No infrastructure package may define business rules. No domain package may import infrastructure.

### SOLID

- Single Responsibility: each package has one reason to change.
- Open/Closed: new exchanges, strategies, notifiers, and persistence backends are added through interfaces.
- Liskov Substitution: implementations must satisfy contracts without surprising callers.
- Interface Segregation: interfaces are small and owned by consumers where practical.
- Dependency Inversion: high-level policy depends on abstractions, not concrete providers.

### Dependency Inversion

Application services depend on interfaces for exchange access, market data, risk decisions, execution, persistence, metrics, logging, and notification. Concrete implementations are wired at startup.

### Composition Over Inheritance

Go has no inheritance hierarchy in this project. Behavior is composed from small interfaces, structs, and constructor-injected dependencies.

### Interface First

Public boundaries are designed before implementation. Every adapter must implement explicit contracts. Every application service must declare what it consumes.

### Explicit Dependencies

Dependencies are passed through constructors. Hidden package-level state is forbidden. Configuration, logger, metrics, database handles, blockchain clients, clocks, and notifiers are explicit dependencies.

### Immutable Domain Models Where Possible

Domain values should be immutable after construction when practical. State transitions should create new values or pass through explicit methods that validate invariants.

### Small Packages

Packages should be cohesive and narrow. A package that mixes market discovery, risk, execution, and persistence is incorrectly designed.

### No Cyclic Dependencies

Package cycles are forbidden. A cycle indicates that boundaries are wrong.

## Project Layers

### Presentation

Presentation owns process entry points and operator-facing surfaces.

Responsibilities:

- CLI commands.
- Daemon startup.
- Configuration loading entry point.
- Dependency graph assembly.
- Signal handling.
- Graceful shutdown coordination.
- Health and readiness endpoints if exposed.

Forbidden dependencies:

- Presentation may depend on application and infrastructure wiring.
- Presentation must not contain trading rules.
- Presentation must not encode exchange-specific logic outside dependency construction.

### Application

Application owns use-case orchestration.

Responsibilities:

- Market discovery workflow.
- Strategy evaluation workflow.
- Risk validation workflow.
- Order execution workflow.
- Position reconciliation workflow.
- Statistics update workflow.
- Paper trading workflow.
- Backtesting workflow.
- Scheduling and cancellation boundaries.

Forbidden dependencies:

- Application must not import concrete exchange implementations except from composition roots.
- Application must not import blockchain client implementations directly.
- Application must not know ABI, contract addresses, private keys, RPC transport details, SQL syntax, or notification provider APIs.

### Domain

Domain owns business language and invariants.

Responsibilities:

- Market models.
- Signal models.
- Order models.
- Position models.
- Portfolio models.
- Risk decisions.
- Trade lifecycle states.
- Errors that represent business rule violations.
- Value objects for money, chain identifiers, expiry, side, direction, slippage, and confidence.

Forbidden dependencies:

- Domain must not import application, infrastructure, blockchain, persistence, logging, metrics, or configuration packages.
- Domain must not perform I/O.

### Infrastructure

Infrastructure owns concrete integrations that are not blockchain-specific.

Responsibilities:

- HTTP clients.
- Notification providers.
- External data providers.
- Clock implementations.
- Serialization adapters.
- Runtime platform integrations.

Forbidden dependencies:

- Infrastructure must not define strategy or risk rules.
- Infrastructure must not bypass application services to mutate trading state.

### Blockchain

Blockchain owns EVM connectivity and transaction mechanics.

Responsibilities:

- RPC clients.
- Chain ID validation.
- Gas estimation.
- Nonce management.
- Transaction signing.
- Transaction broadcasting.
- Receipt polling.
- Reorg-aware confirmation policies.
- Contract binding plumbing.

Forbidden dependencies:

- Blockchain must not evaluate strategies.
- Blockchain must not decide position size.
- Blockchain must not own exchange-independent order policy.

### Persistence

Persistence owns durable storage.

Responsibilities:

- Trade storage.
- Signal storage.
- Market storage.
- Position storage.
- Wallet state storage.
- Statistics storage.
- Idempotency keys.
- Migrations.

Forbidden dependencies:

- Persistence must not call exchanges.
- Persistence must not sign or send transactions.
- Persistence must not contain trading strategy rules.

### Configuration

Configuration owns environment-specific settings.

Responsibilities:

- Loading config from files and environment variables.
- Validating required values.
- Parsing typed values.
- Preventing hardcoded addresses, ABIs, credentials, or network constants.

Forbidden dependencies:

- Configuration must not execute business workflows.
- Configuration must not open long-lived network connections.

### Logging

Logging owns structured operational logs.

Responsibilities:

- Request and workflow correlation IDs.
- Structured fields.
- Error context.
- Redaction of secrets.

Forbidden dependencies:

- Logging must not drive control flow.
- Logging must not leak private keys, mnemonics, API keys, auth tokens, or full signed transactions.

### Metrics

Metrics owns quantitative operational telemetry.

Responsibilities:

- Counters.
- Gauges.
- Histograms.
- Latency measurements.
- Success and failure rates.
- Risk and exposure measurements.

Forbidden dependencies:

- Metrics must not decide behavior directly.
- Metrics must not contain business state that is not stored elsewhere.

### Notification

Notification owns operator alerts.

Responsibilities:

- Critical failure alerts.
- Circuit breaker alerts.
- Trade execution alerts.
- Risk limit alerts.
- Startup and shutdown alerts.

Forbidden dependencies:

- Notification must not be required for core correctness.
- Notification must not contain trading logic.

## Core Components

### Configuration

Loads, validates, and exposes typed settings. Configuration must be immutable after startup. Invalid configuration is a startup failure.

### Logger

Provides structured logging with stable field names. All workflow logs must include component, operation, exchange, chain, market identifier when available, and correlation ID when available.

### Blockchain Client

Provides EVM RPC operations behind an interface. It owns connectivity, request timeouts, chain ID checks, receipt retrieval, block queries, and RPC error normalization.

### Wallet

Owns signing identity and transaction signing. It must not expose private key material after construction. It signs only validated transaction requests from the executor.

### Contract Bindings

Provide typed access to exchange contracts. Bindings are infrastructure details used by exchange adapters and blockchain execution code. ABI files and generated bindings must come from versioned artifacts, not hardcoded strings.

### Exchange Adapter

Translates a specific exchange into the engine's common exchange contract. Buffer Finance, OptionBlitz, Limitless, and future EVM prediction markets are adapters.

### Market Service

Coordinates market discovery, market validation, market normalization, and market cache updates through exchange adapters.

### Strategy Engine

Runs strategies against normalized market snapshots. It produces signals only. It never executes trades.

### Signal Engine

Validates, enriches, deduplicates, persists, and routes strategy signals into risk evaluation.

### Risk Manager

Owns exposure, position sizing, max loss, market eligibility, portfolio constraints, emergency stop, and risk rejection reasons.

### Order Executor

Owns order construction, transaction creation, gas and nonce coordination, signing, broadcasting, receipt confirmation, retry, recovery, and idempotency.

### Position Manager

Owns live position state, reconciliation, lifecycle transitions, settlement status, and exchange-reported state normalization.

### Portfolio Manager

Owns account-level exposure, capital allocation, cross-market limits, cross-exchange limits, drawdown limits, and performance-aware allocation.

### Persistence

Stores durable facts. The engine must be able to restart and recover from persisted signals, orders, transactions, positions, and statistics.

### Statistics

Computes trading and operational statistics from persisted facts. Statistics must not be the source of truth for raw trades.

### Metrics

Exposes runtime telemetry for monitoring and alerting. Metrics must be low-cardinality and safe for production collection.

### Scheduler

Runs periodic and delayed workflows using contexts. It owns timing, not business decisions.

### Notifier

Sends operator notifications for significant events. Notification failure must be logged and measured but must not corrupt trading state.

### Paper Trading

Executes the full strategy, signal, risk, position, portfolio, statistics, logging, metrics, and notification lifecycle without signing or broadcasting transactions.

### Backtesting

Runs strategies and risk policies over historical market data through the same domain and application contracts used by live trading. Backtesting must not import live exchange or wallet implementations.

## Exchange Adapter Contract

Every exchange must implement exactly the same interface so the trading engine can treat Buffer Finance, OptionBlitz, Limitless, and future EVM prediction markets as interchangeable venues.

The strategy engine must never know the exchange implementation. A strategy receives normalized domain market data and emits normalized signals. It does not import Buffer packages, call Buffer contracts, inspect Buffer ABI, or branch on exchange-specific details.

Required adapter interfaces:

### Exchange Identity

Provides immutable metadata:

- Exchange name.
- Supported chain IDs.
- Supported market types.
- Version of the adapter.

### Market Discovery

Discovers active and upcoming markets from the exchange. It returns normalized domain market models with exchange-specific fields captured only in adapter-owned metadata when required.

### Market Validation

Validates whether a discovered market can be traded by the engine. Validation includes status, expiry, liquidity, supported direction, minimum size, maximum size, and exchange-specific constraints translated into normalized rejection reasons.

### Price Feed

Returns current price, implied probability, oracle state, freshness, and confidence for a normalized market.

### Liquidity Query

Returns available liquidity, size constraints, fee information, and expected slippage or payout model.

### Order Builder

Builds an exchange-specific order request from a normalized approved order. It does not sign or broadcast. It returns a transaction request or an execution payload consumed by the executor.

### Position Reader

Reads open, closed, expired, and settled positions for a wallet on that exchange and converts them into normalized domain positions.

### Settlement

Provides settlement eligibility and settlement transaction construction where the exchange requires explicit settlement.

### Health Check

Reports adapter health, contract reachability, oracle availability, and chain compatibility.

The adapter contract must hide exchange differences behind normalized domain objects. If a future exchange needs a capability that does not fit the contract, the architecture should add an explicit optional capability interface rather than leaking exchange-specific logic into strategies or application services.

## Trading Lifecycle

```text
Market discovered
↓
Validated
↓
Strategy evaluated
↓
Signal produced
↓
Risk validated
↓
Order built
↓
Transaction signed
↓
Broadcast
↓
Receipt confirmed
↓
Position updated
↓
Statistics updated
```

### Market Discovered

The market service asks exchange adapters for active and upcoming markets. The result is normalized into domain market models and tagged with exchange identity and chain identity.

### Validated

The market service validates status, expiry, liquidity, oracle freshness, supported direction, size constraints, and operational health. Invalid markets are rejected with explicit reasons and metrics.

### Strategy Evaluated

The strategy engine evaluates normalized market snapshots and historical context. It has no access to wallet, RPC, contracts, gas, nonce, or transaction details.

### Signal Produced

A strategy emits a signal that includes market ID, direction, confidence, time horizon, strategy identity, and optional metadata. Signals are persisted before execution decisions when required for auditability.

### Risk Validated

The risk manager evaluates signal eligibility, position size, exposure, drawdown, max loss, liquidity, market constraints, portfolio state, and emergency-stop status. Rejections are explicit and persisted.

### Order Built

The order executor converts an approved risk decision into a normalized order and asks the exchange adapter to build the exchange-specific execution payload.

### Transaction Signed

The executor coordinates gas and nonce, then asks the wallet to sign. Signing happens only after validation and idempotency checks.

### Broadcast

The signed transaction is sent through the blockchain client. Broadcast result, transaction hash, and errors are persisted.

### Receipt Confirmed

The executor waits for receipt confirmation according to chain-specific confirmation policy. It handles timeout, replacement, reorg, and revert states explicitly.

### Position Updated

The position manager records the new or updated position from transaction receipt data and exchange state reconciliation.

### Statistics Updated

Statistics are recalculated from persisted facts. Metrics are emitted. Notifications are sent for significant lifecycle events.

## Error Handling

Errors must be typed or wrapped with enough context to decide retryability, severity, and operator action. Errors must never be silently ignored.

### Network Errors

Network errors are retryable when transient. Retries must use bounded exponential backoff with jitter and context cancellation.

### RPC Unavailable

RPC unavailability triggers provider failover when configured. Repeated failures trip a circuit breaker for the affected chain/provider.

### Gas Estimation Failure

Gas estimation failure is not automatically fatal. The executor must classify whether the failure indicates contract revert, insufficient funds, unsupported method, RPC issue, or estimation instability.

### Nonce Conflicts

Nonce conflicts are handled by nonce refresh, transaction replacement policy, and idempotency checks. The executor owns all nonce coordination.

### Timeout

Timeouts must respect context deadlines. A timeout does not prove failure for broadcast transactions; the executor must reconcile by transaction hash, nonce, and chain state.

### Contract Revert

Contract reverts are classified into known exchange errors where possible. Reverts are persisted with reason, transaction hash, market ID, and order ID.

### Oracle Unavailable

Oracle unavailability rejects market validation or risk approval depending on where it is detected. Trading must stop for affected markets until freshness returns.

### Retry Policy

Retries must be bounded, observable, idempotent, and safe. No retry may duplicate an order or sign multiple conflicting transactions without executor-owned replacement policy.

### Circuit Breaker

Circuit breakers protect the system from repeated failures. Breakers may exist per RPC provider, exchange adapter, market, strategy, wallet, or entire engine. A tripped breaker must be logged, measured, and optionally notified.

## Concurrency

Concurrency must be explicit, bounded, observable, and cancellable.

### Worker Pools

Worker pools execute market discovery, strategy evaluation, reconciliation, and notification workloads. Pools must have bounded size and backpressure.

### Channels

Channels may connect internal stages when streaming is useful. Channel ownership must be clear. Producers close channels. Consumers must handle closure.

### Contexts

Every external call and long-running workflow accepts `context.Context`. Context carries cancellation and deadlines, not optional dependencies.

### Cancellation

Cancellation must stop new work, unblock waiting operations, and allow in-flight critical sections to finish or roll back safely.

### Graceful Shutdown

Shutdown stops schedulers, closes intake channels, waits for workers, persists in-flight state, flushes logs and metrics, and sends final notifications where configured.

## Folder Structure

The repository should evolve toward this structure:

```text
/cmd
  /engine
/internal
  /app
  /domain
  /exchange
    /buffer
    /optionblitz
    /limitless
  /blockchain
  /wallet
  /persistence
  /config
  /logging
  /metrics
  /notification
  /scheduler
  /paper
  /backtest
/pkg
  /contracts
/configs
/abi
/migrations
/testdata
/docs
```

### `/cmd`

Owns process entry points. Code here wires dependencies and starts applications.

### `/internal/app`

Owns use-case orchestration and application services.

### `/internal/domain`

Owns domain models, value objects, domain errors, and invariants.

### `/internal/exchange`

Owns exchange adapter interfaces and concrete exchange adapters. The shared adapter contract belongs at the top of this area. Buffer, OptionBlitz, and Limitless implementations live in separate subpackages.

### `/internal/blockchain`

Owns EVM RPC, gas, nonce, transaction, receipt, and confirmation services.

### `/internal/wallet`

Owns signing abstractions and wallet implementations.

### `/internal/persistence`

Owns repositories, database models, migrations integration, and durable state mapping.

### `/internal/config`

Owns typed configuration structures and validation.

### `/internal/logging`

Owns logger construction and shared logging field conventions.

### `/internal/metrics`

Owns metric instruments and exporters.

### `/internal/notification`

Owns notifier interfaces and provider implementations.

### `/internal/scheduler`

Owns recurring task orchestration.

### `/internal/paper`

Owns paper-trading execution adapters.

### `/internal/backtest`

Owns historical simulation orchestration.

### `/pkg/contracts`

Owns generated or shared contract-related packages only when they are safe for external import. Prefer `/internal` unless external consumption is required.

### `/configs`

Owns example and environment-specific configuration files without secrets.

### `/abi`

Owns versioned ABI artifacts. ABI content must not be embedded as hardcoded strings in Go code.

### `/migrations`

Owns database schema migrations.

### `/testdata`

Owns deterministic test fixtures.

### `/docs`

Owns supplementary documentation. `ARCHITECTURE.md` remains the authority.

## Coding Standards

### Naming

Names must be precise, boring, and domain-aligned. Avoid abbreviations except standard terms such as RPC, ABI, EVM, ID, URL, and HTTP.

### Interfaces

Interfaces should be small and behavior-focused. Define interfaces at consumer boundaries when possible. Large god interfaces are forbidden.

### Errors

Errors must be returned, not ignored. Wrap errors with operation context. Use typed errors where callers need classification. Do not log and return the same error at every layer; log at workflow boundaries.

### Logging

Use structured logs. Include stable fields. Redact secrets. Logs must explain what failed, where it failed, and the operational impact.

### Configuration

No hardcoded addresses, ABIs, RPC URLs, private keys, chain IDs for runtime behavior, or strategy parameters. Configuration must be validated before the engine starts.

### Tests

Every module requires unit tests for core behavior and error paths. Integration tests must isolate external systems through test fixtures, test networks, mocks, or containers. Tests must be deterministic unless explicitly marked as integration tests.

### Comments

Comments explain why, not what. Public exported Go identifiers require useful documentation. Do not use comments to excuse unclear code.

## Future Roadmap

### Buffer

Implement the first production exchange adapter for Buffer Finance on Arbitrum. Establish the adapter contract, transaction lifecycle, market normalization, risk integration, and position reconciliation.

### OptionBlitz

Add OptionBlitz as a second adapter without changing strategy, risk, execution, or persistence contracts except through deliberate interface extension.

### Limitless

Add Limitless as an adapter and validate whether optional prediction-market capabilities are needed.

### Cross-Exchange Arbitrage

Introduce portfolio-aware multi-exchange execution where strategies can compare normalized markets across venues while execution remains adapter-driven.

### Machine Learning

Add ML-driven strategies as strategy implementations that consume normalized feature sets and produce signals only.

### LLM

Use LLM components only for research, summarization, operator assistance, or offline analysis unless explicit production controls are designed. LLMs must not bypass risk management or executor approval.

### Portfolio Optimization

Add portfolio optimization as a risk and allocation module that controls capital distribution across strategies, exchanges, market types, and time horizons.
