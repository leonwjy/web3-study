const { ethers, upgrades } = require("hardhat")

/**  * 2025/02/15 in sepolia testnet
 * esVault contract deployed to: 0x9f09e574f507BA195b2138474DDA451B4032b96D
     esVault ImplementationAddress: 0xc17268561c1E5b21C9224809DF6456Fa86FD87A1
     esVault AdminAddress: 0xd7733630336a51D14ACadbB8289a1124B4B92b5f
   esDex contract deployed to: 0x0Bbe8D3974A007A0Ae61F19BeCB059F00d95532b
      esDex ImplementationAddress: 0xA35eC7171a97C262d01731A15EC4a7d095ECFcc6
      esDex AdminAddress: 0xd7733630336a51D14ACadbB8289a1124B4B92b5f

  deployer:  0xBb03D7e8380a60665710582E42668cFceF85f457
    esVault contract deployed to: 0x9f09e574f507BA195b2138474DDA451B4032b96D
    0xc17268561c1E5b21C9224809DF6456Fa86FD87A1  esVault getImplementationAddress
    0xd7733630336a51D14ACadbB8289a1124B4B92b5f  esVault getAdminAddress
  esDex contract deployed to: 0x0Bbe8D3974A007A0Ae61F19BeCB059F00d95532b
    0xA35eC7171a97C262d01731A15EC4a7d095ECFcc6  esDex getImplementationAddress
    0xd7733630336a51D14ACadbB8289a1124B4B92b5f  esDex getAdminAddress
  esVault setOrderBook tx: 0x278370f06a8d28387e9dab0cf018a678d03c888575c0a5823465a8e3ba22f4a3
 */

async function main() {
  const [deployer] = await ethers.getSigners()
  console.log("deployer: ", deployer.address)

  let esVault = await ethers.getContractFactory("EasySwapVault")
  esVault = await upgrades.deployProxy(esVault, { initializer: 'initialize' });
  await esVault.deployed()
  console.log("esVault contract deployed to:", esVault.address)
  console.log(await upgrades.erc1967.getImplementationAddress(esVault.address), " esVault getImplementationAddress")
  console.log(await upgrades.erc1967.getAdminAddress(esVault.address), " esVault getAdminAddress")

  newProtocolShare = 200;
  const newESVault = esVault.address;
  // newESVault = "0xaD65f3dEac0Fa9Af4eeDC96E95574AEaba6A2834";
  EIP712Name = "EasySwapOrderBook";
  EIP712Version = "1";
  let esDex = await ethers.getContractFactory("EasySwapOrderBook")
  esDex = await upgrades.deployProxy(
    esDex,
    [newProtocolShare, newESVault, EIP712Name, EIP712Version],
    {
      initializer: 'initialize',
      unsafeAllow: ['state-variable-immutable'],
    }
  );
  await esDex.deployed()
  console.log("esDex contract deployed to:", esDex.address)
  console.log(await upgrades.erc1967.getImplementationAddress(esDex.address), " esDex getImplementationAddress")
  console.log(await upgrades.erc1967.getAdminAddress(esDex.address), " esDex getAdminAddress")

  const esDexAddress = esDex.address
  const esVaultAddress = esVault.address
  // esVaultAddress = "0xaD65f3dEac0Fa9Af4eeDC96E95574AEaba6A2834"
  const esVault_ = await (
    await ethers.getContractFactory("EasySwapVault")
  ).attach(esVaultAddress)
  tx = await esVault_.setOrderBook(esDexAddress)
  await tx.wait()
  console.log("esVault setOrderBook tx:", tx.hash)
}

// We recommend this pattern to be able to use async/await everywhere
// and properly handle errors.
main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error)
    process.exit(1)
  })
