package indexer

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Store is the persistence layer that the Indexer writes to and reads cursor
// state from. Implementations must be safe for concurrent use.
type Store interface {
	// GetLastIndexedBlock returns the last block that was successfully
	// persisted for chainID. The second return value is false when no cursor
	// exists yet (i.e. the chain has never been indexed).
	GetLastIndexedBlock(ctx context.Context, chainID string) (int64, bool, error)

	// SaveEvent persists event and records blockNumber as the new cursor for
	// the event's chain. Both operations should be atomic where possible.
	SaveEvent(ctx context.Context, event ChainEvent, blockNumber int64) error
}

// Indexer is the engine that drives historical catch-up and live subscription
// for a single EventSource / chain combination.
type Indexer struct {
	source      EventSource
	store       Store
	deployBlock int64
	logger      zerolog.Logger
}

// NewIndexer constructs an Indexer. deployBlock is the first block that may
// contain events of interest; it is used as the cursor when no persisted
// cursor is found in store.
func NewIndexer(source EventSource, store Store, deployBlock int64) *Indexer {
	return &Indexer{
		source:      source,
		store:       store,
		deployBlock: deployBlock,
		logger:      log.With().Str("chain_id", source.ChainID()).Logger(),
	}
}

// CatchUp fetches and persists all events between the persisted cursor (or
// deployBlock when no cursor exists) and the current chain head. It returns
// an error only if a fatal operation fails; individual missing ranges are
// logged and skipped.
func (idx *Indexer) CatchUp(ctx context.Context) error {
	chainID := idx.source.ChainID()

	startBlock, hasCursor, err := idx.store.GetLastIndexedBlock(ctx, chainID)
	if err != nil {
		return fmt.Errorf("indexer: get cursor for chain %s: %w", chainID, err)
	}

	if !hasCursor {
		startBlock = idx.deployBlock
	} else {
		// Resume from the block after the last successfully indexed one.
		startBlock++
	}

	latestBlock, err := idx.source.LatestBlock(ctx)
	if err != nil {
		return fmt.Errorf("indexer: get latest block for chain %s: %w", chainID, err)
	}

	if startBlock > latestBlock {
		idx.logger.Info().
			Int64("start_block", startBlock).
			Int64("latest_block", latestBlock).
			Msg("indexer: already caught up, nothing to do")
		return nil
	}

	idx.logger.Info().
		Int64("start_block", startBlock).
		Int64("end_block", latestBlock).
		Msg("indexer: starting catch-up")

	events, err := idx.source.CatchUp(ctx, startBlock, latestBlock)
	if err != nil {
		return fmt.Errorf("indexer: catch-up fetch for chain %s [%d-%d]: %w",
			chainID, startBlock, latestBlock, err)
	}

	for _, event := range events {
		if err := idx.store.SaveEvent(ctx, event, event.BlockNumber); err != nil {
			return fmt.Errorf("indexer: save event tx=%s log=%d: %w",
				event.TxHash, event.LogIndex, err)
		}
	}

	idx.logger.Info().
		Int("events_saved", len(events)).
		Int64("start_block", startBlock).
		Int64("end_block", latestBlock).
		Msg("indexer: catch-up complete")

	return nil
}

// Run performs a full CatchUp and then enters a live subscription loop. It
// blocks until ctx is cancelled or an unrecoverable error occurs. The caller
// should run this in its own goroutine.
func (idx *Indexer) Run(ctx context.Context) error {
	if err := idx.CatchUp(ctx); err != nil {
		return fmt.Errorf("indexer: run catch-up: %w", err)
	}

	eventCh, err := idx.source.Subscribe(ctx)
	if err != nil {
		return fmt.Errorf("indexer: subscribe: %w", err)
	}

	idx.logger.Info().Msg("indexer: live subscription started")

	for {
		select {
		case <-ctx.Done():
			idx.logger.Info().Msg("indexer: context cancelled, shutting down")
			return ctx.Err()

		case event, ok := <-eventCh:
			if !ok {
				return fmt.Errorf("indexer: event channel closed unexpectedly for chain %s",
					idx.source.ChainID())
			}

			if err := idx.store.SaveEvent(ctx, event, event.BlockNumber); err != nil {
				// Log but do not abort; a single bad event should not kill the indexer.
				idx.logger.Error().
					Err(err).
					Str("tx_hash", event.TxHash).
					Int("log_index", event.LogIndex).
					Msg("indexer: failed to save live event")
				continue
			}

			idx.logger.Debug().
				Str("event_type", event.EventType).
				Str("tx_hash", event.TxHash).
				Int64("block_number", event.BlockNumber).
				Msg("indexer: live event saved")
		}
	}
}
