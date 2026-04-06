// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "forge-std/Test.sol";
import "../src/StakingToken.sol";
import "../src/TieredStaking.sol";

contract TieredStakingTest is Test {
    StakingToken internal token;
    TieredStaking internal staking;

    address internal owner;
    address internal user;
    address internal treasury;

    uint256 internal constant REWARD_FUND   = 1_000_000 ether;
    uint256 internal constant USER_BALANCE  = 100_000 ether;
    uint256 internal constant STAKE_AMOUNT  = 10_000 ether;

    function setUp() public {
        owner    = address(this);
        user     = makeAddr("user");
        treasury = makeAddr("treasury");

        // Deploy contracts
        token   = new StakingToken("Staking Token", "STK");
        staking = new TieredStaking(address(token), treasury);

        // Fund staking contract with reward tokens
        token.mint(address(staking), REWARD_FUND);

        // Fund user and approve staking contract
        token.mint(user, USER_BALANCE);
        vm.prank(user);
        token.approve(address(staking), type(uint256).max);
    }

    // ─── Helper ───────────────────────────────────────────────────────────────

    function _stakeAs(address who, uint256 amount, TieredStaking.Tier tier) internal returns (uint256) {
        vm.prank(who);
        return staking.stake(amount, tier);
    }

    // ─── stake ────────────────────────────────────────────────────────────────

    function test_Stake_Bronze_CreatesPosition() public {
        uint256 posId = _stakeAs(user, STAKE_AMOUNT, TieredStaking.Tier.Bronze);
        TieredStaking.Position memory pos = staking.getPosition(posId);

        assertEq(pos.owner,  user);
        assertEq(pos.amount, STAKE_AMOUNT);
        assertEq(uint8(pos.tier), uint8(TieredStaking.Tier.Bronze));
        assertTrue(pos.active);
    }

    function test_Stake_Gold_VerifiesLockUntil() public {
        uint256 before = block.timestamp;
        uint256 posId  = _stakeAs(user, STAKE_AMOUNT, TieredStaking.Tier.Gold);
        TieredStaking.Position memory pos = staking.getPosition(posId);

        assertEq(pos.lockUntil, before + 90 days);
    }

    function test_Stake_ZeroAmount_Reverts() public {
        vm.prank(user);
        vm.expectRevert(TieredStaking.InvalidAmount.selector);
        staking.stake(0, TieredStaking.Tier.Bronze);
    }

    function test_Stake_EmitsEvent() public {
        vm.prank(user);
        vm.expectEmit(true, false, false, true);
        emit TieredStaking.Staked(user, STAKE_AMOUNT, TieredStaking.Tier.Bronze, 0);
        staking.stake(STAKE_AMOUNT, TieredStaking.Tier.Bronze);
    }

    // ─── unstake ──────────────────────────────────────────────────────────────

    function test_Unstake_AfterLock_GetsRewards() public {
        uint256 posId = _stakeAs(user, STAKE_AMOUNT, TieredStaking.Tier.Bronze);

        // Fast-forward past lock period
        vm.warp(block.timestamp + 31 days);

        uint256 balBefore = token.balanceOf(user);
        vm.prank(user);
        staking.unstake(posId);
        uint256 balAfter = token.balanceOf(user);

        uint256 received = balAfter - balBefore;
        // Should receive at least the principal back
        assertGe(received, STAKE_AMOUNT);
        // Should receive some rewards (30+ days at 5% APR)
        assertGt(received, STAKE_AMOUNT);
    }

    function test_Unstake_Early_AppliesPenalty() public {
        uint256 posId = _stakeAs(user, STAKE_AMOUNT, TieredStaking.Tier.Bronze);

        // Warp to halfway through the lock
        vm.warp(block.timestamp + 10 days);

        uint256 userBalBefore     = token.balanceOf(user);
        uint256 treasuryBalBefore = token.balanceOf(treasury);

        vm.prank(user);
        staking.unstake(posId);

        uint256 userReceived    = token.balanceOf(user)     - userBalBefore;
        uint256 treasuryReceived = token.balanceOf(treasury) - treasuryBalBefore;

        uint256 expectedPenalty = (STAKE_AMOUNT * 1000) / 10_000; // 10%
        assertEq(treasuryReceived, expectedPenalty,              "penalty to treasury");
        assertEq(userReceived,     STAKE_AMOUNT - expectedPenalty, "user receives rest");
    }

    function test_Unstake_NotOwner_Reverts() public {
        uint256 posId = _stakeAs(user, STAKE_AMOUNT, TieredStaking.Tier.Bronze);

        address attacker = makeAddr("attacker");
        vm.prank(attacker);
        vm.expectRevert(TieredStaking.NotPositionOwner.selector);
        staking.unstake(posId);
    }

    function test_Unstake_AlreadyUnstaked_Reverts() public {
        uint256 posId = _stakeAs(user, STAKE_AMOUNT, TieredStaking.Tier.Bronze);

        vm.warp(block.timestamp + 31 days);

        vm.prank(user);
        staking.unstake(posId);

        vm.prank(user);
        vm.expectRevert(TieredStaking.PositionNotActive.selector);
        staking.unstake(posId);
    }

    // ─── claimRewards ─────────────────────────────────────────────────────────

    function test_ClaimRewards_AccruesCorrectly() public {
        uint256 posId = _stakeAs(user, STAKE_AMOUNT, TieredStaking.Tier.Bronze);

        // 30 days elapsed
        vm.warp(block.timestamp + 30 days);

        uint256 balBefore = token.balanceOf(user);
        vm.prank(user);
        staking.claimRewards(posId);
        uint256 rewards = token.balanceOf(user) - balBefore;

        // Expected: 10_000e18 * 500 * 30 days / (10_000 * 365 days)
        uint256 expected = (STAKE_AMOUNT * 500 * 30 days) / (10_000 * 365 days);
        assertApproxEqAbs(rewards, expected, 1e15); // within 0.001 token
    }

    // ─── pause / unpause ──────────────────────────────────────────────────────

    function test_Pause_BlocksStaking() public {
        staking.pause();

        vm.prank(user);
        vm.expectRevert();
        staking.stake(STAKE_AMOUNT, TieredStaking.Tier.Bronze);
    }

    function test_Unpause_AllowsStaking() public {
        staking.pause();
        staking.unpause();

        uint256 posId = _stakeAs(user, STAKE_AMOUNT, TieredStaking.Tier.Bronze);
        TieredStaking.Position memory pos = staking.getPosition(posId);
        assertTrue(pos.active);
    }

    // ─── getAPR ───────────────────────────────────────────────────────────────

    function test_GetAPR_Bronze_Returns500() public view {
        assertEq(staking.getAPR(TieredStaking.Tier.Bronze), 500);
    }

    function test_GetAPR_Silver_Returns1000() public view {
        assertEq(staking.getAPR(TieredStaking.Tier.Silver), 1000);
    }

    function test_GetAPR_Gold_Returns1800() public view {
        assertEq(staking.getAPR(TieredStaking.Tier.Gold), 1800);
    }

    // ─── fuzz stake amount ────────────────────────────────────────────────────

    function testFuzz_Stake_Amount(uint256 amount) public {
        // Keep amount within user balance and > 0
        vm.assume(amount > 0 && amount <= USER_BALANCE);

        uint256 posId = _stakeAs(user, amount, TieredStaking.Tier.Silver);
        TieredStaking.Position memory pos = staking.getPosition(posId);

        assertEq(pos.amount, amount);
        assertTrue(pos.active);
    }
}
