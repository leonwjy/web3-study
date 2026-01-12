// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "./Ownable.sol";
import "./Pausable.sol";

contract MyContract is Ownable, Pausable {

    uint256 public value;

    constructor() {
        value = 0;
    }

    function setValue(uint256 _value) public onlyOwner whenNotPaused {
        value = _value;
    }

    function getValue() public whenNotPaused view returns(uint256) {
        return value;
    }
}