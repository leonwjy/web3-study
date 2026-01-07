# NFT 拍卖市场

一个基于 Hardhat 框架开发的 NFT 拍卖市场系统，支持使用 ETH 或 ERC20 代币进行出价，集成 Chainlink 价格预言机进行价格比较，并支持 UUPS 和透明代理两种升级模式。

## 项目概述

本项目实现了一个功能完整的 NFT 拍卖市场，主要特性包括：

- **NFT 合约**：基于 ERC721 标准的 NFT 合约，支持铸造和元数据管理
- **拍卖功能**：支持创建拍卖、ETH/ERC20 出价、结束拍卖、撤回出价
- **价格预言机**：集成 Chainlink 价格预言机，自动将不同代币的出价转换为 USD 进行比较
- **合约升级**：支持 UUPS 和透明代理两种升级模式
- **手续费机制**：固定 2.5% 手续费，自动分配给合约所有者

## 项目结构

```
hardhat-auction/
├── contracts/
│   ├── MyNFT.sol                    # ERC721 NFT 合约
│   ├── Auction.sol                  # 拍卖合约核心实现（可升级基础合约）
│   ├── AuctionUUPS.sol             # UUPS 代理版本
│   ├── AuctionTransparent.sol      # 透明代理版本
│   ├── PriceOracle.sol             # Chainlink 价格预言机封装合约
│   └── mocks/
│       ├── MockERC20.sol           # Mock ERC20 代币（用于测试）
│       └── MockChainlinkAggregator.sol  # Mock Chainlink 聚合器（用于测试）
├── test/
│   ├── MyNFT.test.js               # NFT 合约测试
│   ├── Auction.test.js             # 拍卖核心功能测试
│   ├── AuctionUUPS.test.js         # UUPS 升级测试
│   ├── AuctionTransparent.test.js  # 透明代理升级测试
│   ├── PriceOracle.test.js         # 价格预言机测试
│   └── Integration.test.js         # 集成测试
├── scripts/
│   ├── deploy.js                   # 部署脚本
│   └── upgrade.js                  # 升级脚本
├── ignition/modules/
│   └── deploy.js                   # Hardhat Ignition 部署模块
└── README.md                       # 项目文档
```

## 合约架构

### MyNFT.sol
ERC721 NFT 合约，继承自 OpenZeppelin 的 `ERC721URIStorage` 和 `Ownable`。

**主要功能**：
- `mint(address to, string memory tokenURI)`: 铸造 NFT
- `setTokenURI(uint256 tokenId, string memory tokenURI)`: 设置 NFT URI

### PriceOracle.sol
Chainlink 价格预言机封装合约，用于获取 ETH 和 ERC20 代币的 USD 价格。

**主要功能**：
- `getETHPrice()`: 获取 ETH/USD 价格
- `getERC20Price(address token)`: 获取 ERC20/USD 价格
- `convertToUSD(uint256 amount, address token)`: 将代币数量转换为 USD 价值
- `setTokenPriceFeed(address token, address priceFeed)`: 设置 ERC20 代币的价格源

### Auction.sol
拍卖合约核心实现，继承自 `Initializable`、`OwnableUpgradeable` 和 `ReentrancyGuardUpgradeable`。

**主要功能**：
- `createAuction(...)`: 创建拍卖
- `bidWithETH(uint256 auctionId)`: 使用 ETH 出价
- `bidWithERC20(uint256 auctionId, uint256 amount)`: 使用 ERC20 出价
- `endAuction(uint256 auctionId)`: 结束拍卖
- `withdrawBid(uint256 auctionId)`: 撤回被超过的出价

**拍卖流程**：
1. 卖家创建拍卖，NFT 转移到合约
2. 买家使用 ETH 或 ERC20 出价（自动转换为 USD 比较）
3. 出价被超过的买家可以撤回资金
4. 拍卖结束后，NFT 转移给最高出价者，资金转移给卖家（扣除手续费）

### AuctionUUPS.sol
UUPS 代理版本的拍卖合约，继承自 `Auction` 和 `UUPSUpgradeable`。

**特点**：
- 升级逻辑在实现合约中
- Gas 消耗更低
- 需要实现 `_authorizeUpgrade()` 函数

### AuctionTransparent.sol
透明代理版本的拍卖合约，继承自 `Auction`。

**特点**：
- 升级逻辑在代理合约中
- 更安全，但 Gas 消耗更高
- 需要 ProxyAdmin 合约管理升级

## 安装和配置

### 1. 安装依赖

```bash
npm install
```

### 2. 配置环境变量（可选）

创建 `.env` 文件：

```env
SEPOLIA_URL=https://sepolia.infura.io/v3/YOUR_INFURA_KEY
PRIVATE_KEY=your_private_key_here
```

### 3. 配置 Chainlink 价格源

在 `hardhat.config.js` 中配置不同网络的 Chainlink 价格源地址：

```javascript
chainlink: {
  ethUsdPriceFeed: {
    sepolia: "0x694AA1769357215DE4FAC081bf1f309aDC325306",
    hardhat: "0x5f4eC3Df9cbd43714FE2740f5E3616155c5b8419",
  },
}
```

## 测试

### Hardhat 测试

运行所有测试：

```bash
npx hardhat test
```

运行特定测试文件：

```bash
npx hardhat test test/MyNFT.test.js
npx hardhat test test/Auction.test.js
npx hardhat test test/Integration.test.js
```

查看测试覆盖率（如果配置了）：

```bash
npx hardhat coverage
```

### 本地 Hardhat 环境测试

#### 1. 启动本地 Hardhat 节点

在项目根目录启动本地 Hardhat 节点：

```bash
npx hardhat node
```

节点将在 `http://127.0.0.1:8545` 运行，并自动创建20个测试账户，每个账户有10000 ETH。

#### 2. 使用 Hardhat 控制台测试

启动 Hardhat 控制台连接到本地节点：

```bash
npx hardhat console --network localhost
```

在控制台中可以直接运行 JavaScript 代码来测试合约。

#### 3. 部署合约

**使用部署脚本（推荐）**：

```bash
npx hardhat run scripts/deploy.js --network localhost
```

**或者在控制台中手动部署**：

```javascript
// 在 Hardhat 控制台中运行
const { ethers } = require("hardhat");

// 部署 MyNFT
const MyNFT = await ethers.getContractFactory("MyNFT");
const myNFT = await MyNFT.deploy("MyNFT", "MNFT");
await myNFT.waitForDeployment();
const myNFTAddress = await myNFT.getAddress();
console.log("MyNFT deployed to:", myNFTAddress);

// 部署 PriceOracle
const PriceOracle = await ethers.getContractFactory("PriceOracle");
const priceOracle = await PriceOracle.deploy("0x5f4eC3Df9cbd43714FE2740f5E3616155c5b8419"); // Mock Chainlink aggregator
await priceOracle.waitForDeployment();
const priceOracleAddress = await priceOracle.getAddress();
console.log("PriceOracle deployed to:", priceOracleAddress);

// 部署 AuctionUUPS 实现合约
const AuctionUUPS = await ethers.getContractFactory("AuctionUUPS");
const auctionUUPS = await AuctionUUPS.deploy();
await auctionUUPS.waitForDeployment();
const auctionUUPSAddress = await auctionUUPS.getAddress();
console.log("AuctionUUPS implementation deployed to:", auctionUUPSAddress);

// 部署 ERC1967Proxy
const ERC1967ProxyArtifact = require("@openzeppelin/contracts/build/contracts/ERC1967Proxy.json");
const ERC1967Proxy = await ethers.getContractFactoryFromArtifact(ERC1967ProxyArtifact);
const initData = AuctionUUPS.interface.encodeFunctionData("initialize", [priceOracleAddress]);
const proxy = await ERC1967Proxy.deploy(auctionUUPSAddress, initData);
await proxy.waitForDeployment();
const proxyAddress = await proxy.getAddress();
console.log("ERC1967Proxy deployed to:", proxyAddress);

// 获取代理合约实例
const auction = AuctionUUPS.attach(proxyAddress);
console.log("Auction proxy ready at:", proxyAddress);
```

#### 4. 测试流程示例

**准备工作**：
```javascript
// 在控制台中先定义变量（替换为实际部署的地址）
const myNFTAddress = "0x..."; // MyNFT 合约地址
const proxyAddress = "0x..."; // Auction 代理合约地址
const [seller, buyer1, buyer2] = await ethers.getSigners();
const myNFT = await ethers.getContractAt("MyNFT", myNFTAddress);
const auction = await ethers.getContractAt("AuctionUUPS", proxyAddress);
```

**完整拍卖流程**：

1. **铸造 NFT**：
   ```javascript
   await myNFT.connect(seller).mint(seller.address, "https://example.com/token/1");
   console.log("NFT minted to:", seller.address);
   ```

2. **创建拍卖**：
   ```javascript
   // 授权 NFT 给拍卖合约
   await myNFT.connect(seller).approve(proxyAddress, 1);

   // 创建拍卖（起拍价 $100，持续1小时，使用ETH）
   const oneHour = 3600;
   const startingPriceUSD = ethers.parseEther("100"); // $100 in USD (18 decimals)
   await auction.connect(seller).createAuction(
     myNFTAddress, // NFT 合约地址
     1, // tokenId
     startingPriceUSD,
     oneHour,
     ethers.ZeroAddress // ETH
   );
   console.log("Auction created with ID: 1");
   ```

3. **出价**：
   ```javascript
   // 使用 ETH 出价 $150
   const bidAmount = ethers.parseEther("0.06"); // 假设 ETH 价格为 $2500，0.06 ETH = $150
   await auction.connect(buyer1).bidWithETH(1, { value: bidAmount });
   console.log("Bid placed by buyer1");

   // 另一个买家出价 $200
   const higherBid = ethers.parseEther("0.08"); // 0.08 ETH = $200
   await auction.connect(buyer2).bidWithERC20(1, higherBid);
   console.log("Higher bid placed by buyer2");
   ```

4. **结束拍卖**：
   ```javascript
   // 等待拍卖结束或手动结束
   await auction.connect(seller).endAuction(1);
   console.log("Auction ended");

   // 检查结果
   const auctionInfo = await auction.getAuction(1);
   console.log("Winner:", auctionInfo.highestBidder);
   console.log("Final price:", ethers.formatEther(auctionInfo.currentHighestBid));
   ```

**常见问题排查**：

- **部署失败**：确保本地节点正在运行，账户有足够的 ETH
- **函数调用失败**：检查合约地址是否正确，参数类型是否匹配
- **价格获取失败**：确保 PriceOracle 中设置了正确的价格源地址
- **交易被撤销**：可能是 Gas 不足或网络拥堵

## 部署

### 使用 Hardhat 脚本部署

部署到本地网络：

```bash
npx hardhat run scripts/deploy.js --network hardhat
```

部署到 Sepolia 测试网：

```bash
npx hardhat run scripts/deploy.js --network sepolia
```

部署信息会保存在 `deployments/{network}.json` 文件中。

### 使用 Hardhat Ignition 部署

```bash
npx hardhat ignition deploy ignition/modules/deploy.js --network sepolia
```

### 升级合约

**升级 UUPS 代理**：

```bash
UPGRADE_TYPE=uups npx hardhat run scripts/upgrade.js --network sepolia
```

**升级透明代理**：

```bash
UPGRADE_TYPE=transparent npx hardhat run scripts/upgrade.js --network sepolia
```

## 合约交互示例

### 使用 Hardhat 控制台

```javascript
// 启动 Hardhat 控制台
npx hardhat console --network localhost

// 在控制台中执行以下命令：

// 1. 获取签名者和合约实例
const [deployer, seller, buyer] = await ethers.getSigners();
const myNFT = await ethers.getContractAt("MyNFT", "DEPLOYED_MY_NFT_ADDRESS");
const auction = await ethers.getContractAt("AuctionUUPS", "DEPLOYED_PROXY_ADDRESS");

// 2. 铸造 NFT
await myNFT.connect(seller).mint(seller.address, "https://example.com/nft/1");
console.log("NFT minted to seller");

// 3. 授权 NFT 给拍卖合约
await myNFT.connect(seller).approve("DEPLOYED_PROXY_ADDRESS", 1);

// 4. 创建拍卖（起拍价 $100，持续1小时，使用ETH）
await auction.connect(seller).createAuction(
  "DEPLOYED_MY_NFT_ADDRESS", // NFT 合约地址
  1, // tokenId
  ethers.parseEther("100"), // $100 起拍价
  3600, // 1小时持续时间
  ethers.ZeroAddress // 使用 ETH
);
console.log("Auction created");

// 5. 出价（假设当前 ETH 价格为 $2500/ETH）
const bidAmount = ethers.parseEther("0.06"); // 0.06 ETH = $150
await auction.connect(buyer).bidWithETH(1, { value: bidAmount });
console.log("Bid placed successfully");

// 6. 结束拍卖（等待时间结束或由卖家手动结束）
await auction.connect(seller).endAuction(1);
console.log("Auction ended");

// 7. 检查拍卖结果
const auctionInfo = await auction.getAuction(1);
console.log("Winner:", auctionInfo.highestBidder);
console.log("Final price:", ethers.formatEther(auctionInfo.currentHighestBid));
```

## 升级流程说明

### UUPS 代理升级

1. 部署新版本的实现合约
2. 调用代理合约的 `upgradeToAndCall()` 函数
3. 只有合约 owner 可以执行升级

### 透明代理升级

1. 部署新版本的实现合约
2. 通过 ProxyAdmin 调用 `upgrade()` 函数
3. 只有 ProxyAdmin 的所有者可以执行升级

### 升级注意事项

- **存储布局兼容性**：新版本合约的存储布局必须与旧版本兼容
- **初始化函数**：升级时不应再次调用 `initialize()`，除非使用 `upgradeToAndCall()` 调用其他初始化函数
- **测试**：升级前务必在测试网充分测试

## 安全考虑

1. **重入攻击防护**：使用 `ReentrancyGuard` 保护关键函数
2. **权限控制**：使用 `Ownable` 控制关键操作权限
3. **价格数据验证**：检查 Chainlink 价格数据的有效性和时效性
4. **溢出保护**：Solidity 0.8+ 自动处理溢出
5. **代理模式安全**：遵循 OpenZeppelin 的最佳实践

## 技术栈

- **Solidity**: ^0.8.28
- **Hardhat**: ^2.28.2
- **OpenZeppelin Contracts**: ^5.4.0
- **Chainlink Contracts**: ^1.1.0
- **Hardhat Deploy**: ^1.0.4

## 许可证

MIT