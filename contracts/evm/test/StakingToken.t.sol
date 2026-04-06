// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "forge-std/Test.sol";
import "../src/StakingToken.sol";

contract StakingTokenTest is Test {
    StakingToken internal token;
    address internal owner;
    address internal alice;

    function setUp() public {
        owner = address(this);
        alice = makeAddr("alice");
        token = new StakingToken("Staking Token", "STK");
    }

    // ─── name / symbol ────────────────────────────────────────────────────────

    function test_Name_ReturnsCorrectName() public view {
        assertEq(token.name(), "Staking Token");
    }

    function test_Symbol_ReturnsCorrectSymbol() public view {
        assertEq(token.symbol(), "STK");
    }

    // ─── owner mint ───────────────────────────────────────────────────────────

    function test_Mint_OwnerCanMint() public {
        token.mint(alice, 1_000 ether);
        assertEq(token.balanceOf(alice), 1_000 ether);
    }

    function test_Mint_TotalSupplyIncreasesOnMint() public {
        token.mint(alice, 500 ether);
        assertEq(token.totalSupply(), 500 ether);
    }

    // ─── non-owner revert ─────────────────────────────────────────────────────

    function test_Mint_NonOwnerReverts() public {
        vm.prank(alice);
        vm.expectRevert();
        token.mint(alice, 1_000 ether);
    }

    // ─── fuzz ─────────────────────────────────────────────────────────────────

    function testFuzz_Mint_OwnerCanMintArbitraryAmount(uint256 amount) public {
        // Avoid overflow when comparing total supply
        vm.assume(amount <= type(uint256).max / 2);
        token.mint(alice, amount);
        assertEq(token.balanceOf(alice), amount);
    }
}
