// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "./Auction.sol";

/**
 * @title AuctionTransparent
 * @dev 透明代理版本的拍卖合约 - Remix 简化版本
 * 注意：这个版本移除了 _disableInitializers() 以便在 Remix 中正常工作
 */
contract AuctionTransparent is Auction {
    
    /**
     * @dev 初始化函数（用于可升级合约）
     * @param _priceOracle 价格预言机合约地址
     */
    function initialize(address _priceOracle) public override initializer {
        // 调用父类 Auction 的初始化
        Auction.initialize(_priceOracle);
    }

    /**
     * @dev 获取合约版本
     */
    function version() external pure returns (uint256) {
        return 1;
    }
    
    /**
     * @dev 测试升级功能 - 新增功能
     */
    function newFunction() external pure returns (string memory) {
        return "This is a new function added in version 2";
    }
}

