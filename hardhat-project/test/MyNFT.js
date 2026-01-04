const hre = require("hardhat");
const { expect } = require("chai");

describe("MyNFT Test", async () => {
    const { ethers } = hre;

    let myNFT;
    let owner;
    let addr1;
    let addr2;

    beforeEach(async () => {
        const MyNFT = await ethers.getContractFactory("MyNFT");

        [owner, addr1, addr2] = await ethers.getSigners();

        myNFT = await MyNFT.deploy();
        await myNFT.waitForDeployment();

        const address = await myNFT.getAddress();
        expect(address).to.length.greaterThan(0);
        console.log("MyNFT address:", address);
    });

    it("Should have correct name and symbol", async () => {
        const name = await myNFT.name();
        const symbol = await myNFT.symbol();

        expect(name).to.equal("MyNFT");
        expect(symbol).to.equal("MNFT");
    });

    it("Should mint NFT correctly", async () => {
        const tokenURI = "https://example.com/token/0";
        await myNFT.safeMint(addr1.address, tokenURI);

        expect(await myNFT.ownerOf(0)).to.equal(addr1.address);
        expect(await myNFT.tokenURI(0)).to.equal(tokenURI);
    });

    it("Should increment token ID correctly", async () => {
        await myNFT.safeMint(addr1.address, "https://example.com/token/0");
        await myNFT.safeMint(addr2.address, "https://example.com/token/1");

        expect(await myNFT.ownerOf(0)).to.equal(addr1.address);
        expect(await myNFT.ownerOf(1)).to.equal(addr2.address);
        expect(await myNFT.tokenURI(0)).to.equal("https://example.com/token/0");
        expect(await myNFT.tokenURI(1)).to.equal("https://example.com/token/1");
    });

    it("Should allow owner to burn NFT", async () => {
        await myNFT.safeMint(addr1.address, "https://example.com/token/0");

        // addr1 should be able to burn their own NFT
        await myNFT.connect(addr1).burn(0);

        await expect(myNFT.ownerOf(0)).to.be.reverted;
    });

    it("Should not allow non-owner to burn NFT", async () => {
        await myNFT.safeMint(addr1.address, "https://example.com/token/0");

        // addr2 should not be able to burn addr1's NFT
        await expect(myNFT.connect(addr2).burn(0)).to.be.revertedWith("Only owner can burn");
    });

    it("Should transfer NFT correctly", async () => {
        await myNFT.safeMint(owner.address, "https://example.com/token/0");

        await myNFT.transferFrom(owner.address, addr1.address, 0);

        expect(await myNFT.ownerOf(0)).to.equal(addr1.address);
    });

    it("Should return correct balance", async () => {
        await myNFT.safeMint(addr1.address, "https://example.com/token/0");
        await myNFT.safeMint(addr1.address, "https://example.com/token/1");
        await myNFT.safeMint(addr2.address, "https://example.com/token/2");

        expect(await myNFT.balanceOf(addr1.address)).to.equal(2);
        expect(await myNFT.balanceOf(addr2.address)).to.equal(1);
    });

    it("Should only allow owner to mint", async () => {
        const tokenURI = "https://example.com/token/0";

        // addr1 should not be able to mint
        await expect(
            myNFT.connect(addr1).safeMint(addr1.address, tokenURI)
        ).to.be.reverted;
    });
});
