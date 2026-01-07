# Hardhat Auction 项目在 Remix 中的手工测试指南

## 项目概述

这是一个支持 ETH 和 ERC20 出价的 NFT 拍卖系统，包含以下合约：
- `Auction.sol` - 基础拍卖合约
- `AuctionTransparent.sol` - 透明代理版本的拍卖合约
- `AuctionUUPS.sol` - UUPS代理版本的拍卖合约
- `PriceOracle.sol` - 价格预言机合约
- `MyNFT.sol` - ERC721 NFT 合约
- `MockChainlinkAggregator.sol` - Mock Chainlink 价格聚合器（用于测试）
- `MockERC20.sol` - Mock ERC20 代币（用于测试）

## 第一步：代码复制到 Remix

### 1.1 创建新项目
1. 打开 [Remix IDE](https://remix.ethereum.org/)
2. 点击 "Create New File" 创建新项目
3. **安装必要插件：**
   - 点击左侧插件图标，搜索并安装 "OpenZeppelin Contracts"
   - 这样就可以直接使用 import 语句导入 OZ 合约
4. 在项目根目录创建以下文件夹结构：
```
auction-project/
├── contracts/
│   ├── Auction.sol
│   ├── AuctionTransparent.sol
│   ├── AuctionUUPS.sol
│   ├── PriceOracle.sol
│   ├── MyNFT.sol
│   └── mocks/
│       ├── MockChainlinkAggregator.sol
│       └── MockERC20.sol
```

### 1.2 复制合约代码

**两种方式：**

#### 方式一：使用 OpenZeppelin 插件 + Import 语句（推荐）

**你的合约已经是这种格式了！** 🎉

安装 OpenZeppelin 插件后，你可以直接复制现有的 Hardhat 合约代码，无需任何修改：

```solidity
// 你的 Auction.sol - 已经是最佳实践格式
import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721Receiver.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "./PriceOracle.sol";

contract Auction is Initializable, OwnableUpgradeable, ReentrancyGuardUpgradeable, IERC721Receiver {
    // ... 你的完整代码
}
```

**优势：**
- ✅ 代码与 Hardhat 项目完全一致
- ✅ 无需手动处理依赖
- ✅ 自动获得最新版本的合约
- ✅ Remix 会自动解析所有 import

#### 方式二：手动复制代码（备用方案）
如果插件不可用，可以手动复制所有依赖合约代码到项目中。

#### 1.2.1 复制基础合约

1. **复制 MyNFT.sol**
```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

// 移除导入，使用 Remix 的内置 OpenZeppelin 合约
// import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
// import "@openzeppelin/contracts/access/Ownable.sol";

// 直接粘贴 OpenZeppelin 合约内容或使用 Remix 插件
// 在 Remix 中，你可以：
// 1. 使用 OpenZeppelin 插件 (.openzeppelin)，或者
// 2. 手动复制需要的合约代码

// 这里粘贴 ERC721URIStorage 合约内容
// ... (从 OpenZeppelin 复制完整代码)

// 这里粘贴 Ownable 合约内容
// ... (从 OpenZeppelin 复制完整代码)

contract MyNFT is ERC721URIStorage, Ownable {
    uint256 private _tokenIds;

    constructor(string memory name, string memory symbol) ERC721(name, symbol) Ownable(msg.sender) {}

    function mint(address to, string memory tokenURI) public onlyOwner returns (uint256) {
        unchecked {
            _tokenIds++;
        }
        uint256 newTokenId = _tokenIds;

        _safeMint(to, newTokenId);
        _setTokenURI(newTokenId, tokenURI);

        return newTokenId;
    }

    function setTokenURI(uint256 tokenId, string memory tokenURI) public onlyOwner {
        require(_ownerOf(tokenId) != address(0), "Token does not exist");
        _setTokenURI(tokenId, tokenURI);
    }

    function totalSupply() public view returns (uint256) {
        return _tokenIds;
    }
}
```

2. **复制 MockChainlinkAggregator.sol**
```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

// 在 Remix 中使用完整接口定义
interface AggregatorV3Interface {
    function decimals() external view returns (uint8);
    function description() external view returns (string memory);
    function version() external view returns (uint256);
    function getRoundData(uint80 _roundId) external view returns (
        uint80 roundId,
        int256 answer,
        uint256 startedAt,
        uint256 updatedAt,
        uint80 answeredInRound
    );
    function latestRoundData() external view returns (
        uint80 roundId,
        int256 answer,
        uint256 startedAt,
        uint256 updatedAt,
        uint80 answeredInRound
    );
}

contract MockChainlinkAggregator is AggregatorV3Interface {
    uint8 public constant override decimals = 8;
    string public constant override description = "Mock Price Feed";
    uint256 public constant override version = 0;

    int256 private _price;
    uint256 private _updatedAt;
    uint80 private _roundId;

    constructor(int256 initialPrice) {
        _price = initialPrice;
        _updatedAt = block.timestamp;
        _roundId = 1;
    }

    function setPrice(int256 newPrice) external {
        _price = newPrice;
        _updatedAt = block.timestamp;
        _roundId++;
    }

    function latestRoundData()
        external
        view
        override
        returns (
            uint80 roundId,
            int256 answer,
            uint256 startedAt,
            uint256 updatedAt,
            uint80 answeredInRound
        )
    {
        return (_roundId, _price, _updatedAt, _updatedAt, _roundId);
    }

    function getRoundData(uint80 inputRoundId)
        external
        view
        override
        returns (
            uint80 roundId,
            int256 answer,
            uint256 startedAt,
            uint256 updatedAt,
            uint80 answeredInRound
        )
    {
        return (_roundId, _price, _updatedAt, _updatedAt, _roundId);
    }
}
```

3. **复制 MockERC20.sol**
```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

// 在 Remix 中使用完整 ERC20 实现
// 这里需要复制完整的 ERC20 合约代码
// 或者使用 Remix 的 OpenZeppelin 插件

contract MockERC20 is ERC20 {
    uint8 private _decimals;

    constructor(string memory name, string memory symbol, uint8 decimals_) ERC20(name, symbol) {
        _decimals = decimals_;
        _mint(msg.sender, 1000000 * 10 ** decimals_);
    }

    function decimals() public view virtual override returns (uint8) {
        return _decimals;
    }

    function mint(address to, uint256 amount) external {
        _mint(to, amount);
    }
}
```

4. **复制 PriceOracle.sol**
```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

// 复制完整的 AggregatorV3Interface
interface AggregatorV3Interface {
    function decimals() external view returns (uint8);
    function description() external view returns (string memory);
    function version() external view returns (uint256);
    function getRoundData(uint80 _roundId) external view returns (
        uint80 roundId,
        int256 answer,
        uint256 startedAt,
        uint256 updatedAt,
        uint80 answeredInRound
    );
    function latestRoundData() external view returns (
        uint80 roundId,
        int256 answer,
        uint256 startedAt,
        uint256 updatedAt,
        uint80 answeredInRound
    );
}

contract PriceOracle {
    // ... (复制完整代码)
}

// ERC20 元数据接口
interface IERC20Metadata {
    function decimals() external view returns (uint8);
}
```

5. **复制 Auction.sol**
```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

// 需要导入或直接包含所有依赖的 OpenZeppelin 合约
// import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
// import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
// import "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
// import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
// import "@openzeppelin/contracts/token/ERC721/IERC721Receiver.sol";
// import "@openzeppelin/contracts/token/ERC20/IERC20.sol";

// 在 Remix 中，你需要：
// 1. 使用 OpenZeppelin 插件，或
// 2. 手动复制所有依赖合约

contract Auction is Initializable, OwnableUpgradeable, ReentrancyGuardUpgradeable, IERC721Receiver {
    // ... (复制完整代码)
}
```

## 第二步：设置 Remix 环境

### 2.1 安装必要的插件
1. 在 Remix 中点击插件图标
2. 搜索并安装 "OpenZeppelin" 插件
3. 搜索并安装 "Chainlink" 插件（如果有的话）

### 2.2 编译设置

**使用 OpenZeppelin 插件时：**
1. 选择 Solidity 编译器版本：`0.8.28`
2. **重要：设置 EVM 版本为 Cancun**
   - 在 Remix 的 Solidity 编译器面板中
   - 点击 "Advanced Configurations"
   - 找到 "EVM Version" 选项
   - 从下拉菜单中选择 "cancun"
3. 启用优化：200 runs
4. **直接编译你的合约文件**
   - Remix 会自动解析 @openzeppelin 和 @chainlink 的 import
   - 无需手动编译依赖合约

**如果不使用插件（手动复制方式）：**
1. 编译所有合约，按以下顺序：
   1. `MockChainlinkAggregator.sol`
   2. `MockERC20.sol`
   3. `MyNFT.sol`
   4. `PriceOracle.sol`
   5. `Auction.sol`
   6. `AuctionTransparent.sol`
   7. `AuctionUUPS.sol`

## 第三步：部署合约

### 3.1 账户管理

**为什么需要多个账户：**
在真实的区块链应用中，不同的角色使用不同的账户：
- **部署者账户**：部署合约和管理权限
- **卖家账户**：拥有 NFT 并发起拍卖
- **买家账户**：参与竞价购买

**在 Remix 中设置多个账户：**
1. **生成新账户：**
   - 在部署面板中，点击 "Account" 下拉菜单
   - 选择 "Generate new account"
   - 为账户命名（如 "Seller", "Buyer1", "Buyer2"）

2. **账户分配：**
   - **账户 1（部署者）**：部署所有合约，设置价格源，作为 MyNFT 所有者铸造 NFT
   - **账户 2（卖家）**：接收 NFT，创建拍卖，结束拍卖
   - **账户 3+（买家）**：参与出价竞价

   1：0x5B38Da6a701c568545dCfcB03FcB875f56beddC4
   2：0xAb8483F64d9C6d1EcF9b849Ae677dD3315835cb2
   3：0x4B20993Bc481177ec7E8f571ceCaE8A9e22C02db

3. **账户切换：**
   - 在每个操作前，确认使用正确的账户
   - 在部署面板的 "Account" 菜单中切换账户

### 3.2 部署顺序

**重要提醒：** 在部署任何合约之前，请确保 Solidity 编译器的 EVM 版本设置为 "cancun"。

1. **部署 MockChainlinkAggregator (ETH/USD)**
   - 构造函数参数：`300000000000` (3000 * 10^8，代表 $3000)
   - 在 Remix 中：
     - 选择 `MockChainlinkAggregator.sol`
     - 在构造函数参数输入框中输入：`300000000000`
     - 点击 "Deploy"
   - 记录部署地址：`mockEthPriceFeedAddress`
   - 0xd9145CCE52D386f254917e481eB44e9943F39138

2. **部署 MockChainlinkAggregator (Token/USD)**
   - 构造函数参数：`100000000` (1 * 10^8，代表 $1)
   - 在 Remix 中：
     - 再次选择 `MockChainlinkAggregator.sol`
     - 在构造函数参数输入框中输入：`100000000`
     - 点击 "Deploy"
   - 记录部署地址：`mockTokenPriceFeedAddress`
   - 0xd8b934580fcE35a11B58C6D73aDeE468a2833fa8

### 3.3 参数说明

**MockChainlinkAggregator 构造函数参数详解：**

- **参数类型：** `int256 initialPrice`
- **价格精度：** 8位小数（Chainlink 标准）
- **ETH/USD 价格源：** `300000000000`
  - 含义：3000.00000000 USD（8位小数）
  - 计算：3000 * 10^8 = 300000000000
- **Token/USD 价格源：** `100000000`
  - 含义：1.00000000 USD（8位小数）
  - 计算：1 * 10^8 = 100000000

**验证部署：**
部署后，点击已部署合约的 `latestRoundData()` 函数，应该返回：
- ETH 价格源：`(1, 300000000000, timestamp, timestamp, 1)`
- Token 价格源：`(1, 100000000, timestamp, timestamp, 1)`

3. **部署 PriceOracle**
   - 参数：`mockEthPriceFeedAddress`
   - 记录部署地址：`priceOracleAddress`
   - 0xf8e81D47203A594245E36C48e151709F0C19fBe8

4. **部署 MockERC20**
   - 参数：`"TestToken"`, `"TEST"`, `18`
   - 记录部署地址：`mockERC20Address`
   - 0xD7ACd2a9FD159E69Bb102A1ca21C9a3e3A5F771B
   - 调用 `setTokenPriceFeed(mockERC20Address, mockTokenPriceFeedAddress)` 设置代币价格源

5. **部署 MyNFT**
   - **使用账户：** 账户 1（部署者）
   - 参数：`"MyNFT"`, `"MNFT"`
   - 记录部署地址：`myNFTAddress`
   - 0xDA0bab807633f07f013f94DD0E6A4F96F8742B53
   - **铸造 NFT：** 使用账户 1（部署者）- 因为 MyNFT 的 onlyOwner 限制
     - 保持使用账户 1（当前部署者账户）
     - 调用 `mint(账户2地址, "https://example.com/token/1")` 铸造 NFT 给卖家账户
     - 记录铸造的 tokenId（通常是 1）
     - **注意：** 虽然铸造给账户 2，但铸造动作必须由合约所有者（账户 1）执行
     - 1

6. **部署 Auction 合约**
   - **使用账户：** 账户 1（部署者）
   - **重要：** 确保 EVM 版本设置为 "cancun"
   - **部署环境设置：**
     - 在部署面板中，选择 "Remix VM (Cancun)" 作为环境
     - 如果没有 Cancun 选项，选择 "Remix VM (London)" 或最新版本
     - 确保 Account 有足够的 ETH（至少 1 ETH）
   - **部署步骤：**
     - 选择 `Auction.sol`
     - 确认构造函数参数为空（无需输入参数）
     - 点击 "Deploy" 按钮
   - **如果部署失败：**
     - 检查 Gas Limit（尝试设置更高的值，如 8,000,000）
     - 确认所有依赖合约已正确部署
     - 尝试刷新页面后重新部署
   - 部署成功后，调用 `initialize(priceOracleAddress)` 初始化
   - 记录部署地址：`auctionAddress`
   - 0x9D7f74d0C41E726EC95884E0e97Fa6129e3b5E99

### 3.4 部署故障排除

**常见部署问题及解决方案：**

1. **"evm version: cancun" 错误：**
   - 确认编译器 EVM 版本设置为 "cancun"
   - 确认部署环境选择 "Remix VM (Cancun)"

2. **部署失败 - Gas 不足：**
   - 手动设置 Gas Limit 为 8,000,000 或更高
   - 确保账户有足够的测试 ETH

3. **合约大小超出限制：**
   - 如果提示合约过大，可能是因为包含了太多依赖
   - 考虑使用代理模式（Transparent 或 UUPS）来减小部署大小

4. **初始化失败：**
   - 确保 PriceOracle 地址正确
   - 确保调用 initialize 时使用正确的参数

5. **Remix 页面刷新后丢失部署：**
   - 每次使用前先部署所有依赖合约
   - 或者使用 Remix 的工作区保存功能

### 3.5 参数格式问题

**地址参数格式：**
- **正确格式：** `"0x0000000000000000000000000000000000000000"`
- **错误格式：** `"0"` 或 `0` 或 `"0x0"`
- **ETH 表示：** 使用完整的 40 位零地址，不要使用 `"0"`

**大数字参数格式：**
- **正确格式：** `100000000000000000000` （不带引号）
- **错误格式：** `"100000000000000000000"` （带引号会变成字符串）
- **在 Remix 中：** 直接输入数字，不要加引号

**时间参数格式：**
- **秒数：** `3600` （1小时）
- **不要使用：** `"3600"` 或复杂的时间表达式

### 3.6 createAuction 故障排除

**如果 createAuction 交易 revert，逐步检查：**

1. **账户权限检查：**
   ```javascript
   // 1. 确认你是 NFT 的所有者
   myNFT.ownerOf(1)  // 应该返回你的地址

   // 2. 确认 NFT 已授权给拍卖合约
   myNFT.getApproved(1)  // 应该返回拍卖合约地址

   // 3. 如果没有授权，执行授权
   myNFT.approve(auctionAddress, 1)
   ```

2. **参数验证：**
   ```javascript
   // 确保所有地址都是完整的 40 位格式
   // 确保数字参数不带引号
   // 确保 startingPriceUSD > 0
   // 确保 duration > 0
   ```

3. **合约状态检查：**
   ```javascript
   // 确认 PriceOracle 已正确初始化
   auction.priceOracle()  // 应该返回 PriceOracle 地址

   // 确认拍卖合约已初始化
   auction.owner()  // 应该返回部署者地址
   ```

4. **常见错误原因：**
   - **账户权限问题：**
     - MyNFT.mint 必须由合约所有者（账户 1）调用
     - createAuction 必须由 NFT 实际所有者（账户 2）调用
   - NFT 没有铸造或不属于你
   - 忘记了 approve 授权
   - 参数格式错误（地址不完整、数字带引号等）
   - 合约没有正确初始化

### 3.2 部署验证

部署完成后，验证以下功能：
- PriceOracle 可以获取 ETH 价格
- PriceOracle 可以获取 ERC20 价格
- MockERC20 有足够的代币供应
- MyNFT 成功铸造 NFT

## 第四步：人工测试用例

### 4.0 测试准备提醒

**重要：在 Remix 中调用合约函数时：**
- **地址参数：** 使用完整的 40 位十六进制地址，如 `"0x0000000000000000000000000000000000000000"`
- **数字参数：** 直接输入数字，不要加引号，如 `100000000000000000000`
- **字符串参数：** 使用双引号，如 `"TestToken"`
- **如果遇到 "invalid address" 错误：** 检查地址格式是否正确

**账户角色分配：**
- **账户 1（部署者）**：部署合约，设置价格源
- **账户 2（卖家）**：铸造 NFT，创建拍卖
- **账户 3（买家）**：参与出价
- **在每个测试步骤前确认使用正确的账户！**

### 4.1 Hardhat 部署和升级流程

**完整的开发和测试流程：**

1. **合约开发：**
   - 编写基础合约（Auction.sol）
   - 创建升级版本（AuctionV2.sol）

2. **部署测试：**
   ```bash
   # 部署所有合约
   npx hardhat deploy --tags deploy

   # 升级到新版本
   npx hardhat deploy --tags upgrade
   ```

3. **功能测试：**
   ```bash
   # 运行所有测试
   npx hardhat test

   # 只运行升级测试
   npx hardhat test --grep "合约升级测试"
   ```

4. **审计验证：**
   - 状态变量兼容性
   - 函数签名保持
   - 权限控制完整
   - 升级安全性

### 4.2 基础功能测试

#### 测试 1：创建 ETH 拍卖
1. **准备工作：**
   - **步骤 1：铸造 NFT（使用账户 1）**
     - 切换到账户 1（部署者）
     - 调用 `myNFT.mint(账户2地址, "https://example.com/token/1")`
   - **步骤 2：创建拍卖（使用账户 2）**
     - 切换到账户 2（卖家）- NFT 现在属于账户 2
     - 调用 `myNFT.approve(auctionAddress, 1)` 授权拍卖合约
     - 验证授权：调用 `myNFT.getApproved(1)` 确认返回 auctionAddress
     - 0x9D7f74d0C41E726EC95884E0e97Fa6129e3b5E99

2. **创建拍卖：**
   - **在 Remix 中的输入参数：**
     - `nftContract`: `myNFTAddress`（替换为实际的 MyNFT 合约地址）
     - `tokenId`: `1`
     - `startingPriceUSD`: `100000000000000000000` （100 * 10^18，代表 $100）
     - `duration`: `3600` （1小时 = 3600秒）
     - `paymentToken`: `0x0000000000000000000000000000000000000000` （address(0)，表示使用 ETH）
   - **完整调用示例：**
     ```javascript
     auction.createAuction(
       "0xDA0bab807633f07f013f94DD0E6A4F96F8742B53", // myNFTAddress
       1,                                            // tokenId
       100000000000000000000,                       // $100 in 18 decimals
       3600,                                         // 1 hour
       "0x0000000000000000000000000000000000000000"  // address(0) for ETH
     )
     ```
   - 成功

3. **故障排除：**
   - **如果 revert，检查以下各项：**
     - 使用正确的账户（MyNFT owner）
     - NFT 是否已铸造：`myNFT.ownerOf(1)` 应该返回你的地址
     - NFT 是否已授权：`myNFT.getApproved(1)` 应该返回 auction 地址
     - 所有参数格式是否正确

3. **验证：**
   - 事件 `AuctionCreated` 被触发
   - 调用 `auction.getAuction(1)` 查看拍卖信息
   - 确认 NFT 已转移到拍卖合约

#### 测试 2：使用 ETH 出价
1. **使用账户：** 账户 3（买家）
2. **出价：**
   - 计算出价金额：假设 ETH 价格 $3000，要出价 $150，需要发送 0.05 ETH
   - 调用 `auction.bidWithETH(1)` 并发送 0.05 ETH

3. **验证：**
   - 事件 `BidPlaced` 被触发
   - 调用 `auction.getAuction(1)` 确认最高出价更新
   - 确认 `auction.bids(1, bidder1Address)` 返回正确的 USD 值

#### 测试 3：多个出价者竞价
1. **使用账户：** 账户 3（买家）或其他买家账户
2. **更高出价：**
   - 出价 $200，需要发送约 0.0667 ETH
   - 调用 `auction.bidWithETH(1)` 发送更高金额

3. **验证：**
   - bidder1 的出价被退回
   - bidder2 成为最高出价者
   - 调用 `auction.withdrawableBids(1, bidder1Address)` 查看可提取金额

#### 测试 4：撤回被超过的出价
1. **切换到 bidder1 账户**
2. **撤回出价：**
   - 调用 `auction.withdrawBid(1)`

3. **验证：**
   - ETH 被退回到 bidder1 账户
   - 事件 `BidWithdrawn` 被触发

### 4.3 ERC20 出价测试

#### 测试 5：创建 ERC20 拍卖
1. **创建新拍卖：**
   - 调用 `myNFT.mint(yourAddress, "https://example.com/token/2")` 铸造新 NFT
   - 授权：`myNFT.approve(auctionAddress, 2)`
   - 创建拍卖：`auction.createAuction(myNFTAddress, 2, ethers.utils.parseEther("50"), 3600, mockERC20Address)`

#### 测试 6：使用 ERC20 出价
1. **切换到 bidder1**
2. **授权代币：**
   - 调用 `mockERC20.approve(auctionAddress, ethers.utils.parseEther("100"))`

3. **出价：**
   - 调用 `auction.bidWithERC20(1, ethers.utils.parseEther("60"))` （出价 60 代币 = $60）

4. **验证：**
   - 代币从 bidder1 转移到拍卖合约
   - 事件 `BidPlaced` 被触发

### 4.4 拍卖结束测试

#### 测试 7：时间结束后结束拍卖
1. **等待拍卖结束或加速时间**
   - 在 Remix 中，可以使用 "Time Travel" 功能加速时间

2. **结束拍卖：**
   - 调用 `auction.endAuction(1)`

3. **验证：**
   - NFT 转移给最高出价者
   - 拍卖收入分配正确（97.5% 给卖家，2.5% 给合约所有者）
   - 事件 `AuctionEnded` 被触发

#### 测试 8：提前结束拍卖（卖家）
1. **创建新拍卖**
2. **卖家提前结束：**
   - 调用 `auction.endAuction(auctionId)`

3. **验证：**
   - 拍卖立即结束
   - NFT 退还给卖家

### 4.5 Hardhat Upgrades 测试

#### 完整的升级测试流程

1. **初始部署：**
   ```bash
   npx hardhat deploy --tags deploy
   ```

2. **运行升级脚本：**
   ```bash
   npx hardhat deploy --tags upgrade
   ```

3. **运行升级测试：**
   ```bash
   npx hardhat test test/Upgrade.test.js
   ```

#### 测试覆盖

- ✅ 透明代理部署和升级 (`upgrades.deployProxy` / `upgrades.upgradeProxy`)
- ✅ UUPS 代理部署和升级
- ✅ 状态保持性验证
- ✅ 新功能测试（暂停、批量操作）
- ✅ 权限控制验证
- ✅ 升级安全性检查

#### 升级脚本说明

**deploy.js**: 使用 `upgrades.deployProxy()` 部署透明代理和 UUPS 代理
**upgrade.js**: 使用 `upgrades.upgradeProxy()` 执行合约升级
**缓存文件**: `.cache/proxies.json` 保存代理地址信息

## 第五步：错误情况测试

### 5.1 边界条件测试

#### 测试 11：出价低于起拍价
- 预期：交易 revert，错误信息 "Bid below starting price"

#### 测试 12：拍卖结束后出价
- 预期：交易 revert，错误信息 "Auction already ended"

#### 测试 13：重复结束拍卖
- 预期：交易 revert，错误信息 "Auction already ended"

#### 测试 14：未授权的 NFT
- 预期：交易 revert，错误信息 "NFT not approved"

#### 测试 15：非所有者创建拍卖
- 预期：交易 revert

### 5.2 权限测试

#### 测试 16：非所有者升级合约
- 预期：透明代理 revert，UUPS 代理 revert

#### 测试 17：非最高出价者撤回出价
- 预期：交易 revert，错误信息 "Cannot withdraw current highest bid"

## 第六步：测试覆盖检查

确保测试覆盖以下场景：
- ✅ ETH 出价流程
- ✅ ERC20 出价流程
- ✅ 多个竞价者
- ✅ 出价撤回
- ✅ 拍卖结束和资金分配
- ✅ 提前结束拍卖
- ✅ 代理合约升级
- ✅ 错误处理
- ✅ 权限控制
- ✅ 边界条件

## 第七步：清理和总结

1. **记录所有测试结果**
2. **截图重要交易和事件**
3. **验证所有合约功能正常**
4. **确认 gas 使用合理**
5. **总结发现的问题和改进建议**

## 附录：常用工具函数

### 在 Remix 中快速转换单位
```javascript
// ETH to Wei
ethers.utils.parseEther("1") // 1 ETH = 1000000000000000000 Wei

// USD to Contract units (18 decimals)
ethers.utils.parseEther("100") // $100 = 100000000000000000000

// 计算出价金额
// 如果 ETH = $3000，要出价 $150，则需要：150 / 3000 = 0.05 ETH
ethers.utils.parseEther("0.05")
```

### 常用地址
- Zero Address: `0x0000000000000000000000000000000000000000`
- 你的地址：通过 `web3.eth.accounts[0]` 获取

### 时间操作
- 当前时间：`Math.floor(Date.now() / 1000)`
- 1小时后：`Math.floor(Date.now() / 1000) + 3600`

这个指南涵盖了项目的完整测试流程。在 Remix 中进行手工测试时，请按照步骤仔细执行，确保每个功能都经过验证。
