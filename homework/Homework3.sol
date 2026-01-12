// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

error VoteClosed();
error AlreadyVoted();
error NotOwner();

contract Task {
    bool public canVote;
    uint8 public voteCount;
    address public immutable owner;  // immutable 节省读取 gas
    mapping(address => bool) public voteMap;

    constructor() {
        canVote = true;
        // voteCount 默认就是 0，不需要显式赋值
        owner = msg.sender;
    }

    function vote() public {
        if (!canVote) revert VoteClosed();
        if (voteMap[msg.sender]) revert AlreadyVoted();

        voteMap[msg.sender] = true;
        unchecked {
            voteCount++;  // 如果确定不会超过 255
        }
    }

    // 加上权限控制！
    function toggleVote() public {
        if (msg.sender != owner) revert NotOwner();
        canVote = !canVote;
    }

    // 删除 getOwner() 和 getOwnerVoted()，public 变量已有自动 getter
}