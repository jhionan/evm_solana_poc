use anchor_lang::prelude::*;

#[error_code]
pub enum StakingError {
    #[msg("Invalid staking tier")]
    InvalidTier,
    #[msg("Stake amount must be greater than zero")]
    InvalidAmount,
    #[msg("Position is not active")]
    PositionNotActive,
    #[msg("Unauthorized: not the position owner")]
    Unauthorized,
    #[msg("Arithmetic overflow")]
    Overflow,
}
