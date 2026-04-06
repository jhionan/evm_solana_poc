# Multi-Chain Staking POC

A production-grade proof-of-concept demonstrating cross-chain staking across **Ethereum (EVM)** and **Solana**, with tiered rewards, early withdrawal penalties, catastrophe recovery, and a unified gRPC API.

---

## Architecture

```
                         ┌─────────────────────────────────────────────────┐
                         │                  Client Layer                    │
                         │      grpcui / buf curl / web / mobile            │
                         └──────────────────────┬──────────────────────────┘
                                                │  ConnectRPC (gRPC + HTTP/JSON)
                         ┌──────────────────────▼──────────────────────────┐
                         │          Interceptor Stack                       │
                         │  1. Security Headers (HSTS, CSP, XFO)           │
                         │  2. Rate Limiter (Valkey-backed, per wallet)     │
                         │  3. JWT Auth (HS256, skip public endpoints)      │
                         │  4. Audit Logger (hash chain, mutating only)     │
                         └──────────────────────┬──────────────────────────┘
                                                │
                         ┌──────────────────────▼──────────────────────────┐
                         │              StakingService (Go)                 │
                         │  • Routes by chain field                         │
                         │  • Constructor-injected adapters                 │
                         │  • Typed errors → ConnectRPC codes               │
                         └───────────┬───────────────────┬─────────────────┘
                                     │                   │
              ┌──────────────────────▼──┐       ┌────────▼──────────────────┐
              │       EVMStaker          │       │      SolanaStaker          │
              │  (go-ethereum bindings)  │       │  (gagliardetto/solana-go)  │
              │  • approve → stake tx    │       │  • PDA-based positions     │
              │  • Parse Staked events   │       │  • CPI token transfers     │
              │  • Tx signing (isolated) │       │  • 2 tiers (Bronze/Gold)   │
              └──────────┬──────────────┘       └─────────┬──────────────────┘
                         │                                 │
              ┌──────────▼─────────────────────────────────▼──────────────┐
              │                    EventIndexer (Go)                        │
              │  • Catch-up from last indexed block on startup             │
              │  • Live WebSocket subscription for new events              │
              │  • Idempotent writes (tx_hash + log_index dedup)           │
              │  • Chain is source of truth — DB fully rebuildable         │
              └──────────────────────────┬─────────────────────────────────┘
                                         │
                         ┌───────────────▼───────────────┐
                         │       PostgreSQL 18            │
                         │  positions • chain_events      │
                         │  block_cursors • rewards        │
                         │  audit_log (hash chain)         │
                         └───────────────────────────────┘
```

---

## Live Testnet Deployments

### Ethereum Sepolia (Chain ID `11155111`)

| Contract | Address | Explorer |
|----------|---------|----------|
| **StakingToken (LLSTK)** | `0xD8315d4a129d6dEb6A54D8C42445911AF3aD8F64` | [View on Etherscan](https://sepolia.etherscan.io/address/0xD8315d4a129d6dEb6A54D8C42445911AF3aD8F64) |
| **TieredStaking** | `0x1D1Ca3770e7366EA6A5431a6185eD4952F570470` | [View on Etherscan](https://sepolia.etherscan.io/address/0x1D1Ca3770e7366EA6A5431a6185eD4952F570470) |
| **Deployer** | `0x239f3d3F6FFBC5A01d1109d33BDdE0Bb12325FB7` | [View on Etherscan](https://sepolia.etherscan.io/address/0x239f3d3F6FFBC5A01d1109d33BDdE0Bb12325FB7) |

### Solana Devnet

| Component | Address | Explorer |
|-----------|---------|----------|
| **Tiered Staking Program** | `zaEDLTCw2b5zGrunQEDKvHuVEEgFoAwjHViMH357QGn` | [View on Solscan](https://solscan.io/account/zaEDLTCw2b5zGrunQEDKvHuVEEgFoAwjHViMH357QGn?cluster=devnet) |
| **IDL Account** | `AERgQFi4B5wrgg2NbMQyS49TKtV3b53pJ96wECGBst3P` | [View on Solscan](https://solscan.io/account/AERgQFi4B5wrgg2NbMQyS49TKtV3b53pJ96wECGBst3P?cluster=devnet) |
| **Deploy Signature** | `3b3iyyrDefPHx7Gv444QiYS3ybsxgzhSuzRp45zE5DkZYgFUToYqU9iSUwxhcXtXiUhucf9Kayk2cAxBNeNmCsow` | [View on Solscan](https://solscan.io/tx/3b3iyyrDefPHx7Gv444QiYS3ybsxgzhSuzRp45zE5DkZYgFUToYqU9iSUwxhcXtXiUhucf9Kayk2cAxBNeNmCsow?cluster=devnet) |
| **Upgrade Authority** | `8ZAapig8xZLbTKDSi2ehMseu6eYcPqzZJX7BkFXKMfmK` | [View on Solscan](https://solscan.io/account/8ZAapig8xZLbTKDSi2ehMseu6eYcPqzZJX7BkFXKMfmK?cluster=devnet) |

### Sepolia E2E — Stake Transaction (verified on-chain)

```json
{
  "position": {
    "id": "0",
    "chain": "CHAIN_EVM",
    "wallet": "0x239f3d3F6FFBC5A01d1109d33BDdE0Bb12325FB7",
    "amount": "100000000000000000000",
    "tier": "TIER_GOLD",
    "status": "POSITION_STATUS_ACTIVE",
    "stakedAt": "1775493480",
    "lockUntil": "1783269480"
  },
  "txHash": "0xd4dba2b950f60d807c4174a1f08615320e21281d8791f9f666c0bebcb44d3f57"
}
```

Verify: [View tx on Etherscan](https://sepolia.etherscan.io/tx/0xd4dba2b950f60d807c4174a1f08615320e21281d8791f9f666c0bebcb44d3f57)

### How to Reproduce on Sepolia

```bash
# 1. Get Sepolia ETH from faucet
#    https://cloud.google.com/application/web3/faucet/ethereum/sepolia

# 2. Deploy contracts
PRIVATE_KEY=0xYOUR_KEY make deploy-sepolia

# 3. Update .env with contract addresses and run
EVM_RPC_URL=https://ethereum-sepolia-rpc.publicnode.com \
EVM_DEPLOY_BLOCK=10603490 \
EVM_STAKING_CONTRACT=0x1D1Ca3770e7366EA6A5431a6185eD4952F570470 \
EVM_TOKEN_CONTRACT=0xD8315d4a129d6dEb6A54D8C42445911AF3aD8F64 \
make build && ./backend/bin/server
```

---

## E2E Test Results (Local Anvil)

The following is a verified end-to-end test run against a local Anvil chain with real on-chain transactions:

### 1. GetTiers (public, no auth required)

```json
{"tier":"TIER_BRONZE","aprBps":500}
{"tier":"TIER_SILVER","aprBps":1000}
{"tier":"TIER_GOLD","aprBps":1800}
```

### 2. Stake 1000 STK in Gold tier

```json
{
  "positionId": "1",
  "txHash": "0x69d9877217379327b116fe4ca47265842959905a4cf74320de8804bf2e060158",
  "tier": "TIER_GOLD",
  "status": "POSITION_STATUS_ACTIVE"
}
```

### 3. GetPosition (reads from chain)

```json
{
  "id": "0",
  "wallet": "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
  "amount": "1000000000000000000000",
  "tier": "TIER_GOLD",
  "status": "POSITION_STATUS_ACTIVE"
}
```

### 4. ClaimRewards (accrued in ~3 seconds at 18% APR)

```json
{
  "rewardsClaimed": "1638127853881278",
  "txHash": "0xb97201c0cbf5c21b795e6e26c6048fccc606909409f56ea15c834db8c8cada27"
}
```

### 5. Unstake (early withdrawal — 10% penalty to treasury)

```json
{
  "amountReturned": "1000000000000000000000",
  "rewards": "0",
  "penalty": "100000000000000000000",
  "txHash": "0xdc39063ff67047d14db5f0721dc3f12e80937f683b3314607fec5f2a6f23684a"
}
```

### 6. Indexed Events (PostgreSQL)

```
   event_type   | block_number
----------------+--------------
 Staked         |         2771
 Staked         |         2912
 RewardsClaimed |         2914
 Unstaked       |         2915
```

All 4 events caught by the live WebSocket subscription and persisted idempotently.

### 7. Audit Log with Hash Chain

```
    action    |                   actor                    | chain_id |   hash_prefix    | prev_hash_prefix
--------------+--------------------------------------------+----------+------------------+------------------
 Stake        | 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 | evm      | 5c0ca4f702343337 |
 ClaimRewards | 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 | evm      | e6100c5b359d3d55 | 5c0ca4f702343337
 Unstake      | 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266 | evm      | 5e0b12a1174e2731 | e6100c5b359d3d55
```

Each entry's `prev_hash` matches the previous entry's `hash` — tamper-evident chain verified.

### 8. Test Suite

```
=== Go Backend (11 packages) ===
ok  internal/api          — ConnectRPC handler + proto conversion
ok  internal/audit        — Hash chain logger + interceptor (9 tests)
ok  internal/auth         — JWT sign/verify + RBAC (4 tests)
ok  internal/chain/evm    — EVMStaker adapter (10 tests)
ok  internal/chain/solana — SolanaStaker adapter (9 tests)
ok  internal/config       — Viper config + validation (5 tests)
ok  internal/indexer      — Catch-up + live subscription (4 tests)
ok  internal/security     — Validation + rate limiter (8 tests)
ok  internal/staking      — Service + chain routing (10 tests)
ok  pkg/errors            — AppError + ConnectRPC mapping (8 tests)
ok  pkg/middleware         — Security headers (2 tests)

=== Solidity (Foundry) ===
StakingTokenTest:  6 passed, 0 failed (includes fuzz)
TieredStakingTest: 15 passed, 0 failed (includes fuzz)

=== Solana (Anchor) ===
cargo check: 0 errors
```

---

## Quick Start

### Prerequisites

| Tool | Version | Install |
|------|---------|---------|
| Docker + Compose | v2.x | https://docs.docker.com/get-docker/ |
| Go | 1.26+ | https://go.dev/dl/ |
| Foundry (forge, anvil) | latest | `curl -L https://foundry.paradigm.xyz \| bash` |
| Buf | latest | `brew install bufbuild/buf/buf` |
| sqlc | latest | `brew install sqlc` |
| Goose | latest | `go install github.com/pressly/goose/v3/cmd/goose@latest` |

### 1. Start infrastructure

```bash
make up     # PostgreSQL, Valkey, Anvil, Solana validator + migrations
```

### 2. Deploy contracts

```bash
make deploy-local   # Deploy to local Anvil
```

Copy the printed contract addresses into your `.env` file (see `.env.example`).

### 3. Run the server

```bash
cp .env.example .env
# Edit .env with contract addresses from step 2
make build
./backend/bin/server
```

### 4. Test with grpcui

```bash
make demo   # Opens browser-based RPC explorer
```

---

## Tech Stack

| Layer | Technology | Purpose |
|-------|-----------|---------|
| API | ConnectRPC + Buf | gRPC + gRPC-Web + HTTP/JSON from single handler |
| Backend | Go 1.26 | Business logic, chain adapters, indexer |
| EVM Contracts | Solidity 0.8.28 + Foundry | Tiered staking with OpenZeppelin security |
| Solana Program | Anchor 0.30 / Rust | PDA-based staking with CPI token transfers |
| Database | PostgreSQL 18 + sqlc | Type-safe queries, event store, audit log |
| Cache | Valkey 8 | Rate limiting counters (Redis-compatible) |
| Auth | JWT (HS256) | Wallet-based authentication |
| Migrations | Goose | Sequential, versioned schema changes |
| Testing | Forge (fuzz) + testify | Property-based + unit tests |
| Logging | zerolog | Structured JSON logs |

---

## Key Features

- **Tiered staking** — Bronze (30d / 5% APR), Silver (60d / 10%), Gold (90d / 18%) with lock periods enforced on-chain.
- **Early withdrawal penalty** — 10% of staked amount sent to treasury on early unstake, enforced in Solidity with `SafeERC20`.
- **Chain adapter pattern** — `ChainStaker` interface decouples business logic from chain specifics. Adding a new chain = implement one interface.
- **Event indexer** — Catches up from last indexed block on startup, then subscribes via WebSocket for live events. Idempotent writes via `(tx_hash, log_index)` unique constraint.
- **Catastrophe recovery** — DB is a cache of on-chain state. Truncate all tables, restart → indexer rebuilds everything from chain events.
- **Audit log with hash chain** — Every mutating operation (Stake, Unstake, ClaimRewards) is logged with a SHA-256 hash chained to the previous entry. Tamper-evident.
- **Security from day 0** — ReentrancyGuard, Ownable2Step, Pausable, Checks-Effects-Interactions (Solidity). JWT auth, RBAC, rate limiting, security headers, typed errors (Go).
- **Transaction signing isolation** — Private keys never leave the signer module. API layer has no access to key material.

---

## Smart Contracts

### EVM — TieredStaking.sol

| Tier | Lock Period | APR (bps) | Early Penalty |
|------|-----------|-----------|---------------|
| Bronze | 30 days | 500 (5%) | 10% to treasury |
| Silver | 60 days | 1000 (10%) | 10% to treasury |
| Gold | 90 days | 1800 (18%) | 10% to treasury |

Security: `ReentrancyGuard` + `Ownable2Step` + `Pausable` + `SafeERC20` + Checks-Effects-Interactions

### Solana — Tiered Staking (Anchor)

| Tier | Lock Period | APR (bps) |
|------|-----------|-----------|
| Bronze | 30 days | 500 (5%) |
| Gold | 90 days | 1800 (18%) |

PDA-based position tracking, CPI token transfers, event emission for indexer.

**Adapter status:** The Go `SolanaStaker` adapter implements the `ChainStaker` interface with `GetTiers` and `HealthCheck` fully wired. Write operations (Stake, Unstake, ClaimRewards) are scaffolded — the Anchor program is deployed and the adapter pattern is proven via EVM; wiring the Solana RPC calls follows the same mechanical pattern. See [scope tradeoffs ADR](../etherium_poc_docs/03-solana-integration/decisions/scope-tradeoffs.md).

---

## Running Tests

```bash
make test              # Go unit tests (11 packages)
make test-cover        # Go tests with HTML coverage report
make test-contracts    # Solidity tests (21 tests, includes fuzz)
make test-all          # Everything
```

---

## Decision Records

Architecture decisions documented in ADR format:

```
../etherium_poc_docs/
├── 00-strategy/          — Why this POC, interview strategy
├── 01-smart-contracts/   — Foundry vs Hardhat, tiered design, security
├── 02-go-backend/        — ConnectRPC vs REST, chain adapter pattern, indexing
├── 03-solana-integration/— Anchor, EVM vs Solana model, scope tradeoffs
├── 04-multi-chain-orchestration/ — Unified interface, catastrophe recovery
└── 05-retrospective/     — Learnings, production recommendations
```

---

## Project Layout

```
etherium_poc/
├── backend/
│   ├── cmd/server/            # Entry point, full wiring
│   ├── db/
│   │   ├── migrations/        # Goose SQL migrations (5 tables)
│   │   ├── queries/           # sqlc query definitions
│   │   └── sqlc/              # Generated type-safe Go
│   ├── gen/staking/v1/        # Generated ConnectRPC + protobuf
│   ├── internal/
│   │   ├── api/               # ConnectRPC handler (6 RPCs)
│   │   ├── audit/             # Hash chain logger + PG adapter + interceptor
│   │   ├── auth/              # JWT, RBAC, ConnectRPC interceptor
│   │   ├── chain/
│   │   │   ├── evm/           # EVMStaker (real contract bindings)
│   │   │   └── solana/        # SolanaStaker
│   │   ├── config/            # Viper config with validation
│   │   ├── indexer/           # Engine + EVM source + PG store
│   │   ├── security/          # Validation + rate limiter
│   │   ├── signer/            # Tx signing (isolated key material)
│   │   └── staking/           # Service with chain routing
│   ├── pkg/
│   │   ├── errors/            # AppError sentinels → ConnectRPC codes
│   │   └── middleware/        # Security headers
│   └── proto/staking/v1/      # Protobuf service definition
├── contracts/
│   ├── evm/                   # Foundry (Solidity 0.8.28)
│   │   ├── src/               # StakingToken + TieredStaking
│   │   ├── test/              # 21 Forge tests (unit + fuzz)
│   │   └── script/            # Deploy script
│   └── solana/                # Anchor (Rust)
│       └── programs/tiered-staking/
├── docker-compose.yml         # PostgreSQL 18, Valkey 8, Anvil, Solana
├── Makefile                   # 18 dev targets
└── README.md
```
