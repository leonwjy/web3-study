// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Homework2 {
    uint[] public data;

    // 这个函数有很多优化空间
    function processData(uint[] calldata input) public {
        data = new uint[](input.length);
        for(uint i = 0; i < input.length; i++) {
            data[i] = input[i] * 2;
        }
    }

    // 这个函数也可以优化
    function getSum() public view returns (uint) {
        uint sum = 0;
        uint len = data.length;
        for(uint i = 0; i < len; i++) {
            sum += data[i];
        }
        return sum;
    }
}