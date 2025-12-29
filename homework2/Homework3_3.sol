// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract RoleMananger {

    enum Role {
        Owner,
        Admin,
        User
    }

    struct User {
        string name;
        string email;
        uint256 balance;
        address addr;
        Role role;
    }

    mapping(address => Role) public roles;
    mapping(uint256 => User) public users;
    
    uint256 public userCount;
    uint256 public constant MAX_USERS = 1000;

    constructor() { 
        roles[msg.sender] = Role.Owner;
        userCount++;
        users[userCount] = User({
            name: "Owner",
            email: "owner@example.com",
            balance: 0 ether,
            addr: msg.sender,
            role: Role.Owner
        });
    }

    modifier onlyOwner(){
        require(roles[msg.sender] == Role.Owner, "role must be owner");
        _;
    }

    modifier onlyAdmin(){
        require(roles[msg.sender] == Role.Admin, "role must be admin");
        _;
    }

    modifier onlyUser(){
        require(roles[msg.sender] == Role.User, "role must be user");
        _;
    }

    // 管理员注册用户
    function registerUser(string calldata name, string calldata email) public onlyAdmin returns(uint256) {
        require(bytes(name).length > 0, "name required");
        require(userCount < MAX_USERS, "max users reached");
        
        uint256 userId = userCount++;
        users[userId] = User({
            name: name,
            email: email,
            balance: 0 ether,
            addr: msg.sender,
            role: Role.User
        });
        roles[msg.sender] = Role.User;
        return userId;
    }

    // 创建人注册管理员
    function registerAdmin(string calldata name, string calldata email) public onlyOwner returns(uint256) {
        require(bytes(name).length > 0, "name required");
        require(userCount < MAX_USERS, "max users reached");
        uint256 userId = userCount++;
        users[userId] = User({
            name: name,
            email: email,
            balance: 0 ether,
            addr: msg.sender,
            role: Role.Admin
        });
        roles[msg.sender] = Role.Admin;
        return userId;
    }

    // 存款
    function deposit(uint256 userId) public onlyOwner payable {
        require(userId <= userCount, "user not found");
        require(msg.value > 0 ether, "amount must be greater than 0");
        User storage user = users[userId];
        user.balance += msg.value;
    }

    // 查询信息
    function getUserInfo(uint256 userId) public view returns(string memory name, string memory email) { 
        require(userId <= userCount, "user not found");
        User storage user = users[userId];
        require(user.role == Role.User, "user is not a user");
        return (user.name, user.email);
    }

}