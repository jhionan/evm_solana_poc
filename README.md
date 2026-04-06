# Multi-Chain Staking POC

A production-grade proof-of-concept demonstrating cross-chain staking across **Ethereum (EVM)** and **Solana**, with tiered rewards, slashing penalties, catastrophe recovery, and a unified gRPC API surface.

---

## Architecture

```
                         ┌─────────────────────────────────────────────────┐
                         │                  Client Layer                    │
                         │     grpcui / mobile app / CLI / web frontend     │
                         └──────────────────────┬──────────────────────────┘
                                                │  HTTP/2 + gRPC (ConnectRPC)
                         ┌──────────────────────▼──────────────────────────┐
                         │              StakingService (Go)                 │
                         │  • Validate requests                             │
                         │  • Enforce tiered staking rules (Bronze→Diamond) │
                         │  • Apply penalty / slashing logic                │
                         │  • Publish domain events to Valkey               │
                         └───────────┬───────────────────┬─────────────────┘
                                     │                   │
              ┌──────────────────────▼──┐       ┌────────▼──────────────────┐
              │       EVMStaker          │       │      SolanaStaker          │
              │  (go-ethereum / ethclient)│      │  (gagliardetto/solana-go)  │
              │  • Stake / unstake txs   │       │  • Stake / unstake txs     │
              │  • Gas estimation        │       │  • Compute unit budgeting  │
              │  • EIP-1559 fee strategy │       │  • Priority fee strategy   │
              └──────────┬──────────────┘       └─────────┬──────────────────┘
                         │                                 │
              ┌──────────▼─────────────────────────────────▼──────────────┐
              │                    EventIndexer (Go)                        │
              │  • Polls Anvil / Solana validator for new blocks            │
              │  • Decodes StakeDeposited / Unstaked / Slashed events       │
              │  • Writes canonical records into PostgreSQL                 │
              └──────────────────────────┬─────────────────────────────────┘
                                         │
                         ┌───────────────▼───────────────┐
                         │       PostgreSQL 16            │
                         │  stakes • events • snapshots   │
                         │  positions • audit_log         │
                         └───────────────────────────────┘
```

**Supporting infrastructure:**

| Component | Role |
|-----------|------|
| Anvil (Foundry) | Local EVM node, 2-second block time |
| Solana Test Validator | Local Solana validator |
| Valkey 8 | Pub/sub for domain events; rate-limit counters |
| sqlc | Compile-time safe SQL → Go |
| buf / protoc | Proto compilation & linting |
| Goose | Deterministic DB migrations |

---

## Quick Start

### Prerequisites

| Tool | Version | Install |
|------|---------|---------|
| Docker + Compose | v2.x | https://docs.docker.com/get-docker/ |
| Go | 1.23+ | https://go.dev/dl/ |
| Foundry (forge, anvil) | latest | `curl -L https://foundry.paradigm.xyz \| bash` |
| Anchor CLI | 0.30+ | https://www.anchor-lang.com/docs/installation |
| Buf | 1.x | https://buf.build/docs/installation |
| sqlc | 1.x | `go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest` |
| Goose | 3.x | `go install github.com/pressly/goose/v3/cmd/goose@latest` |
| golangci-lint | 1.x | https://golangci-lint.run/usage/install/ |
| grpcui | latest | `go install github.com/fullstorydev/grpcui/cmd/grpcui@latest` |

### 1. Start the full stack

```bash
make up
```

This starts PostgreSQL, Valkey, Anvil, and the Solana validator, waits for health checks, then applies all DB migrations automatically.

### 2. Generate code

```bash
make proto    # compile .proto → Go stubs
make sqlc     # compile SQL → type-safe Go
```

### 3. Build the server

```bash
make build
```

Binary is written to `backend/bin/server`.

### 4. Deploy contracts to local nodes

```bash
# Ethereum (Anvil)
make deploy-local

# Solana
make test-solana   # anchor test deploys + exercises the program
```

### 5. Run the server

```bash
./backend/bin/server
```

The gRPC server listens on `:8080` by default.

### 6. Explore with the browser UI

```bash
make demo
```

Opens `grpcui` in your browser — no client code required.

---

## Tech Stack

| Layer | Technology | Purpose |
|-------|-----------|---------|
| API | ConnectRPC + protobuf | Unified gRPC + HTTP/1.1 + gRPC-Web |
| Backend | Go 1.23 | Business logic, indexer, adapters |
| EVM contracts | Solidity 0.8 + Foundry | Staking vault, penalty ledger |
| Solana program | Anchor / Rust | Native Solana staking program |
| Database | PostgreSQL 16 | Canonical event store, position ledger |
| Cache / events | Valkey 8 | Pub/sub, idempotency keys, rate limits |
| Migrations | Goose | Sequential, versioned schema changes |
| SQL codegen | sqlc | Zero-reflection, compile-time DB access |
| Proto toolchain | buf | Linting, breaking-change detection |
| Testing (EVM) | Forge + fuzz | Unit + property-based contract tests |
| Testing (Solana) | Anchor test | TypeScript integration tests |
| Observability | zerolog + structured JSON | Machine-readable logs, trace IDs |

---

## Key Features

- **Tiered staking** — Bronze / Silver / Gold / Diamond tiers with distinct APY multipliers and lock-up rules enforced on-chain and mirrored in the service layer.
- **Slashing & penalties** — Configurable penalty rates applied automatically on early withdrawal or validator misbehaviour, with full audit trail in PostgreSQL.
- **Chain adapter pattern** — `ChainStaker` interface decouples business logic from EVM and Solana specifics; adding a new chain requires implementing one interface.
- **Event indexer** — Background goroutine polls each chain for new blocks, decodes contract events, and writes idempotent records to the event store.
- **Catastrophe recovery** — Snapshot-based recovery mechanism; the indexer can rebuild all position state from on-chain events from genesis, providing a verifiable audit trail even after a complete DB loss.
- **Idempotent operations** — Every stake/unstake request carries a client-generated idempotency key stored in Valkey, preventing duplicate on-chain transactions from retries.
- **Security first** — Input validation at the proto layer, replay protection, rate limiting per address via Valkey, and no private keys stored in the service (external signer / HSM integration points provided).

---

## Running Tests

```bash
# Go unit tests
make test

# Go tests with HTML coverage report
make test-cover

# EVM contract tests (verbose)
make test-contracts

# EVM fuzz tests (1 000 runs)
make test-contracts-fuzz

# Solana program tests
make test-solana

# Everything at once
make test-all
```

---

## Demo Flow

The following walkthrough exercises the full system end-to-end:

1. **Spin up** — `make up` starts all four services and applies migrations.
2. **Deploy** — `make deploy-local` deploys the Solidity staking vault to Anvil; `make test-solana` deploys the Anchor program to the local validator.
3. **Stake (Bronze)** — Call `StakingService.Stake` with 0.1 ETH via `grpcui`; observe the transaction mined on Anvil within 2 seconds.
4. **Upgrade tier** — Stake additional ETH to cross the Gold threshold; confirm the tier upgrade event appears in `events` table.
5. **Index events** — Watch the EventIndexer goroutine decode `StakeDeposited` logs and populate `positions` in real time.
6. **Unstake with penalty** — Unstake before the lock-up expires; verify the penalty is deducted on-chain and recorded in `audit_log`.
7. **Catastrophe recovery** — Stop the service, drop the `positions` table, restart; the indexer replays all on-chain events from block 0 and fully reconstructs the ledger.
8. **Solana round-trip** — Repeat steps 3–5 against the Solana validator via the same `StakingService` RPC — zero client-side changes required.

---

## Decision Records

Architecture decisions, trade-offs, and chain selection rationale are documented in:

```
../etherium_poc_docs/
```

---

## Project Layout

```
etherium_poc/
├── backend/
│   ├── cmd/server/         # Entrypoint
│   ├── db/migrations/      # Goose SQL migrations
│   ├── gen/                # sqlc-generated Go (do not edit)
│   ├── internal/           # Domain logic, adapters, indexer
│   ├── pkg/                # Shared utilities (errors, logger, config)
│   └── proto/              # Protobuf definitions + buf config
├── contracts/
│   ├── evm/                # Foundry project (Solidity)
│   └── solana/             # Anchor project (Rust)
├── docker-compose.yml
├── Makefile
└── README.md
```

---

*Built as a portfolio piece demonstrating production-grade multi-chain infrastructure design.*
