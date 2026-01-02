// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

contract MyToken {

    // error
    error InsufficientBalance(address account, uint256 available, uint256 required);
    error InvalidRecipient(address recipient);
    error InvalidAmount(uint256 amount);

    // event
    event Transfer(address indexed from, address indexed to, uint256 value);

    mapping(address => uint256) public balanceOf;

    constructor() {
        balanceOf[msg.sender] = 1000;
    }

    function transfer(address to, uint256 amount) external returns (bool) {
        if (to == address(0)) revert InvalidRecipient(to);
        if (amount <= 0) revert InvalidAmount(amount);

        // 转账前金额
        uint256 beforeFormBalance = balanceOf[msg.sender];
        uint256 beforeToBalance = balanceOf[to];

        // 处理交易
        balanceOf[msg.sender] -= amount;
        balanceOf[to] += amount;

        // 交易后验证
        assert(
            balanceOf[msg.sender] + balanceOf[to] == beforeFormBalance + beforeToBalance
        );

        // 事件
        emit Transfer(msg.sender, to, amount);

        return true;
    }
}


contract AuctionContract {
    // 自定义错误
    error BidTooLow(uint256 currentBid, uint256 newBid);
    error LowMinBid(uint256 amount, uint256 minBid);
    error TimeEnded();
    error TimeNotEnded();
    error NotHighestBidder();
    error WithdrawalFailed();
    error AlreadyWithdrawn();
    
    address public owner;
    uint256 public timeEnd;
    uint256 public highestBid;
    address public highestBidder;
    
    mapping(address => uint256) public pendingReturns; // 记录待退款
    mapping(address => bool) public hasWithdrawn; //已拍人

    uint256 public constant MIN_BID =  0.1 ether;  // 最小加价
    
    event NewBid(address indexed bidder, uint256 amount);
    event TimeIsEnded(address winner, uint256 amount);
    event Withdrawal(address indexed bidder, uint256 amount);

    modifier onlyOwner(){
        require(msg.sender == owner, unicode"只有所有者可以结束拍卖");
        _;
    }
    
    constructor(uint256 _duration) {
        owner = msg.sender;
        timeEnd = block.timestamp + _duration;
    }
    
    /**
     * 出价函数
     */
    function bid() public payable returns(bool) {
        // 检查拍卖是否还在进行
        if (block.timestamp >= timeEnd) {
            revert TimeEnded();
        }

        // 小于最小加价
        if(msg.value < highestBid + MIN_BID) {
            revert LowMinBid(msg.value, MIN_BID);
        }
        
        // 检查出价是否高于当前最高价
        if (msg.value <= highestBid) {
            revert BidTooLow(highestBid, msg.value);
        }
        
        // 如果有之前的最高出价者,记录待退款
        if (highestBidder != address(0)) {
            pendingReturns[highestBidder] += highestBid;
        }
        
        // 更新最高出价
        highestBidder = msg.sender;
        highestBid = msg.value;
        
        emit NewBid(msg.sender, msg.value);

        return true;
    }
    
    /**
     * 提取未中标的出价
     */
    function withdraw() public returns (bool) {
        // 检查是否有待退款
        uint256 amount = pendingReturns[msg.sender];
        require(amount > 0, unicode"没有待退款");
        
        // 检查是否已经提取过
        if (hasWithdrawn[msg.sender]) {
            revert AlreadyWithdrawn();
        }
        
        // 先更新状态
        pendingReturns[msg.sender] = 0;
        hasWithdrawn[msg.sender] = true;
        
        // 处理转账
        (bool success, ) = msg.sender.call{value: amount}("");
        if (!success) {
            // 如果转账失败,恢复状态
            pendingReturns[msg.sender] = amount;
            hasWithdrawn[msg.sender] = false;
            revert WithdrawalFailed();
        }
        
        emit Withdrawal(msg.sender, amount);
        return true;
    }
    
    /**
     * 结束拍卖（仅所有者）
     */
    function endAuction() public onlyOwner {
        
        if (block.timestamp < timeEnd) {
            revert TimeNotEnded();
        }
        
        emit TimeIsEnded(highestBidder, highestBid);
        
        // 转账给所有者
        (bool success, ) = owner.call{value: highestBid}("");
        require(success, unicode"转账失败");
    }
}