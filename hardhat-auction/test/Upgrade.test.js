const { expect } = require("chai");
const { ethers, upgrades } = require("hardhat");

describe("合约升级测试 (Hardhat Upgrades)", function () {
  let deployer, seller, bidder;
  let myNFT, priceOracle, mockERC20;
  let auctionTransparent, auctionTransparentV2;
  let auctionUUPS, auctionUUPSV2;

  before(async function () {
    [deployer, seller, bidder] = await ethers.getSigners();

    // 部署基础合约
    const MyNFT = await ethers.getContractFactory("MyNFT");
    myNFT = await MyNFT.deploy("MyNFT", "MNFT");
    await myNFT.waitForDeployment();

    const MockChainlinkAggregator = await ethers.getContractFactory("MockChainlinkAggregator");
    const mockEthPriceFeed = await MockChainlinkAggregator.deploy(300000000000); // $3000
    await mockEthPriceFeed.waitForDeployment();

    const PriceOracle = await ethers.getContractFactory("PriceOracle");
    priceOracle = await PriceOracle.deploy(await mockEthPriceFeed.getAddress());
    await priceOracle.waitForDeployment();

    const MockERC20 = await ethers.getContractFactory("MockERC20");
    mockERC20 = await MockERC20.deploy("TestToken", "TEST", 18);
    await mockERC20.waitForDeployment();

    await priceOracle.setTokenPriceFeed(await mockERC20.getAddress(), await mockEthPriceFeed.getAddress());
  });

  describe("透明代理升级测试", function () {
    it("应该成功部署透明代理", async function () {
      const AuctionTransparent = await ethers.getContractFactory("AuctionTransparent");

      auctionTransparent = await upgrades.deployProxy(AuctionTransparent, [await priceOracle.getAddress()], {
        initializer: "initialize",
        kind: "transparent",
      });

      await auctionTransparent.waitForDeployment();

      expect(await auctionTransparent.priceOracle()).to.equal(await priceOracle.getAddress());
    });

    it("应该成功升级到 V2", async function () {
      const AuctionV2 = await ethers.getContractFactory("AuctionV2");

      auctionTransparentV2 = await upgrades.upgradeProxy(await auctionTransparent.getAddress(), AuctionV2, {
        call: { fn: "initializeV2" },
      });

      await auctionTransparentV2.waitForDeployment();

      expect(await auctionTransparentV2.getVersion()).to.equal("AuctionV2");
      expect(await auctionTransparentV2.getAuctionStats()).to.not.be.undefined;
    });

    it("应该保持原有状态", async function () {
      // 验证状态保持
      expect(await auctionTransparentV2.priceOracle()).to.equal(await priceOracle.getAddress());
      expect(await auctionTransparentV2.owner()).to.equal(deployer.address);
    });

    it("应该支持新功能", async function () {
      // 测试新功能
      const stats = await auctionTransparentV2.getAuctionStats();
      expect(stats.totalAuctions).to.equal(0);
      expect(stats.activeAuctions).to.equal(0);

      // 测试 V2 的增强验证（起拍价至少 1 USD）
      await expect(
        auctionTransparentV2.connect(seller).createAuction(
          await myNFT.getAddress(),
          1,
          ethers.parseEther("0.5"), // 小于 1 USD
          3600,
          ethers.ZeroAddress
        )
      ).to.be.revertedWith("Starting price must be at least 1 USD");
    });
  });

  describe("UUPS代理升级测试", function () {
    it("应该成功部署 UUPS 代理", async function () {
      const AuctionUUPS = await ethers.getContractFactory("AuctionUUPS");

      auctionUUPS = await upgrades.deployProxy(AuctionUUPS, [await priceOracle.getAddress()], {
        initializer: "initialize",
        kind: "uups",
      });

      await auctionUUPS.waitForDeployment();

      expect(await auctionUUPS.priceOracle()).to.equal(await priceOracle.getAddress());
    });

    it("应该成功升级到 V2", async function () {
      const AuctionUUPSV2 = await ethers.getContractFactory("AuctionUUPSV2");

      auctionUUPSV2 = await upgrades.upgradeProxy(await auctionUUPS.getAddress(), AuctionUUPSV2, {
        call: { fn: "initializeV2" },
        kind: "uups",
      });

      await auctionUUPSV2.waitForDeployment();

      expect(await auctionUUPSV2.getVersion()).to.equal("AuctionUUPSV2");
      expect(await auctionUUPSV2.getAuctionStats()).to.not.be.undefined;
    });

    it("应该保持原有状态", async function () {
      expect(await auctionUUPSV2.priceOracle()).to.equal(await priceOracle.getAddress());
      expect(await auctionUUPSV2.owner()).to.equal(deployer.address);
    });

    it("应该支持新功能", async function () {
      // 测试 V2 功能
      const stats = await auctionUUPSV2.getAuctionStats();
      expect(stats.totalAuctions).to.equal(0);
      expect(stats.activeAuctions).to.equal(0);

      expect(await auctionUUPSV2.getVersion()).to.equal("AuctionUUPSV2");
    });
  });

  describe("升级安全性测试", function () {
    it("应该防止非所有者升级", async function () {
      const AuctionV2 = await ethers.getContractFactory("AuctionV2");

      // 尝试从非所有者账户升级
      await expect(
        upgrades.upgradeProxy(await auctionTransparent.getAddress(), AuctionV2.connect(seller))
      ).to.be.reverted;
    });

    it("应该正确处理初始化函数", async function () {
      // initializeV2 已经在升级时通过 call 参数调用过了
      // 这里验证升级后的合约状态是正确的
      const version = await auctionTransparentV2.getVersion();
      expect(version).to.equal("AuctionV2");
    });

    it("应该保持函数选择器兼容性", async function () {
      // 测试原有函数仍然可用
      const currentId = await auctionTransparentV2.getCurrentAuctionId();
      expect(currentId).to.equal(0);
    });
  });
});
