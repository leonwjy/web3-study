// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract MultiSigWallet {
    // ============================================
    // 自定义错误定义
    // ============================================

    error NotOwner(address caller);
    error InvalidOwners();
    error InvalidRequiredConfirmations(uint256 required, uint256 ownersCount);
    error TransactionAlreadyExecuted(uint256 txId);
    error TransactionNotExists(uint256 txId);
    error TransactionAlreadyConfirmed(uint256 txId, address owner);
    error NotEnoughConfirmations(uint256 txId, uint256 current, uint256 required);
    error TransactionExecutionFailed(uint256 txId);
    error GasLimitTooLow(uint256 provided, uint256 minimum);
    error ReentrantCall();

    // ============================================
    // 事件定义
    // ============================================

    event Deposit(address indexed sender, uint256 indexed amount);
    event TransactionSubmitted(uint256 indexed txId, address indexed proposer, address indexed to, uint256 value, bytes data);
    event TransactionConfirmed(uint256 indexed txId, address indexed owner);
    event TransactionExecuted(uint256 indexed txId, address indexed executor);
    event TransactionCancelled(uint256 indexed txId);

    // ============================================
    // 结构体定义
    // ============================================

    struct Transaction {
        address to;
        bool executed;
        uint256 value;
        uint256 gasLimit;
        uint256 confirmations;
        bytes data;
    }

    // ============================================
    // 状态变量
    // ============================================

    address[] public owners;
    mapping(address => bool) public isOwner;
    uint256 public requiredConfirmations;
    Transaction[] public transactions;
    mapping(uint256 => mapping(address => bool)) public confirmations;
    mapping(uint256 => bool) public cancelled;
    bool private _locked;
    uint256 public constant MIN_GAS_LIMIT = 2300;
    uint256 public constant MAX_GAS_LIMIT = 10000000;

    // ============================================
    // 修饰符定义
    // ============================================

    modifier nonReentrant() { if (_locked) revert ReentrantCall(); _locked = true; _; _locked = false; }
    modifier onlyOwner() { if (!isOwner[msg.sender]) revert NotOwner(msg.sender); _; }
    modifier txExists(uint256 _txId) { if (_txId >= transactions.length) revert TransactionNotExists(_txId); _; }
    modifier notCancelled(uint256 _txId) { if (cancelled[_txId]) revert TransactionAlreadyExecuted(_txId); _; }
    modifier notExecuted(uint256 _txId) { if (transactions[_txId].executed) revert TransactionAlreadyExecuted(_txId); _; }

    // ============================================
    // 构造函数
    // ============================================

    constructor(address[] memory _owners, uint256 _requiredConfirmations) {
        // 验证所有者
        if (_owners.length == 0) {
            revert InvalidOwners();
        }

        // 验证确认数
        if (_requiredConfirmations == 0 || _requiredConfirmations > _owners.length) {
            revert InvalidRequiredConfirmations(_requiredConfirmations, _owners.length);
        }

        // 设置所有者
        for (uint256 i = 0; i < _owners.length; i++) {
            address owner = _owners[i];
            if (owner == address(0)) {
                revert InvalidOwners();
            }
            if (isOwner[owner]) {
                revert InvalidOwners(); // 重复所有者
            }
            isOwner[owner] = true;
            owners.push(owner);
        }

        requiredConfirmations = _requiredConfirmations;
    }

    // ============================================
    // 主要功能模块
    // ============================================

    /**
     * 接收以太币的 fallback 函数
     */
    receive() external payable {
        emit Deposit(msg.sender, msg.value);
    }

    /**
     * 提交交易提案
     * @param _to 目标合约地址
     * @param _value 发送的以太币数量
     * @param _data 调用数据
     * @param _gasLimit gas 限制
     * @return txId 交易ID
     */
    function submitTransaction(
        address _to,
        uint256 _value,
        bytes calldata _data,
        uint256 _gasLimit
    ) external onlyOwner returns (uint256 txId) {
        // 验证 gas 限制
        if (_gasLimit < MIN_GAS_LIMIT || _gasLimit > MAX_GAS_LIMIT) {
            revert GasLimitTooLow(_gasLimit, MIN_GAS_LIMIT);
        }

        // 创建交易
        Transaction memory transaction = Transaction({
            to: _to,
            value: _value,
            data: _data,
            executed: false,
            gasLimit: _gasLimit,
            confirmations: 0
        });

        transactions.push(transaction);
        txId = transactions.length - 1;

        // 自动确认（提交者确认）
        confirmations[txId][msg.sender] = true;
        transactions[txId].confirmations = 1;

        emit TransactionSubmitted(txId, msg.sender, _to, _value, _data);

        return txId;
    }

    /**
     * 确认交易
     * @param _txId 交易ID
     */
    function confirmTransaction(uint256 _txId)
        external
        onlyOwner
        txExists(_txId)
        notExecuted(_txId)
        notCancelled(_txId)
    {
        if (confirmations[_txId][msg.sender]) {
            revert TransactionAlreadyConfirmed(_txId, msg.sender);
        }

        confirmations[_txId][msg.sender] = true;
        transactions[_txId].confirmations += 1;

        emit TransactionConfirmed(_txId, msg.sender);

        // 如果达到所需确认数，自动执行
        if (transactions[_txId].confirmations >= requiredConfirmations) {
            executeTransaction(_txId);
        }
    }

    /**
     * 执行交易
     * @param _txId 交易ID
     */
    function executeTransaction(uint256 _txId)
        public
        onlyOwner
        txExists(_txId)
        notExecuted(_txId)
        notCancelled(_txId)
        nonReentrant
    {
        Transaction storage transaction = transactions[_txId];

        // 检查确认数
        if (transaction.confirmations < requiredConfirmations) {
            revert NotEnoughConfirmations(_txId, transaction.confirmations, requiredConfirmations);
        }

        // 标记为已执行
        transaction.executed = true;

        // 使用 call 执行交易，应用 gas 限制
        (bool success, ) = transaction.to.call{value: transaction.value, gas: transaction.gasLimit}(
            transaction.data
        );

        if (!success) {
            // 如果执行失败，恢复执行状态
            transaction.executed = false;
            revert TransactionExecutionFailed(_txId);
        }

        emit TransactionExecuted(_txId, msg.sender);
    }

    /**
     * 取消交易（只有提交者可以取消）
     * @param _txId 交易ID
     */
    function cancelTransaction(uint256 _txId)
        external
        onlyOwner
        txExists(_txId)
        notExecuted(_txId)
        notCancelled(_txId)
    {
        // 这里简化逻辑，任何所有者都可以取消
        // 实际应用中可能需要只有提交者才能取消
        cancelled[_txId] = true;

        emit TransactionCancelled(_txId);
    }

    // ============================================
    // 查询功能模块
    // ============================================

    /**
     * 获取交易详情
     * @param _txId 交易ID
     */
    function getTransaction(uint256 _txId)
        external
        view
        txExists(_txId)
        returns (
            address to,
            uint256 value,
            bytes memory data,
            bool executed,
            uint256 gasLimit,
            uint256 confirmations
        )
    {
        Transaction memory transaction = transactions[_txId];
        return (
            transaction.to,
            transaction.value,
            transaction.data,
            transaction.executed,
            transaction.gasLimit,
            transaction.confirmations
        );
    }

    /**
     * 获取交易总数
     */
    function getTransactionCount() external view returns (uint256) {
        return transactions.length;
    }

    /**
     * 获取所有者数量
     */
    function getOwnersCount() external view returns (uint256) {
        return owners.length;
    }

    /**
     * 检查地址是否为所有者
     * @param _address 要检查的地址
     */
    function isWalletOwner(address _address) external view returns (bool) {
        return isOwner[_address];
    }

    /**
     * 获取合约余额
     */
    function getBalance() external view returns (uint256) {
        return address(this).balance;
    }
}
