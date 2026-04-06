// Package api implements the ConnectRPC StakingService handler.
package api

import (
	"context"
	"math/big"

	"connectrpc.com/connect"

	stakingv1 "github.com/jhionan/multichain-staking/gen/staking/v1"
	"github.com/jhionan/multichain-staking/gen/staking/v1/stakingv1connect"
	"github.com/jhionan/multichain-staking/internal/chain"
	"github.com/jhionan/multichain-staking/internal/staking"
	apperrors "github.com/jhionan/multichain-staking/pkg/errors"
)

// Ensure Handler implements the generated interface at compile time.
var _ stakingv1connect.StakingServiceHandler = (*Handler)(nil)

// Handler is the ConnectRPC implementation of StakingService.
type Handler struct {
	svc *staking.Service
}

// NewHandler constructs a Handler backed by the given staking.Service.
func NewHandler(svc *staking.Service) *Handler {
	return &Handler{svc: svc}
}

// ---------------------------------------------------------------------------
// Conversion helpers
// ---------------------------------------------------------------------------

// toChainType converts the proto Chain enum to the domain ChainType.
func toChainType(c stakingv1.Chain) chain.ChainType {
	switch c {
	case stakingv1.Chain_CHAIN_SOLANA:
		return chain.ChainSolana
	default:
		return chain.ChainEVM
	}
}

// toProtoTier converts the domain TierType to the proto Tier enum.
func toProtoTier(t chain.TierType) stakingv1.Tier {
	switch t {
	case chain.TierSilver:
		return stakingv1.Tier_TIER_SILVER
	case chain.TierGold:
		return stakingv1.Tier_TIER_GOLD
	default:
		return stakingv1.Tier_TIER_BRONZE
	}
}

// fromProtoTier converts the proto Tier enum to the domain TierType.
func fromProtoTier(t stakingv1.Tier) chain.TierType {
	switch t {
	case stakingv1.Tier_TIER_SILVER:
		return chain.TierSilver
	case stakingv1.Tier_TIER_GOLD:
		return chain.TierGold
	default:
		return chain.TierBronze
	}
}

// toProtoPositionStatus converts the domain PositionStatus to the proto enum.
func toProtoPositionStatus(s chain.PositionStatus) stakingv1.PositionStatus {
	switch s {
	case chain.StatusActive:
		return stakingv1.PositionStatus_POSITION_STATUS_ACTIVE
	case chain.StatusUnstaked:
		return stakingv1.PositionStatus_POSITION_STATUS_UNSTAKED
	case chain.StatusPenalty:
		return stakingv1.PositionStatus_POSITION_STATUS_PENALTY
	default:
		return stakingv1.PositionStatus_POSITION_STATUS_UNSPECIFIED
	}
}

// toProtoChain converts the domain ChainType to the proto Chain enum.
func toProtoChain(c chain.ChainType) stakingv1.Chain {
	switch c {
	case chain.ChainSolana:
		return stakingv1.Chain_CHAIN_SOLANA
	default:
		return stakingv1.Chain_CHAIN_EVM
	}
}

// bigIntStr safely converts a *big.Int to its decimal string representation.
// A nil pointer is returned as "0".
func bigIntStr(n *big.Int) string {
	if n == nil {
		return "0"
	}
	return n.String()
}

// toProtoPosition converts a domain StakePosition to its proto equivalent.
func toProtoPosition(p chain.StakePosition) *stakingv1.StakePosition {
	return &stakingv1.StakePosition{
		Id:             p.ID,
		Chain:          toProtoChain(p.Chain),
		Wallet:         bigIntStr(p.Wallet),
		Amount:         bigIntStr(p.Amount),
		Tier:           toProtoTier(p.Tier),
		Status:         toProtoPositionStatus(p.Status),
		StakedAt:       p.StakedAt.Unix(),
		LockUntil:      p.LockUntil.Unix(),
		AccruedRewards: bigIntStr(p.AccruedRewards),
		TxHash:         p.TxHash,
	}
}

// ---------------------------------------------------------------------------
// RPC methods
// ---------------------------------------------------------------------------

// GetTiers returns the staking tiers for the requested chain (public endpoint).
func (h *Handler) GetTiers(
	ctx context.Context,
	req *connect.Request[stakingv1.GetTiersRequest],
) (*connect.Response[stakingv1.GetTiersResponse], error) {
	chainType := toChainType(req.Msg.GetChain())

	tiers, err := h.svc.GetTiers(ctx, chainType)
	if err != nil {
		return nil, apperrors.ToConnectError(err)
	}

	protoTiers := make([]*stakingv1.TierInfo, 0, len(tiers))
	for _, t := range tiers {
		protoTiers = append(protoTiers, &stakingv1.TierInfo{
			Tier:     toProtoTier(t.Type),
			LockDays: t.LockDays,
			AprBps:   t.APRBps,
			MinStake: bigIntStr(t.MinStake),
		})
	}

	return connect.NewResponse(&stakingv1.GetTiersResponse{Tiers: protoTiers}), nil
}

// Stake creates a new staking position.
func (h *Handler) Stake(
	ctx context.Context,
	req *connect.Request[stakingv1.StakeRequest],
) (*connect.Response[stakingv1.StakeResponse], error) {
	chainType := toChainType(req.Msg.GetChain())

	amount, ok := new(big.Int).SetString(req.Msg.GetAmount(), 10)
	if !ok {
		return nil, apperrors.ToConnectError(apperrors.ErrValidation.Wrap("amount is not a valid integer"))
	}

	wallet, ok := new(big.Int).SetString(req.Msg.GetWallet(), 10)
	if !ok {
		return nil, apperrors.ToConnectError(apperrors.ErrValidation.Wrap("wallet is not a valid integer"))
	}

	domainReq := chain.StakeRequest{
		Wallet: wallet,
		Amount: amount,
		Tier:   fromProtoTier(req.Msg.GetTier()),
	}

	receipt, err := h.svc.Stake(ctx, chainType, domainReq)
	if err != nil {
		return nil, apperrors.ToConnectError(err)
	}

	pos, err := h.svc.GetPosition(ctx, chainType, receipt.PositionID)
	if err != nil {
		return nil, apperrors.ToConnectError(err)
	}

	return connect.NewResponse(&stakingv1.StakeResponse{
		Position: toProtoPosition(pos),
		TxHash:   receipt.TxHash,
	}), nil
}

// Unstake withdraws a staking position.
func (h *Handler) Unstake(
	ctx context.Context,
	req *connect.Request[stakingv1.UnstakeRequest],
) (*connect.Response[stakingv1.UnstakeResponse], error) {
	chainType := toChainType(req.Msg.GetChain())

	receipt, err := h.svc.Unstake(ctx, chainType, req.Msg.GetPositionId())
	if err != nil {
		return nil, apperrors.ToConnectError(err)
	}

	return connect.NewResponse(&stakingv1.UnstakeResponse{
		AmountReturned: bigIntStr(receipt.AmountReturned),
		Rewards:        bigIntStr(receipt.Rewards),
		Penalty:        bigIntStr(receipt.Penalty),
		TxHash:         receipt.TxHash,
	}), nil
}

// ClaimRewards claims accrued rewards for a staking position.
func (h *Handler) ClaimRewards(
	ctx context.Context,
	req *connect.Request[stakingv1.ClaimRewardsRequest],
) (*connect.Response[stakingv1.ClaimRewardsResponse], error) {
	chainType := toChainType(req.Msg.GetChain())

	receipt, err := h.svc.ClaimRewards(ctx, chainType, req.Msg.GetPositionId())
	if err != nil {
		return nil, apperrors.ToConnectError(err)
	}

	return connect.NewResponse(&stakingv1.ClaimRewardsResponse{
		RewardsClaimed: bigIntStr(receipt.RewardsClaimed),
		TxHash:         receipt.TxHash,
	}), nil
}

// GetPosition returns the current state of a staking position.
func (h *Handler) GetPosition(
	ctx context.Context,
	req *connect.Request[stakingv1.GetPositionRequest],
) (*connect.Response[stakingv1.GetPositionResponse], error) {
	chainType := toChainType(req.Msg.GetChain())

	pos, err := h.svc.GetPosition(ctx, chainType, req.Msg.GetPositionId())
	if err != nil {
		return nil, apperrors.ToConnectError(err)
	}

	return connect.NewResponse(&stakingv1.GetPositionResponse{
		Position: toProtoPosition(pos),
	}), nil
}

// ListPositions returns all positions for a wallet on a chain.
func (h *Handler) ListPositions(
	ctx context.Context,
	req *connect.Request[stakingv1.ListPositionsRequest],
) (*connect.Response[stakingv1.ListPositionsResponse], error) {
	chainType := toChainType(req.Msg.GetChain())

	positions, err := h.svc.ListPositions(ctx, chainType, req.Msg.GetWallet())
	if err != nil {
		return nil, apperrors.ToConnectError(err)
	}

	protoPositions := make([]*stakingv1.StakePosition, 0, len(positions))
	for _, p := range positions {
		protoPositions = append(protoPositions, toProtoPosition(p))
	}

	return connect.NewResponse(&stakingv1.ListPositionsResponse{
		Positions: protoPositions,
	}), nil
}
