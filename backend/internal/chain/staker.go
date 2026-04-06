package chain

import "context"

// ChainStaker is the abstraction every chain adapter must implement.
// All methods are context-aware to support cancellation and deadlines.
type ChainStaker interface {
	// ChainID returns the chain family this adapter belongs to.
	ChainID() ChainType

	// HealthCheck verifies the adapter can reach its underlying node/RPC.
	HealthCheck(ctx context.Context) error

	// GetTiers returns the staking tiers supported by this chain.
	GetTiers(ctx context.Context) ([]Tier, error)

	// Stake creates a new staking position for the given request.
	Stake(ctx context.Context, req StakeRequest) (StakeReceipt, error)

	// Unstake withdraws a staking position identified by positionID.
	Unstake(ctx context.Context, positionID string) (UnstakeReceipt, error)

	// ClaimRewards claims accrued rewards for positionID without unstaking.
	ClaimRewards(ctx context.Context, positionID string) (ClaimReceipt, error)

	// GetPosition returns the current state of a single staking position.
	GetPosition(ctx context.Context, positionID string) (StakePosition, error)

	// ListPositions returns all staking positions owned by wallet.
	ListPositions(ctx context.Context, wallet string) ([]StakePosition, error)
}
