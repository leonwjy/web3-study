require("@nomicfoundation/hardhat-toolbox");
require('hardhat-deploy');
require("@openzeppelin/hardhat-upgrades");

/** @type import('hardhat/config').HardhatUserConfig */
module.exports = {
  solidity: {
    version: "0.8.28",
    settings: {
      optimizer: {
        enabled: true,
        runs: 200,
      },
    },
  },
  namedAccounts: {
    deployer: {
      default: 0,
    },
    user1: {
      default: 1,
    },
    user2: {
      default: 2,
    },
  },
  networks: {
    hardhat: {
      chainId: 31337,
    },
    sepolia: {
      url: process.env.SEPOLIA_URL || "",
      accounts: process.env.PRIVATE_KEY ? [process.env.PRIVATE_KEY] : [],
      chainId: 11155111,
    },
  },
  // Chainlink 价格源地址配置
  chainlink: {
    ethUsdPriceFeed: {
      sepolia: "0x694AA1769357215DE4FAC081bf1f309aDC325306",
      hardhat: "0x5f4eC3Df9cbd43714FE2740f5E3616155c5b8419", // Mainnet address for fork testing
    },
  },
};
