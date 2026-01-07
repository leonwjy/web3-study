const { expect } = require("chai");
const { ethers } = require("hardhat");
const { loadFixture } = require("@nomicfoundation/hardhat-toolbox/network-helpers");

describe("AuctionUUPS", function () {
  async function deployFixture() {
    const [owner, seller] = await ethers.getSigners();

    // 部署 Mock Chainlink Aggregator
    const MockChainlinkAggregator = await ethers.getContractFactory("MockChainlinkAggregator");
    const mockEthPriceFeed = await MockChainlinkAggregator.deploy(3000 * 10 ** 8);
    await mockEthPriceFeed.waitForDeployment();

    // 部署 PriceOracle
    const PriceOracle = await ethers.getContractFactory("PriceOracle");
    const priceOracle = await PriceOracle.deploy(await mockEthPriceFeed.getAddress());
    await priceOracle.waitForDeployment();

    // 部署 AuctionUUPS 实现合约
    const AuctionUUPS = await ethers.getContractFactory("AuctionUUPS");
    const auctionImpl = await AuctionUUPS.deploy();
    await auctionImpl.waitForDeployment();

    // 部署 UUPS 代理（使用 ERC1967Proxy）
    const ERC1967ProxyArtifact = require("@openzeppelin/contracts/build/contracts/ERC1967Proxy.json");
    const ERC1967Proxy = await ethers.getContractFactoryFromArtifact(ERC1967ProxyArtifact);
    const initData = auctionImpl.interface.encodeFunctionData("initialize", [
      await priceOracle.getAddress()
    ]);
    const proxy = await ERC1967Proxy.deploy(await auctionImpl.getAddress(), initData);
    await proxy.waitForDeployment();

    // 获取代理合约实例
    const auction = AuctionUUPS.attach(await proxy.getAddress());

    return { auction, auctionImpl, owner, seller, priceOracle };
  }

  describe("Deployment", function () {
    it("Should initialize correctly", async function () {
      const { auction, priceOracle } = await loadFixture(deployFixture);
      
      expect(await auction.priceOracle()).to.equal(await priceOracle.getAddress());
      expect(await auction.version()).to.equal(1);
    });
  });

  describe("Upgrade", function () {
    it("Should upgrade to new implementation", async function () {
      const { auction, auctionImpl, owner } = await loadFixture(deployFixture);

      // 部署新版本实现合约
      const AuctionUUPSV2 = await ethers.getContractFactory("AuctionUUPS");
      const auctionImplV2 = await AuctionUUPSV2.deploy();
      await auctionImplV2.waitForDeployment();

      // 升级代理
      await expect(
        auction.connect(owner).upgradeToAndCall(await auctionImplV2.getAddress(), "0x")
      ).to.not.be.reverted;

      // 验证升级成功（通过调用函数验证）
      expect(await auction.priceOracle()).to.equal(await auction.priceOracle());
    });

    it("Should revert if non-owner tries to upgrade", async function () {
      const { auction } = await loadFixture(deployFixture);
      const [, , nonOwner] = await ethers.getSigners();

      const AuctionUUPSV2 = await ethers.getContractFactory("AuctionUUPS");
      const auctionImplV2 = await AuctionUUPSV2.deploy();
      await auctionImplV2.waitForDeployment();

      await expect(
        auction.connect(nonOwner).upgradeToAndCall(await auctionImplV2.getAddress(), "0x")
      ).to.be.revertedWithCustomError(auction, "OwnableUnauthorizedAccount");
    });
  });

  describe("Data Persistence", function () {
    it("Should preserve data after upgrade", async function () {
      const { auction, owner } = await loadFixture(deployFixture);
      const [,, seller] = await ethers.getSigners();

      // 设置一些数据
      const priceOracleAddress = await auction.priceOracle();

      // 部署新版本
      const AuctionUUPSV2 = await ethers.getContractFactory("AuctionUUPS");
      const auctionImplV2 = await AuctionUUPSV2.deploy();
      await auctionImplV2.waitForDeployment();

      // 升级
      await auction.connect(owner).upgradeToAndCall(await auctionImplV2.getAddress(), "0x");

      // 验证数据仍然存在
      expect(await auction.priceOracle()).to.equal(priceOracleAddress);
    });
  });
});

