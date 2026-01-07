// SPDX-License-Identifier: MIT
pragma solidity ^0.8.28;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/utils/ReentrancyGuardUpgradeable.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721Receiver.sol";
import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "./PriceOracle.sol";

/**
 * @title Auction
 * @dev NFT 拍卖合约核心实现，支持 ETH 和 ERC20 出价
 */
contract Auction is Initializable, OwnableUpgradeable, ReentrancyGuardUpgradeable, IERC721Receiver {
    // 手续费比例（基点，250 = 2.5%）
    uint256 public constant FEE_BASIS_POINTS = 250; // 2.5%
    uint256 public constant BASIS_POINTS = 10000;
    
    // 价格预言机合约
    PriceOracle public priceOracle;
    
    // 拍卖结构体
    struct AuctionInfo {
        address nftContract;      // NFT 合约地址
        uint256 tokenId;          // NFT tokenId
        address seller;            // 卖家地址
        uint256 startingPrice;    // 起拍价（USD，18 位小数）
        uint256 currentHighestBid; // 当前最高出价（USD，18 位小数）
        address highestBidder;   // 当前最高出价者
        uint256 endTime;          // 拍卖结束时间
        address paymentToken;     // 支付代币地址（address(0) 表示 ETH）
        bool ended;               // 是否已结束
    }
    
    // 拍卖 ID 到拍卖信息的映射
    mapping(uint256 => AuctionInfo) public auctions;
    
    // 拍卖 ID 计数器
    uint256 internal _auctionIdCounter;
    
    // 用户出价记录：auctionId => bidder => bidAmount (USD, 18 decimals)
    mapping(uint256 => mapping(address => uint256)) public bids;
    
    // 用户可提取的出价金额：auctionId => bidder => withdrawableAmount (原始代币单位)
    mapping(uint256 => mapping(address => uint256)) public withdrawableBids;
    
    // 事件
    event AuctionCreated(
        uint256 indexed auctionId,
        address indexed nftContract,
        uint256 indexed tokenId,
        address seller,
        uint256 startingPrice,
        uint256 endTime,
        address paymentToken
    );
    
    event BidPlaced(
        uint256 indexed auctionId,
        address indexed bidder,
        uint256 bidAmountUSD,
        uint256 originalAmount,
        address paymentToken
    );
    
    event AuctionEnded(
        uint256 indexed auctionId,
        address indexed winner,
        uint256 finalPriceUSD,
        address paymentToken
    );
    
    event BidWithdrawn(
        uint256 indexed auctionId,
        address indexed bidder,
        uint256 amount
    );
    
    /**
     * @dev 初始化函数（用于可升级合约）
     * @param _priceOracle 价格预言机合约地址
     */
    function initialize(address _priceOracle) public virtual initializer {
        __Ownable_init(msg.sender);
        __ReentrancyGuard_init();
        require(_priceOracle != address(0), "Invalid price oracle address");
        priceOracle = PriceOracle(_priceOracle);
    }
    
    /**
     * @dev 创建拍卖
     * @param nftContract NFT 合约地址
     * @param tokenId NFT tokenId
     * @param startingPriceUSD 起拍价（USD，18 位小数）
     * @param duration 拍卖持续时间（秒）
     * @param paymentToken 支付代币地址（address(0) 表示 ETH）
     */
    function createAuction(
        address nftContract,
        uint256 tokenId,
        uint256 startingPriceUSD,
        uint256 duration,
        address paymentToken
    ) public virtual nonReentrant {
        require(nftContract != address(0), "Invalid NFT contract");
        require(duration > 0, "Duration must be greater than 0");
        require(startingPriceUSD > 0, "Starting price must be greater than 0");
        
        IERC721 nft = IERC721(nftContract);
        require(nft.ownerOf(tokenId) == msg.sender, "Not the owner of NFT");
        require(nft.getApproved(tokenId) == address(this) || 
                nft.isApprovedForAll(msg.sender, address(this)), 
                "NFT not approved");
        
        // 转移 NFT 到合约
        nft.safeTransferFrom(msg.sender, address(this), tokenId);
        
        _auctionIdCounter++;
        uint256 auctionId = _auctionIdCounter;
        
        auctions[auctionId] = AuctionInfo({
            nftContract: nftContract,
            tokenId: tokenId,
            seller: msg.sender,
            startingPrice: startingPriceUSD,
            currentHighestBid: 0,
            highestBidder: address(0),
            endTime: block.timestamp + duration,
            paymentToken: paymentToken,
            ended: false
        });
        
        emit AuctionCreated(
            auctionId,
            nftContract,
            tokenId,
            msg.sender,
            startingPriceUSD,
            block.timestamp + duration,
            paymentToken
        );
    }
    
    /**
     * @dev 使用 ETH 出价
     * @param auctionId 拍卖 ID
     */
    function bidWithETH(uint256 auctionId) external payable nonReentrant {
        AuctionInfo storage auction = auctions[auctionId];
        require(!auction.ended, "Auction already ended");
        require(block.timestamp < auction.endTime, "Auction has ended");
        require(auction.paymentToken == address(0), "This auction does not accept ETH");
        require(msg.value > 0, "Bid amount must be greater than 0");
        
        // 将 ETH 出价转换为 USD
        uint256 bidAmountUSD = priceOracle.convertToUSD(msg.value, address(0));

        // 如果已有出价，新出价必须高于当前最高出价
        // 如果没有出价，新出价必须不低于起拍价
        if (auction.currentHighestBid > 0) {
            require(bidAmountUSD > auction.currentHighestBid, "Bid must be higher than current highest bid");

            // 退还前一个最高出价者的出价
            address previousBidder = auction.highestBidder;
            uint256 previousBidAmount = withdrawableBids[auctionId][previousBidder];
            if (previousBidAmount > 0) {
                withdrawableBids[auctionId][previousBidder] = previousBidAmount;
            }
        } else {
            require(bidAmountUSD >= auction.startingPrice, "Bid below starting price");
        }
        
        // 记录当前出价者的出价
        bids[auctionId][msg.sender] = bidAmountUSD;
        withdrawableBids[auctionId][msg.sender] = msg.value;
        
        // 更新最高出价
        auction.currentHighestBid = bidAmountUSD;
        auction.highestBidder = msg.sender;
        
        emit BidPlaced(auctionId, msg.sender, bidAmountUSD, msg.value, address(0));
    }
    
    /**
     * @dev 使用 ERC20 代币出价
     * @param auctionId 拍卖 ID
     * @param amount 出价金额（代币最小单位）
     */
    function bidWithERC20(uint256 auctionId, uint256 amount) external nonReentrant {
        AuctionInfo storage auction = auctions[auctionId];
        require(!auction.ended, "Auction already ended");
        require(block.timestamp < auction.endTime, "Auction has ended");
        require(auction.paymentToken != address(0), "This auction does not accept ERC20");
        require(amount > 0, "Bid amount must be greater than 0");
        
        IERC20 token = IERC20(auction.paymentToken);
        require(token.balanceOf(msg.sender) >= amount, "Insufficient token balance");
        require(token.allowance(msg.sender, address(this)) >= amount, "Token allowance insufficient");
        
        // 转移代币到合约
        token.transferFrom(msg.sender, address(this), amount);
        
        // 将代币出价转换为 USD
        uint256 bidAmountUSD = priceOracle.convertToUSD(amount, auction.paymentToken);

        // 如果已有出价，新出价必须高于当前最高出价
        // 如果没有出价，新出价必须不低于起拍价
        if (auction.currentHighestBid > 0) {
            require(bidAmountUSD > auction.currentHighestBid, "Bid must be higher than current highest bid");

            // 退还前一个最高出价者的出价
            address previousBidder = auction.highestBidder;
            uint256 previousBidAmount = withdrawableBids[auctionId][previousBidder];
            if (previousBidAmount > 0) {
                withdrawableBids[auctionId][previousBidder] = previousBidAmount;
            }
        } else {
            require(bidAmountUSD >= auction.startingPrice, "Bid below starting price");
        }
        
        // 记录当前出价者的出价
        bids[auctionId][msg.sender] = bidAmountUSD;
        withdrawableBids[auctionId][msg.sender] = amount;
        
        // 更新最高出价
        auction.currentHighestBid = bidAmountUSD;
        auction.highestBidder = msg.sender;
        
        emit BidPlaced(auctionId, msg.sender, bidAmountUSD, amount, auction.paymentToken);
    }
    
    /**
     * @dev 结束拍卖
     * @param auctionId 拍卖 ID
     */
    function endAuction(uint256 auctionId) external nonReentrant {
        AuctionInfo storage auction = auctions[auctionId];
        require(!auction.ended, "Auction already ended");
        require(block.timestamp >= auction.endTime || msg.sender == auction.seller, 
                "Auction not ended and not seller");
        
        auction.ended = true;
        
        IERC721 nft = IERC721(auction.nftContract);
        
        if (auction.highestBidder != address(0)) {
            // 有出价者，转移 NFT 给最高出价者
            nft.safeTransferFrom(address(this), auction.highestBidder, auction.tokenId);
            
            // 计算手续费和卖家收益
            uint256 totalAmount = withdrawableBids[auctionId][auction.highestBidder];
            uint256 fee = (totalAmount * FEE_BASIS_POINTS) / BASIS_POINTS;
            uint256 sellerAmount = totalAmount - fee;
            
            // 转移资金
            if (auction.paymentToken == address(0)) {
                // ETH
                (bool sellerSuccess,) = payable(auction.seller).call{value: sellerAmount}("");
                require(sellerSuccess, "ETH transfer to seller failed");
                (bool ownerSuccess,) = payable(owner()).call{value: fee}("");
                require(ownerSuccess, "ETH transfer to owner failed");
            } else {
                // ERC20
                IERC20 token = IERC20(auction.paymentToken);
                token.transfer(auction.seller, sellerAmount);
                token.transfer(owner(), fee);
            }
            
            emit AuctionEnded(auctionId, auction.highestBidder, auction.currentHighestBid, auction.paymentToken);
        } else {
            // 没有出价者，退还 NFT 给卖家
            nft.safeTransferFrom(address(this), auction.seller, auction.tokenId);
            emit AuctionEnded(auctionId, address(0), 0, auction.paymentToken);
        }
    }
    
    /**
     * @dev 撤回出价（当出价被超过时）
     * @param auctionId 拍卖 ID
     */
    function withdrawBid(uint256 auctionId) external nonReentrant {
        AuctionInfo storage auction = auctions[auctionId];
        require(!auction.ended, "Auction already ended");
        require(msg.sender != auction.highestBidder, "Cannot withdraw current highest bid");
        
        uint256 amount = withdrawableBids[auctionId][msg.sender];
        require(amount > 0, "No bid to withdraw");
        
        withdrawableBids[auctionId][msg.sender] = 0;
        
        if (auction.paymentToken == address(0)) {
            // ETH
            (bool success,) = payable(msg.sender).call{value: amount}("");
            require(success, "ETH transfer failed");
        } else {
            // ERC20
            IERC20 token = IERC20(auction.paymentToken);
            token.transfer(msg.sender, amount);
        }
        
        emit BidWithdrawn(auctionId, msg.sender, amount);
    }
    
    /**
     * @dev 获取拍卖信息
     * @param auctionId 拍卖 ID
     * @return 拍卖信息结构体
     */
    function getAuction(uint256 auctionId) external view returns (AuctionInfo memory) {
        return auctions[auctionId];
    }
    
    /**
     * @dev 获取当前拍卖 ID 计数器
     * @return 当前拍卖 ID
     */
    function getCurrentAuctionId() external view returns (uint256) {
        return _auctionIdCounter;
    }

    /**
     * @dev 实现 IERC721Receiver 接口，支持接收 NFT
     * @return bytes4 返回 IERC721Receiver.onERC721Received.selector
     */
    function onERC721Received(
        address /*operator*/,
        address /*from*/,
        uint256 /*tokenId*/,
        bytes calldata /*data*/
    ) external pure override returns (bytes4) {
        return IERC721Receiver.onERC721Received.selector;
    }
}

