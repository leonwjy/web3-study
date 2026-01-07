// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

/**
 * @title PriceOracle
 * @dev Chainlink 价格预言机封装合约，用于获取 ETH 和 ERC20 代币的 USD 价格
 */
contract PriceOracle {
    // ETH/USD 价格源地址
    AggregatorV3Interface public ethUsdPriceFeed;
    
    // ERC20 token 地址到 Chainlink feed 地址的映射
    mapping(address => address) public tokenPriceFeeds;
    
    // 价格过期时间（秒），默认 1 小时
    uint256 public constant PRICE_STALENESS_THRESHOLD = 3600;
    
    event TokenPriceFeedSet(address indexed token, address indexed priceFeed);
    
    /**
     * @dev 构造函数
     * @param _ethUsdPriceFeed ETH/USD Chainlink 价格源地址
     */
    constructor(address _ethUsdPriceFeed) {
        require(_ethUsdPriceFeed != address(0), "Invalid price feed address");
        ethUsdPriceFeed = AggregatorV3Interface(_ethUsdPriceFeed);
    }
    
    /**
     * @dev 设置 ERC20 token 的价格源地址
     * @param token ERC20 token 地址
     * @param priceFeed Chainlink 价格源地址
     */
    function setTokenPriceFeed(address token, address priceFeed) external {
        require(token != address(0), "Invalid token address");
        require(priceFeed != address(0), "Invalid price feed address");
        tokenPriceFeeds[token] = priceFeed;
        emit TokenPriceFeedSet(token, priceFeed);
    }
    
    /**
     * @dev 获取 ETH/USD 价格
     * @return price ETH 价格（8 位小数）
     * @return decimals 价格精度（通常为 8）
     */
    function getETHPrice() public view returns (int256 price, uint8 decimals) {
        (
            /* uint80 roundId */,
            int256 answer,
            /* uint256 startedAt */,
            uint256 updatedAt,
            /* uint80 answeredInRound */
        ) = ethUsdPriceFeed.latestRoundData();
        
        require(answer > 0, "Invalid price data");
        require(updatedAt > 0, "Round not complete");
        require(
            block.timestamp - updatedAt <= PRICE_STALENESS_THRESHOLD,
            "Price data is stale"
        );
        
        decimals = ethUsdPriceFeed.decimals();
        price = answer;
    }
    
    /**
     * @dev 获取 ERC20 token/USD 价格
     * @param token ERC20 token 地址
     * @return price Token 价格（8 位小数）
     * @return decimals 价格精度（通常为 8）
     */
    function getERC20Price(address token) public view returns (int256 price, uint8 decimals) {
        address priceFeedAddress = tokenPriceFeeds[token];
        require(priceFeedAddress != address(0), "Price feed not set for token");
        
        AggregatorV3Interface priceFeed = AggregatorV3Interface(priceFeedAddress);
        
        (
            /* uint80 roundId */,
            int256 answer,
            /* uint256 startedAt */,
            uint256 updatedAt,
            /* uint80 answeredInRound */
        ) = priceFeed.latestRoundData();
        
        require(answer > 0, "Invalid price data");
        require(updatedAt > 0, "Round not complete");
        require(
            block.timestamp - updatedAt <= PRICE_STALENESS_THRESHOLD,
            "Price data is stale"
        );
        
        decimals = priceFeed.decimals();
        price = answer;
    }
    
    /**
     * @dev 将代币数量转换为美元价值
     * @param amount 代币数量（以代币的最小单位计算，例如 wei）
     * @param token 代币地址（address(0) 表示 ETH）
     * @return usdValue 美元价值（以最小单位计算，需要除以 1e18 得到实际美元数）
     */
    function convertToUSD(uint256 amount, address token) external view returns (uint256 usdValue) {
        int256 price;
        uint8 priceDecimals;
        uint8 tokenDecimals = 18; // 默认 18 位小数
        
        if (token == address(0)) {
            // ETH
            (price, priceDecimals) = getETHPrice();
        } else {
            // ERC20 token
            (price, priceDecimals) = getERC20Price(token);
            // 尝试获取 token decimals，如果失败则使用默认值 18
            try IERC20Metadata(token).decimals() returns (uint8 decimals) {
                tokenDecimals = decimals;
            } catch {}
        }
        
        // 计算 USD 价值
        // amount * price / (10^tokenDecimals) * (10^priceDecimals) / (10^priceDecimals)
        // = amount * price / (10^tokenDecimals)
        // 为了保持精度，我们使用更大的精度进行计算
        // 最终结果以 18 位小数返回
        
        uint256 priceUint = uint256(price);
        
        // 如果 priceDecimals 是 8，tokenDecimals 是 18
        // 我们需要：amount * price / 10^tokenDecimals * 10^18 / 10^priceDecimals
        // = amount * price * 10^(18 - priceDecimals) / 10^tokenDecimals
        
        if (priceDecimals <= 18 && tokenDecimals <= 18) {
            uint256 adjustment = 10 ** (18 - priceDecimals);
            usdValue = (amount * priceUint * adjustment) / (10 ** tokenDecimals);
        } else {
            // 处理其他情况，使用更安全的计算方式
            usdValue = (amount * priceUint) / (10 ** tokenDecimals);
            if (priceDecimals < 18) {
                usdValue = usdValue * (10 ** (18 - priceDecimals));
            } else if (priceDecimals > 18) {
                usdValue = usdValue / (10 ** (priceDecimals - 18));
            }
        }
    }
}

// ERC20 元数据接口，用于获取代币精度
interface IERC20Metadata {
    function decimals() external view returns (uint8);
}

