// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "./Auction.sol";

/**
 * @title AuctionTransparent
 * @dev 透明代理版本的拍卖合约
 * 注意：透明代理模式不需要在实现合约中继承任何代理相关的合约
 * 代理逻辑由部署时使用的 TransparentUpgradeableProxy 合约处理
 */
contract AuctionTransparent is Auction {
    // 版本号，用于测试升级
    uint256 public version;
    
    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }
    
    /**
     * @dev 初始化函数（用于可升级合约）
     * @param _priceOracle 价格预言机合约地址
     */
    function initialize(address _priceOracle) public override initializer {
        __Ownable_init(msg.sender);
        __ReentrancyGuard_init();
        require(_priceOracle != address(0), "Invalid price oracle address");
        priceOracle = PriceOracle(_priceOracle);
        version = 1;
    }
    
    /**
     * @dev 升级后调用函数（可选）
     */
    function updateVersion() external {
        version++;
    }
}

