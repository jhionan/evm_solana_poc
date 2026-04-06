use anchor_lang::prelude::*;

// ---------------------------------------------------------------------------
// Tier enum
// ---------------------------------------------------------------------------

#[derive(AnchorSerialize, AnchorDeserialize, Clone, Copy, PartialEq, Eq, Default)]
pub enum Tier {
    #[default]
    Bronze,
    Gold,
}

// ---------------------------------------------------------------------------
// StakingPool
// ---------------------------------------------------------------------------

/// On-chain account that holds global pool configuration.
///
/// PDA seeds: [b"pool"]
#[account]
#[derive(Default)]
pub struct StakingPool {
    /// Authority that initialized the pool.
    pub authority: Pubkey,

    /// Treasury wallet that receives early-exit penalty tokens.
    pub treasury: Pubkey,

    /// SPL token mint being staked.
    pub token_mint: Pubkey,

    /// Token vault PDA that custodies all staked tokens.
    pub token_vault: Pubkey,

    /// Total tokens currently held in the vault.
    pub total_staked: u64,

    /// Monotonically-increasing counter used as the next position ID.
    pub next_position_id: u64,

    /// Early-exit penalty in basis points (default 1000 = 10%).
    pub penalty_bps: u16,

    /// Bronze tier APR in basis points (default 500 = 5%).
    pub bronze_apr_bps: u16,

    /// Bronze tier lock duration in days (default 30).
    pub bronze_lock_days: u16,

    /// Gold tier APR in basis points (default 1800 = 18%).
    pub gold_apr_bps: u16,

    /// Gold tier lock duration in days (default 90).
    pub gold_lock_days: u16,

    /// Bump for the pool PDA.
    pub bump: u8,
}

impl StakingPool {
    /// Discriminator (8) + authority (32) + treasury (32) + token_mint (32) +
    /// token_vault (32) + total_staked (8) + next_position_id (8) +
    /// penalty_bps (2) + bronze_apr_bps (2) + bronze_lock_days (2) +
    /// gold_apr_bps (2) + gold_lock_days (2) + bump (1) = 163
    pub const SIZE: usize = 8 + 32 + 32 + 32 + 32 + 8 + 8 + 2 + 2 + 2 + 2 + 2 + 1;
}

// ---------------------------------------------------------------------------
// UserStake
// ---------------------------------------------------------------------------

/// On-chain account representing a single staking position.
///
/// PDA seeds: [b"stake", pool.key().as_ref(), position_id.to_le_bytes()]
#[account]
#[derive(Default)]
pub struct UserStake {
    /// Wallet that owns this staking position.
    pub owner: Pubkey,

    /// The pool this position belongs to.
    pub pool: Pubkey,

    /// Unique monotonic position identifier (copied from pool at stake time).
    pub position_id: u64,

    /// Tokens staked in this position.
    pub amount: u64,

    /// Tier chosen when the position was opened.
    pub tier: Tier,

    /// Unix timestamp when tokens were staked.
    pub staked_at: i64,

    /// Unix timestamp after which penalty-free withdrawal is allowed.
    pub lock_until: i64,

    /// Unix timestamp of the most recent reward claim (starts equal to staked_at).
    pub last_claimed_at: i64,

    /// Whether the position is still active (false after unstake).
    pub active: bool,

    /// Bump for this position's PDA.
    pub bump: u8,
}

impl UserStake {
    /// Discriminator (8) + owner (32) + pool (32) + position_id (8) +
    /// amount (8) + tier (1) + staked_at (8) + lock_until (8) +
    /// last_claimed_at (8) + active (1) + bump (1) = 115
    pub const SIZE: usize = 8 + 32 + 32 + 8 + 8 + 1 + 8 + 8 + 8 + 1 + 1;
}
