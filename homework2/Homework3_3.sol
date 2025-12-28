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

    // 用户注册
    function registerUser(string calldata name, string calldata email) public returns(uint256) {
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

    // 设置权限
    function setAdminRole(uint256 userId) public onlyOwner {
        require(userId <= userCount, "user not found");
        require(users[userId].role != Role.Admin, "user already has admin role");
        User storage user = users[userId];
        user.role = Role.Admin;
        roles[user.addr] = Role.Admin;
    }

    // 存款
    function deposit(uint256 userId) public payable {
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