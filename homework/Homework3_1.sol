// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Homework3_1 {

    uint[] public data; 

    uint public constant MAX_LENGTH = 100;

    // 安全添加
    function safePush(uint num) public {
        require(data.length < MAX_LENGTH, "Array is full");
        data.push(num);
    }

    // 快速删除
    function quickRemove(uint index) public {
        require(index < data.length, "Index out of bounds");

        data[index] = data[data.length - 1];

        data.pop();
    }

    // 保持顺序
    function removeWithOrder(uint index) public {
        require(index < data.length, "Index out of bounds");

        // 需要把数组从删除位置逐个向前移动一个
        uint len = data.length;
        for (uint i = index; i < len - 1; i ++) {
            data[i] = data[i + 1];
        }

        // 删除最后一个元素
        data.pop();
    }

    // 分批求和
    function sumRange(uint start, uint end) public view returns (uint) {
        require(start < end, "Invalid range");
        require(end < data.length, "End exceeds array length");

        uint total = 0;
        for(uint i = start; i < end; i++ ) {
            total += data[i];
        }

        return total;
    } 

    // 查找元素
    function selectIndex(uint index) public view returns (uint) {
        require(index < data.length, "Index out of bounds");

        return data[index];
    }

    // 获取所有元素
    function getArray() public view returns(uint[] memory) {
        require(data.length <= 100, "Array too large, use getArrayPaged instead");
        return data;
    }


    // 示例：优化下面方法的 ges
    function process(uint[] memory values) public {
        for(uint i= 0;i < values.length; i++){
            if(values[i]>10){
                data.push(values[i]);
            }
        }
    }

    // 优化版本
    function processOptimized (uint[] calldata values) public {
        for (uint i = 0; i < values.length; i++ ) {
            if(values[i] > 10) {
                data.push(values[i]);
            }
        }
    }

}