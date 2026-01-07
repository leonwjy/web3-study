// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "./Auction.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

/**
 * @title AuctionUUPS
 * @dev UUPS 代理版本的拍卖合约
 */
contract AuctionUUPS is Auction, UUPSUpgradeable {
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
        __UUPSUpgradeable_init();
        require(_priceOracle != address(0), "Invalid price oracle address");
        priceOracle = PriceOracle(_priceOracle);
        version = 1;
    }
    
    /**
     * @dev 授权升级函数（仅 owner 可升级）
     * @param newImplementation 新实现合约地址
     */
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {
        // 可以在这里添加额外的升级验证逻辑
    }
    
    /**
     * @dev 升级后调用函数（可选）
     */
    function _updateVersion() internal {
        version++;
    }
}

