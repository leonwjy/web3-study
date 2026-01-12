// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;


// 众筹合约
contract Crowdfunding { 

    // 众筹成功
    event CrowdSuccess(address indexed contributor);

    // 退款成功
    event RefundSuccess(address indexed contributor);

    // 退款失败
    error TransferFailed(address contributor, uint256 amount);

    // 众筹失败
    event CrowdFailed();

    // 众筹状态
    enum CrowdStatus {
        Pending,    // 待开始
        Running,    // 进行中
        Success,    // 成功
        Failed      // 失败
    }

    // 众筹信息
    struct CrowdInfo {
        string name;       // 名称
        string description; // 描述
        uint256 targetAmount; // 目标金额
        uint256 currentAmount; // 当前金额
        uint256 startTime; // 开始时间
        uint256 deadline; // 截止时间
    }

    // 众筹信息
    CrowdInfo public crowdInfo;

    // 众筹状态
    CrowdStatus public status;

    // 是否暂停
    bool public paused = false;

    // 众筹参与者
    mapping(address => uint256) public contributions;

    // 众筹参与者地址
    address[] public userAddresses;

    // 众筹参与者数量
    uint256 public userCount;

    // 所有者
    address public owner;

    // 最小支持金额
    uint256 public constant MIN_SUPPORT_AMOUNT = 0.01 ether;

    // 批量退款大小
    uint256 private constant BATCH_SIZE = 10;

    // 一天的秒数
    uint256 private constant SECONDS_PER_DAY = 24 * 60 * 60;

    // 构造器
    constructor(string memory _name, string memory _description, uint256 _targetAmount, uint256 _durationDays) {
        crowdInfo = CrowdInfo({
            name: _name,
            description: _description,
            targetAmount: _targetAmount,
            currentAmount: 0,
            startTime: block.timestamp,
            deadline: block.timestamp + _durationDays * SECONDS_PER_DAY
        });
        owner = msg.sender;
        status = CrowdStatus.Pending;
    }

    // modifier
    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }

    modifier whenNotPaused() {
        require(!paused, "Crowdfunding is paused");
        _;
    }

    modifier whenPaused() {
        require(paused, "Crowdfunding is not paused");
        _;
    }

    modifier whenRunning() {
        require(status == CrowdStatus.Running, "Crowdfunding is not running");
        _;
    }

    modifier whenSuccess() {
        require(status == CrowdStatus.Success, "Crowdfunding is not successful");
        _;
    }
    
    modifier whenFailed() {
        require(status == CrowdStatus.Failed, "Crowdfunding is failed");
        _;
    }

    modifier whenPending() {
        require(status == CrowdStatus.Pending, "Crowdfunding is not pending");
        _;
    }

    // 开始众筹
    function startCrowd() public onlyOwner whenPending {
        status = CrowdStatus.Running;
    }

    // 暂停众筹
    function pauseCrowd() public onlyOwner whenRunning {
        paused = true;
    }

    // 恢复众筹
    function resumeCrowd() public onlyOwner whenPaused {
        paused = false;
    }

    // 内部方法状态转化为 string
    function _statusToString(CrowdStatus _status) internal pure returns(string memory) {
        string[4] memory statusStrings = ["Pending", "Running", "Success", "Failed"];
        return statusStrings[uint256(_status)];
    }

    // 内部方法，快速删除删除捐赠者信息退款
    function _removeContributor(address contributor) internal {
        for (uint i = 0; i < userCount; i++) {
            if (userAddresses[i] == contributor) {
                userAddresses[i] = userAddresses[userCount - 1];
                userCount--;
                break;
            }
        }
    }

    // 内部方法，分页退款捐款者
    function _refundContributorsBatch(uint256 startIndex, uint256 endIndex, bool updateTotal) internal {
        require(endIndex <= userCount, "End index out of bounds");

        for (uint i = startIndex; i < endIndex; i++) {
            address contributor = userAddresses[i];
            if (contributions[contributor] > 0) {
                _processRefund(contributor, updateTotal);
            }
        }
    }

    // 内部方法，清理已退款的捐款者地址
    function _cleanupRefundedContributors(uint256 startIndex, uint256 endIndex) internal {
        // 从后往前清理，避免索引变化问题
        for (uint i = endIndex; i > startIndex; i--) {
            address contributor = userAddresses[i - 1];
            if (contributions[contributor] == 0) {
                _removeContributor(contributor);
            }
        }
    }

    // 内部方法，执行单个捐款者的退款
    function _processRefund(address contributor, bool updateTotal) internal {
        uint256 amount = contributions[contributor];
        require(amount > 0, "Contributor has no funds to refund");

        // 清空捐款者捐款
        contributions[contributor] = 0;

        // 退款，使用 call 模式退款
        (bool success, ) = payable(contributor).call{value: amount}("");
        // 判断退款是否成功
        if (!success) {
            revert TransferFailed(contributor, amount);
        }

        // 通知退款成功
        emit RefundSuccess(contributor);

        // 条件更新众筹总金额（owner可能不希望更新）
        if (updateTotal) {
            crowdInfo.currentAmount -= amount;
        }
    }

    // 内部方法，判断众筹是否还在进行中（未结束）
    function _isActive() internal view returns(bool) {
        return block.timestamp < crowdInfo.deadline;
    }

    // 内部方法，判断是否是众筹失败
    function _isFailed() internal view returns(bool) {
        return !_isActive() && crowdInfo.currentAmount < crowdInfo.targetAmount;
    }

    // 捐款
    function contribute() public payable whenRunning {
        require(msg.value >= MIN_SUPPORT_AMOUNT, "Amount is less than the minimum support amount");
        require(_isActive(), "Crowdfunding has ended");

        // 处理捐款者地址（只在第一次捐款时添加）
        if (contributions[msg.sender] == 0) {
            userAddresses.push(msg.sender);
            userCount++;
        }

        contributions[msg.sender] += msg.value;
        crowdInfo.currentAmount += msg.value;

        if (crowdInfo.currentAmount >= crowdInfo.targetAmount) {
            status = CrowdStatus.Success;
            // 通知众筹成功
            emit CrowdSuccess(msg.sender);
        }
    }

    // 获取众筹信息
    function getCrowdInfo() public view returns(CrowdInfo memory, string memory statusString) {
        return (crowdInfo, _statusToString(status));
    }

    // 获取我的捐款信息
    function getMyContributions() public view returns(uint256) {
        return contributions[msg.sender];
    }

    // 获取众筹参与者信息
    function getContributors() public view onlyOwner returns(address[] memory, uint256[] memory) {
        address[] memory contributors = new address[](userCount);
        uint256[] memory contributionsAmount = new uint256[](userCount);
        for (uint i = 0; i < userCount; i++) {
            contributors[i] = address(userAddresses[i]);
            contributionsAmount[i] = contributions[userAddresses[i]];
        }
        return (contributors, contributionsAmount);
    }

    // 捐款者退款
    function refundContributor() public payable whenFailed whenRunning {
        // 执行退款（更新总金额）
        _processRefund(msg.sender, true);
        // 清理捐款者地址
        _removeContributor(msg.sender);
    }

    // 众筹失败，全部退款
    function refundAllContributors() public onlyOwner whenFailed whenRunning payable {
        // 判断是否是众筹失败，众筹时间众筹金额
        require(_isFailed(), "Crowdfunding is not failed");

        // 分页处理退款
        uint256 totalContributors = userCount;
        for (uint256 startIndex = 0; startIndex < totalContributors; startIndex += BATCH_SIZE) {
            uint256 endIndex = startIndex + BATCH_SIZE;
            if (endIndex > totalContributors) {
                endIndex = totalContributors;
            }
            // 先退款（不更新总金额，因为是owner操作）
            _refundContributorsBatch(startIndex, endIndex, false);
            // 再清理已退款的捐款者
            _cleanupRefundedContributors(startIndex, endIndex);
        }

        // 修改众筹状态
        status = CrowdStatus.Failed;
        // 通知众筹失败
        emit CrowdFailed();
    }

}