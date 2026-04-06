// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "forge-std/Script.sol";
import "../src/StakingToken.sol";
import "../src/TieredStaking.sol";

/// @notice Deploys StakingToken and TieredStaking, then:
///         - Funds the staking contract with 1 000 000 tokens (reward pool)
///         - Mints 100 000 tokens to the deployer
contract Deploy is Script {
    function run() external {
        uint256 deployerKey = vm.envUint("PRIVATE_KEY");
        address deployer    = vm.addr(deployerKey);

        // Use deployer as initial treasury (can be changed via updateTreasury)
        address treasury = deployer;

        vm.startBroadcast(deployerKey);

        // 1. Deploy token
        StakingToken token = new StakingToken("Staking Token", "STK");
        console2.log("StakingToken deployed at:", address(token));

        // 2. Deploy staking contract
        TieredStaking staking = new TieredStaking(address(token), treasury);
        console2.log("TieredStaking deployed at:", address(staking));

        // 3. Fund staking contract with 1 000 000 tokens
        token.mint(address(staking), 1_000_000 ether);
        console2.log("Funded staking contract with 1 000 000 STK");

        // 4. Mint 100 000 tokens to deployer
        token.mint(deployer, 100_000 ether);
        console2.log("Minted 100 000 STK to deployer:", deployer);

        vm.stopBroadcast();
    }
}
