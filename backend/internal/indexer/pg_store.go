package indexer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	db "github.com/jhionan/multichain-staking/db/sqlc"
)

// PGStore is a PostgreSQL-backed implementation of the indexer Store interface.
// It uses sqlc-generated queries and a pgxpool for connection management.
// All writes within SaveEvent are wrapped in a single database transaction.
type PGStore struct {
	pool *pgxpool.Pool
	q    *db.Queries
}

// NewPGStore constructs a PGStore backed by the given connection pool.
func NewPGStore(pool *pgxpool.Pool) *PGStore {
	return &PGStore{
		pool: pool,
		q:    db.New(pool),
	}
}

// GetLastIndexedBlock returns the last block that was successfully persisted
// for chainID. The second return value is false when no cursor exists yet.
func (s *PGStore) GetLastIndexedBlock(ctx context.Context, chainID string) (int64, bool, error) {
	cursor, err := s.q.GetBlockCursor(ctx, chainID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("pg_store: GetBlockCursor chain=%s: %w", chainID, err)
	}
	return cursor.LastBlock, true, nil
}

// SaveEvent persists the chain event and advances the block cursor for its
// chain atomically within a single transaction.
//
// Event-type-specific behaviour:
//   - "Staked":   inserts a new position row derived from the event's RawData.
//   - "Unstaked": finds the active position for the staker wallet and marks it
//     as "unstaked". If no matching position is found the status update is
//     skipped — this can happen when the contract was already active before the
//     indexer started and the position was never recorded locally.
//
// The chain event insert uses ON CONFLICT DO NOTHING, making SaveEvent safe to
// call more than once for the same (tx_hash, log_index) pair.
func (s *PGStore) SaveEvent(ctx context.Context, event ChainEvent, blockNumber int64) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("pg_store: begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.q.WithTx(tx)

	// -----------------------------------------------------------------------
	// 1. Persist the raw chain event.
	// -----------------------------------------------------------------------
	rawJSON, err := json.Marshal(event.RawData)
	if err != nil {
		return fmt.Errorf("pg_store: marshal raw_data for tx=%s: %w", event.TxHash, err)
	}

	_, err = qtx.InsertChainEvent(ctx, db.InsertChainEventParams{
		ChainID:     event.ChainID,
		EventType:   event.EventType,
		TxHash:      event.TxHash,
		LogIndex:    int32(event.LogIndex), //nolint:gosec // log index fits comfortably in int32
		BlockNumber: blockNumber,
		RawData:     rawJSON,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		// ErrNoRows is returned by pgx when the INSERT hit ON CONFLICT DO NOTHING
		// and produced no row; that is idempotent and expected — any other error
		// is a genuine failure.
		return fmt.Errorf("pg_store: InsertChainEvent tx=%s log=%d: %w",
			event.TxHash, event.LogIndex, err)
	}

	// -----------------------------------------------------------------------
	// 2. Apply event-type-specific side effects.
	// -----------------------------------------------------------------------
	switch event.EventType {
	case "Staked":
		if err := s.handleStaked(ctx, qtx, event); err != nil {
			return err
		}

	case "Unstaked":
		if err := s.handleUnstaked(ctx, qtx, event); err != nil {
			return err
		}
	}

	// -----------------------------------------------------------------------
	// 3. Advance the block cursor.
	// -----------------------------------------------------------------------
	if err := qtx.UpsertBlockCursor(ctx, db.UpsertBlockCursorParams{
		ChainID:   event.ChainID,
		LastBlock: blockNumber,
	}); err != nil {
		return fmt.Errorf("pg_store: UpsertBlockCursor chain=%s block=%d: %w",
			event.ChainID, blockNumber, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("pg_store: commit tx for tx=%s: %w", event.TxHash, err)
	}

	return nil
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// handleStaked inserts a new position row from a "Staked" event's RawData.
// Expected RawData keys: "staker" (address hex), "amount" (decimal string),
// "tier" (uint64), "lock_end" (unix timestamp string).
func (s *PGStore) handleStaked(ctx context.Context, qtx *db.Queries, event ChainEvent) error {
	wallet, _ := event.RawData["staker"].(string)
	amountStr, _ := event.RawData["amount"].(string)
	lockEndStr, _ := event.RawData["lock_end"].(string)

	// Tier may arrive as float64 when decoded from JSON.
	var tierNum uint64
	switch v := event.RawData["tier"].(type) {
	case float64:
		tierNum = uint64(v)
	case uint64:
		tierNum = v
	}

	// Build pgtype.Numeric from the decimal string.
	amountBig := new(big.Int)
	if amountStr != "" {
		amountBig.SetString(amountStr, 10)
	}
	var amount pgtype.Numeric
	if err := amount.Scan(amountBig.String()); err != nil {
		return fmt.Errorf("pg_store: parse amount %q for tx=%s: %w", amountStr, event.TxHash, err)
	}

	// lock_end is a Unix timestamp (seconds) encoded as a decimal string.
	var lockUntil pgtype.Timestamptz
	if lockEndStr != "" {
		lockEndBig := new(big.Int)
		lockEndBig.SetString(lockEndStr, 10)
		ts := time.Unix(lockEndBig.Int64(), 0).UTC()
		lockUntil = pgtype.Timestamptz{Time: ts, Valid: true}
	}

	stakedAt := pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true}

	tierStr := fmt.Sprintf("%d", tierNum)

	_, err := qtx.InsertPosition(ctx, db.InsertPositionParams{
		ChainID:   event.ChainID,
		Wallet:    wallet,
		Amount:    amount,
		Tier:      tierStr,
		StakedAt:  stakedAt,
		LockUntil: lockUntil,
		Status:    "active",
		TxHash:    event.TxHash,
	})
	if err != nil {
		return fmt.Errorf("pg_store: InsertPosition tx=%s: %w", event.TxHash, err)
	}

	return nil
}

// handleUnstaked updates the status of the active position belonging to the
// staker wallet on this chain. It looks up positions by (chain_id, wallet)
// and marks the first "active" one as "unstaked".
//
// If no matching active position is found the update is silently skipped; this
// is expected when the contract had prior activity before the indexer started.
func (s *PGStore) handleUnstaked(ctx context.Context, qtx *db.Queries, event ChainEvent) error {
	wallet, _ := event.RawData["staker"].(string)
	if wallet == "" {
		// Nothing we can do without a wallet address.
		return nil
	}

	positions, err := qtx.ListPositionsByWallet(ctx, db.ListPositionsByWalletParams{
		ChainID: event.ChainID,
		Wallet:  wallet,
	})
	if err != nil {
		return fmt.Errorf("pg_store: ListPositionsByWallet chain=%s wallet=%s: %w",
			event.ChainID, wallet, err)
	}

	for _, pos := range positions {
		if pos.Status == "active" {
			if err := qtx.UpdatePositionStatus(ctx, db.UpdatePositionStatusParams{
				ID:     pos.ID,
				Status: "unstaked",
			}); err != nil {
				return fmt.Errorf("pg_store: UpdatePositionStatus pos=%v: %w", pos.ID, err)
			}
			// Update the first active position found and stop.
			return nil
		}
	}

	// No active position found — skip silently.
	return nil
}
