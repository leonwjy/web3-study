// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract VoteSystem {

    struct Proposal {
        string content;                     //提案内容
        uint voteCount;                     //投票数量
        uint endTime;                       //结束时间
        mapping(address => bool) voters;    //投票人信息
    }

    uint public proposalCount;              //提案总数
    mapping(uint => Proposal) idToProposal; //提案映射

    // 可以定义一个变量，用于维护获胜者 id，每次投票的时候进行更新

    // 创建提案
    function createProposal(string calldata content, uint time) public returns(uint) {
        require(bytes(content).length > 0, "content required");
        require(time > 0, "time required");

        uint pId = proposalCount++;

        Proposal storage p = idToProposal[pId];
        p.content = content;
        p.voteCount = 0;
        p.endTime = block.timestamp + time;

        return pId;
    }

    // 投票
    function vote(uint pId) public {
        require(pId < proposalCount, "pid is not exists");
        

        Proposal storage p = idToProposal[pId];

        require(block.timestamp < p.endTime, "vote is over");
        require(!p.voters[msg.sender], "already vote");

        p.voters[msg.sender] = true;
        p.voteCount++;

    }

    // 查询提案信息
    function getProposal(uint pId) public view returns(string memory content, uint endTime, uint voteCount) {
        require(pId < proposalCount, "pid is not exists");

        Proposal storage p = idToProposal[pId];

        return (p.content, p.endTime, p.voteCount);
    }

    // 获取获胜提案
    function getWinProposal() public view returns(uint pId) {
        uint maxVotes = 0;
        for(uint i = 0; i < proposalCount; i ++) {
            Proposal storage p = idToProposal[i];
            if(p.voteCount > maxVotes) {
                maxVotes = p.voteCount;
                pId = i;
            }
        }

        return pId;
    }

}