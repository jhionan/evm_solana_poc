// Package solana provides a ChainStaker adapter for the Solana blockchain.
// Transaction methods are placeholders until the on-chain program IDL is
// compiled and the instruction builders are wired in.
package solana

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	solanago "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/jhionan/multichain-staking/internal/chain"
	"github.com/rs/zerolog"
)

// ErrNotYetImplemented is returned by mutation methods until the Solana
// program client is integrated.
var ErrNotYetImplemented = errors.New("solana: not yet implemented — program client not yet integrated")

// ErrUseIndexedDB is returned by ListPositions because position listing is
// delegated to the off-chain indexer, not direct account scanning.
var ErrUseIndexedDB = errors.New("solana: use indexed DB for listing positions")

// hardcodedTiers mirrors the on-chain staking program's tier configuration.
// Solana variant supports Bronze and Gold only (no Silver tier).
var hardcodedTiers = []chain.Tier{
	{
		Type:     chain.TierBronze,
		LockDays: 30,
		APRBps:   500,
		MinStake: big.NewInt(0), // placeholder until program IDL is parsed
	},
	{
		Type:     chain.TierGold,
		LockDays: 90,
		APRBps:   1800,
		MinStake: big.NewInt(0),
	},
}

// SolanaStaker implements chain.ChainStaker for the Solana blockchain.
type SolanaStaker struct {
	client    *rpc.Client
	programID solanago.PublicKey
	authority solanago.PrivateKey
	logger    zerolog.Logger
}

// NewSolanaStaker creates a SolanaStaker.
// client is required; authority may be a zero-value PrivateKey for read-only
// operation until signing is wired in.
func NewSolanaStaker(
	client *rpc.Client,
	programID solanago.PublicKey,
	authority solanago.PrivateKey,
	logger zerolog.Logger,
) (*SolanaStaker, error) {
	if client == nil {
		return nil, fmt.Errorf("solana: rpc.Client must not be nil")
	}
	return &SolanaStaker{
		client:    client,
		programID: programID,
		authority: authority,
		logger:    logger,
	}, nil
}

// ChainID returns ChainSolana to identify this adapter family.
func (s *SolanaStaker) ChainID() chain.ChainType {
	return chain.ChainSolana
}

// HealthCheck verifies the adapter can reach the Solana RPC node.
func (s *SolanaStaker) HealthCheck(ctx context.Context) error {
	out, err := s.client.GetHealth(ctx)
	if err != nil {
		return fmt.Errorf("solana: health check failed: %w", err)
	}
	if out != "ok" {
		return fmt.Errorf("solana: node unhealthy: %s", out)
	}
	return nil
}

// GetTiers returns the hardcoded tier configuration for the Solana program.
func (s *SolanaStaker) GetTiers(_ context.Context) ([]chain.Tier, error) {
	return hardcodedTiers, nil
}

// Stake is a placeholder pending program client integration.
func (s *SolanaStaker) Stake(_ context.Context, _ chain.StakeRequest) (chain.StakeReceipt, error) {
	return chain.StakeReceipt{}, ErrNotYetImplemented
}

// Unstake is a placeholder pending program client integration.
func (s *SolanaStaker) Unstake(_ context.Context, _ string) (chain.UnstakeReceipt, error) {
	return chain.UnstakeReceipt{}, ErrNotYetImplemented
}

// ClaimRewards is a placeholder pending program client integration.
func (s *SolanaStaker) ClaimRewards(_ context.Context, _ string) (chain.ClaimReceipt, error) {
	return chain.ClaimReceipt{}, ErrNotYetImplemented
}

// GetPosition is a placeholder pending program client integration.
func (s *SolanaStaker) GetPosition(_ context.Context, _ string) (chain.StakePosition, error) {
	return chain.StakePosition{}, ErrNotYetImplemented
}

// ListPositions always returns ErrUseIndexedDB. Position listing is handled
// by the off-chain indexer.
func (s *SolanaStaker) ListPositions(_ context.Context, _ string) ([]chain.StakePosition, error) {
	return nil, ErrUseIndexedDB
}

// compile-time interface assertion
var _ chain.ChainStaker = (*SolanaStaker)(nil)
