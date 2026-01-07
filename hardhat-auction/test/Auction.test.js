const { expect } = require("chai");
const { ethers, time } = require("hardhat");

describe("Auction", function () {
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
    // ETH/USD: $3000
    mockEthPriceFeed = await MockChainlinkAggregator.deploy(3000 * 10 ** 8);
    await mockEthPriceFeed.waitForDeployment();
    
    // Token/USD: $1
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

    // 铸造 NFT 给 seller
    await myNFT.mint(seller.address, "https://example.com/token/1");
    await myNFT.connect(seller).approve(await auction.getAddress(), 1);
  });

  describe("Deployment", function () {
    it("Should set the right price oracle", async function () {
      expect(await auction.priceOracle()).to.equal(await priceOracle.getAddress());
    });

    it("Should set the right owner", async function () {
      expect(await auction.owner()).to.equal(owner.address);
    });
  });

  describe("Create Auction", function () {
    it("Should create auction with ETH", async function () {
      const startingPriceUSD = ethers.parseEther("100"); // $100
      const duration = 3600; // 1 hour

      await expect(
        auction.connect(seller).createAuction(
          await myNFT.getAddress(),
          1,
          startingPriceUSD,
          duration,
          ethers.ZeroAddress
        )
      ).to.emit(auction, "AuctionCreated");

      const auctionInfo = await auction.getAuction(1);
      expect(auctionInfo.nftContract).to.equal(await myNFT.getAddress());
      expect(auctionInfo.tokenId).to.equal(1);
      expect(auctionInfo.seller).to.equal(seller.address);
      expect(auctionInfo.startingPrice).to.equal(startingPriceUSD);
      expect(auctionInfo.paymentToken).to.equal(ethers.ZeroAddress);
    });

    it("Should create auction with ERC20", async function () {
      const startingPriceUSD = ethers.parseEther("100");
      const duration = 3600;

      await expect(
        auction.connect(seller).createAuction(
          await myNFT.getAddress(),
          1,
          startingPriceUSD,
          duration,
          await mockERC20.getAddress()
        )
      ).to.emit(auction, "AuctionCreated");

      const auctionInfo = await auction.getAuction(1);
      expect(auctionInfo.paymentToken).to.equal(await mockERC20.getAddress());
    });

    it("Should transfer NFT to auction contract", async function () {
      await auction.connect(seller).createAuction(
        await myNFT.getAddress(),
        1,
        ethers.parseEther("100"),
        3600,
        ethers.ZeroAddress
      );

      expect(await myNFT.ownerOf(1)).to.equal(await auction.getAddress());
    });

    it("Should revert if NFT not approved", async function () {
      await myNFT.mint(seller.address, "uri2");
      
      await expect(
        auction.connect(seller).createAuction(
          await myNFT.getAddress(),
          2,
          ethers.parseEther("100"),
          3600,
          ethers.ZeroAddress
        )
      ).to.be.revertedWith("NFT not approved");
    });
  });

  describe("Bid with ETH", function () {
    beforeEach(async function () {
      await auction.connect(seller).createAuction(
        await myNFT.getAddress(),
        1,
        ethers.parseEther("100"), // $100 starting price
        3600,
        ethers.ZeroAddress
      );
    });

    it("Should place ETH bid", async function () {
      // $100 = 0.0333 ETH (at $3000/ETH)
      const bidAmount = ethers.parseEther("0.04"); // ~$120
      
      await expect(
        auction.connect(bidder1).bidWithETH(1, { value: bidAmount })
      ).to.emit(auction, "BidPlaced");

      const auctionInfo = await auction.getAuction(1);
      expect(auctionInfo.highestBidder).to.equal(bidder1.address);
    });

    it("Should reject bid below starting price", async function () {
      // $50 < $100 starting price
      const bidAmount = ethers.parseEther("0.01"); // ~$30
      
      await expect(
        auction.connect(bidder1).bidWithETH(1, { value: bidAmount })
      ).to.be.revertedWith("Bid below starting price");
    });

    it("Should reject bid lower than current highest", async function () {
      const bid1 = ethers.parseEther("0.04"); // ~$120
      const bid2 = ethers.parseEther("0.03"); // ~$90
      
      await auction.connect(bidder1).bidWithETH(1, { value: bid1 });
      
      await expect(
        auction.connect(bidder2).bidWithETH(1, { value: bid2 })
      ).to.be.revertedWith("Bid must be higher than current highest bid");
    });

    it("Should update highest bidder", async function () {
      const bid1 = ethers.parseEther("0.04");
      const bid2 = ethers.parseEther("0.05");
      
      await auction.connect(bidder1).bidWithETH(1, { value: bid1 });
      await auction.connect(bidder2).bidWithETH(1, { value: bid2 });

      const auctionInfo = await auction.getAuction(1);
      expect(auctionInfo.highestBidder).to.equal(bidder2.address);
    });
  });

  describe("Bid with ERC20", function () {
    beforeEach(async function () {
      await auction.connect(seller).createAuction(
        await myNFT.getAddress(),
        1,
        ethers.parseEther("100"),
        3600,
        await mockERC20.getAddress()
      );
    });

    it("Should place ERC20 bid", async function () {
      // $100 = 100 tokens (at $1/token)
      const bidAmount = ethers.parseEther("120");
      
      await mockERC20.connect(bidder1).approve(await auction.getAddress(), bidAmount);
      
      await expect(
        auction.connect(bidder1).bidWithERC20(1, bidAmount)
      ).to.emit(auction, "BidPlaced");

      const auctionInfo = await auction.getAuction(1);
      expect(auctionInfo.highestBidder).to.equal(bidder1.address);
    });

    it("Should transfer tokens to auction contract", async function () {
      const bidAmount = ethers.parseEther("120");
      await mockERC20.connect(bidder1).approve(await auction.getAddress(), bidAmount);
      
      await auction.connect(bidder1).bidWithERC20(1, bidAmount);

      expect(await mockERC20.balanceOf(await auction.getAddress())).to.equal(bidAmount);
    });
  });

  describe("End Auction", function () {
    beforeEach(async function () {
      await auction.connect(seller).createAuction(
        await myNFT.getAddress(),
        1,
        ethers.parseEther("100"),
        3600,
        ethers.ZeroAddress
      );
    });

    it("Should end auction and transfer NFT to winner", async function () {
      const bidAmount = ethers.parseEther("0.04");
      await auction.connect(bidder1).bidWithETH(1, { value: bidAmount });
      
      await ethers.provider.send("evm_increaseTime", [3601]); // 超过结束时间
      await ethers.provider.send("evm_mine");
      
      await expect(
        auction.endAuction(1)
      ).to.emit(auction, "AuctionEnded");

      expect(await myNFT.ownerOf(1)).to.equal(bidder1.address);
    });

    it("Should transfer funds to seller (minus fee)", async function () {
      const bidAmount = ethers.parseEther("0.04");
      await auction.connect(bidder1).bidWithETH(1, { value: bidAmount });
      
      const sellerBalanceBefore = await ethers.provider.getBalance(seller.address);
      await ethers.provider.send("evm_increaseTime", [3601]);
      await ethers.provider.send("evm_mine");
      await auction.endAuction(1);
      const sellerBalanceAfter = await ethers.provider.getBalance(seller.address);
      
      // 手续费 2.5%，卖家获得 97.5%
      const expectedSellerAmount = bidAmount * BigInt(9750) / BigInt(10000);
      expect(sellerBalanceAfter - sellerBalanceBefore).to.be.closeTo(
        expectedSellerAmount,
        ethers.parseEther("0.001") // 允许 gas 费用误差
      );
    });

    it("Should refund NFT if no bids", async function () {
      await ethers.provider.send("evm_increaseTime", [3601]);
      await ethers.provider.send("evm_mine");
      await auction.endAuction(1);

      expect(await myNFT.ownerOf(1)).to.equal(seller.address);
    });

    it("Should revert if auction not ended", async function () {
      await expect(
        auction.endAuction(1)
      ).to.be.revertedWith("Auction not ended and not seller");
    });
  });

  describe("Withdraw Bid", function () {
    beforeEach(async function () {
      await auction.connect(seller).createAuction(
        await myNFT.getAddress(),
        1,
        ethers.parseEther("100"),
        3600,
        ethers.ZeroAddress
      );
    });

    it("Should allow bidder to withdraw outbid bid", async function () {
      const bid1 = ethers.parseEther("0.04");
      const bid2 = ethers.parseEther("0.05");
      
      await auction.connect(bidder1).bidWithETH(1, { value: bid1 });
      await auction.connect(bidder2).bidWithETH(1, { value: bid2 });

      const bidder1BalanceBefore = await ethers.provider.getBalance(bidder1.address);
      await auction.connect(bidder1).withdrawBid(1);
      const bidder1BalanceAfter = await ethers.provider.getBalance(bidder1.address);

      expect(bidder1BalanceAfter - bidder1BalanceBefore).to.be.closeTo(
        bid1,
        ethers.parseEther("0.001")
      );
    });

    it("Should revert if trying to withdraw current highest bid", async function () {
      const bidAmount = ethers.parseEther("0.04");
      await auction.connect(bidder1).bidWithETH(1, { value: bidAmount });

      await expect(
        auction.connect(bidder1).withdrawBid(1)
      ).to.be.revertedWith("Cannot withdraw current highest bid");
    });
  });
});

