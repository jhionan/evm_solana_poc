// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "forge-std/Script.sol";
import "../src/StakingToken.sol";
import "../src/TieredStaking.sol";

/// @notice Deploys to a public testnet (Sepolia, LightLink Pegasus, etc.)
///         Usage: PRIVATE_KEY=0x... forge script script/DeployTestnet.s.sol \
///                --rpc-url <RPC_URL> --broadcast --verify
contract DeployTestnet is Script {
    function run() external {
        uint256 deployerKey = vm.envUint("PRIVATE_KEY");
        address deployer = vm.addr(deployerKey);
        address treasury = deployer;

        console2.log("Deployer:", deployer);
        console2.log("Chain ID:", block.chainid);

        vm.startBroadcast(deployerKey);

        // 1. Deploy token
        StakingToken token = new StakingToken("LightLink Staking Token", "LLSTK");
        console2.log("StakingToken:", address(token));

        // 2. Deploy staking contract
        TieredStaking staking = new TieredStaking(address(token), treasury);
        console2.log("TieredStaking:", address(staking));

        // 3. Fund staking contract with reward tokens (100K for testnet)
        token.mint(address(staking), 100_000 ether);
        console2.log("Funded staking contract with 100,000 LLSTK");

        // 4. Mint test tokens to deployer (10K)
        token.mint(deployer, 10_000 ether);
        console2.log("Minted 10,000 LLSTK to deployer");

        vm.stopBroadcast();

        // Summary for .env configuration
        console2.log("---");
        console2.log("Add to .env:");
        console2.log("EVM_TOKEN_CONTRACT=", address(token));
        console2.log("EVM_STAKING_CONTRACT=", address(staking));
    }
}
