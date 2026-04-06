//go:build integration

package indexer

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	db "github.com/jhionan/multichain-staking/db/sqlc"
)

const testDSN = "postgres://staking:staking@localhost:5432/staking?sslmode=disable"

// newTestPool opens a connection pool for the integration test database and
// registers cleanup on t to close it automatically.
func newTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	pool, err := pgxpool.New(context.Background(), testDSN)
	require.NoError(t, err, "connect to test DB")
	t.Cleanup(pool.Close)
	return pool
}

// cleanDB truncates all relevant tables so each test starts from a known
// empty state. Positions must be truncated before chain_events because of
// foreign keys (rewards cascades from positions).
func cleanDB(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx := context.Background()
	q := db.New(pool)

	require.NoError(t, q.TruncateRewards(ctx))
	require.NoError(t, q.TruncatePositions(ctx))
	require.NoError(t, q.TruncateChainEvents(ctx))
	require.NoError(t, q.ResetAllBlockCursors(ctx))
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

// TestPGStore_GetLastIndexedBlock_NoCursor verifies that querying a chain that
// has never been indexed returns (0, false, nil).
func TestPGStore_GetLastIndexedBlock_NoCursor(t *testing.T) {
	pool := newTestPool(t)
	cleanDB(t, pool)

	store := NewPGStore(pool)

	block, hasCursor, err := store.GetLastIndexedBlock(context.Background(), "test-chain-1")
	require.NoError(t, err)
	assert.False(t, hasCursor, "expected no cursor for unseen chain")
	assert.Equal(t, int64(0), block)
}

// TestPGStore_SaveEvent_CreatesEventAndCursor verifies that SaveEvent persists
// the chain event and creates a block cursor for the chain.
func TestPGStore_SaveEvent_CreatesEventAndCursor(t *testing.T) {
	pool := newTestPool(t)
	cleanDB(t, pool)

	store := NewPGStore(pool)
	ctx := context.Background()

	event := ChainEvent{
		ChainID:     "test-chain-2",
		EventType:   "Staked",
		TxHash:      "0xdeadbeef000000000000000000000000000000000000000000000000000000aa",
		LogIndex:    0,
		BlockNumber: 100,
		RawData: map[string]interface{}{
			"staker":   "0xAbCdEf0000000000000000000000000000000001",
			"amount":   "1000000000000000000",
			"tier":     float64(1),
			"lock_end": "9999999999",
		},
	}

	err := store.SaveEvent(ctx, event, 100)
	require.NoError(t, err)

	// Cursor should now exist and point at block 100.
	block, hasCursor, err := store.GetLastIndexedBlock(ctx, "test-chain-2")
	require.NoError(t, err)
	assert.True(t, hasCursor)
	assert.Equal(t, int64(100), block)

	// The chain event should be queryable by tx hash.
	q := db.New(pool)
	events, err := q.GetEventsByTxHash(ctx, event.TxHash)
	require.NoError(t, err)
	require.Len(t, events, 1)
	assert.Equal(t, "test-chain-2", events[0].ChainID)
	assert.Equal(t, "Staked", events[0].EventType)
	assert.Equal(t, int64(100), events[0].BlockNumber)
}

// TestPGStore_SaveEvent_Idempotent verifies that calling SaveEvent twice with
// the same (tx_hash, log_index) does not create duplicate chain event rows.
func TestPGStore_SaveEvent_Idempotent(t *testing.T) {
	pool := newTestPool(t)
	cleanDB(t, pool)

	store := NewPGStore(pool)
	ctx := context.Background()

	event := ChainEvent{
		ChainID:     "test-chain-3",
		EventType:   "RewardsClaimed",
		TxHash:      "0xdeadbeef000000000000000000000000000000000000000000000000000000bb",
		LogIndex:    2,
		BlockNumber: 200,
		RawData: map[string]interface{}{
			"staker":      "0xAbCdEf0000000000000000000000000000000002",
			"position_id": "42",
			"rewards":     "500000000000000000",
		},
	}

	// Save the same event twice.
	require.NoError(t, store.SaveEvent(ctx, event, 200))
	require.NoError(t, store.SaveEvent(ctx, event, 200))

	// Exactly one row should exist.
	q := db.New(pool)
	events, err := q.GetEventsByTxHash(ctx, event.TxHash)
	require.NoError(t, err)
	assert.Len(t, events, 1, "duplicate saves must not create extra rows")
}

// TestPGStore_GetLastIndexedBlock_AfterSave verifies that the block cursor is
// correctly updated after each SaveEvent call.
func TestPGStore_GetLastIndexedBlock_AfterSave(t *testing.T) {
	pool := newTestPool(t)
	cleanDB(t, pool)

	store := NewPGStore(pool)
	ctx := context.Background()

	chainID := "test-chain-4"

	events := []ChainEvent{
		{
			ChainID:     chainID,
			EventType:   "Staked",
			TxHash:      "0xdeadbeef000000000000000000000000000000000000000000000000000000cc",
			LogIndex:    0,
			BlockNumber: 300,
			RawData: map[string]interface{}{
				"staker":   "0xAbCdEf0000000000000000000000000000000003",
				"amount":   "2000000000000000000",
				"tier":     float64(2),
				"lock_end": "9999999999",
			},
		},
		{
			ChainID:     chainID,
			EventType:   "Staked",
			TxHash:      "0xdeadbeef000000000000000000000000000000000000000000000000000000dd",
			LogIndex:    0,
			BlockNumber: 350,
			RawData: map[string]interface{}{
				"staker":   "0xAbCdEf0000000000000000000000000000000004",
				"amount":   "3000000000000000000",
				"tier":     float64(0),
				"lock_end": "9999999999",
			},
		},
	}

	for _, e := range events {
		require.NoError(t, store.SaveEvent(ctx, e, e.BlockNumber))
	}

	block, hasCursor, err := store.GetLastIndexedBlock(ctx, chainID)
	require.NoError(t, err)
	assert.True(t, hasCursor)
	assert.Equal(t, int64(350), block, "cursor should reflect the last saved block")
}

// TestPGStore_SaveEvent_UnstakedUpdatesPosition verifies that an "Unstaked"
// event updates the status of an existing "active" position to "unstaked".
func TestPGStore_SaveEvent_UnstakedUpdatesPosition(t *testing.T) {
	pool := newTestPool(t)
	cleanDB(t, pool)

	store := NewPGStore(pool)
	ctx := context.Background()

	chainID := "test-chain-5"
	wallet := "0xAbCdEf0000000000000000000000000000000005"
	stakedTxHash := "0xdeadbeef000000000000000000000000000000000000000000000000000000ee"
	unstakedTxHash := "0xdeadbeef000000000000000000000000000000000000000000000000000000ff"

	// First, stake to create an active position.
	stakeEvent := ChainEvent{
		ChainID:     chainID,
		EventType:   "Staked",
		TxHash:      stakedTxHash,
		LogIndex:    0,
		BlockNumber: 400,
		RawData: map[string]interface{}{
			"staker":   wallet,
			"amount":   "1000000000000000000",
			"tier":     float64(0),
			"lock_end": "9999999999",
		},
	}
	require.NoError(t, store.SaveEvent(ctx, stakeEvent, 400))

	// Verify position is active.
	q := db.New(pool)
	positions, err := q.ListPositionsByWallet(ctx, db.ListPositionsByWalletParams{
		ChainID: chainID,
		Wallet:  wallet,
	})
	require.NoError(t, err)
	require.Len(t, positions, 1)
	assert.Equal(t, "active", positions[0].Status)

	// Now unstake.
	unstakeEvent := ChainEvent{
		ChainID:     chainID,
		EventType:   "Unstaked",
		TxHash:      unstakedTxHash,
		LogIndex:    0,
		BlockNumber: 450,
		RawData: map[string]interface{}{
			"staker":      wallet,
			"position_id": "1",
			"amount":      "1000000000000000000",
			"penalty":     "0",
		},
	}
	require.NoError(t, store.SaveEvent(ctx, unstakeEvent, 450))

	// Position status should now be "unstaked".
	positions, err = q.ListPositionsByWallet(ctx, db.ListPositionsByWalletParams{
		ChainID: chainID,
		Wallet:  wallet,
	})
	require.NoError(t, err)
	require.Len(t, positions, 1)
	assert.Equal(t, "unstaked", positions[0].Status)

	// Cursor should be at block 450.
	block, hasCursor, err := store.GetLastIndexedBlock(ctx, chainID)
	require.NoError(t, err)
	assert.True(t, hasCursor)
	assert.Equal(t, int64(450), block)
}
