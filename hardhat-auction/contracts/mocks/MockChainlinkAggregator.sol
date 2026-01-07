// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

/**
 * @title MockChainlinkAggregator
 * @dev Mock Chainlink 价格聚合器，用于测试
 */
contract MockChainlinkAggregator is AggregatorV3Interface {
    uint8 public constant override decimals = 8;
    string public constant override description = "Mock Price Feed";
    uint256 public constant override version = 0;
    
    int256 private _price;
    uint256 private _updatedAt;
    uint80 private _roundId;
    
    constructor(int256 initialPrice) {
        _price = initialPrice;
        _updatedAt = block.timestamp;
        _roundId = 1;
    }
    
    function setPrice(int256 newPrice) external {
        _price = newPrice;
        _updatedAt = block.timestamp;
        _roundId++;
    }
    
    function latestRoundData()
        external
        view
        override
        returns (
            uint80 roundId,
            int256 answer,
            uint256 startedAt,
            uint256 updatedAt,
            uint80 answeredInRound
        )
    {
        return (_roundId, _price, _updatedAt, _updatedAt, _roundId);
    }
    
    function getRoundData(uint80 inputRoundId)
        external
        view
        override
        returns (
            uint80 roundId,
            int256 answer,
            uint256 startedAt,
            uint256 updatedAt,
            uint80 answeredInRound
        )
    {
        return (_roundId, _price, _updatedAt, _updatedAt, _roundId);
    }
}

