const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  const [deployer] = await hre.ethers.getSigners();
  console.log("Deploying contracts with account:", deployer.address);
  console.log("Account balance:", (await hre.ethers.provider.getBalance(deployer.address)).toString());

  // 获取 Chainlink 价格源地址
  const networkName = hre.network.name;
  const ethUsdPriceFeed = hre.config.chainlink?.ethUsdPriceFeed?.[networkName] || 
                         "0x694AA1769357215DE4FAC081bf1f309aDC325306"; // Sepolia default

  console.log("\n=== Deploying MyNFT ===");
  const MyNFT = await hre.ethers.getContractFactory("MyNFT");
  const myNFT = await MyNFT.deploy("MyNFT", "MNFT");
  await myNFT.waitForDeployment();
  const myNFTAddress = await myNFT.getAddress();
  console.log("MyNFT deployed to:", myNFTAddress);

  console.log("\n=== Deploying PriceOracle ===");
  const PriceOracle = await hre.ethers.getContractFactory("PriceOracle");
  const priceOracle = await PriceOracle.deploy(ethUsdPriceFeed);
  await priceOracle.waitForDeployment();
  const priceOracleAddress = await priceOracle.getAddress();
  console.log("PriceOracle deployed to:", priceOracleAddress);
  console.log("ETH/USD Price Feed:", ethUsdPriceFeed);

  console.log("\n=== Deploying AuctionUUPS Implementation ===");
  const AuctionUUPS = await hre.ethers.getContractFactory("AuctionUUPS");
  const auctionUUPSImpl = await AuctionUUPS.deploy();
  await auctionUUPSImpl.waitForDeployment();
  const auctionUUPSImplAddress = await auctionUUPSImpl.getAddress();
  console.log("AuctionUUPS Implementation deployed to:", auctionUUPSImplAddress);

  console.log("\n=== Deploying AuctionUUPS Proxy ===");
  const ERC1967ProxyArtifact = require("@openzeppelin/contracts/build/contracts/ERC1967Proxy.json");
  const ERC1967Proxy = await hre.ethers.getContractFactoryFromArtifact(ERC1967ProxyArtifact);
  const auctionUUPSInitData = AuctionUUPS.interface.encodeFunctionData("initialize", [
    priceOracleAddress
  ]);
  const auctionUUPSProxy = await ERC1967Proxy.deploy(auctionUUPSImplAddress, auctionUUPSInitData);
  await auctionUUPSProxy.waitForDeployment();
  const auctionUUPSProxyAddress = await auctionUUPSProxy.getAddress();
  console.log("AuctionUUPS Proxy deployed to:", auctionUUPSProxyAddress);

  // 获取代理合约实例并验证初始化
  const auctionUUPS = AuctionUUPS.attach(auctionUUPSProxyAddress);
  console.log("AuctionUUPS Owner:", await auctionUUPS.owner());
  console.log("AuctionUUPS Version:", await auctionUUPS.version());

  console.log("\n=== Deploying AuctionTransparent Implementation ===");
  const AuctionTransparent = await hre.ethers.getContractFactory("AuctionTransparent");
  const auctionTransparentImpl = await AuctionTransparent.deploy();
  await auctionTransparentImpl.waitForDeployment();
  const auctionTransparentImplAddress = await auctionTransparentImpl.getAddress();
  console.log("AuctionTransparent Implementation deployed to:", auctionTransparentImplAddress);

  console.log("\n=== Deploying ProxyAdmin ===");
  const ProxyAdminArtifact = require("@openzeppelin/contracts/build/contracts/ProxyAdmin.json");
  const ProxyAdmin = await hre.ethers.getContractFactoryFromArtifact(ProxyAdminArtifact);
  const proxyAdmin = await ProxyAdmin.deploy(deployer.address);
  await proxyAdmin.waitForDeployment();
  const proxyAdminAddress = await proxyAdmin.getAddress();
  console.log("ProxyAdmin deployed to:", proxyAdminAddress);

  console.log("\n=== Deploying AuctionTransparent Proxy ===");
  const TransparentUpgradeableProxyArtifact = require("@openzeppelin/contracts/build/contracts/TransparentUpgradeableProxy.json");
  const TransparentUpgradeableProxy = await hre.ethers.getContractFactoryFromArtifact(TransparentUpgradeableProxyArtifact);
  const auctionTransparentInitData = AuctionTransparent.interface.encodeFunctionData("initialize", [
    priceOracleAddress
  ]);
  const auctionTransparentProxy = await TransparentUpgradeableProxy.deploy(
    auctionTransparentImplAddress,
    proxyAdminAddress,
    auctionTransparentInitData
  );
  await auctionTransparentProxy.waitForDeployment();
  const auctionTransparentProxyAddress = await auctionTransparentProxy.getAddress();
  console.log("AuctionTransparent Proxy deployed to:", auctionTransparentProxyAddress);

  // 获取代理合约实例并验证初始化
  const auctionTransparent = AuctionTransparent.attach(auctionTransparentProxyAddress);
  console.log("AuctionTransparent Owner:", await auctionTransparent.owner());
  console.log("AuctionTransparent Version:", await auctionTransparent.version());

  // 保存部署地址
  const deploymentInfo = {
    network: networkName,
    deployer: deployer.address,
    contracts: {
      MyNFT: myNFTAddress,
      PriceOracle: priceOracleAddress,
      AuctionUUPS: {
        implementation: auctionUUPSImplAddress,
        proxy: auctionUUPSProxyAddress,
      },
      AuctionTransparent: {
        implementation: auctionTransparentImplAddress,
        proxy: auctionTransparentProxyAddress,
        proxyAdmin: proxyAdminAddress,
      },
    },
    chainlink: {
      ethUsdPriceFeed: ethUsdPriceFeed,
    },
    timestamp: new Date().toISOString(),
  };

  const deploymentPath = path.join(__dirname, "..", "deployments", `${networkName}.json`);
  fs.mkdirSync(path.dirname(deploymentPath), { recursive: true });
  fs.writeFileSync(deploymentPath, JSON.stringify(deploymentInfo, null, 2));
  console.log("\n=== Deployment info saved to:", deploymentPath);

  console.log("\n=== Deployment Summary ===");
  console.log(JSON.stringify(deploymentInfo, null, 2));
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });

