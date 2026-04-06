package api_test

import (
	"context"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	stakingv1 "github.com/jhionan/multichain-staking/gen/staking/v1"
	"github.com/jhionan/multichain-staking/gen/staking/v1/stakingv1connect"
	"github.com/jhionan/multichain-staking/internal/api"
	"github.com/jhionan/multichain-staking/internal/chain"
	"github.com/jhionan/multichain-staking/internal/staking"
)

// ---------------------------------------------------------------------------
// mockStaker — inline implementation of chain.ChainStaker for tests
// ---------------------------------------------------------------------------

type mockStaker struct {
	chainID chain.ChainType
}

func (m *mockStaker) ChainID() chain.ChainType { return m.chainID }

func (m *mockStaker) HealthCheck(_ context.Context) error { return nil }

func (m *mockStaker) GetTiers(_ context.Context) ([]chain.Tier, error) {
	return []chain.Tier{
		{
			Type:     chain.TierBronze,
			LockDays: 30,
			APRBps:   500,
			MinStake: big.NewInt(100),
		},
		{
			Type:     chain.TierSilver,
			LockDays: 60,
			APRBps:   800,
			MinStake: big.NewInt(500),
		},
		{
			Type:     chain.TierGold,
			LockDays: 90,
			APRBps:   1200,
			MinStake: big.NewInt(1000),
		},
	}, nil
}

func (m *mockStaker) Stake(_ context.Context, req chain.StakeRequest) (chain.StakeReceipt, error) {
	return chain.StakeReceipt{PositionID: "pos-1", TxHash: "0xabc"}, nil
}

func (m *mockStaker) Unstake(_ context.Context, positionID string) (chain.UnstakeReceipt, error) {
	return chain.UnstakeReceipt{
		AmountReturned: big.NewInt(100),
		Rewards:        big.NewInt(10),
		Penalty:        big.NewInt(0),
		TxHash:         "0xdef",
	}, nil
}

func (m *mockStaker) ClaimRewards(_ context.Context, positionID string) (chain.ClaimReceipt, error) {
	return chain.ClaimReceipt{RewardsClaimed: big.NewInt(5), TxHash: "0xghi"}, nil
}

func (m *mockStaker) GetPosition(_ context.Context, positionID string) (chain.StakePosition, error) {
	return chain.StakePosition{
		ID:             positionID,
		Chain:          m.chainID,
		Wallet:         "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		Amount:         big.NewInt(100),
		Tier:           chain.TierBronze,
		Status:         chain.StatusActive,
		StakedAt:       time.Now(),
		LockUntil:      time.Now().Add(30 * 24 * time.Hour),
		AccruedRewards: big.NewInt(0),
		TxHash:         "0xabc",
	}, nil
}

func (m *mockStaker) ListPositions(_ context.Context, wallet string) ([]chain.StakePosition, error) {
	return []chain.StakePosition{}, nil
}

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

func newTestServer(t *testing.T) (*httptest.Server, stakingv1connect.StakingServiceClient) {
	t.Helper()

	mock := &mockStaker{chainID: chain.ChainEVM}
	svc := staking.NewService([]chain.ChainStaker{mock})
	h := api.NewHandler(svc)

	mux := http.NewServeMux()
	path, handler := stakingv1connect.NewStakingServiceHandler(h)
	mux.Handle(path, handler)

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	client := stakingv1connect.NewStakingServiceClient(http.DefaultClient, srv.URL)
	return srv, client
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestHandler_GetTiers_ReturnsThreeTiersWithBronzeFirst(t *testing.T) {
	_, client := newTestServer(t)

	resp, err := client.GetTiers(context.Background(), connect.NewRequest(&stakingv1.GetTiersRequest{
		Chain: stakingv1.Chain_CHAIN_EVM,
	}))
	require.NoError(t, err)

	tiers := resp.Msg.GetTiers()
	assert.Len(t, tiers, 3, "expected 3 tiers")
	assert.Equal(t, stakingv1.Tier_TIER_BRONZE, tiers[0].GetTier(), "first tier should be TIER_BRONZE")
}
