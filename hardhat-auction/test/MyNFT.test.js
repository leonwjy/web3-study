const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("MyNFT", function () {
  let myNFT;
  let owner;
  let user1;
  let user2;

  beforeEach(async function () {
    [owner, user1, user2] = await ethers.getSigners();

    const MyNFT = await ethers.getContractFactory("MyNFT");
    myNFT = await MyNFT.deploy("MyNFT", "MNFT");
    await myNFT.waitForDeployment();
  });

  describe("Deployment", function () {
    it("Should set the right name and symbol", async function () {
      expect(await myNFT.name()).to.equal("MyNFT");
      expect(await myNFT.symbol()).to.equal("MNFT");
    });

    it("Should set the right owner", async function () {
      expect(await myNFT.owner()).to.equal(owner.address);
    });
  });

  describe("Minting", function () {
    it("Should mint NFT to specified address", async function () {
      const tokenURI = "https://example.com/token/1";
      await myNFT.mint(user1.address, tokenURI);
      
      expect(await myNFT.ownerOf(1)).to.equal(user1.address);
      expect(await myNFT.tokenURI(1)).to.equal(tokenURI);
      expect(await myNFT.totalSupply()).to.equal(1);
    });

    it("Should increment token ID for each mint", async function () {
      await myNFT.mint(user1.address, "uri1");
      await myNFT.mint(user2.address, "uri2");
      
      expect(await myNFT.ownerOf(1)).to.equal(user1.address);
      expect(await myNFT.ownerOf(2)).to.equal(user2.address);
      expect(await myNFT.totalSupply()).to.equal(2);
    });

    it("Should revert if non-owner tries to mint", async function () {
      await expect(
        myNFT.connect(user1).mint(user2.address, "uri")
      ).to.be.revertedWithCustomError(myNFT, "OwnableUnauthorizedAccount");
    });
  });

  describe("Token URI", function () {
    it("Should set token URI", async function () {
      const tokenId = 1;
      const tokenURI = "https://example.com/token/1";
      
      await myNFT.mint(user1.address, tokenURI);
      expect(await myNFT.tokenURI(tokenId)).to.equal(tokenURI);
    });

    it("Should update token URI", async function () {
      const tokenId = 1;
      const initialURI = "https://example.com/token/1";
      const newURI = "https://example.com/token/1-updated";
      
      await myNFT.mint(user1.address, initialURI);
      await myNFT.setTokenURI(tokenId, newURI);
      
      expect(await myNFT.tokenURI(tokenId)).to.equal(newURI);
    });

    it("Should revert if setting URI for non-existent token", async function () {
      await expect(
        myNFT.setTokenURI(999, "uri")
      ).to.be.revertedWith("Token does not exist");
    });

    it("Should revert if non-owner tries to set URI", async function () {
      await myNFT.mint(user1.address, "uri");
      
      await expect(
        myNFT.connect(user1).setTokenURI(1, "new-uri")
      ).to.be.revertedWithCustomError(myNFT, "OwnableUnauthorizedAccount");
    });
  });

  describe("Transfers", function () {
    it("Should transfer NFT", async function () {
      await myNFT.mint(user1.address, "uri");
      await myNFT.connect(user1).transferFrom(user1.address, user2.address, 1);
      
      expect(await myNFT.ownerOf(1)).to.equal(user2.address);
    });
  });
});

