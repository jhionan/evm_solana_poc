package indexer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// mockEventSource
// ---------------------------------------------------------------------------

// mockEventSource simulates an on-chain event source backed by a fixed set of
// events. CatchUp filters them by block range; Subscribe immediately closes
// the channel so tests can assert on catch-up alone without blocking.
type mockEventSource struct {
	chainID string
	events  []ChainEvent
	latest  int64
}

func newMockEventSource(chainID string, latest int64, events []ChainEvent) *mockEventSource {
	return &mockEventSource{chainID: chainID, events: events, latest: latest}
}

func (m *mockEventSource) ChainID() string { return m.chainID }

func (m *mockEventSource) LatestBlock(_ context.Context) (int64, error) {
	return m.latest, nil
}

func (m *mockEventSource) CatchUp(_ context.Context, startBlock, endBlock int64) ([]ChainEvent, error) {
	var result []ChainEvent
	for _, e := range m.events {
		if e.BlockNumber >= startBlock && e.BlockNumber <= endBlock {
			result = append(result, e)
		}
	}
	return result, nil
}

func (m *mockEventSource) Subscribe(ctx context.Context) (<-chan ChainEvent, error) {
	ch := make(chan ChainEvent)
	// Close the channel immediately so Run's select loop terminates cleanly
	// once the context is cancelled in tests.
	close(ch)
	return ch, nil
}

// ---------------------------------------------------------------------------
// mockStore
// ---------------------------------------------------------------------------

// mockStore is an in-memory Store implementation used during tests. It keeps
// track of the cursor per chain and all saved events.
type mockStore struct {
	cursors map[string]int64
	events  []ChainEvent
}

func newMockStore() *mockStore {
	return &mockStore{cursors: make(map[string]int64)}
}

func (s *mockStore) GetLastIndexedBlock(_ context.Context, chainID string) (int64, bool, error) {
	block, ok := s.cursors[chainID]
	return block, ok, nil
}

func (s *mockStore) SaveEvent(_ context.Context, event ChainEvent, blockNumber int64) error {
	s.cursors[event.ChainID] = blockNumber
	s.events = append(s.events, event)
	return nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func buildTestEvents() []ChainEvent {
	return []ChainEvent{
		{ChainID: "1", EventType: "Staked", TxHash: "0xaaa", LogIndex: 0, BlockNumber: 10},
		{ChainID: "1", EventType: "Staked", TxHash: "0xbbb", LogIndex: 0, BlockNumber: 20},
		{ChainID: "1", EventType: "Unstaked", TxHash: "0xccc", LogIndex: 0, BlockNumber: 30},
		{ChainID: "1", EventType: "RewardsClaimed", TxHash: "0xddd", LogIndex: 0, BlockNumber: 40},
		{ChainID: "1", EventType: "Staked", TxHash: "0xeee", LogIndex: 0, BlockNumber: 50},
	}
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

// TestIndexer_CatchUp_FromZero verifies that when no cursor is stored the
// indexer uses deployBlock as the start and indexes every event up to the
// chain head.
func TestIndexer_CatchUp_FromZero(t *testing.T) {
	ctx := context.Background()
	events := buildTestEvents()

	src := newMockEventSource("1", 50, events)
	store := newMockStore() // no cursor set
	idx := NewIndexer(src, store, 0)

	err := idx.CatchUp(ctx)
	require.NoError(t, err)

	// All 5 events should have been saved.
	assert.Len(t, store.events, 5, "expected all events to be indexed from block 0")

	// Cursor should be advanced to the block of the last event.
	cursor, hasCursor, _ := store.GetLastIndexedBlock(ctx, "1")
	assert.True(t, hasCursor)
	assert.Equal(t, int64(50), cursor)
}

// TestIndexer_CatchUp_FromCursor verifies that when a cursor is present the
// indexer resumes from cursor+1 and only persists events in the new range.
func TestIndexer_CatchUp_FromCursor(t *testing.T) {
	ctx := context.Background()
	events := buildTestEvents()

	src := newMockEventSource("1", 50, events)
	store := newMockStore()

	// Pre-set cursor at block 30 — events at blocks 10, 20, 30 are already
	// indexed; the indexer should only pick up blocks 31–50.
	store.cursors["1"] = 30

	idx := NewIndexer(src, store, 0)

	err := idx.CatchUp(ctx)
	require.NoError(t, err)

	// Only events at blocks 40 and 50 should have been saved in this run.
	require.Len(t, store.events, 2, "expected only events after cursor block to be indexed")
	assert.Equal(t, int64(40), store.events[0].BlockNumber)
	assert.Equal(t, int64(50), store.events[1].BlockNumber)

	// Cursor should be advanced to 50.
	cursor, hasCursor, _ := store.GetLastIndexedBlock(ctx, "1")
	assert.True(t, hasCursor)
	assert.Equal(t, int64(50), cursor)
}

// TestIndexer_CatchUp_AlreadyCaughtUp ensures CatchUp is a no-op when the
// cursor is at or beyond the chain head.
func TestIndexer_CatchUp_AlreadyCaughtUp(t *testing.T) {
	ctx := context.Background()
	events := buildTestEvents()

	src := newMockEventSource("1", 50, events)
	store := newMockStore()
	store.cursors["1"] = 50 // cursor == head

	idx := NewIndexer(src, store, 0)

	err := idx.CatchUp(ctx)
	require.NoError(t, err)

	// No new events should have been saved.
	assert.Empty(t, store.events)
}

// TestIndexer_CatchUp_EmptyChain ensures CatchUp works correctly when no
// events fall within the range (contract freshly deployed, no activity yet).
func TestIndexer_CatchUp_EmptyChain(t *testing.T) {
	ctx := context.Background()

	src := newMockEventSource("1", 100, nil) // no events at all
	store := newMockStore()

	idx := NewIndexer(src, store, 90)

	err := idx.CatchUp(ctx)
	require.NoError(t, err)

	assert.Empty(t, store.events)
}
