const hre = require("hardhat");
const { expect } = require("chai");

describe("MyTokenTest", async () => {

    const { ethers } = hre;

    // 10000 MTK
    const initialSupply = 10000;

    let myToken;

    beforeEach(async () => {
        const MyToken = await ethers.getContractFactory("MyToken");

        myToken = await MyToken.deploy(initialSupply);

        myToken.waitForDeployment();

        const address = await myToken.getAddress();

        expect(address).to.length.greaterThan(0);

        console.log("myToken address:", address);

    });

    it("test contract name、symbol、decimals", async () => {

        const name = await myToken.name();
        const symbol = await myToken.symbol();
        const decimals = await myToken.decimals();

        expect(name).to.equal("MyToken");
        expect(symbol).to.equal("MTK");
        expect(decimals).to.equal(18);
    })
});