// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract UserManager {

    struct User {
        string name;
        string email;
        uint balance;
        uint registeredAt;
        bool exists;
    }

    mapping(address => User) public users;
    address[] public userAddresses;

    uint public constant MAX_USERS = 1000;

    // 注册用户
    function registerUser(string calldata name, string calldata email) public {
        require(!users[msg.sender].exists, unicode"用户已注册");
        require(userAddresses.length < MAX_USERS, unicode"超过最大人数");
        require(bytes(name).length > 0, unicode"名字不能为空");
        users[msg.sender] = User({
            name: name,
            email: email,
            balance: 0,
            registeredAt: block.timestamp,
            exists: true
        });
        userAddresses.push(msg.sender);
    }

    // 更新用户信息
    function updateProfile(string calldata name, string calldata email) public {
        User storage user = users[msg.sender];
        require(user.exists, unicode"用户未注册");
        user.name = name;
        user.email = email;
    }

    // 存款功能
    function deposit(uint amount) public {
        User storage user = users[msg.sender];
        require(user.exists, unicode"用户未注册");
        require(amount > 0, unicode"存款金额不能为0");
        user.balance += amount;
    }

    // 查询用户信息
    function getUserInfo(address add) public view returns(User memory) {
        return users[add];
    }

    // 获取所有用户
    function getAllUsers() public view returns(User[] memory) {
        uint len = userAddresses.length;
        User[] memory temp_users = new User[](len);
        for (uint i = 0; i < len; i++) {
            temp_users[i] = users[userAddresses[i]];
        }
        return temp_users;
    }

    // 分批获取用户
    function getUsersByRange(uint start, uint end) public view returns(User[] memory) {
        require(start < end, unicode"无效范围");
        require(end < userAddresses.length, unicode"超出范围");

        User[] memory temp_users = new User[](end - start);
        for(uint i = start; i < end; i ++) {
            temp_users[i - start] = users[userAddresses[i]];
        }
        return temp_users;
    }
}