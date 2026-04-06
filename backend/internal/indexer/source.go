package indexer

import "context"

// ChainEvent represents a decoded on-chain event from any supported chain.
type ChainEvent struct {
	ChainID     string                 `json:"chain_id"`
	EventType   string                 `json:"event_type"`
	TxHash      string                 `json:"tx_hash"`
	LogIndex    int                    `json:"log_index"`
	BlockNumber int64                  `json:"block_number"`
	RawData     map[string]interface{} `json:"raw_data"`
}

// EventSource is the abstraction every chain adapter must implement to supply
// events to the indexer engine. Implementations are expected to be safe for
// concurrent use.
type EventSource interface {
	// CatchUp returns all events emitted between startBlock and endBlock
	// (both inclusive). The caller is responsible for providing a sensible
	// range; implementations may split it internally to respect RPC limits.
	CatchUp(ctx context.Context, startBlock, endBlock int64) ([]ChainEvent, error)

	// Subscribe returns a channel that delivers live events as they are mined.
	// The channel is closed when ctx is cancelled or an unrecoverable error
	// occurs. Implementations should not block the caller after returning.
	Subscribe(ctx context.Context) (<-chan ChainEvent, error)

	// LatestBlock returns the current head block number reported by the node.
	LatestBlock(ctx context.Context) (int64, error)

	// ChainID returns the canonical identifier for this chain (e.g. "1" for
	// Ethereum mainnet, "137" for Polygon).
	ChainID() string
}
