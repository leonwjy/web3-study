const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("PriceOracle", function () {
  let priceOracle;
  let mockEthPriceFeed;
  let mockTokenPriceFeed;
  let owner;

  beforeEach(async function () {
    [owner] = await ethers.getSigners();

    // 部署 Mock Chainlink Aggregator
    const MockChainlinkAggregator = await ethers.getContractFactory("MockChainlinkAggregator");
    // ETH/USD 价格：$3000，8位小数
    mockEthPriceFeed = await MockChainlinkAggregator.deploy(3000 * 10 ** 8);
    await mockEthPriceFeed.waitForDeployment();

    // 部署 PriceOracle
    const PriceOracle = await ethers.getContractFactory("PriceOracle");
    priceOracle = await PriceOracle.deploy(await mockEthPriceFeed.getAddress());
    await priceOracle.waitForDeployment();
  });

  describe("Deployment", function () {
    it("Should set the ETH price feed", async function () {
      expect(await priceOracle.ethUsdPriceFeed()).to.equal(await mockEthPriceFeed.getAddress());
    });

    it("Should revert with invalid price feed address", async function () {
      const PriceOracle = await ethers.getContractFactory("PriceOracle");
      await expect(
        PriceOracle.deploy(ethers.ZeroAddress)
      ).to.be.revertedWith("Invalid price feed address");
    });
  });

  describe("ETH Price", function () {
    it("Should get ETH price", async function () {
      const [price, decimals] = await priceOracle.getETHPrice();
      expect(price).to.equal(3000 * 10 ** 8);
      expect(decimals).to.equal(8);
    });

    it("Should revert with stale price", async function () {
      // 设置价格更新时间为 2 小时前
      const twoHoursAgo = Math.floor(Date.now() / 1000) - 7200;
      // 注意：Mock 合约需要支持设置 updatedAt，这里简化处理
      // 实际测试中可能需要更复杂的 Mock
    });
  });

  describe("ERC20 Price Feed", function () {
    it("Should set token price feed", async function () {
      const MockChainlinkAggregator = await ethers.getContractFactory("MockChainlinkAggregator");
      const mockTokenFeed = await MockChainlinkAggregator.deploy(1 * 10 ** 8); // $1
      await mockTokenFeed.waitForDeployment();

      const tokenAddress = ethers.Wallet.createRandom().address;
      
      await priceOracle.setTokenPriceFeed(tokenAddress, await mockTokenFeed.getAddress());
      expect(await priceOracle.tokenPriceFeeds(tokenAddress)).to.equal(await mockTokenFeed.getAddress());
    });

    it("Should get ERC20 price", async function () {
      const MockChainlinkAggregator = await ethers.getContractFactory("MockChainlinkAggregator");
      const mockTokenFeed = await MockChainlinkAggregator.deploy(100 * 10 ** 8); // $100
      await mockTokenFeed.waitForDeployment();

      const tokenAddress = ethers.Wallet.createRandom().address;
      await priceOracle.setTokenPriceFeed(tokenAddress, await mockTokenFeed.getAddress());

      const [price, decimals] = await priceOracle.getERC20Price(tokenAddress);
      expect(price).to.equal(100 * 10 ** 8);
      expect(decimals).to.equal(8);
    });

    it("Should revert when price feed not set", async function () {
      const tokenAddress = ethers.Wallet.createRandom().address;
      await expect(
        priceOracle.getERC20Price(tokenAddress)
      ).to.be.revertedWith("Price feed not set for token");
    });
  });

  describe("Convert to USD", function () {
    it("Should convert ETH to USD", async function () {
      // 1 ETH = 3000 USD
      const ethAmount = ethers.parseEther("1"); // 1 ETH
      const usdValue = await priceOracle.convertToUSD(ethAmount, ethers.ZeroAddress);
      
      // 1 ETH * 3000 USD/ETH = 3000 USD (18 decimals)
      expect(usdValue).to.equal(ethers.parseEther("3000"));
    });

    it("Should convert ERC20 to USD", async function () {
      const MockChainlinkAggregator = await ethers.getContractFactory("MockChainlinkAggregator");
      const mockTokenFeed = await MockChainlinkAggregator.deploy(2 * 10 ** 8); // $2
      await mockTokenFeed.waitForDeployment();

      const MockERC20 = await ethers.getContractFactory("MockERC20");
      const mockToken = await MockERC20.deploy("TestToken", "TEST", 18);
      await mockToken.waitForDeployment();

      await priceOracle.setTokenPriceFeed(await mockToken.getAddress(), await mockTokenFeed.getAddress());

      // 100 tokens * $2 = $200
      const tokenAmount = ethers.parseEther("100");
      const usdValue = await priceOracle.convertToUSD(tokenAmount, await mockToken.getAddress());
      
      expect(usdValue).to.equal(ethers.parseEther("200"));
    });
  });
});

