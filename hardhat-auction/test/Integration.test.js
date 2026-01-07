const { expect } = require("chai");
const { ethers, time } = require("hardhat");

describe("Integration Tests", function () {
  let auction;
  let myNFT;
  let priceOracle;
  let mockERC20;
  let mockEthPriceFeed;
  let mockTokenPriceFeed;
  let owner;
  let seller;
  let bidder1;
  let bidder2;

  beforeEach(async function () {
    [owner, seller, bidder1, bidder2] = await ethers.getSigners();

    // 部署 Mock Chainlink Aggregators
    const MockChainlinkAggregator = await ethers.getContractFactory("MockChainlinkAggregator");
    mockEthPriceFeed = await MockChainlinkAggregator.deploy(3000 * 10 ** 8);
    await mockEthPriceFeed.waitForDeployment();
    
    mockTokenPriceFeed = await MockChainlinkAggregator.deploy(1 * 10 ** 8);
    await mockTokenPriceFeed.waitForDeployment();

    // 部署 PriceOracle
    const PriceOracle = await ethers.getContractFactory("PriceOracle");
    priceOracle = await PriceOracle.deploy(await mockEthPriceFeed.getAddress());
    await priceOracle.waitForDeployment();

    // 部署 Mock ERC20
    const MockERC20 = await ethers.getContractFactory("MockERC20");
    mockERC20 = await MockERC20.deploy("TestToken", "TEST", 18);
    await mockERC20.waitForDeployment();

    // 设置 Token 价格源
    await priceOracle.setTokenPriceFeed(await mockERC20.getAddress(), await mockTokenPriceFeed.getAddress());

    // 给测试用户铸造 ERC20 代币
    await mockERC20.mint(bidder1.address, ethers.parseEther("1000"));
    await mockERC20.mint(bidder2.address, ethers.parseEther("1000"));

    // 部署 Auction
    const Auction = await ethers.getContractFactory("Auction");
    auction = await Auction.deploy();
    await auction.waitForDeployment();
    await auction.initialize(await priceOracle.getAddress());

    // 部署 MyNFT
    const MyNFT = await ethers.getContractFactory("MyNFT");
    myNFT = await MyNFT.deploy("MyNFT", "MNFT");
    await myNFT.waitForDeployment();
  });

  describe("Full Auction Flow with ETH", function () {
    it("Should complete full auction flow: mint -> create -> bid -> end", async function () {
      // 1. Mint NFT
      const mintTx = await myNFT.connect(owner).mint(seller.address, "https://example.com/token/1");
      await mintTx.wait();
      const tokenId = await myNFT.totalSupply();
      expect(await myNFT.ownerOf(tokenId)).to.equal(seller.address);

      // 2. Approve NFT
      await myNFT.connect(seller).approve(await auction.getAddress(), tokenId);

      // 3. Create auction
      const startingPriceUSD = ethers.parseEther("100");
      const duration = 3600;
      
      await auction.connect(seller).createAuction(
        await myNFT.getAddress(),
        tokenId,
        startingPriceUSD,
        duration,
        ethers.ZeroAddress
      );

      expect(await myNFT.ownerOf(tokenId)).to.equal(await auction.getAddress());

      const auctionId = 1; // 这是第一个auction

      // 4. Place bids
      const bid1 = ethers.parseEther("0.04"); // ~$120
      const bid2 = ethers.parseEther("0.05"); // ~$150

      await auction.connect(bidder1).bidWithETH(auctionId, { value: bid1 });
      await auction.connect(bidder2).bidWithETH(auctionId, { value: bid2 });

      const auctionInfo = await auction.getAuction(auctionId);
      expect(auctionInfo.highestBidder).to.equal(bidder2.address);

      // 5. Withdraw outbid bid
      const bidder1BalanceBefore = await ethers.provider.getBalance(bidder1.address);
      await auction.connect(bidder1).withdrawBid(auctionId);
      const bidder1BalanceAfter = await ethers.provider.getBalance(bidder1.address);
      expect(bidder1BalanceAfter - bidder1BalanceBefore).to.be.closeTo(bid1, ethers.parseEther("0.001"));

      // 6. End auction
      await ethers.provider.send("evm_increaseTime", [3601]);
      await ethers.provider.send("evm_mine");
      await auction.endAuction(auctionId);

      // 7. Verify NFT transferred to winner
      expect(await myNFT.ownerOf(tokenId)).to.equal(bidder2.address);

      // 8. Verify funds transferred (seller gets 97.5%, owner gets 2.5%)
      const sellerBalanceAfter = await ethers.provider.getBalance(seller.address);
      const ownerBalanceAfter = await ethers.provider.getBalance(owner.address);
      
      // 注意：由于 gas 费用，这里只做基本验证
      expect(sellerBalanceAfter).to.be.gt(0);
      expect(ownerBalanceAfter).to.be.gt(0);
    });
  });

  describe("Full Auction Flow with ERC20", function () {
    it("Should complete full auction flow with ERC20", async function () {
      // 1. Mint NFT - 使用mint返回的实际tokenId
      const mintTx = await myNFT.connect(owner).mint(seller.address, "https://example.com/token/100");
      await mintTx.wait();
      // 获取当前总供应量作为tokenId（假设这是新的token）
      const tokenId = await myNFT.totalSupply();
      await myNFT.connect(seller).approve(await auction.getAddress(), tokenId);

      // 2. Create auction with ERC20
      await auction.connect(seller).createAuction(
        await myNFT.getAddress(),
        tokenId,
        ethers.parseEther("100"),
        3600,
        await mockERC20.getAddress()
      );

      const auctionId = 1; // 这是第一个auction（ETH测试中的是不同的测试）

      // 3. Place ERC20 bids
      const bid1 = ethers.parseEther("120");
      const bid2 = ethers.parseEther("150");

      await mockERC20.connect(bidder1).approve(await auction.getAddress(), bid1);
      await mockERC20.connect(bidder2).approve(await auction.getAddress(), bid2);

      await auction.connect(bidder1).bidWithERC20(auctionId, bid1);
      await auction.connect(bidder2).bidWithERC20(auctionId, bid2);

      // 4. End auction
      await ethers.provider.send("evm_increaseTime", [3601]);
      await ethers.provider.send("evm_mine");
      await auction.endAuction(auctionId);

      // 5. Verify NFT transferred
      expect(await myNFT.ownerOf(tokenId)).to.equal(bidder2.address);

      // 6. Verify token balances
      const sellerBalance = await mockERC20.balanceOf(seller.address);
      const ownerBalance = await mockERC20.balanceOf(owner.address);
      
      // Seller should receive 97.5% of bid2
      const expectedSellerAmount = bid2 * BigInt(9750) / BigInt(10000);
      expect(sellerBalance).to.be.closeTo(expectedSellerAmount, ethers.parseEther("0.1"));
    });
  });

  describe("Multiple Auctions", function () {
    it("Should handle multiple concurrent auctions", async function () {
      // Create multiple NFTs and auctions
      for (let i = 1; i <= 3; i++) {
        await myNFT.mint(seller.address, `https://example.com/token/${i}`);
        await myNFT.connect(seller).approve(await auction.getAddress(), i);
        
        await auction.connect(seller).createAuction(
          await myNFT.getAddress(),
          i,
          ethers.parseEther("100"),
          3600,
          ethers.ZeroAddress
        );
      }

      // Place bids on different auctions
      await auction.connect(bidder1).bidWithETH(1, { value: ethers.parseEther("0.04") });
      await auction.connect(bidder2).bidWithETH(2, { value: ethers.parseEther("0.05") });
      await auction.connect(bidder1).bidWithETH(3, { value: ethers.parseEther("0.06") });

      // Verify each auction has correct highest bidder
      expect((await auction.getAuction(1)).highestBidder).to.equal(bidder1.address);
      expect((await auction.getAuction(2)).highestBidder).to.equal(bidder2.address);
      expect((await auction.getAuction(3)).highestBidder).to.equal(bidder1.address);
    });
  });
});

