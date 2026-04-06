package chain

import (
	"math/big"
	"time"
)

// ChainType identifies a blockchain family.
type ChainType string

const (
	ChainEVM    ChainType = "evm"
	ChainSolana ChainType = "solana"
)

// TierType identifies a staking tier.
type TierType string

const (
	TierBronze TierType = "bronze"
	TierSilver TierType = "silver"
	TierGold   TierType = "gold"
)

// PositionStatus represents the lifecycle state of a staking position.
type PositionStatus string

const (
	StatusActive   PositionStatus = "active"
	StatusUnstaked PositionStatus = "unstaked"
	StatusPenalty  PositionStatus = "penalty"
)

// Tier describes the parameters for a single staking tier.
type Tier struct {
	Type     TierType
	LockDays uint32
	APRBps   uint32   // annual percentage rate in basis points (1 bps = 0.01%)
	MinStake *big.Int // minimum stake amount in smallest denomination
}

// StakeRequest is the input for a staking operation.
type StakeRequest struct {
	Wallet string
	Amount *big.Int
	Tier   TierType
}

// StakeReceipt is returned after a successful stake.
type StakeReceipt struct {
	PositionID string
	TxHash     string
}

// UnstakeReceipt is returned after a successful unstake.
type UnstakeReceipt struct {
	AmountReturned *big.Int
	Rewards        *big.Int
	Penalty        *big.Int
	TxHash         string
}

// ClaimReceipt is returned after rewards are claimed.
type ClaimReceipt struct {
	RewardsClaimed *big.Int
	TxHash         string
}

// StakePosition represents the full state of a staking position.
type StakePosition struct {
	ID             string
	Chain          ChainType
	Wallet         string
	Amount         *big.Int
	Tier           TierType
	Status         PositionStatus
	StakedAt       time.Time
	LockUntil      time.Time
	AccruedRewards *big.Int
	TxHash         string
}
