const { ethers, upgrades } = require("hardhat");

module.exports = async function ({ getNamedAccounts, deployments }) {
  const { deploy, save } = deployments;
  const { deployer } = await getNamedAccounts();

  console.log("部署用户地址:", deployer);

  // 获取 Chainlink 价格源地址
  const networkName = await ethers.provider.getNetwork().then(n => n.name);
  const ethUsdPriceFeed = networkName === "sepolia"
    ? "0x694AA1769357215DE4FAC081bf1f309aDC325306"
    : "0x5f4eC3Df9cbd43714FE2740f5E3616155c5b8419"; // Mainnet

  // 部署基础合约
  console.log("\n=== 部署 MyNFT ===");
  const myNFT = await deploy("MyNFT", {
    from: deployer,
    args: ["MyNFT", "MNFT"],
    log: true,
  });

  console.log("\n=== 部署 PriceOracle ===");
  const priceOracle = await deploy("PriceOracle", {
    from: deployer,
    args: [ethUsdPriceFeed],
    log: true,
  });

  console.log("\n=== 部署 MockERC20 ===");
  const mockERC20 = await deploy("MockERC20", {
    from: deployer,
    args: ["TestToken", "TEST", 18],
    log: true,
  });

  // 设置代币价格源
  const priceOracleContract = await ethers.getContractAt("PriceOracle", priceOracle.address);
  await priceOracleContract.setTokenPriceFeed(mockERC20.address, ethUsdPriceFeed);

  console.log("\n=== 部署 AuctionTransparent (透明代理) ===");
  const AuctionTransparent = await ethers.getContractFactory("AuctionTransparent");
  const auctionTransparent = await upgrades.deployProxy(AuctionTransparent, [priceOracle.address], {
    initializer: "initialize",
    kind: "transparent",
  });
  await auctionTransparent.waitForDeployment();

  // 保存部署信息
  await save("AuctionTransparent", {
    abi: AuctionTransparent.interface.format(),
    address: await auctionTransparent.getAddress(),
  });

  console.log("\n=== 部署 AuctionUUPS (UUPS代理) ===");
  const AuctionUUPS = await ethers.getContractFactory("AuctionUUPS");
  const auctionUUPS = await upgrades.deployProxy(AuctionUUPS, [priceOracle.address], {
    initializer: "initialize",
    kind: "uups",
  });
  await auctionUUPS.waitForDeployment();

  // 保存部署信息
  await save("AuctionUUPS", {
    abi: AuctionUUPS.interface.format(),
    address: await auctionUUPS.getAddress(),
  });

  console.log("\n=== 部署总结 ===");
  console.log("MyNFT:", myNFT.address);
  console.log("PriceOracle:", priceOracle.address);
  console.log("MockERC20:", mockERC20.address);
  console.log("AuctionTransparent Proxy:", await auctionTransparent.getAddress());
  console.log("AuctionUUPS Proxy:", await auctionUUPS.getAddress());

  // 保存代理地址到文件（用于升级脚本）
  const fs = require("fs");
  const path = require("path");

  const proxyData = {
    network: networkName,
    deployer,
    proxies: {
      AuctionTransparent: await auctionTransparent.getAddress(),
      AuctionUUPS: await auctionUUPS.getAddress(),
    },
    timestamp: new Date().toISOString(),
  };

  const cacheDir = path.join(__dirname, "..", ".cache");
  if (!fs.existsSync(cacheDir)) {
    fs.mkdirSync(cacheDir, { recursive: true });
  }

  fs.writeFileSync(
    path.join(cacheDir, "proxies.json"),
    JSON.stringify(proxyData, null, 2)
  );

  console.log("代理地址已保存到 .cache/proxies.json");
}

module.exports.tags = ["deploy"];

