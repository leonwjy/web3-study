// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "./Auction.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

/**
 * @title AuctionUUPS
 * @dev UUPS 代理版本的拍卖合约 - Remix 简化版本
 */
contract AuctionUUPS is Auction, UUPSUpgradeable {
    
    /**
     * @dev 初始化函数（用于可升级合约）
     * @param _priceOracle 价格预言机合约地址
     */
    function initialize(address _priceOracle) public override initializer {
        // 调用父类 Auction 的初始化
        Auction.initialize(_priceOracle);
        __UUPSUpgradeable_init();
    }

    /**
     * @dev 获取合约版本
     */
    function version() external pure returns (uint256) {
        return 1;
    }
    
    /**
     * @dev 授权升级函数（仅 owner 可升级）
     * @param newImplementation 新实现合约地址
     */
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {
        // 可以在这里添加额外的升级验证逻辑
    }
}

