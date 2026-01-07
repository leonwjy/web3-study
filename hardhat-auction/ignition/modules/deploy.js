const { buildModule } = require("@nomicfoundation/hardhat-ignition/modules");

module.exports = buildModule("AuctionModule", (m) => {
  const deployer = m.getAccount(0);

  // 部署 MyNFT
  const myNFT = m.contract("MyNFT", ["MyNFT", "MNFT"], {
    from: deployer,
  });

  // 获取 Chainlink ETH/USD 价格源地址（从配置中读取）
  // 对于 Sepolia: 0x694AA1769357215DE4FAC081bf1f309aDC325306
  // 对于本地测试，可以使用 Mock
  const ethUsdPriceFeed = m.getParameter(
    "ethUsdPriceFeed",
    "0x694AA1769357215DE4FAC081bf1f309aDC325306" // Sepolia default
  );

  // 部署 PriceOracle
  const priceOracle = m.contract("PriceOracle", [ethUsdPriceFeed], {
    from: deployer,
  });

  // 部署 AuctionUUPS 实现合约
  const auctionUUPSImpl = m.contract("AuctionUUPS", [], {
    from: deployer,
  });

  // 部署 AuctionTransparent 实现合约
  const auctionTransparentImpl = m.contract("AuctionTransparent", [], {
    from: deployer,
  });

  // 部署 UUPS 代理
  const auctionUUPSInitData = m.encodeFunctionCall(
    auctionUUPSImpl,
    "initialize",
    [priceOracle]
  );

  const auctionUUPSProxy = m.contract(
    "ERC1967Proxy",
    [auctionUUPSImpl, auctionUUPSInitData],
    {
      from: deployer,
      id: "AuctionUUPSProxy",
    }
  );

  // 部署 ProxyAdmin（用于透明代理）
  const proxyAdmin = m.contract("ProxyAdmin", [], {
    from: deployer,
  });

  // 部署透明代理
  const auctionTransparentInitData = m.encodeFunctionCall(
    auctionTransparentImpl,
    "initialize",
    [priceOracle]
  );

  const auctionTransparentProxy = m.contract(
    "TransparentUpgradeableProxy",
    [auctionTransparentImpl, proxyAdmin, auctionTransparentInitData],
    {
      from: deployer,
      id: "AuctionTransparentProxy",
    }
  );

  return {
    myNFT,
    priceOracle,
    auctionUUPSImpl,
    auctionUUPSProxy,
    auctionTransparentImpl,
    auctionTransparentProxy,
    proxyAdmin,
  };
});

