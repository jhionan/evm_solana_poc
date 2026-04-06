// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/utils/ReentrancyGuard.sol";
import "@openzeppelin/contracts/access/Ownable2Step.sol";
import "@openzeppelin/contracts/utils/Pausable.sol";

/// @title TieredStaking
/// @notice Multi-tier staking contract with Bronze / Silver / Gold tiers, lock periods,
///         early-exit penalty, and pro-rata reward accrual.
contract TieredStaking is ReentrancyGuard, Ownable2Step, Pausable {
    using SafeERC20 for IERC20;

    // ─── Custom errors ────────────────────────────────────────────────────────

    error InvalidAmount();
    error NotPositionOwner();
    error PositionNotActive();
    error InvalidTier();

    // ─── Constants ───────────────────────────────────────────────────────────

    uint256 public constant PENALTY_BPS = 1000; // 10 %
    uint256 public constant BPS_DENOMINATOR = 10_000;
    uint256 public constant SECONDS_PER_YEAR = 365 days;

    // ─── Enums ────────────────────────────────────────────────────────────────

    enum Tier {
        Bronze, // 0
        Silver, // 1
        Gold    // 2
    }

    // ─── Structs ──────────────────────────────────────────────────────────────

    struct TierConfig {
        uint256 lockDuration; // seconds
        uint256 aprBps;       // annual percentage rate in basis points
    }

    struct Position {
        address owner;
        uint256 amount;
        Tier tier;
        uint256 stakedAt;
        uint256 lockUntil;
        uint256 lastClaimedAt;
        bool active;
    }

    // ─── State ────────────────────────────────────────────────────────────────

    IERC20 public immutable stakingToken;
    address public treasury;

    uint256 public nextPositionId;

    mapping(uint256 => Position) private _positions;
    mapping(Tier => TierConfig) private _tierConfigs;

    // ─── Events ───────────────────────────────────────────────────────────────

    event Staked(address indexed user, uint256 amount, Tier tier, uint256 positionId);
    event Unstaked(address indexed user, uint256 amount, uint256 rewards, uint256 penalty);
    event RewardsClaimed(address indexed user, uint256 amount, uint256 positionId);
    event TierUpdated(Tier tier, uint256 newAprBps);
    event TreasuryUpdated(address newTreasury);

    // ─── Constructor ─────────────────────────────────────────────────────────

    constructor(address _stakingToken, address _treasury) Ownable(msg.sender) {
        require(_stakingToken != address(0), "zero token");
        require(_treasury != address(0), "zero treasury");

        stakingToken = IERC20(_stakingToken);
        treasury = _treasury;

        // Tier configurations
        _tierConfigs[Tier.Bronze] = TierConfig({lockDuration: 30 days,  aprBps: 500});
        _tierConfigs[Tier.Silver] = TierConfig({lockDuration: 60 days,  aprBps: 1000});
        _tierConfigs[Tier.Gold]   = TierConfig({lockDuration: 90 days,  aprBps: 1800});
    }

    // ─── External functions ───────────────────────────────────────────────────

    /// @notice Stake `amount` tokens in the chosen `tier`.
    /// @return positionId The ID of the newly created position.
    function stake(uint256 amount, Tier tier) external nonReentrant whenNotPaused returns (uint256 positionId) {
        if (amount == 0) revert InvalidAmount();
        if (uint8(tier) > 2) revert InvalidTier();

        TierConfig storage config = _tierConfigs[tier];

        positionId = nextPositionId++;

        _positions[positionId] = Position({
            owner:         msg.sender,
            amount:        amount,
            tier:          tier,
            stakedAt:      block.timestamp,
            lockUntil:     block.timestamp + config.lockDuration,
            lastClaimedAt: block.timestamp,
            active:        true
        });

        stakingToken.safeTransferFrom(msg.sender, address(this), amount);

        emit Staked(msg.sender, amount, tier, positionId);
    }

    /// @notice Unstake a position. Returns principal + rewards if lock expired;
    ///         otherwise applies a 10 % early-exit penalty (sent to treasury).
    function unstake(uint256 positionId) external nonReentrant whenNotPaused {
        Position storage pos = _positions[positionId];

        if (pos.owner != msg.sender) revert NotPositionOwner();
        if (!pos.active)             revert PositionNotActive();

        // Checks-Effects-Interactions: mark inactive before any transfer
        pos.active = false;

        uint256 principal = pos.amount;
        uint256 rewards   = _calculateRewards(pos);
        uint256 penalty   = 0;
        uint256 payout;

        if (block.timestamp >= pos.lockUntil) {
            // Lock expired — full principal + accumulated rewards
            payout = principal + rewards;
        } else {
            // Early exit — 10 % penalty on principal, rewards are forfeit
            penalty = (principal * PENALTY_BPS) / BPS_DENOMINATOR;
            payout  = principal - penalty;
            rewards = 0; // forfeited on early exit
        }

        // Transfer penalty to treasury first (if any)
        if (penalty > 0) {
            stakingToken.safeTransfer(treasury, penalty);
        }

        stakingToken.safeTransfer(msg.sender, payout);

        emit Unstaked(msg.sender, principal, rewards, penalty);
    }

    /// @notice Claim accrued rewards for a position without unstaking.
    function claimRewards(uint256 positionId) external nonReentrant whenNotPaused {
        Position storage pos = _positions[positionId];

        if (pos.owner != msg.sender) revert NotPositionOwner();
        if (!pos.active)             revert PositionNotActive();

        uint256 rewards = _calculateRewards(pos);

        // Update checkpoint before transfer (CEI)
        pos.lastClaimedAt = block.timestamp;

        if (rewards > 0) {
            stakingToken.safeTransfer(msg.sender, rewards);
        }

        emit RewardsClaimed(msg.sender, rewards, positionId);
    }

    // ─── View functions ───────────────────────────────────────────────────────

    /// @notice Returns the full Position struct for a given positionId.
    function getPosition(uint256 positionId) external view returns (Position memory) {
        return _positions[positionId];
    }

    /// @notice Returns the APR in basis points for a given tier.
    function getAPR(Tier tier) external view returns (uint256) {
        if (uint8(tier) > 2) revert InvalidTier();
        return _tierConfigs[tier].aprBps;
    }

    // ─── Owner functions ──────────────────────────────────────────────────────

    function pause() external onlyOwner {
        _pause();
    }

    function unpause() external onlyOwner {
        _unpause();
    }

    /// @notice Update the APR for a tier. Does not affect existing positions.
    function updateTierAPR(Tier tier, uint256 newAprBps) external onlyOwner {
        if (uint8(tier) > 2) revert InvalidTier();
        _tierConfigs[tier].aprBps = newAprBps;
        emit TierUpdated(tier, newAprBps);
    }

    /// @notice Update the treasury address.
    function updateTreasury(address newTreasury) external onlyOwner {
        require(newTreasury != address(0), "zero treasury");
        treasury = newTreasury;
        emit TreasuryUpdated(newTreasury);
    }

    // ─── Internal helpers ─────────────────────────────────────────────────────

    /// @dev Simple linear reward: amount * aprBps * elapsed / (BPS_DENOMINATOR * SECONDS_PER_YEAR)
    function _calculateRewards(Position storage pos) internal view returns (uint256) {
        uint256 elapsed = block.timestamp - pos.lastClaimedAt;
        TierConfig storage config = _tierConfigs[pos.tier];
        return (pos.amount * config.aprBps * elapsed) / (BPS_DENOMINATOR * SECONDS_PER_YEAR);
    }
}
