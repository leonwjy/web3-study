const { ethers, upgrades } = require("hardhat");
const fs = require("fs");
const path = require("path");

module.exports = async function ({ getNamedAccounts, deployments }) {
  const { save } = deployments;
  const { deployer } = await getNamedAccounts();

  console.log("升级用户地址:", deployer);

  // 读取代理地址
  const cachePath = path.join(__dirname, "..", ".cache", "proxies.json");
  if (!fs.existsSync(cachePath)) {
    throw new Error("未找到代理地址文件，请先运行部署脚本");
  }

  const proxyData = JSON.parse(fs.readFileSync(cachePath, "utf-8"));
  const auctionTransparentAddress = proxyData.proxies.AuctionTransparent;
  const auctionUUPSAddress = proxyData.proxies.AuctionUUPS;

  console.log("\n=== 升级 AuctionTransparent 到 V2 ===");
  const AuctionV2 = await ethers.getContractFactory("AuctionV2");

  // 升级透明代理
  const auctionTransparentV2 = await upgrades.upgradeProxy(auctionTransparentAddress, AuctionV2, {
    call: { fn: "initializeV2" },
    kind: "transparent",
  });

  await auctionTransparentV2.waitForDeployment();
  const transparentAddress = await auctionTransparentV2.getAddress();

  console.log("AuctionTransparent 升级完成:", transparentAddress);

  console.log("\n=== 升级 AuctionUUPS 到 V2 ===");

  // 升级 UUPS 代理
  const auctionUUPSV2 = await upgrades.upgradeProxy(auctionUUPSAddress, AuctionV2, {
    call: { fn: "initializeV2" },
    kind: "uups",
  });

  await auctionUUPSV2.waitForDeployment();
  const uupsAddress = await auctionUUPSV2.getAddress();

  console.log("AuctionUUPS 升级完成:", uupsAddress);

  // 保存升级后的合约信息
  await save("AuctionTransparentV2", {
    abi: AuctionV2.interface.format(),
    address: transparentAddress,
  });

  await save("AuctionUUPSV2", {
    abi: AuctionV2.interface.format(),
    address: uupsAddress,
  });

  console.log("\n=== 升级总结 ===");
  console.log("AuctionTransparent V2:", transparentAddress);
  console.log("AuctionUUPS V2:", uupsAddress);

  // 验证升级
  console.log("\n=== 验证升级 ===");

  const transparentVersion = await auctionTransparentV2.getVersion();
  const uupsVersion = await auctionUUPSV2.getVersion();

  console.log("AuctionTransparent 版本:", transparentVersion);
  console.log("AuctionUUPS 版本:", uupsVersion);

  // 测试新功能
  const stats = await auctionTransparentV2.getAuctionStats();
  console.log("新功能测试 - 统计信息:", stats);

  // 测试暂停功能
  await auctionTransparentV2.pause();
  const isPaused = await auctionTransparentV2.paused();
  console.log("暂停功能测试:", isPaused ? "成功" : "失败");

  // 更新缓存文件
  proxyData.upgraded = {
    AuctionTransparentV2: transparentAddress,
    AuctionUUPSV2: uupsAddress,
  };
  proxyData.upgradedAt = new Date().toISOString();

  fs.writeFileSync(cachePath, JSON.stringify(proxyData, null, 2));
  console.log("升级信息已保存到 .cache/proxies.json");
}

module.exports.tags = ["upgrade"];