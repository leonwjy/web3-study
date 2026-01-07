const hre = require("hardhat");
const fs = require("fs");
const path = require("path");

async function main() {
  const [deployer] = await hre.ethers.getSigners();
  const networkName = hre.network.name;

  // 读取部署信息
  const deploymentPath = path.join(__dirname, "..", "deployments", `${networkName}.json`);
  if (!fs.existsSync(deploymentPath)) {
    throw new Error(`Deployment file not found: ${deploymentPath}`);
  }

  const deploymentInfo = JSON.parse(fs.readFileSync(deploymentPath, "utf8"));
  console.log("Deployment info loaded from:", deploymentPath);

  const upgradeType = process.env.UPGRADE_TYPE || "uups"; // "uups" or "transparent"
  console.log(`\n=== Upgrading ${upgradeType.toUpperCase()} Proxy ===`);

  if (upgradeType === "uups") {
    // UUPS 升级
    const proxyAddress = deploymentInfo.contracts.AuctionUUPS.proxy;
    console.log("Proxy address:", proxyAddress);

    // 部署新版本实现合约
    console.log("\nDeploying new implementation...");
    const AuctionUUPS = await hre.ethers.getContractFactory("AuctionUUPS");
    const newImpl = await AuctionUUPS.deploy();
    await newImpl.waitForDeployment();
    const newImplAddress = await newImpl.getAddress();
    console.log("New implementation deployed to:", newImplAddress);

    // 获取代理合约实例
    const proxy = AuctionUUPS.attach(proxyAddress);
    const oldVersion = await proxy.version();
    console.log("Current version:", oldVersion.toString());

    // 执行升级
    console.log("\nUpgrading proxy...");
    const tx = await proxy.upgradeToAndCall(newImplAddress, "0x");
    await tx.wait();
    console.log("Upgrade transaction:", tx.hash);

    // 验证升级
    const newVersion = await proxy.version();
    console.log("New version:", newVersion.toString());

    // 更新部署信息
    deploymentInfo.contracts.AuctionUUPS.implementation = newImplAddress;
    deploymentInfo.contracts.AuctionUUPS.lastUpgrade = {
      timestamp: new Date().toISOString(),
      newImplementation: newImplAddress,
      transactionHash: tx.hash,
    };

  } else if (upgradeType === "transparent") {
    // 透明代理升级
    const proxyAddress = deploymentInfo.contracts.AuctionTransparent.proxy;
    const proxyAdminAddress = deploymentInfo.contracts.AuctionTransparent.proxyAdmin;
    console.log("Proxy address:", proxyAddress);
    console.log("ProxyAdmin address:", proxyAdminAddress);

    // 部署新版本实现合约
    console.log("\nDeploying new implementation...");
    const AuctionTransparent = await hre.ethers.getContractFactory("AuctionTransparent");
    const newImpl = await AuctionTransparent.deploy();
    await newImpl.waitForDeployment();
    const newImplAddress = await newImpl.getAddress();
    console.log("New implementation deployed to:", newImplAddress);

    // 获取 ProxyAdmin 实例
    const ProxyAdmin = await hre.ethers.getContractFactory("ProxyAdmin");
    const proxyAdmin = ProxyAdmin.attach(proxyAdminAddress);

    // 获取代理合约实例以检查当前版本
    const proxy = AuctionTransparent.attach(proxyAddress);
    const oldVersion = await proxy.version();
    console.log("Current version:", oldVersion.toString());

    // 执行升级
    console.log("\nUpgrading proxy via ProxyAdmin...");
    const tx = await proxyAdmin.upgrade(proxyAddress, newImplAddress);
    await tx.wait();
    console.log("Upgrade transaction:", tx.hash);

    // 验证升级
    const newVersion = await proxy.version();
    console.log("New version:", newVersion.toString());

    // 更新部署信息
    deploymentInfo.contracts.AuctionTransparent.implementation = newImplAddress;
    deploymentInfo.contracts.AuctionTransparent.lastUpgrade = {
      timestamp: new Date().toISOString(),
      newImplementation: newImplAddress,
      transactionHash: tx.hash,
    };
  } else {
    throw new Error("Invalid UPGRADE_TYPE. Use 'uups' or 'transparent'");
  }

  // 保存更新的部署信息
  fs.writeFileSync(deploymentPath, JSON.stringify(deploymentInfo, null, 2));
  console.log("\n=== Deployment info updated ===");
}

main()
  .then(() => process.exit(0))
  .catch((error) => {
    console.error(error);
    process.exit(1);
  });

