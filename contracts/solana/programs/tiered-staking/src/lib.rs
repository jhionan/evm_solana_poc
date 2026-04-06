use anchor_lang::prelude::*;
use anchor_spl::token::{self, Mint, Token, TokenAccount, Transfer};

pub mod errors;
pub mod state;

use errors::StakingError;
use state::{StakingPool, Tier, UserStake};

// Placeholder program ID — replace with the actual key from `anchor keys list`
// before deploying to any cluster.
declare_id!("Fg6PaFpoGXkYsidMpWTK6W2BeZ7FEfcYkg476zPFsLnS");

// ---------------------------------------------------------------------------
// Constants
// ---------------------------------------------------------------------------

pub const SECONDS_PER_DAY: i64 = 86_400;
pub const SECONDS_PER_YEAR: i64 = 365 * SECONDS_PER_DAY;
pub const BPS_DENOMINATOR: u64 = 10_000;

// ---------------------------------------------------------------------------
// Program entrypoint
// ---------------------------------------------------------------------------

#[program]
pub mod tiered_staking {
    use super::*;

    // -----------------------------------------------------------------------
    // initialize
    // -----------------------------------------------------------------------

    /// Create the StakingPool PDA and its associated token vault.
    pub fn initialize(
        ctx: Context<Initialize>,
        penalty_bps: u16,
        bronze_apr_bps: u16,
        bronze_lock_days: u16,
        gold_apr_bps: u16,
        gold_lock_days: u16,
    ) -> Result<()> {
        let pool = &mut ctx.accounts.pool;

        pool.authority = ctx.accounts.authority.key();
        pool.treasury = ctx.accounts.treasury.key();
        pool.token_mint = ctx.accounts.token_mint.key();
        pool.token_vault = ctx.accounts.token_vault.key();
        pool.total_staked = 0;
        pool.next_position_id = 0;
        pool.penalty_bps = penalty_bps;
        pool.bronze_apr_bps = bronze_apr_bps;
        pool.bronze_lock_days = bronze_lock_days;
        pool.gold_apr_bps = gold_apr_bps;
        pool.gold_lock_days = gold_lock_days;
        pool.bump = ctx.bumps.pool;

        emit!(PoolInitializedEvent {
            authority: pool.authority,
            treasury: pool.treasury,
            token_mint: pool.token_mint,
        });

        Ok(())
    }

    // -----------------------------------------------------------------------
    // stake
    // -----------------------------------------------------------------------

    /// Open a new staking position by locking tokens in the vault.
    pub fn stake(ctx: Context<Stake>, amount: u64, tier: Tier) -> Result<()> {
        require!(amount > 0, StakingError::InvalidAmount);

        let pool = &mut ctx.accounts.pool;
        let clock = Clock::get()?;
        let now = clock.unix_timestamp;

        // Calculate lock duration based on tier.
        let lock_days: i64 = match tier {
            Tier::Bronze => pool.bronze_lock_days as i64,
            Tier::Gold => pool.gold_lock_days as i64,
        };
        let lock_until = now
            .checked_add(lock_days * SECONDS_PER_DAY)
            .ok_or(StakingError::Overflow)?;

        // Assign and increment position ID.
        let position_id = pool.next_position_id;
        pool.next_position_id = pool
            .next_position_id
            .checked_add(1)
            .ok_or(StakingError::Overflow)?;

        // Update pool total.
        pool.total_staked = pool
            .total_staked
            .checked_add(amount)
            .ok_or(StakingError::Overflow)?;

        // Populate the UserStake account.
        let user_stake = &mut ctx.accounts.user_stake;
        user_stake.owner = ctx.accounts.user.key();
        user_stake.pool = pool.key();
        user_stake.position_id = position_id;
        user_stake.amount = amount;
        user_stake.tier = tier;
        user_stake.staked_at = now;
        user_stake.lock_until = lock_until;
        user_stake.last_claimed_at = now;
        user_stake.active = true;
        user_stake.bump = ctx.bumps.user_stake;

        // CPI: transfer tokens from user's wallet to vault.
        let cpi_accounts = Transfer {
            from: ctx.accounts.user_token.to_account_info(),
            to: ctx.accounts.token_vault.to_account_info(),
            authority: ctx.accounts.user.to_account_info(),
        };
        token::transfer(
            CpiContext::new(ctx.accounts.token_program.to_account_info(), cpi_accounts),
            amount,
        )?;

        emit!(StakedEvent {
            owner: user_stake.owner,
            position_id,
            amount,
            tier,
            lock_until,
        });

        Ok(())
    }

    // -----------------------------------------------------------------------
    // unstake
    // -----------------------------------------------------------------------

    /// Close a staking position and return tokens (with penalty if early exit).
    pub fn unstake(ctx: Context<Unstake>) -> Result<()> {
        let user_stake = &mut ctx.accounts.user_stake;
        let pool = &mut ctx.accounts.pool;
        let clock = Clock::get()?;
        let now = clock.unix_timestamp;

        require!(user_stake.active, StakingError::PositionNotActive);
        require!(
            user_stake.owner == ctx.accounts.user.key(),
            StakingError::Unauthorized
        );

        // Determine pending rewards accumulated since last claim.
        let elapsed_since_claim = now
            .checked_sub(user_stake.last_claimed_at)
            .ok_or(StakingError::Overflow)?;
        let apr_bps = match user_stake.tier {
            Tier::Bronze => pool.bronze_apr_bps as u64,
            Tier::Gold => pool.gold_apr_bps as u64,
        };
        let rewards = calculate_rewards(user_stake.amount, apr_bps, elapsed_since_claim)?;

        // Determine principal amount and any penalty.
        let principal = user_stake.amount;
        let (to_user, to_treasury) = if now >= user_stake.lock_until {
            // Past lock: full principal + rewards to user, nothing to treasury.
            let total = principal
                .checked_add(rewards)
                .ok_or(StakingError::Overflow)?;
            (total, 0u64)
        } else {
            // Early exit: apply penalty to principal; unclaimed rewards are forfeited.
            let penalty = (principal as u128)
                .checked_mul(pool.penalty_bps as u128)
                .ok_or(StakingError::Overflow)?
                / BPS_DENOMINATOR as u128;
            let penalty = penalty as u64;
            let net = principal
                .checked_sub(penalty)
                .ok_or(StakingError::Overflow)?;
            (net, penalty)
        };

        // Mark position inactive and reduce pool total.
        user_stake.active = false;
        pool.total_staked = pool
            .total_staked
            .checked_sub(principal)
            .ok_or(StakingError::Overflow)?;

        // PDA signer seeds for the vault.
        let pool_key = pool.key();
        let vault_seeds: &[&[&[u8]]] = &[&[b"vault", pool_key.as_ref(), &[ctx.bumps.token_vault]]];

        // CPI: transfer tokens to user.
        if to_user > 0 {
            let cpi_accounts = Transfer {
                from: ctx.accounts.token_vault.to_account_info(),
                to: ctx.accounts.user_token.to_account_info(),
                authority: ctx.accounts.token_vault.to_account_info(),
            };
            token::transfer(
                CpiContext::new_with_signer(
                    ctx.accounts.token_program.to_account_info(),
                    cpi_accounts,
                    vault_seeds,
                ),
                to_user,
            )?;
        }

        // CPI: transfer penalty to treasury (if any).
        if to_treasury > 0 {
            let cpi_accounts = Transfer {
                from: ctx.accounts.token_vault.to_account_info(),
                to: ctx.accounts.treasury_token.to_account_info(),
                authority: ctx.accounts.token_vault.to_account_info(),
            };
            token::transfer(
                CpiContext::new_with_signer(
                    ctx.accounts.token_program.to_account_info(),
                    cpi_accounts,
                    vault_seeds,
                ),
                to_treasury,
            )?;
        }

        emit!(UnstakedEvent {
            owner: user_stake.owner,
            position_id: user_stake.position_id,
            amount_returned: to_user,
            penalty: to_treasury,
            early_exit: now < user_stake.lock_until,
        });

        Ok(())
    }

    // -----------------------------------------------------------------------
    // claim_rewards
    // -----------------------------------------------------------------------

    /// Claim accrued rewards without closing the position.
    pub fn claim_rewards(ctx: Context<ClaimRewards>) -> Result<()> {
        let user_stake = &mut ctx.accounts.user_stake;
        let pool = &ctx.accounts.pool;
        let clock = Clock::get()?;
        let now = clock.unix_timestamp;

        require!(user_stake.active, StakingError::PositionNotActive);
        require!(
            user_stake.owner == ctx.accounts.user.key(),
            StakingError::Unauthorized
        );

        let elapsed = now
            .checked_sub(user_stake.last_claimed_at)
            .ok_or(StakingError::Overflow)?;

        let apr_bps = match user_stake.tier {
            Tier::Bronze => pool.bronze_apr_bps as u64,
            Tier::Gold => pool.gold_apr_bps as u64,
        };

        let rewards = calculate_rewards(user_stake.amount, apr_bps, elapsed)?;

        // Update claim timestamp before any external call (checks-effects-interactions).
        user_stake.last_claimed_at = now;

        if rewards > 0 {
            let pool_key = pool.key();
            let vault_seeds: &[&[&[u8]]] =
                &[&[b"vault", pool_key.as_ref(), &[ctx.bumps.token_vault]]];

            let cpi_accounts = Transfer {
                from: ctx.accounts.token_vault.to_account_info(),
                to: ctx.accounts.user_token.to_account_info(),
                authority: ctx.accounts.token_vault.to_account_info(),
            };
            token::transfer(
                CpiContext::new_with_signer(
                    ctx.accounts.token_program.to_account_info(),
                    cpi_accounts,
                    vault_seeds,
                ),
                rewards,
            )?;
        }

        emit!(RewardsClaimedEvent {
            owner: user_stake.owner,
            position_id: user_stake.position_id,
            rewards,
        });

        Ok(())
    }
}

// ---------------------------------------------------------------------------
// Reward helper
// ---------------------------------------------------------------------------

/// Compute: (amount * apr_bps * elapsed_secs) / (BPS_DENOMINATOR * SECONDS_PER_YEAR)
fn calculate_rewards(amount: u64, apr_bps: u64, elapsed_secs: i64) -> Result<u64> {
    if elapsed_secs <= 0 {
        return Ok(0);
    }
    let elapsed = elapsed_secs as u128;
    let numerator = (amount as u128)
        .checked_mul(apr_bps as u128)
        .ok_or(StakingError::Overflow)?
        .checked_mul(elapsed)
        .ok_or(StakingError::Overflow)?;
    let denominator = (BPS_DENOMINATOR as u128) * (SECONDS_PER_YEAR as u128);
    Ok((numerator / denominator) as u64)
}

// ---------------------------------------------------------------------------
// Account contexts
// ---------------------------------------------------------------------------

#[derive(Accounts)]
pub struct Initialize<'info> {
    /// The StakingPool PDA (created here).
    #[account(
        init,
        payer = authority,
        space = StakingPool::SIZE,
        seeds = [b"pool"],
        bump
    )]
    pub pool: Account<'info, StakingPool>,

    /// SPL token mint for the staked asset.
    pub token_mint: Account<'info, Mint>,

    /// Token vault PDA — custodies all staked tokens.
    /// Owned by the pool PDA via a PDA authority.
    #[account(
        init,
        payer = authority,
        token::mint = token_mint,
        token::authority = pool,
        seeds = [b"vault", pool.key().as_ref()],
        bump
    )]
    pub token_vault: Account<'info, TokenAccount>,

    /// Treasury wallet that receives penalty tokens (unchecked — any account).
    /// CHECK: Caller is responsible for supplying the correct treasury address.
    pub treasury: UncheckedAccount<'info>,

    /// Pool creator and fee payer.
    #[account(mut)]
    pub authority: Signer<'info>,

    pub system_program: Program<'info, System>,
    pub token_program: Program<'info, Token>,
    pub rent: Sysvar<'info, Rent>,
}

#[derive(Accounts)]
pub struct Stake<'info> {
    /// Pool state — increments counters and total_staked.
    #[account(
        mut,
        seeds = [b"pool"],
        bump = pool.bump
    )]
    pub pool: Account<'info, StakingPool>,

    /// New UserStake PDA for this position.
    #[account(
        init,
        payer = user,
        space = UserStake::SIZE,
        seeds = [
            b"stake",
            pool.key().as_ref(),
            &pool.next_position_id.to_le_bytes()
        ],
        bump
    )]
    pub user_stake: Account<'info, UserStake>,

    /// Vault receives the staked tokens.
    #[account(
        mut,
        seeds = [b"vault", pool.key().as_ref()],
        bump
    )]
    pub token_vault: Account<'info, TokenAccount>,

    /// User's token account — source of staked tokens.
    #[account(mut)]
    pub user_token: Account<'info, TokenAccount>,

    /// Wallet signing the transaction; pays for the new UserStake account.
    #[account(mut)]
    pub user: Signer<'info>,

    pub token_program: Program<'info, Token>,
    pub system_program: Program<'info, System>,
}

#[derive(Accounts)]
pub struct Unstake<'info> {
    #[account(
        mut,
        seeds = [b"pool"],
        bump = pool.bump
    )]
    pub pool: Account<'info, StakingPool>,

    /// The staking position being closed.
    #[account(
        mut,
        seeds = [
            b"stake",
            pool.key().as_ref(),
            &user_stake.position_id.to_le_bytes()
        ],
        bump = user_stake.bump
    )]
    pub user_stake: Account<'info, UserStake>,

    /// Vault from which tokens are returned.
    #[account(
        mut,
        seeds = [b"vault", pool.key().as_ref()],
        bump
    )]
    pub token_vault: Account<'info, TokenAccount>,

    /// User's token account — receives principal (and rewards if matured).
    #[account(mut)]
    pub user_token: Account<'info, TokenAccount>,

    /// Treasury token account — receives penalty tokens on early exit.
    #[account(mut)]
    pub treasury_token: Account<'info, TokenAccount>,

    #[account(mut)]
    pub user: Signer<'info>,

    pub token_program: Program<'info, Token>,
}

#[derive(Accounts)]
pub struct ClaimRewards<'info> {
    #[account(
        seeds = [b"pool"],
        bump = pool.bump
    )]
    pub pool: Account<'info, StakingPool>,

    #[account(
        mut,
        seeds = [
            b"stake",
            pool.key().as_ref(),
            &user_stake.position_id.to_le_bytes()
        ],
        bump = user_stake.bump
    )]
    pub user_stake: Account<'info, UserStake>,

    /// Vault pays out rewards.
    #[account(
        mut,
        seeds = [b"vault", pool.key().as_ref()],
        bump
    )]
    pub token_vault: Account<'info, TokenAccount>,

    /// User receives claimed rewards here.
    #[account(mut)]
    pub user_token: Account<'info, TokenAccount>,

    pub user: Signer<'info>,

    pub token_program: Program<'info, Token>,
}

// ---------------------------------------------------------------------------
// Events
// ---------------------------------------------------------------------------

#[event]
pub struct PoolInitializedEvent {
    pub authority: Pubkey,
    pub treasury: Pubkey,
    pub token_mint: Pubkey,
}

#[event]
pub struct StakedEvent {
    pub owner: Pubkey,
    pub position_id: u64,
    pub amount: u64,
    pub tier: Tier,
    pub lock_until: i64,
}

#[event]
pub struct UnstakedEvent {
    pub owner: Pubkey,
    pub position_id: u64,
    pub amount_returned: u64,
    pub penalty: u64,
    pub early_exit: bool,
}

#[event]
pub struct RewardsClaimedEvent {
    pub owner: Pubkey,
    pub position_id: u64,
    pub rewards: u64,
}
