// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/security/Pausable.sol";

contract MyToken is ERC20, Ownable, Pausable {

    // 构造函数
    constructor(
        string memory name,
        string memory symbol,
        uint256 initialSupply
    ) ERC20(name, symbol) {
        // 铸造初始供应量给合约部署者
        _mint(msg.sender, initialSupply * 10**decimals());
    }

    // 只有owner可以铸造代币
    function mint(address to, uint256 amount) public onlyOwner {
        _mint(to, amount);
    }

    // 用户可以销毁自己的代币
    function burn(uint256 amount) public {
        _burn(msg.sender, amount);
    }

    // owner可以销毁任意地址的代币（用于治理）
    function burnFrom(address account, uint256 amount) public onlyOwner {
        _burn(account, amount);
    }

    // 重写transfer函数，添加暂停检查
    function transfer(address to, uint256 amount)
        public
        override
        whenNotPaused
        returns (bool)
    {
        return super.transfer(to, amount);
    }

    // 重写transferFrom函数，添加暂停检查
    function transferFrom(address from, address to, uint256 amount)
        public
        override
        whenNotPaused
        returns (bool)
    {
        return super.transferFrom(from, to, amount);
    }
}