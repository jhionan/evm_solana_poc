# Multi-Chain Staking POC — Design Spec

**Date:** 2026-04-06
**Author:** Jhionan Rian Santos
**Status:** Approved
**Timeline:** 7 days (Apr 7–13, 2026)

---

## 1. Purpose

Build a production-quality proof of concept that demonstrates senior-level capability across the full LightLink stack: EVM smart contracts (Solidity), Go backend, and multi-chain expansion to Rust-based chains (Solana). This serves as a portfolio piece for an active interview with LightLink (lightlink.io) for a Senior Software Engineer role.

### What This Proves to LightLink

1. Production-quality Solidity with security best practices (their current stack)
2. Go backend with clean abstractions for chain interaction (their backend language)
3. Multi-chain extensibility to Rust-based chains like Solana (their roadmap)
4. Architectural thinking: chain adapter pattern, catastrophe recovery, security from day 0

---

## 2. Tech Stack (All Open Source, Latest Stable)

| Layer | Tool | Version | Why |
|-------|------|---------|-----|
| **EVM Contracts** | Solidity | 0.8.28 | Latest stable, transient storage support |
| **EVM Toolchain** | Foundry | Latest | Rust-based, fastest Solidity tooling |
| **EVM Libraries** | OpenZeppelin 5.x | 5.1+ | Audited ERC-20, ReentrancyGuard, Pausable |
| **Solana Contracts** | Anchor | 0.30+ | Dominant Solana framework, Rust-based |
| **Solana CLI** | Solana CLI | 1.18+ | Latest stable toolchain |
| **Go** | 1.26 | Latest | Primary backend language |
| **RPC Framework** | ConnectRPC (connect-go) | Latest | gRPC + gRPC-Web + Connect protocol |
| **Protobuf Tooling** | Buf CLI | Latest | Modern proto management, linting, codegen |
| **EVM Go Client** | go-ethereum (geth) | Latest | Official Ethereum Go library |
| **Solana Go Client** | gagliardetto/solana-go | Latest | Most active open-source Solana Go client |
| **Logging** | zerolog | Latest | Structured logging |
| **Config** | Viper | Latest | Standard Go config management |
| **Database** | PostgreSQL | 16+ | Indexed events, positions, block cursors, audit log |
| **DB Queries** | sqlc | Latest | Type-safe generated Go from SQL |
| **Migrations** | Goose | Latest | SQL migration management |
| **Cache / Rate Limiting** | Valkey | 8+ | Open-source Redis fork, rate limits + caching |
| **Auth** | golang-jwt/jwt/v5 | Latest | JWT authentication |
| **API Testing** | grpcui | Latest | Interactive browser UI for gRPC |
| **Testing (Go)** | testify | Latest | Assertions + mocks |
| **Testing (Solidity)** | Forge test | Built-in | Foundry native, tests in Solidity |
| **Local EVM** | Anvil | Built-in | Foundry's local chain |
| **Local Solana** | solana-test-validator | Built-in | Local Solana cluster |
| **Containers** | Docker Compose | Latest | Orchestrate local chains + backend |
| **CI** | GitHub Actions | N/A | Free, open-source workflows |

---

## 3. Architecture — Chain Adapter Pattern

### Overview

The Go backend uses a `ChainStaker` interface. Each chain gets its own adapter implementing the interface. A `StakingService` orchestrates across chains and routes requests based on the `chain` field in each RPC call.

### Project Structure

```
etherium_poc/
├── contracts/
│   ├── evm/              # Foundry project (Solidity)
│   │   ├── src/
│   │   │   ├── StakingToken.sol
│   │   │   └── TieredStaking.sol
│   │   ├── test/
│   │   ├── script/
│   │   └── foundry.toml
│   └── solana/            # Anchor project (Rust)
│       ├── programs/
│       │   └── tiered-staking/
│       ├── tests/
│       └── Anchor.toml
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go
│   │   ├── auth/
│   │   │   ├── jwt.go
│   │   │   ├── interceptor.go
│   │   │   └── rbac.go
│   │   ├── security/
│   │   │   ├── ratelimit.go
│   │   │   ├── headers.go
│   │   │   └── validation.go
│   │   ├── signer/
│   │   │   ├── signer.go
│   │   │   ├── evm_signer.go
│   │   │   └── solana_signer.go
│   │   ├── audit/
│   │   │   ├── logger.go
│   │   │   └── middleware.go
│   │   ├── chain/
│   │   │   ├── staker.go
│   │   │   ├── types.go
│   │   │   ├── evm/
│   │   │   │   ├── adapter.go
│   │   │   │   ├── bindings/
│   │   │   │   └── adapter_test.go
│   │   │   └── solana/
│   │   │       ├── adapter.go
│   │   │       ├── idl/
│   │   │       └── adapter_test.go
│   │   ├── indexer/
│   │   │   ├── indexer.go
│   │   │   ├── evm_source.go
│   │   │   ├── solana_source.go
│   │   │   └── indexer_test.go
│   │   ├── staking/
│   │   │   ├── service.go
│   │   │   └── service_test.go
│   │   └── api/
│   │       ├── handler.go
│   │       └── handler_test.go
│   ├── db/
│   │   ├── migrations/
│   │   ├── queries/
│   │   └── sqlc/
│   ├── proto/
│   │   ├── buf.yaml
│   │   ├── buf.gen.yaml
│   │   └── staking/v1/
│   │       └── staking.proto
│   └── pkg/
│       ├── errors/
│       └── middleware/
├── Makefile
├── docker-compose.yml
└── README.md
```

### Core Interface

```go
type ChainStaker interface {
    // Write operations
    Stake(ctx context.Context, req StakeRequest) (*StakeReceipt, error)
    Unstake(ctx context.Context, positionID string) (*UnstakeReceipt, error)
    ClaimRewards(ctx context.Context, positionID string) (*ClaimReceipt, error)

    // Read operations
    GetPosition(ctx context.Context, positionID string) (*StakePosition, error)
    ListPositions(ctx context.Context, wallet string) ([]StakePosition, error)
    GetTiers(ctx context.Context) ([]Tier, error)

    // Chain info
    ChainID() ChainType
    HealthCheck(ctx context.Context) error
}
```

### Request Flow

```
Client (grpcui / buf curl / frontend)
    │
    ▼
ConnectRPC Handler (api/handler.go)
    │  deserializes proto, validates input
    ▼
StakingService (staking/service.go)
    │  reads chain field, selects adapter
    │  applies business logic
    ▼
ChainStaker (chain/staker.go interface)
    ├── EVMStaker ──► go-ethereum ──► Anvil / EVM chain
    └── SolanaStaker ──► solana-go ──► solana-test-validator
    │
    ▼
EventIndexer (indexer/indexer.go)
    │  subscribes to chain events
    │  writes to PostgreSQL
    ▼
PostgreSQL (db/)
    │  positions, events, block cursors
    ▼
StakingService reads DB for queries
    │  (chain is source of truth, DB is fast cache)
    ▼
Response back to client
```

### Protobuf Service Definition

```protobuf
service StakingService {
  rpc Stake(StakeRequest) returns (StakeResponse);
  rpc Unstake(UnstakeRequest) returns (UnstakeResponse);
  rpc ClaimRewards(ClaimRewardsRequest) returns (ClaimRewardsResponse);
  rpc GetPosition(GetPositionRequest) returns (GetPositionResponse);
  rpc ListPositions(ListPositionsRequest) returns (ListPositionsResponse);
  rpc GetTiers(GetTiersRequest) returns (GetTiersResponse);
}
```

Each RPC accepts a `chain` field so the caller specifies which chain to interact with.

---

## 4. Smart Contracts

### EVM — Solidity (Foundry)

**StakingToken.sol** — Simple ERC-20 token (mintable by owner, for testing).

**TieredStaking.sol** — Core staking logic:

| Tier | Lock Period | APR | 
|------|-----------|-----|
| Bronze | 30 days | 5% |
| Silver | 60 days | 10% |
| Gold | 90 days | 18% |

**Functions:**
- `stake(amount, tier)` — locks tokens, records position
- `unstake(positionId)` — returns tokens + rewards (if lock expired) or tokens - penalty (if early)
- `claimRewards(positionId)` — claim accrued rewards without unstaking
- `getPosition(positionId)` — view stake details
- `getAPR(tier)` — view current tier rates

**Early withdrawal penalty:** 10% of staked amount, sent to treasury address.

**Events:**
- `Staked(user, amount, tier, positionId)`
- `Unstaked(user, amount, rewards, penalty)`
- `RewardsClaimed(user, amount, positionId)`
- `TierUpdated(tier, newAPR)`

**Security integrations:**
- OpenZeppelin `ReentrancyGuard` on all state-changing functions
- OpenZeppelin `Ownable2Step` for admin (two-step ownership transfer)
- OpenZeppelin `Pausable` for emergency stops
- OpenZeppelin `SafeERC20` for token transfers
- Checks-Effects-Interactions pattern throughout

### Solana — Rust/Anchor

Mirrored staking logic, simplified for POC scope:

| Tier | Lock Period | APR |
|------|-----------|-----|
| Bronze | 30 days | 5% |
| Gold | 90 days | 18% |

**Accounts:**
- `StakingPool` — program state (authority, reward rate, treasury)
- `UserStake` — PDA per user per position (amount, tier, timestamp)
- `TokenVault` — PDA-owned token account holding staked tokens

**Instructions:**
- `initialize()` — create pool with tier configs
- `stake(amount, tier)` — transfer tokens to vault, create UserStake PDA
- `unstake()` — calculate rewards/penalty, transfer back
- `claim_rewards()` — claim without unstaking

**Simplified vs EVM:**
- 2 tiers instead of 3 (Bronze + Gold)
- No pausable (added complexity for POC)
- Fixed penalty rate (no admin update)

---

## 5. Event Indexer & Catastrophe Recovery

### Strategy: Chain-as-Source-of-Truth + DB Rebuild (Primary) + WAL Archiving (Fast Restore)

The database is treated as a cache of on-chain state. If lost, the re-indexer replays all contract events from the deployment block and rebuilds completely.

### Indexer Behavior

**Startup:**
1. Read `last_indexed_block` from `block_cursors` table per chain
2. If null → full re-index from contract deployment block
3. If behind chain head → catch-up replay
4. Once caught up → switch to websocket subscription (EVM) / polling (Solana)

**Per Event:**
1. Parse event (Staked, Unstaked, RewardsClaimed)
2. Upsert position in DB (idempotent by `tx_hash` + `log_index`)
3. Update `last_indexed_block`
4. All in a single DB transaction

**Recovery:**
1. Truncate positions/events tables
2. Reset block cursors
3. Restart indexer → full re-index from chain
4. DB is rebuilt entirely from on-chain events

### Production Recommendations (Documented, Not Built)

- PostgreSQL WAL archiving for point-in-time recovery (fast restore shortcut)
- WAL shipped to object storage (S3/GCS)
- Re-index validates restored data integrity

### Database Schema

```sql
-- Block cursor per chain (for indexer recovery)
CREATE TABLE block_cursors (
    chain_id    TEXT PRIMARY KEY,
    last_block  BIGINT NOT NULL,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Staking positions (rebuilt from chain events)
CREATE TABLE positions (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chain_id    TEXT NOT NULL,
    wallet      TEXT NOT NULL,
    amount      NUMERIC NOT NULL,
    tier        TEXT NOT NULL,
    staked_at   TIMESTAMPTZ NOT NULL,
    lock_until  TIMESTAMPTZ NOT NULL,
    status      TEXT NOT NULL DEFAULT 'active',
    tx_hash     TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Raw indexed events (audit trail, re-indexing validation)
CREATE TABLE chain_events (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chain_id     TEXT NOT NULL,
    event_type   TEXT NOT NULL,
    tx_hash      TEXT NOT NULL,
    log_index    INTEGER NOT NULL,
    block_number BIGINT NOT NULL,
    raw_data     JSONB NOT NULL,
    indexed_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (tx_hash, log_index)
);

-- Reward snapshots (computed from events)
CREATE TABLE rewards (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    position_id         UUID NOT NULL REFERENCES positions(id),
    accrued_amount      NUMERIC NOT NULL DEFAULT 0,
    last_calculated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Audit log (security)
CREATE TABLE audit_log (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    action      TEXT NOT NULL,
    actor       TEXT NOT NULL,
    chain_id    TEXT,
    details     JSONB,
    prev_hash   TEXT,
    hash        TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

## 6. Security Architecture

### Carried From Infinite Brain

| Pattern | Implementation |
|---------|---------------|
| JWT auth + ConnectRPC interceptor | Wallet owner proves identity before chain txs |
| RBAC (`user` + `admin`) | Admin: pause contracts, update tiers. Users: own positions only |
| Rate limiting (Valkey-backed) | Prevent tx spam, protect RPC node quotas |
| Typed errors (AppError sentinels) | Never expose chain RPC errors, private keys, or internal state |
| Security headers middleware | XSS, clickjacking, HSTS, CSP |
| Audit log with hash chain | Every staking action logged immutably |
| Input validation | Wallet address format, amount bounds, tier enum |
| Secrets via env vars | Private keys, JWT secret never in code |

### DeFi-Specific Security

| Pattern | Implementation |
|---------|---------------|
| Reentrancy guard | OpenZeppelin `ReentrancyGuard` on all state-changing Solidity functions |
| Checks-Effects-Interactions | Validate → update state → external call (defense in depth) |
| Tx signing isolation | Private keys in dedicated signer module, never exposed to API layer |
| Nonce management | Sequential nonce tracker per chain per wallet |
| Amount bounds | Min/max stake enforced in contract AND backend (double validation) |
| Pausable contracts | Emergency freeze via OpenZeppelin `Pausable` |
| Ownable2Step | Two-step ownership transfer (prevents accidental loss) |
| Event integrity | Indexer validates event signatures match expected ABI |
| Wallet validation | EVM: checksum address. Solana: base58 + valid pubkey |

### Interceptor Stack (Order)

```
1. Security Headers
2. Rate Limiter (Valkey)
3. Auth Interceptor (JWT)
4. RBAC Check
5. Input Validation
6. Audit Logger
7. Handler → Tx Signer (isolated) → Chain Adapter
```

### Documented But Not Built (Production Recommendations)

- Honeypot endpoints
- Field-level encryption
- mTLS between services
- OpenBao/Vault secret rotation
- Pre-commit secret scanning hooks

---

## 7. Phased Build Plan (7 Days)

### Phase 1 — Foundation (Day 1)
- Project scaffold (Go module, Foundry init, Anchor init)
- Docker Compose (PostgreSQL, Anvil, solana-test-validator, Valkey)
- Protobuf definitions + Buf code generation
- DB migrations (Goose)
- Config + logging setup
- Makefile with all dev commands

### Phase 2 — EVM Smart Contracts (Day 2)
- `StakingToken.sol` (ERC-20, mintable)
- `TieredStaking.sol` (3 tiers, penalties, security patterns)
- Forge tests (unit + fuzz)
- Deploy script to Anvil
- Generate Go bindings from ABI

### Phase 3 — Go Backend Core (Day 3–4)
- `ChainStaker` interface + types
- `EVMStaker` adapter
- Security layer (auth, rate limiting, validation, audit)
- `StakingService` with chain routing
- ConnectRPC handlers
- Unit + integration tests

### Phase 4 — Event Indexer + Recovery (Day 4–5)
- Event indexer engine (catch-up + subscription)
- EVM event source
- Idempotent DB writes
- Full re-index recovery
- Recovery tests

### Phase 5 — Solana Integration (Day 5–6)
- Anchor staking program (2 tiers)
- Anchor tests
- `SolanaStaker` Go adapter
- Solana event source for indexer
- Integration tests

### Phase 6 — Polish + Demo (Day 7)
- grpcui setup
- End-to-end demo flow
- README with architecture diagram
- Retrospective doc

### Demo Flow (What "Done" Looks Like)

```
1. make up                          → start all services
2. Open grpcui                      → interactive RPC browser
3. GetTiers                         → Bronze/Silver/Gold (EVM), Bronze/Gold (Solana)
4. Stake(chain=EVM, amount=1000, tier=GOLD)
5. Stake(chain=SOLANA, amount=500, tier=BRONZE)
6. ListPositions(wallet=...)        → both positions unified
7. Kill PostgreSQL, restart         → indexer rebuilds from chain
8. ListPositions again              → same data, zero loss
```

---

## 8. Documentation Strategy

All decision docs live at `/Users/rian/focaApp/etherium_poc_docs/` in ADR format:

```
etherium_poc_docs/
├── README.md
├── 00-strategy/
│   └── why-this-poc.md
├── 01-smart-contracts/
│   ├── decisions/
│   │   ├── foundry-over-hardhat.md
│   │   ├── tiered-staking-design.md
│   │   └── security-considerations.md
│   └── phase-summary.md
├── 02-go-backend/
│   ├── decisions/
│   │   ├── chain-adapter-pattern.md
│   │   ├── connect-rpc-over-rest.md
│   │   └── event-indexing-strategy.md
│   └── phase-summary.md
├── 03-solana-integration/
│   ├── decisions/
│   │   ├── anchor-framework.md
│   │   ├── evm-vs-solana-model.md
│   │   └── scope-tradeoffs.md
│   └── phase-summary.md
├── 04-multi-chain-orchestration/
│   ├── decisions/
│   │   ├── unified-interface-design.md
│   │   ├── catastrophe-recovery-strategy.md
│   │   └── cross-chain-state.md
│   └── phase-summary.md
└── 05-retrospective/
    └── what-i-learned.md
```

Each ADR follows: **Context → Options Considered → Decision → Consequences**.
