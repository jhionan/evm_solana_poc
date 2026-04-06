package solana_test

import (
	"context"
	"testing"

	"github.com/jhionan/multichain-staking/internal/chain"
	"github.com/jhionan/multichain-staking/internal/chain/solana"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// compile-time assertion: SolanaStaker must satisfy the ChainStaker interface.
var _ chain.ChainStaker = (*solana.SolanaStaker)(nil)

func newTestStaker(t *testing.T) *solana.SolanaStaker {
	t.Helper()
	s, err := solana.NewSolanaStakerForTest(zerolog.Nop())
	require.NoError(t, err)
	return s
}

func TestSolanaStaker_ChainID_ReturnsChainSolana(t *testing.T) {
	s := newTestStaker(t)
	assert.Equal(t, chain.ChainSolana, s.ChainID())
}

func TestSolanaStaker_GetTiers_ReturnsTwoTiers(t *testing.T) {
	s := newTestStaker(t)

	tiers, err := s.GetTiers(context.Background())
	require.NoError(t, err)
	require.Len(t, tiers, 2, "expected exactly 2 tiers (Bronze, Gold)")
}

func TestSolanaStaker_GetTiers_BronzeValues(t *testing.T) {
	s := newTestStaker(t)

	tiers, err := s.GetTiers(context.Background())
	require.NoError(t, err)

	bronze := tiers[0]
	assert.Equal(t, chain.TierBronze, bronze.Type)
	assert.Equal(t, uint32(30), bronze.LockDays)
	assert.Equal(t, uint32(500), bronze.APRBps)
}

func TestSolanaStaker_GetTiers_GoldValues(t *testing.T) {
	s := newTestStaker(t)

	tiers, err := s.GetTiers(context.Background())
	require.NoError(t, err)

	gold := tiers[1]
	assert.Equal(t, chain.TierGold, gold.Type)
	assert.Equal(t, uint32(90), gold.LockDays)
	assert.Equal(t, uint32(1800), gold.APRBps)
}

func TestSolanaStaker_Stake_ReturnsNotYetImplementedError(t *testing.T) {
	s := newTestStaker(t)

	_, err := s.Stake(context.Background(), chain.StakeRequest{})
	assert.ErrorIs(t, err, solana.ErrNotYetImplemented)
}

func TestSolanaStaker_Unstake_ReturnsNotYetImplementedError(t *testing.T) {
	s := newTestStaker(t)

	_, err := s.Unstake(context.Background(), "pos-1")
	assert.ErrorIs(t, err, solana.ErrNotYetImplemented)
}

func TestSolanaStaker_ClaimRewards_ReturnsNotYetImplementedError(t *testing.T) {
	s := newTestStaker(t)

	_, err := s.ClaimRewards(context.Background(), "pos-1")
	assert.ErrorIs(t, err, solana.ErrNotYetImplemented)
}

func TestSolanaStaker_GetPosition_ReturnsNotYetImplementedError(t *testing.T) {
	s := newTestStaker(t)

	_, err := s.GetPosition(context.Background(), "pos-1")
	assert.ErrorIs(t, err, solana.ErrNotYetImplemented)
}

func TestSolanaStaker_ListPositions_ReturnsUseIndexedDBError(t *testing.T) {
	s := newTestStaker(t)

	_, err := s.ListPositions(context.Background(), "some-wallet")
	assert.ErrorIs(t, err, solana.ErrUseIndexedDB)
}
