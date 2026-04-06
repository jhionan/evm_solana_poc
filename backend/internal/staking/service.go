package staking

import (
	"context"
	"fmt"

	"github.com/jhionan/multichain-staking/internal/chain"
)

// Service routes staking operations to the correct chain adapter.
type Service struct {
	adapters map[chain.ChainType]chain.ChainStaker
}

// NewService constructs a Service from a slice of chain adapters.
// Duplicate chain IDs are silently overwritten by the last adapter.
func NewService(adapters []chain.ChainStaker) *Service {
	m := make(map[chain.ChainType]chain.ChainStaker, len(adapters))
	for _, a := range adapters {
		m[a.ChainID()] = a
	}
	return &Service{adapters: m}
}

// adapter returns the ChainStaker for chainType or a descriptive error.
func (s *Service) adapter(chainType chain.ChainType) (chain.ChainStaker, error) {
	a, ok := s.adapters[chainType]
	if !ok {
		return nil, fmt.Errorf("staking: no adapter registered for chain %q", chainType)
	}
	return a, nil
}

// GetTiers returns the staking tiers available on chainType.
func (s *Service) GetTiers(ctx context.Context, chainType chain.ChainType) ([]chain.Tier, error) {
	a, err := s.adapter(chainType)
	if err != nil {
		return nil, err
	}
	return a.GetTiers(ctx)
}

// Stake creates a new staking position on chainType.
func (s *Service) Stake(ctx context.Context, chainType chain.ChainType, req chain.StakeRequest) (chain.StakeReceipt, error) {
	a, err := s.adapter(chainType)
	if err != nil {
		return chain.StakeReceipt{}, err
	}
	return a.Stake(ctx, req)
}

// Unstake withdraws the position identified by positionID on chainType.
func (s *Service) Unstake(ctx context.Context, chainType chain.ChainType, positionID string) (chain.UnstakeReceipt, error) {
	a, err := s.adapter(chainType)
	if err != nil {
		return chain.UnstakeReceipt{}, err
	}
	return a.Unstake(ctx, positionID)
}

// ClaimRewards claims accrued rewards for positionID on chainType.
func (s *Service) ClaimRewards(ctx context.Context, chainType chain.ChainType, positionID string) (chain.ClaimReceipt, error) {
	a, err := s.adapter(chainType)
	if err != nil {
		return chain.ClaimReceipt{}, err
	}
	return a.ClaimRewards(ctx, positionID)
}

// GetPosition returns the current state of positionID on chainType.
func (s *Service) GetPosition(ctx context.Context, chainType chain.ChainType, positionID string) (chain.StakePosition, error) {
	a, err := s.adapter(chainType)
	if err != nil {
		return chain.StakePosition{}, err
	}
	return a.GetPosition(ctx, positionID)
}

// ListPositions returns all staking positions owned by wallet on chainType.
func (s *Service) ListPositions(ctx context.Context, chainType chain.ChainType, wallet string) ([]chain.StakePosition, error) {
	a, err := s.adapter(chainType)
	if err != nil {
		return nil, err
	}
	return a.ListPositions(ctx, wallet)
}

// HealthCheck pings all registered adapters and returns a map of chain → error.
// A nil error means the adapter is healthy.
func (s *Service) HealthCheck(ctx context.Context) map[chain.ChainType]error {
	results := make(map[chain.ChainType]error, len(s.adapters))
	for id, a := range s.adapters {
		results[id] = a.HealthCheck(ctx)
	}
	return results
}
