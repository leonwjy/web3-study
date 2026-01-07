// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "./Auction.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";

/**
 * @title AuctionV2
 * @dev 拍卖合约升级版本 - 透明代理升级用
 * 注意：为保持存储兼容性，不添加新的存储变量
 */
contract AuctionV2 is Auction {
    // V2 初始化标志（避免存储布局冲突，使用常量）
    bool private _v2Initialized;

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /**
     * @dev 升级时调用的初始化函数
     */
    function initializeV2() public {
        require(!_v2Initialized, "Already initialized");
        _v2Initialized = true;
        // 升级时的初始化逻辑
        // 注意：不添加新的存储变量以保持兼容性
    }

    /**
     * @dev 新功能：获取 V2 标识
     */
    function getVersion() external pure returns (string memory) {
        return "AuctionV2";
    }

    /**
     * @dev 新功能：获取合约统计信息
     */
    function getAuctionStats() external view returns (
        uint256 totalAuctions,
        uint256 activeAuctions
    ) {
        uint256 total = _auctionIdCounter;
        uint256 active = 0;

        for (uint256 i = 1; i <= total; i++) {
            if (!auctions[i].ended && block.timestamp < auctions[i].endTime) {
                active++;
            }
        }

        return (total, active);
    }

    /**
     * @dev 重写创建拍卖函数，添加额外验证
     */
    function createAuction(
        address nftContract,
        uint256 tokenId,
        uint256 startingPriceUSD,
        uint256 duration,
        address paymentToken
    ) public override {
        // 添加 V2 的额外验证
        require(startingPriceUSD >= 10**18, "Starting price must be at least 1 USD");

        // 调用父类函数
        super.createAuction(nftContract, tokenId, startingPriceUSD, duration, paymentToken);
    }
}

/**
 * @title AuctionUUPSV2
 * @dev UUPS 升级版本的拍卖合约
 * 注意：UUPS 升级需要实现合约也继承 UUPSUpgradeable
 */
contract AuctionUUPSV2 is Auction, UUPSUpgradeable {
    // V2 初始化标志
    bool private _v2Initialized;

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    /**
     * @dev 初始化函数（用于可升级合约）
     */
    function initialize(address _priceOracle) public override initializer {
        Auction.initialize(_priceOracle);
        __UUPSUpgradeable_init();
    }

    /**
     * @dev 升级时调用的初始化函数
     */
    function initializeV2() public {
        require(!_v2Initialized, "Already initialized");
        _v2Initialized = true;
        // UUPS 升级时的初始化逻辑
    }

    /**
     * @dev 授权升级函数（仅 owner 可升级）
     */
    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {
        // UUPS 升级授权逻辑
    }

    /**
     * @dev 新功能：获取 V2 标识
     */
    function getVersion() external pure returns (string memory) {
        return "AuctionUUPSV2";
    }

    /**
     * @dev 新功能：获取合约统计信息
     */
    function getAuctionStats() external view returns (
        uint256 totalAuctions,
        uint256 activeAuctions
    ) {
        uint256 total = _auctionIdCounter;
        uint256 active = 0;

        for (uint256 i = 1; i <= total; i++) {
            if (!auctions[i].ended && block.timestamp < auctions[i].endTime) {
                active++;
            }
        }

        return (total, active);
    }

    /**
     * @dev 重写创建拍卖函数，添加额外验证
     */
    function createAuction(
        address nftContract,
        uint256 tokenId,
        uint256 startingPriceUSD,
        uint256 duration,
        address paymentToken
    ) public override {
        // 添加 V2 的额外验证
        require(startingPriceUSD >= 10**18, "Starting price must be at least 1 USD");

        // 调用父类函数
        super.createAuction(nftContract, tokenId, startingPriceUSD, duration, paymentToken);
    }
}
