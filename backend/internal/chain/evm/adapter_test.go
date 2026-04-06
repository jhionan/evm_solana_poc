package evm_test

import (
	"context"
	"testing"

	"github.com/jhionan/multichain-staking/internal/chain"
	"github.com/jhionan/multichain-staking/internal/chain/evm"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// compile-time assertion: EVMStaker must satisfy the ChainStaker interface.
// If this line fails to compile the adapter is out of sync with the interface.
var _ chain.ChainStaker = (*evm.EVMStaker)(nil)

// newTestStaker builds an EVMStaker without a live RPC client.
// We pass a nil ethclient to trigger NewEVMStaker's guard and then bypass it
// for tests that only exercise in-memory logic (GetTiers, ChainID).
// For that we expose a test constructor that skips the nil check.
func newTestStaker(t *testing.T) *evm.EVMStaker {
	t.Helper()
	s, err := evm.NewEVMStakerForTest(zerolog.Nop())
	require.NoError(t, err)
	return s
}

func TestEVMStaker_ChainID_ReturnsChainEVM(t *testing.T) {
	s := newTestStaker(t)
	assert.Equal(t, chain.ChainEVM, s.ChainID())
}

func TestEVMStaker_GetTiers_ReturnsThreeTiers(t *testing.T) {
	s := newTestStaker(t)

	tiers, err := s.GetTiers(context.Background())
	require.NoError(t, err)
	require.Len(t, tiers, 3, "expected exactly 3 tiers (Bronze, Silver, Gold)")
}

func TestEVMStaker_GetTiers_BronzeValues(t *testing.T) {
	s := newTestStaker(t)

	tiers, err := s.GetTiers(context.Background())
	require.NoError(t, err)

	bronze := tiers[0]
	assert.Equal(t, chain.TierBronze, bronze.Type)
	assert.Equal(t, uint32(30), bronze.LockDays)
	assert.Equal(t, uint32(500), bronze.APRBps)
}

func TestEVMStaker_GetTiers_SilverValues(t *testing.T) {
	s := newTestStaker(t)

	tiers, err := s.GetTiers(context.Background())
	require.NoError(t, err)

	silver := tiers[1]
	assert.Equal(t, chain.TierSilver, silver.Type)
	assert.Equal(t, uint32(60), silver.LockDays)
	assert.Equal(t, uint32(1000), silver.APRBps)
}

func TestEVMStaker_GetTiers_GoldValues(t *testing.T) {
	s := newTestStaker(t)

	tiers, err := s.GetTiers(context.Background())
	require.NoError(t, err)

	gold := tiers[2]
	assert.Equal(t, chain.TierGold, gold.Type)
	assert.Equal(t, uint32(90), gold.LockDays)
	assert.Equal(t, uint32(1800), gold.APRBps)
}

func TestEVMStaker_Stake_ReturnsContractNotConnectedError(t *testing.T) {
	s := newTestStaker(t)

	_, err := s.Stake(context.Background(), chain.StakeRequest{})
	assert.ErrorIs(t, err, evm.ErrContractNotConnected)
}

func TestEVMStaker_Unstake_ReturnsContractNotConnectedError(t *testing.T) {
	s := newTestStaker(t)

	_, err := s.Unstake(context.Background(), "pos-1")
	assert.ErrorIs(t, err, evm.ErrContractNotConnected)
}

func TestEVMStaker_ClaimRewards_ReturnsContractNotConnectedError(t *testing.T) {
	s := newTestStaker(t)

	_, err := s.ClaimRewards(context.Background(), "pos-1")
	assert.ErrorIs(t, err, evm.ErrContractNotConnected)
}

func TestEVMStaker_GetPosition_ReturnsContractNotConnectedError(t *testing.T) {
	s := newTestStaker(t)

	_, err := s.GetPosition(context.Background(), "pos-1")
	assert.ErrorIs(t, err, evm.ErrContractNotConnected)
}

func TestEVMStaker_ListPositions_ReturnsUseIndexedDBError(t *testing.T) {
	s := newTestStaker(t)

	_, err := s.ListPositions(context.Background(), "0xdeadbeef")
	assert.ErrorIs(t, err, evm.ErrUseIndexedDB)
}
