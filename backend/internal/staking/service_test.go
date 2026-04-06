package staking_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/jhionan/multichain-staking/internal/chain"
	"github.com/jhionan/multichain-staking/internal/staking"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// mockStaker — inline test double for chain.ChainStaker
// ---------------------------------------------------------------------------

type mockStaker struct {
	chainID chain.ChainType
	tiers   []chain.Tier
}

func (m *mockStaker) ChainID() chain.ChainType { return m.chainID }

func (m *mockStaker) HealthCheck(_ context.Context) error { return nil }

func (m *mockStaker) GetTiers(_ context.Context) ([]chain.Tier, error) {
	return m.tiers, nil
}

func (m *mockStaker) Stake(_ context.Context, req chain.StakeRequest) (chain.StakeReceipt, error) {
	return chain.StakeReceipt{
		PositionID: "pos-001",
		TxHash:     "0xdeadbeef",
	}, nil
}

func (m *mockStaker) Unstake(_ context.Context, positionID string) (chain.UnstakeReceipt, error) {
	return chain.UnstakeReceipt{
		AmountReturned: big.NewInt(1000),
		Rewards:        big.NewInt(50),
		Penalty:        big.NewInt(0),
		TxHash:         "0xunstake",
	}, nil
}

func (m *mockStaker) ClaimRewards(_ context.Context, positionID string) (chain.ClaimReceipt, error) {
	return chain.ClaimReceipt{
		RewardsClaimed: big.NewInt(25),
		TxHash:         "0xclaim",
	}, nil
}

func (m *mockStaker) GetPosition(_ context.Context, positionID string) (chain.StakePosition, error) {
	return chain.StakePosition{ID: positionID, Chain: m.chainID}, nil
}

func (m *mockStaker) ListPositions(_ context.Context, wallet string) ([]chain.StakePosition, error) {
	return []chain.StakePosition{{ID: "pos-001", Chain: m.chainID}}, nil
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func newEVMTiers() []chain.Tier {
	return []chain.Tier{
		{Type: chain.TierBronze, LockDays: 30, APRBps: 500, MinStake: big.NewInt(100)},
		{Type: chain.TierSilver, LockDays: 90, APRBps: 1000, MinStake: big.NewInt(500)},
		{Type: chain.TierGold, LockDays: 180, APRBps: 1500, MinStake: big.NewInt(1000)},
	}
}

func newSolanaTiers() []chain.Tier {
	return []chain.Tier{
		{Type: chain.TierBronze, LockDays: 14, APRBps: 400, MinStake: big.NewInt(10)},
		{Type: chain.TierSilver, LockDays: 60, APRBps: 900, MinStake: big.NewInt(100)},
	}
}

func newService() *staking.Service {
	evm := &mockStaker{chainID: chain.ChainEVM, tiers: newEVMTiers()}
	sol := &mockStaker{chainID: chain.ChainSolana, tiers: newSolanaTiers()}
	return staking.NewService([]chain.ChainStaker{evm, sol})
}

// ---------------------------------------------------------------------------
// tests
// ---------------------------------------------------------------------------

func TestService_GetTiers_EVM_Returns3(t *testing.T) {
	svc := newService()
	tiers, err := svc.GetTiers(context.Background(), chain.ChainEVM)
	require.NoError(t, err)
	assert.Len(t, tiers, 3, "EVM adapter should return 3 tiers")
}

func TestService_GetTiers_Solana_Returns2(t *testing.T) {
	svc := newService()
	tiers, err := svc.GetTiers(context.Background(), chain.ChainSolana)
	require.NoError(t, err)
	assert.Len(t, tiers, 2, "Solana adapter should return 2 tiers")
}

func TestService_GetTiers_UnknownChain_Errors(t *testing.T) {
	svc := newService()
	_, err := svc.GetTiers(context.Background(), chain.ChainType("bitcoin"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no adapter registered")
}

func TestService_Stake_ReturnsReceipt(t *testing.T) {
	svc := newService()
	req := chain.StakeRequest{
		Wallet: big.NewInt(12345),
		Amount: big.NewInt(500),
		Tier:   chain.TierSilver,
	}
	receipt, err := svc.Stake(context.Background(), chain.ChainEVM, req)
	require.NoError(t, err)
	assert.Equal(t, "pos-001", receipt.PositionID)
	assert.Equal(t, "0xdeadbeef", receipt.TxHash)
}

func TestService_Stake_UnknownChain_Errors(t *testing.T) {
	svc := newService()
	req := chain.StakeRequest{Amount: big.NewInt(100), Tier: chain.TierBronze}
	_, err := svc.Stake(context.Background(), chain.ChainType("polygon"), req)
	assert.Error(t, err)
}

func TestService_HealthCheck_AllHealthy(t *testing.T) {
	svc := newService()
	results := svc.HealthCheck(context.Background())
	assert.Len(t, results, 2)
	for ch, err := range results {
		assert.NoError(t, err, "chain %s should be healthy", ch)
	}
}

func TestService_Unstake_ReturnsReceipt(t *testing.T) {
	svc := newService()
	receipt, err := svc.Unstake(context.Background(), chain.ChainEVM, "pos-001")
	require.NoError(t, err)
	assert.Equal(t, "0xunstake", receipt.TxHash)
	assert.Equal(t, big.NewInt(1000), receipt.AmountReturned)
}

func TestService_ClaimRewards_ReturnsReceipt(t *testing.T) {
	svc := newService()
	receipt, err := svc.ClaimRewards(context.Background(), chain.ChainSolana, "pos-001")
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(25), receipt.RewardsClaimed)
}

func TestService_GetPosition_ReturnsPosition(t *testing.T) {
	svc := newService()
	pos, err := svc.GetPosition(context.Background(), chain.ChainEVM, "pos-42")
	require.NoError(t, err)
	assert.Equal(t, "pos-42", pos.ID)
	assert.Equal(t, chain.ChainEVM, pos.Chain)
}

func TestService_ListPositions_ReturnsList(t *testing.T) {
	svc := newService()
	positions, err := svc.ListPositions(context.Background(), chain.ChainSolana, "wallet-abc")
	require.NoError(t, err)
	assert.Len(t, positions, 1)
}
