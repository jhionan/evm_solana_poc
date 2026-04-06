package indexer

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ---------------------------------------------------------------------------
// Event signature topics
// ---------------------------------------------------------------------------

// Canonical Solidity event signatures for the staking contract.
// These must match exactly what is compiled into the contract ABI.
const (
	sigStaked         = "Staked(address,uint256,uint8,uint256)"
	sigUnstaked       = "Unstaked(address,uint256,uint256,uint256)"
	sigRewardsClaimed = "RewardsClaimed(address,uint256,uint256)"
)

var (
	topicStaked         = crypto.Keccak256Hash([]byte(sigStaked))
	topicUnstaked       = crypto.Keccak256Hash([]byte(sigUnstaked))
	topicRewardsClaimed = crypto.Keccak256Hash([]byte(sigRewardsClaimed))

	// allTopics is the set of topic[0] values the EVMEventSource subscribes to.
	allTopics = []common.Hash{topicStaked, topicUnstaked, topicRewardsClaimed}
)

// ---------------------------------------------------------------------------
// EVMEventSource
// ---------------------------------------------------------------------------

// EVMEventSource implements EventSource against an EVM-compatible chain using
// two go-ethereum clients:
//   - httpClient  — used for historical log filtering (FilterLogs)
//   - wsClient    — used for live log subscriptions (SubscribeFilterLogs)
//
// The contract address restricts the filter so only the staking contract's
// events are returned.
type EVMEventSource struct {
	chainID         string
	contractAddress common.Address
	httpClient      *ethclient.Client
	wsClient        *ethclient.Client
	logger          zerolog.Logger
}

// NewEVMEventSource constructs an EVMEventSource.
//
//   - chainID          — canonical chain identifier (e.g. "1", "137")
//   - contractAddress  — address of the deployed staking contract
//   - httpClient       — connected HTTP/IPC ethclient (used for FilterLogs)
//   - wsClient         — connected WebSocket ethclient (used for Subscribe)
func NewEVMEventSource(
	chainID string,
	contractAddress common.Address,
	httpClient *ethclient.Client,
	wsClient *ethclient.Client,
) *EVMEventSource {
	return &EVMEventSource{
		chainID:         chainID,
		contractAddress: contractAddress,
		httpClient:      httpClient,
		wsClient:        wsClient,
		logger:          log.With().Str("chain_id", chainID).Logger(),
	}
}

// ChainID implements EventSource.
func (s *EVMEventSource) ChainID() string { return s.chainID }

// LatestBlock implements EventSource. It queries the HTTP client for the
// current chain head block number.
func (s *EVMEventSource) LatestBlock(ctx context.Context) (int64, error) {
	n, err := s.httpClient.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("evm_source: BlockNumber: %w", err)
	}
	return int64(n), nil
}

// maxBlockRange is the maximum number of blocks to query in a single FilterLogs
// call. Public RPC endpoints typically cap this at 50 000.
const maxBlockRange int64 = 10_000

// CatchUp implements EventSource. It queries logs in chunks of maxBlockRange
// to respect RPC provider limits.
func (s *EVMEventSource) CatchUp(ctx context.Context, startBlock, endBlock int64) ([]ChainEvent, error) {
	var allEvents []ChainEvent

	for from := startBlock; from <= endBlock; from += maxBlockRange + 1 {
		to := from + maxBlockRange
		if to > endBlock {
			to = endBlock
		}

		chunk, err := s.fetchLogs(ctx, from, to)
		if err != nil {
			return nil, err
		}
		allEvents = append(allEvents, chunk...)
	}

	s.logger.Info().
		Str("chain_id", s.chainID).
		Int64("start_block", startBlock).
		Int64("end_block", endBlock).
		Int("events_found", len(allEvents)).
		Msg("evm_source: catch-up complete")

	return allEvents, nil
}

// fetchLogs performs a single FilterLogs call for a block range.
func (s *EVMEventSource) fetchLogs(ctx context.Context, startBlock, endBlock int64) ([]ChainEvent, error) {
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(startBlock),
		ToBlock:   big.NewInt(endBlock),
		Addresses: []common.Address{s.contractAddress},
		Topics:    [][]common.Hash{allTopics},
	}

	logs, err := s.httpClient.FilterLogs(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("evm_source: FilterLogs [%d-%d]: %w", startBlock, endBlock, err)
	}

	events := make([]ChainEvent, 0, len(logs))
	for _, l := range logs {
		event, err := s.parseLog(l)
		if err != nil {
			s.logger.Warn().
				Err(err).
				Str("tx_hash", l.TxHash.Hex()).
				Uint("log_index", l.Index).
				Msg("evm_source: skipping unparseable log")
			continue
		}
		events = append(events, event)
	}

	s.logger.Info().
		Int64("start_block", startBlock).
		Int64("end_block", endBlock).
		Int("events_found", len(events)).
		Msg("evm_source: catch-up complete")

	return events, nil
}

// Subscribe implements EventSource. It opens a WebSocket subscription for
// live logs emitted by the staking contract and fans them out onto the
// returned channel. The channel is closed when ctx is cancelled.
func (s *EVMEventSource) Subscribe(ctx context.Context) (<-chan ChainEvent, error) {
	if s.wsClient == nil {
		// No WebSocket client — return a channel that blocks until context is done.
		ch := make(chan ChainEvent)
		go func() {
			<-ctx.Done()
			close(ch)
		}()
		s.logger.Warn().Str("chain_id", s.chainID).Msg("evm_source: no WS client — subscription disabled, catch-up only")
		return ch, nil
	}

	query := ethereum.FilterQuery{
		Addresses: []common.Address{s.contractAddress},
		Topics:    [][]common.Hash{allTopics},
	}

	rawCh := make(chan types.Log, 64)

	sub, err := s.wsClient.SubscribeFilterLogs(ctx, query, rawCh)
	if err != nil {
		return nil, fmt.Errorf("evm_source: SubscribeFilterLogs: %w", err)
	}

	outCh := make(chan ChainEvent, 64)

	go func() {
		defer close(outCh)
		for {
			select {
			case <-ctx.Done():
				sub.Unsubscribe()
				return

			case err := <-sub.Err():
				if err != nil {
					s.logger.Error().Err(err).Msg("evm_source: subscription error")
				}
				return

			case l, ok := <-rawCh:
				if !ok {
					return
				}
				event, err := s.parseLog(l)
				if err != nil {
					s.logger.Warn().
						Err(err).
						Str("tx_hash", l.TxHash.Hex()).
						Uint("log_index", l.Index).
						Msg("evm_source: skipping unparseable live log")
					continue
				}
				outCh <- event
			}
		}
	}()

	return outCh, nil
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// parseLog converts a raw EVM log into a ChainEvent. It returns an error for
// logs that do not match any of the known event signatures.
func (s *EVMEventSource) parseLog(l types.Log) (ChainEvent, error) {
	if len(l.Topics) == 0 {
		return ChainEvent{}, fmt.Errorf("log has no topics")
	}

	sig := l.Topics[0]

	var eventType string
	var rawData map[string]interface{}
	var err error

	switch sig {
	case topicStaked:
		eventType = "Staked"
		rawData, err = parseStaked(l)
	case topicUnstaked:
		eventType = "Unstaked"
		rawData, err = parseUnstaked(l)
	case topicRewardsClaimed:
		eventType = "RewardsClaimed"
		rawData, err = parseRewardsClaimed(l)
	default:
		return ChainEvent{}, fmt.Errorf("unknown event topic %s", sig.Hex())
	}

	if err != nil {
		return ChainEvent{}, fmt.Errorf("parse %s: %w", eventType, err)
	}

	return ChainEvent{
		ChainID:     s.chainID,
		EventType:   eventType,
		TxHash:      l.TxHash.Hex(),
		LogIndex:    int(l.Index),
		BlockNumber: int64(l.BlockNumber),
		RawData:     rawData,
	}, nil
}

// parseStaked decodes a Staked(address,uint256,uint8,uint256) log.
//
// ABI encoding for non-indexed params packs them into Data in order:
//   amount uint256  (32 bytes)
//   tier   uint8    (32 bytes — padded)
//   lockEnd uint256 (32 bytes)
//
// The indexed param (address) is in Topics[1].
func parseStaked(l types.Log) (map[string]interface{}, error) {
	if len(l.Topics) < 2 {
		return nil, fmt.Errorf("Staked: expected ≥2 topics, got %d", len(l.Topics))
	}
	if len(l.Data) < 96 {
		return nil, fmt.Errorf("Staked: data too short (%d bytes)", len(l.Data))
	}

	staker := common.BytesToAddress(l.Topics[1].Bytes())
	amount := new(big.Int).SetBytes(l.Data[0:32])
	tier := new(big.Int).SetBytes(l.Data[32:64])
	lockEnd := new(big.Int).SetBytes(l.Data[64:96])

	return map[string]interface{}{
		"staker":  staker.Hex(),
		"amount":  amount.String(),
		"tier":    tier.Uint64(),
		"lock_end": lockEnd.String(),
	}, nil
}

// parseUnstaked decodes an Unstaked(address,uint256,uint256,uint256) log.
//
// Indexed: address (Topics[1])
// Data: positionID uint256 | amount uint256 | penalty uint256
func parseUnstaked(l types.Log) (map[string]interface{}, error) {
	if len(l.Topics) < 2 {
		return nil, fmt.Errorf("Unstaked: expected ≥2 topics, got %d", len(l.Topics))
	}
	if len(l.Data) < 96 {
		return nil, fmt.Errorf("Unstaked: data too short (%d bytes)", len(l.Data))
	}

	staker := common.BytesToAddress(l.Topics[1].Bytes())
	positionID := new(big.Int).SetBytes(l.Data[0:32])
	amount := new(big.Int).SetBytes(l.Data[32:64])
	penalty := new(big.Int).SetBytes(l.Data[64:96])

	return map[string]interface{}{
		"staker":      staker.Hex(),
		"position_id": positionID.String(),
		"amount":      amount.String(),
		"penalty":     penalty.String(),
	}, nil
}

// parseRewardsClaimed decodes a RewardsClaimed(address,uint256,uint256) log.
//
// Indexed: address (Topics[1])
// Data: positionID uint256 | rewards uint256
func parseRewardsClaimed(l types.Log) (map[string]interface{}, error) {
	if len(l.Topics) < 2 {
		return nil, fmt.Errorf("RewardsClaimed: expected ≥2 topics, got %d", len(l.Topics))
	}
	if len(l.Data) < 64 {
		return nil, fmt.Errorf("RewardsClaimed: data too short (%d bytes)", len(l.Data))
	}

	staker := common.BytesToAddress(l.Topics[1].Bytes())
	positionID := new(big.Int).SetBytes(l.Data[0:32])
	rewards := new(big.Int).SetBytes(l.Data[32:64])

	return map[string]interface{}{
		"staker":      staker.Hex(),
		"position_id": positionID.String(),
		"rewards":     rewards.String(),
	}, nil
}

// RawDataJSON serialises RawData to a JSON byte slice suitable for DB
// storage. Returns an error if marshalling fails.
func RawDataJSON(rawData map[string]interface{}) ([]byte, error) {
	b, err := json.Marshal(rawData)
	if err != nil {
		return nil, fmt.Errorf("RawDataJSON: %w", err)
	}
	return b, nil
}
