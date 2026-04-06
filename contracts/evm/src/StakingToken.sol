// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/// @title StakingToken
/// @notice Simple ERC-20 token mintable by the owner, used as the staking / reward token
contract StakingToken is ERC20, Ownable {
    constructor(string memory name, string memory symbol)
        ERC20(name, symbol)
        Ownable(msg.sender)
    {}

    /// @notice Mint `amount` tokens to `to`. Only callable by the owner.
    function mint(address to, uint256 amount) external onlyOwner {
        _mint(to, amount);
    }
}
