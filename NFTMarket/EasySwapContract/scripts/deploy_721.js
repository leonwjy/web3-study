const { ethers, upgrades } = require("hardhat")

async function main() {
  const [deployer] = await ethers.getSigners()
  console.log("deployer: ", deployer.address)  // 0xBb03D7e8380a60665710582E42668cFceF85f457

  // let TestERC721 = await ethers.getContractFactory("Troll")
  // const testERC721 = await TestERC721.deploy()
  // await testERC721.deployed()
  // console.log("testERC721 contract deployed to:", testERC721.address)  // 0x5AD7186d8cf84091b4823DB49198b03207E5247E

  //mint
  // let testERC721Address = "0x5AD7186d8cf84091b4823DB49198b03207E5247E";
  // let testERC721 = await (await ethers.getContractFactory("Troll")).attach(testERC721Address)
  // tx = await testERC721.mint(deployer.address, 50);
  // await tx.wait()
  // console.log("mint tx:", tx.hash) // 0x105ed13f5febe9256f063afe14304a3c6e8ea24aa9b48e7bb288d67148c5f616
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error)
    process.exit(1)
  })
