// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

// ============================================
// ERC20 接口定义
// ============================================

interface IERC20 {
    function transfer(address to, uint256 amount) external returns (bool);
    function transferFrom(address from, address to, uint256 amount) external returns (bool);
    function balanceOf(address account) external view returns (uint256);
    function approve(address spender, uint256 amount) external returns (bool);
    function allowance(address owner, address spender) external view returns (uint256);
}

contract SwapContract {
    // ============================================
    // 自定义错误定义
    // ============================================

    error InsufficientAllowance(address token, address owner, address spender, uint256 required, uint256 available);
    error TransferFailed(address token, address from, address to, uint256 amount);
    error InsufficientBalance(address token, address account, uint256 required, uint256 available);

    // ============================================
    // 事件定义
    // ============================================

    event TokenSwap(
        address indexed user,
        address indexed tokenIn,
        address indexed tokenOut,
        uint256 amountIn,
        uint256 amountOut
    );

    event Deposit(
        address indexed user,
        address indexed token,
        uint256 amount
    );

    event Withdraw(
        address indexed user,
        address indexed token,
        uint256 amount
    );

    // ============================================
    // 状态变量
    // ============================================

    // 合约地址映射（用于存储用户存入的代币）
    mapping(address => mapping(address => uint256)) public userBalances;

    // ============================================
    // 主要功能模块
    // ============================================

    /**
     * 存款代币到合约
     * @param token 代币合约地址
     * @param amount 存款数量
     */
    function deposit(address token, uint256 amount) external {
        // 检查用户余额
        uint256 userBalance = IERC20(token).balanceOf(msg.sender);
        if (userBalance < amount) {
            revert InsufficientBalance(token, msg.sender, amount, userBalance);
        }

        // 检查授权额度
        uint256 allowance = IERC20(token).allowance(msg.sender, address(this));
        if (allowance < amount) {
            revert InsufficientAllowance(token, msg.sender, address(this), amount, allowance);
        }

        // 从用户转移代币到合约
        bool success = IERC20(token).transferFrom(msg.sender, address(this), amount);
        if (!success) {
            revert TransferFailed(token, msg.sender, address(this), amount);
        }

        // 更新用户在合约中的余额
        userBalances[msg.sender][token] += amount;

        emit Deposit(msg.sender, token, amount);
    }

    /**
     * 从合约提取代币
     * @param token 代币合约地址
     * @param amount 提取数量
     */
    function withdraw(address token, uint256 amount) external {
        // 检查用户在合约中的余额
        uint256 contractBalance = userBalances[msg.sender][token];
        if (contractBalance < amount) {
            revert InsufficientBalance(token, msg.sender, amount, contractBalance);
        }

        // 更新用户余额（先更新状态）
        userBalances[msg.sender][token] -= amount;

        // 转移代币给用户
        bool success = IERC20(token).transfer(msg.sender, amount);
        if (!success) {
            // 如果转账失败，恢复状态
            userBalances[msg.sender][token] += amount;
            revert TransferFailed(token, address(this), msg.sender, amount);
        }

        emit Withdraw(msg.sender, token, amount);
    }

    /**
     * 代币交换功能
     * 将用户存放在合约中的 tokenIn 交换为 tokenOut
     * 这里使用简单的 1:1 交换比例作为示例
     *
     * @param tokenIn 输入代币合约地址
     * @param tokenOut 输出代币合约地址
     * @param amountIn 输入代币数量
     */
    function swap(
        address tokenIn,
        address tokenOut,
        uint256 amountIn
    ) external {
        // 检查用户在合约中的 tokenIn 余额
        uint256 userTokenInBalance = userBalances[msg.sender][tokenIn];
        if (userTokenInBalance < amountIn) {
            revert InsufficientBalance(tokenIn, msg.sender, amountIn, userTokenInBalance);
        }

        // 检查合约是否有足够的 tokenOut 来交换
        uint256 contractTokenOutBalance = userBalances[address(this)][tokenOut];
        if (contractTokenOutBalance < amountIn) {
            revert InsufficientBalance(tokenOut, address(this), amountIn, contractTokenOutBalance);
        }

        // 计算输出数量（这里使用 1:1 的简单比例）
        uint256 amountOut = amountIn;

        // 先更新状态（防止重入攻击）
        userBalances[msg.sender][tokenIn] -= amountIn;
        userBalances[address(this)][tokenOut] -= amountOut;

        // 增加用户的 tokenOut 余额
        userBalances[msg.sender][tokenOut] += amountOut;

        // 增加合约的 tokenIn 余额（作为储备）
        userBalances[address(this)][tokenIn] += amountIn;

        emit TokenSwap(msg.sender, tokenIn, tokenOut, amountIn, amountOut);
    }

    // ============================================
    // 查询功能模块
    // ============================================

    /**
     * 获取用户在合约中的代币余额
     * @param user 用户地址
     * @param token 代币合约地址
     */
    function getUserBalance(address user, address token) external view returns (uint256) {
        return userBalances[user][token];
    }
}
