// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./Ownable.sol";

contract Pausable is Ownable {

    bool public paused;

    event Paused(address indexed account);
    event Unpaused(address indexed account);

    constructor() {
        paused = false;
    }
    
    modifier whenNotPaused() {
        require(!paused, "Pausable is paused");
        _;
    }

    modifier whenPaused() {
        require(paused, "Pausable is not paused");
        _;
    }
    
    function pause() public onlyOwner whenNotPaused {
        paused = true;
        emit Paused(msg.sender);
    }

    function unpause() public onlyOwner whenPaused {
        paused = false;
        emit Unpaused(msg.sender);
    }
}