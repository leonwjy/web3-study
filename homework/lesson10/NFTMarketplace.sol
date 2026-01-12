// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/IERC721.sol";
import "@openzeppelin/contracts/token/ERC721/IERC721Receiver.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/utils/introspection/IERC165.sol";

/**
 * @title IERC2981
 * @dev ERC2981版税标准接口
 */
interface IERC2981 is IERC165 {
    function royaltyInfo(
        uint256 tokenId,
        uint256 salePrice
    ) external view returns (
        address receiver,
        uint256 royaltyAmount
    );
}

/**
 * @title NFTMarketplace
 * @dev 完整的NFT交易市场合约，支持上架、购买、版税和拍卖功能
 * @notice 使用ReentrancyGuard防止重入攻击
 */
contract NFTMarketplace is ReentrancyGuard {
    
    /**
     * @dev 挂单结构体
     */
    struct Listing {
        address seller;        // 卖家地址
        address nftContract;    // NFT合约地址
        uint256 tokenId;         // Token ID
        uint256 price;          // 售价（wei）
        bool active;            // 是否激活
        bool isPrivate;         // 是否为私密销售
        address[] whitelist;    // 白名单买家地址数组
    }
    
    /**
     * @dev 拍卖结构体
     */
    struct Auction {
        address seller;           // 卖家地址
        address nftContract;      // NFT合约地址
        uint256 tokenId;          // Token ID
        uint256 startPrice;       // 起拍价
        uint256 highestBid;       // 当前最高出价
        address highestBidder;    // 当前最高出价者
        uint256 endTime;          // 拍卖结束时间
        bool active;              // 是否激活
    }

    /**
     * @dev 邀约结构体
     */
    struct Offer {
        address buyer;            // 买家地址
        address nftContract;      // NFT合约地址
        uint256 tokenId;          // Token ID
        uint256 price;            // 出价
        uint256 deadline;         // 截止时间
        bool active;              // 是否激活
    }
    
    // 挂单映射
    mapping(uint256 => Listing) public listings;
    uint256 public listingCounter;
    
    // 拍卖映射
    mapping(uint256 => Auction) public auctions;
    uint256 public auctionCounter;

    // 邀约映射
    mapping(uint256 => Offer) public offers;
    uint256 public offerCounter;

    // 待退款映射（用于拍卖）
    mapping(uint256 => mapping(address => uint256)) public pendingReturns;
    
    // 平台手续费（基点，10000 = 100%）
    uint256 public platformFee = 250; // 2.5%
    
    // 手续费接收地址
    address public feeRecipient;
    
    /**
     * @dev NFT上架事件
     */
    event NFTListed(
        uint256 indexed listingId,
        address indexed seller,
        address indexed nftContract,
        uint256 tokenId,
        uint256 price
    );
    
    /**
     * @dev NFT下架事件
     */
    event NFTDelisted(
        uint256 indexed listingId
    );
    
    /**
     * @dev 价格更新事件
     */
    event PriceUpdated(
        uint256 indexed listingId,
        uint256 newPrice
    );
    
    /**
     * @dev NFT售出事件
     */
    event NFTSold(
        uint256 indexed listingId,
        address indexed buyer,
        address indexed seller,
        uint256 price
    );
    
    /**
     * @dev 拍卖创建事件
     */
    event AuctionCreated(
        uint256 indexed auctionId,
        address indexed seller,
        address indexed nftContract,
        uint256 tokenId,
        uint256 startPrice,
        uint256 endTime
    );
    
    /**
     * @dev 出价事件
     */
    event BidPlaced(
        uint256 indexed auctionId,
        address indexed bidder,
        uint256 amount
    );
    
    /**
     * @dev 拍卖结束事件
     */
    event AuctionEnded(
        uint256 indexed auctionId,
        address indexed winner,
        uint256 finalPrice
    );

    /**
     * @dev 邀约创建事件
     */
    event OfferCreated(
        uint256 indexed offerId,
        address indexed buyer,
        address indexed nftContract,
        uint256 tokenId,
        uint256 price,
        uint256 deadline
    );

    /**
     * @dev 邀约接受事件
     */
    event OfferAccepted(
        uint256 indexed offerId,
        address indexed seller,
        address indexed buyer,
        uint256 price
    );

    /**
     * @dev 邀约拒绝事件
     */
    event OfferRejected(
        uint256 indexed offerId,
        address indexed seller
    );

    /**
     * @dev 邀约取消事件
     */
    event OfferCancelled(
        uint256 indexed offerId,
        address indexed buyer
    );
    
    /**
     * @dev 构造函数
     * @param _feeRecipient 手续费接收地址
     */
    constructor(address _feeRecipient) {
        require(_feeRecipient != address(0), "Invalid fee recipient");
        feeRecipient = _feeRecipient;
    }
    
    /**
     * @dev 上架NFT
     * @param nftContract NFT合约地址
     * @param tokenId Token ID
     * @param price 售价（wei）
     * @return listingId 挂单ID
     */
    function listNFT(
        address nftContract,
        uint256 tokenId,
        uint256 price
    ) external returns (uint256) {
        return _listNFT(nftContract, tokenId, price, false, new address[](0));
    }

    /**
     * @dev 上架NFT（支持私密销售）
     * @param nftContract NFT合约地址
     * @param tokenId Token ID
     * @param price 售价（wei）
     * @param isPrivate 是否为私密销售
     * @param whitelist 白名单买家地址数组
     * @return listingId 挂单ID
     */
    function listNFTWithWhitelist(
        address nftContract,
        uint256 tokenId,
        uint256 price,
        bool isPrivate,
        address[] calldata whitelist
    ) external returns (uint256) {
        return _listNFT(nftContract, tokenId, price, isPrivate, whitelist);
    }

    /**
     * @dev 内部上架NFT函数
     */
    function _listNFT(
        address nftContract,
        uint256 tokenId,
        uint256 price,
        bool isPrivate,
        address[] memory whitelist
    ) internal returns (uint256) {
        require(price > 0, "Price must be greater than 0");
        require(nftContract != address(0), "Invalid NFT contract");

        IERC721 nft = IERC721(nftContract);

        // 验证所有权
        require(nft.ownerOf(tokenId) == msg.sender, "Not the owner");

        // 验证授权
        require(
            nft.getApproved(tokenId) == address(this) ||
            nft.isApprovedForAll(msg.sender, address(this)),
            "Marketplace not approved"
        );

        // 创建挂单
        listingCounter++;
        listings[listingCounter] = Listing({
            seller: msg.sender,
            nftContract: nftContract,
            tokenId: tokenId,
            price: price,
            active: true,
            isPrivate: isPrivate,
            whitelist: whitelist
        });

        emit NFTListed(
            listingCounter,
            msg.sender,
            nftContract,
            tokenId,
            price
        );

        return listingCounter;
    }
    
    /**
     * @dev 下架NFT
     * @param listingId 挂单ID
     */
    function delistNFT(uint256 listingId) external {
        Listing storage listing = listings[listingId];
        
        require(listing.active, "Listing not active");
        require(listing.seller == msg.sender, "Not the seller");
        
        listing.active = false;
        
        emit NFTDelisted(listingId);
    }
    
    /**
     * @dev 更新挂单价格
     * @param listingId 挂单ID
     * @param newPrice 新价格（wei）
     */
    function updatePrice(uint256 listingId, uint256 newPrice) external {
        require(newPrice > 0, "Price must be greater than 0");
        
        Listing storage listing = listings[listingId];
        require(listing.active, "Listing not active");
        require(listing.seller == msg.sender, "Not the seller");
        
        listing.price = newPrice;
        
        emit PriceUpdated(listingId, newPrice);
    }
    
    /**
     * @dev 购买NFT
     * @param listingId 挂单ID
     * @notice 需要支付足够的ETH，多余部分会自动退还
     */
    function buyNFT(uint256 listingId) external payable nonReentrant {
        Listing storage listing = listings[listingId];

        // 检查挂单状态
        require(listing.active, "Listing not active");
        require(msg.value >= listing.price, "Insufficient payment");
        require(msg.sender != listing.seller, "Cannot buy your own NFT");

        // 检查私密销售白名单
        if (listing.isPrivate) {
            require(_isWhitelisted(listing.whitelist, msg.sender), "Not whitelisted for private sale");
        }
        
        // 先更新状态（CEI原则）
        listing.active = false;
        
        // 计算手续费
        uint256 fee = (listing.price * platformFee) / 10000;
        
        // 获取版税信息
        (address royaltyReceiver, uint256 royaltyAmount) = _getRoyaltyInfo(
            listing.nftContract,
            listing.tokenId,
            listing.price
        );
        
        // 计算卖家收益
        uint256 sellerAmount = listing.price - fee - royaltyAmount;
        
        // 转移NFT
        IERC721(listing.nftContract).safeTransferFrom(
            listing.seller,
            msg.sender,
            listing.tokenId
        );
        
        // 资金分配：版税 -> 平台手续费 -> 卖家收益
        if (royaltyAmount > 0 && royaltyReceiver != address(0)) {
            (bool successRoyalty, ) = royaltyReceiver.call{value: royaltyAmount}("");
            require(successRoyalty, "Royalty transfer failed");
        }
        
        (bool successSeller, ) = listing.seller.call{value: sellerAmount}("");
        require(successSeller, "Transfer to seller failed");
        
        (bool successFee, ) = feeRecipient.call{value: fee}("");
        require(successFee, "Transfer fee failed");
        
        // 退还多余资金
        if (msg.value > listing.price) {
            (bool successRefund, ) = msg.sender.call{
                value: msg.value - listing.price
            }("");
            require(successRefund, "Refund failed");
        }
        
        emit NFTSold(listingId, msg.sender, listing.seller, listing.price);
    }
    
    /**
     * @dev 创建拍卖
     * @param nftContract NFT合约地址
     * @param tokenId Token ID
     * @param startPrice 起拍价（wei）
     * @param durationHours 拍卖时长（小时）
     * @return auctionId 拍卖ID
     */
    function createAuction(
        address nftContract,
        uint256 tokenId,
        uint256 startPrice,
        uint256 durationHours
    ) external returns (uint256) {
        require(startPrice > 0, "Start price must be greater than 0");
        require(durationHours >= 1, "Duration must be at least 1 hour");
        require(nftContract != address(0), "Invalid NFT contract");
        
        IERC721 nft = IERC721(nftContract);
        
        // 验证所有权
        require(nft.ownerOf(tokenId) == msg.sender, "Not the owner");
        
        // 验证授权
        require(
            nft.getApproved(tokenId) == address(this) ||
            nft.isApprovedForAll(msg.sender, address(this)),
            "Marketplace not approved"
        );
        
        // 创建拍卖
        auctionCounter++;
        auctions[auctionCounter] = Auction({
            seller: msg.sender,
            nftContract: nftContract,
            tokenId: tokenId,
            startPrice: startPrice,
            highestBid: 0,
            highestBidder: address(0),
            endTime: block.timestamp + (durationHours * 1 hours),
            active: true
        });
        
        emit AuctionCreated(
            auctionCounter,
            msg.sender,
            nftContract,
            tokenId,
            startPrice,
            auctions[auctionCounter].endTime
        );
        
        return auctionCounter;
    }
    
    /**
     * @dev 出价
     * @param auctionId 拍卖ID
     * @notice 需要支付足够的ETH，出价必须高于当前最高出价的5%
     */
    function placeBid(uint256 auctionId) external payable {
        Auction storage auction = auctions[auctionId];
        
        require(auction.active, "Auction not active");
        require(block.timestamp < auction.endTime, "Auction ended");
        require(msg.sender != auction.seller, "Seller cannot bid");
        
        // 计算最低出价
        uint256 minBid;
        if (auction.highestBid == 0) {
            minBid = auction.startPrice;
        } else {
            minBid = auction.highestBid + (auction.highestBid * 5 / 100); // 5% increment
        }
        
        require(msg.value >= minBid, "Bid too low");
        
        // 如果有之前的出价者，记录他们的待退款金额
        if (auction.highestBidder != address(0)) {
            pendingReturns[auctionId][auction.highestBidder] += auction.highestBid;
        }
        
        // 更新最高出价
        auction.highestBid = msg.value;
        auction.highestBidder = msg.sender;
        
        emit BidPlaced(auctionId, msg.sender, msg.value);
    }
    
    /**
     * @dev 提取出价退款
     * @param auctionId 拍卖ID
     * @notice 被超越的出价者可以提取他们的资金
     */
    function withdrawBid(uint256 auctionId) external {
        uint256 amount = pendingReturns[auctionId][msg.sender];
        require(amount > 0, "No pending return");
        
        pendingReturns[auctionId][msg.sender] = 0;
        
        (bool success, ) = msg.sender.call{value: amount}("");
        require(success, "Transfer failed");
    }
    
    /**
     * @dev 结束拍卖
     * @param auctionId 拍卖ID
     * @notice 任何人都可以在拍卖结束后调用此函数进行结算
     */
    function endAuction(uint256 auctionId) external nonReentrant {
        Auction storage auction = auctions[auctionId];
        
        require(auction.active, "Auction not active");
        require(block.timestamp >= auction.endTime, "Auction not ended");
        
        auction.active = false;
        
        if (auction.highestBidder != address(0)) {
            // 有人出价，进行结算
            uint256 fee = (auction.highestBid * platformFee) / 10000;
            
            (address royaltyReceiver, uint256 royaltyAmount) = _getRoyaltyInfo(
                auction.nftContract,
                auction.tokenId,
                auction.highestBid
            );
            
            uint256 sellerAmount = auction.highestBid - fee - royaltyAmount;
            
            // 转移NFT
            IERC721(auction.nftContract).safeTransferFrom(
                auction.seller,
                auction.highestBidder,
                auction.tokenId
            );
            
            // 资金分配
            if (royaltyAmount > 0 && royaltyReceiver != address(0)) {
                (bool successRoyalty, ) = royaltyReceiver.call{value: royaltyAmount}("");
                require(successRoyalty, "Royalty transfer failed");
            }
            
            (bool successSeller, ) = auction.seller.call{value: sellerAmount}("");
            require(successSeller, "Transfer to seller failed");
            
            (bool successFee, ) = feeRecipient.call{value: fee}("");
            require(successFee, "Transfer fee failed");
            
            emit AuctionEnded(
                auctionId,
                auction.highestBidder,
                auction.highestBid
            );
        } else {
            // 没有人出价，拍卖流拍
            emit AuctionEnded(auctionId, address(0), 0);
        }
    }
    
    /**
     * @dev 获取版税信息
     * @param nftContract NFT合约地址
     * @param tokenId Token ID
     * @param salePrice 售价
     * @return receiver 版税接收地址
     * @return royaltyAmount 版税金额
     * @notice 内部函数，检查NFT合约是否支持ERC2981标准
     */
    function _getRoyaltyInfo(
        address nftContract,
        uint256 tokenId,
        uint256 salePrice
    ) internal view returns (address receiver, uint256 royaltyAmount) {
        // 检查NFT合约是否支持ERC2981
        if (IERC165(nftContract).supportsInterface(type(IERC2981).interfaceId)) {
            (receiver, royaltyAmount) = IERC2981(nftContract).royaltyInfo(
                tokenId,
                salePrice
            );
        } else {
            // 不支持版税，返回零地址和零金额
            receiver = address(0);
            royaltyAmount = 0;
        }
    }
    
    /**
     * @dev 查询挂单信息
     * @param listingId 挂单ID
     * @return seller 卖家地址
     * @return nftContract NFT合约地址
     * @return tokenId Token ID
     * @return price 价格
     * @return active 是否激活
     * @return isPrivate 是否为私密销售
     */
    function getListing(uint256 listingId) external view returns (
        address seller,
        address nftContract,
        uint256 tokenId,
        uint256 price,
        bool active,
        bool isPrivate
    ) {
        Listing memory listing = listings[listingId];
        return (
            listing.seller,
            listing.nftContract,
            listing.tokenId,
            listing.price,
            listing.active,
            listing.isPrivate
        );
    }
    
    /**
     * @dev 查询拍卖信息
     * @param auctionId 拍卖ID
     * @return seller 卖家地址
     * @return nftContract NFT合约地址
     * @return tokenId Token ID
     * @return startPrice 起拍价
     * @return highestBid 当前最高出价
     * @return highestBidder 当前最高出价者
     * @return endTime 结束时间
     * @return active 是否激活
     */
    function getAuction(uint256 auctionId) external view returns (
        address seller,
        address nftContract,
        uint256 tokenId,
        uint256 startPrice,
        uint256 highestBid,
        address highestBidder,
        uint256 endTime,
        bool active
    ) {
        Auction memory auction = auctions[auctionId];
        return (
            auction.seller,
            auction.nftContract,
            auction.tokenId,
            auction.startPrice,
            auction.highestBid,
            auction.highestBidder,
            auction.endTime,
            auction.active
        );
    }
    
    /**
     * @dev 设置平台手续费
     * @param newFee 新的手续费（基点）
     * @notice 只有手续费接收地址可以调用
     */
    function setPlatformFee(uint256 newFee) external {
        require(msg.sender == feeRecipient, "Not fee recipient");
        require(newFee <= 1000, "Fee too high"); // 最大10%
        platformFee = newFee;
    }
    
    /**
     * @dev 更新手续费接收地址
     * @param newRecipient 新的接收地址
     * @notice 只有当前手续费接收地址可以调用
     */
    function updateFeeRecipient(address newRecipient) external {
        require(msg.sender == feeRecipient, "Not fee recipient");
        require(newRecipient != address(0), "Invalid address");
        feeRecipient = newRecipient;
    }

    /**
     * @dev 创建邀约
     * @param nftContract NFT合约地址
     * @param tokenId Token ID
     * @param price 出价（wei）
     * @param durationHours 邀约有效期（小时）
     * @return offerId 邀约ID
     * @notice 买家可以对任意NFT出价
     */
    function makeOffer(
        address nftContract,
        uint256 tokenId,
        uint256 price,
        uint256 durationHours
    ) external payable returns (uint256) {
        require(price > 0, "Price must be greater than 0");
        require(durationHours >= 1, "Duration must be at least 1 hour");
        require(nftContract != address(0), "Invalid NFT contract");
        require(msg.value >= price, "Insufficient payment");

        IERC721 nft = IERC721(nftContract);

        // 验证NFT存在
        require(nft.ownerOf(tokenId) != address(0), "NFT does not exist");

        // 买家不能给自己出价
        require(nft.ownerOf(tokenId) != msg.sender, "Cannot make offer to yourself");

        // 创建邀约
        offerCounter++;
        offers[offerCounter] = Offer({
            buyer: msg.sender,
            nftContract: nftContract,
            tokenId: tokenId,
            price: price,
            deadline: block.timestamp + (durationHours * 1 hours),
            active: true
        });

        emit OfferCreated(
            offerCounter,
            msg.sender,
            nftContract,
            tokenId,
            price,
            offers[offerCounter].deadline
        );

        return offerCounter;
    }

    /**
     * @dev 接受邀约
     * @param offerId 邀约ID
     * @notice 只有NFT所有者可以接受邀约
     */
    function acceptOffer(uint256 offerId) external payable nonReentrant {
        Offer storage offer = offers[offerId];

        require(offer.active, "Offer not active");
        require(block.timestamp <= offer.deadline, "Offer expired");

        IERC721 nft = IERC721(offer.nftContract);

        // 验证调用者是NFT所有者
        require(nft.ownerOf(offer.tokenId) == msg.sender, "Not the owner");

        // 验证授权
        require(
            nft.getApproved(offer.tokenId) == address(this) ||
            nft.isApprovedForAll(msg.sender, address(this)),
            "Marketplace not approved"
        );

        // 先更新状态（CEI原则）
        offer.active = false;

        // 计算手续费
        uint256 fee = (offer.price * platformFee) / 10000;

        // 获取版税信息
        (address royaltyReceiver, uint256 royaltyAmount) = _getRoyaltyInfo(
            offer.nftContract,
            offer.tokenId,
            offer.price
        );

        // 计算卖家收益
        uint256 sellerAmount = offer.price - fee - royaltyAmount;

        // 转移NFT
        nft.safeTransferFrom(msg.sender, offer.buyer, offer.tokenId);

        // 资金分配：版税 -> 平台手续费 -> 卖家收益
        if (royaltyAmount > 0 && royaltyReceiver != address(0)) {
            (bool successRoyalty, ) = royaltyReceiver.call{value: royaltyAmount}("");
            require(successRoyalty, "Royalty transfer failed");
        }

        (bool successSeller, ) = msg.sender.call{value: sellerAmount}("");
        require(successSeller, "Transfer to seller failed");

        (bool successFee, ) = feeRecipient.call{value: fee}("");
        require(successFee, "Transfer fee failed");

        // 退还未使用的资金
        if (msg.value > offer.price) {
            uint256 refundAmount = msg.value - offer.price;
            (bool successRefund, ) = offer.buyer.call{value: refundAmount}("");
            require(successRefund, "Refund failed");
        }

        emit OfferAccepted(offerId, msg.sender, offer.buyer, offer.price);
    }

    /**
     * @dev 拒绝邀约
     * @param offerId 邀约ID
     * @notice 只有NFT所有者可以拒绝邀约
     */
    function rejectOffer(uint256 offerId) external {
        Offer storage offer = offers[offerId];

        require(offer.active, "Offer not active");

        IERC721 nft = IERC721(offer.nftContract);

        // 验证调用者是NFT所有者
        require(nft.ownerOf(offer.tokenId) == msg.sender, "Not the owner");

        offer.active = false;

        // 退还买家的资金
        (bool success, ) = offer.buyer.call{value: offer.price}("");
        require(success, "Refund failed");

        emit OfferRejected(offerId, msg.sender);
    }

    /**
     * @dev 取消邀约
     * @param offerId 邀约ID
     * @notice 只有邀约创建者可以取消邀约
     */
    function cancelOffer(uint256 offerId) external {
        Offer storage offer = offers[offerId];

        require(offer.active, "Offer not active");
        require(offer.buyer == msg.sender, "Not the buyer");

        offer.active = false;

        // 退还买家的资金
        (bool success, ) = msg.sender.call{value: offer.price}("");
        require(success, "Refund failed");

        emit OfferCancelled(offerId, msg.sender);
    }

    /**
     * @dev 查询邀约信息
     * @param offerId 邀约ID
     * @return buyer 买家地址
     * @return nftContract NFT合约地址
     * @return tokenId Token ID
     * @return price 价格
     * @return deadline 截止时间
     * @return active 是否激活
     */
    function getOffer(uint256 offerId) external view returns (
        address buyer,
        address nftContract,
        uint256 tokenId,
        uint256 price,
        uint256 deadline,
        bool active
    ) {
        Offer memory offer = offers[offerId];
        return (
            offer.buyer,
            offer.nftContract,
            offer.tokenId,
            offer.price,
            offer.deadline,
            offer.active
        );
    }

    /**
     * @dev 批量上架NFT
     * @param nftContract NFT合约地址
     * @param tokenIds Token ID数组
     * @param prices 价格数组（wei）
     * @return listingIds 挂单ID数组
     * @notice 批量上架多个NFT，tokenIds和prices数组长度必须相等
     */
    function batchListNFTs(
        address nftContract,
        uint256[] calldata tokenIds,
        uint256[] calldata prices
    ) external returns (uint256[] memory) {
        require(tokenIds.length == prices.length, "Arrays length mismatch");
        require(tokenIds.length > 0, "Empty arrays");
        require(tokenIds.length <= 50, "Too many NFTs"); // 限制批量操作数量
        require(nftContract != address(0), "Invalid NFT contract");

        uint256[] memory listingIds = new uint256[](tokenIds.length);
        IERC721 nft = IERC721(nftContract);

        for (uint256 i = 0; i < tokenIds.length; i++) {
            uint256 tokenId = tokenIds[i];
            uint256 price = prices[i];

            require(price > 0, "Price must be greater than 0");

            // 验证所有权
            require(nft.ownerOf(tokenId) == msg.sender, "Not the owner");

            // 验证授权
            require(
                nft.getApproved(tokenId) == address(this) ||
                nft.isApprovedForAll(msg.sender, address(this)),
                "Marketplace not approved"
            );

            // 创建挂单
            listingCounter++;
            listings[listingCounter] = Listing({
                seller: msg.sender,
                nftContract: nftContract,
                tokenId: tokenId,
                price: price,
                active: true,
                isPrivate: false,
                whitelist: new address[](0)
            });

            listingIds[i] = listingCounter;

            emit NFTListed(
                listingCounter,
                msg.sender,
                nftContract,
                tokenId,
                price
            );
        }

        return listingIds;
    }

    /**
     * @dev 批量购买NFT
     * @param listingIds 挂单ID数组
     * @notice 批量购买多个NFT，需要为每个NFT支付足够的ETH
     */
    function batchBuyNFTs(uint256[] calldata listingIds) external payable nonReentrant {
        require(listingIds.length > 0, "Empty array");
        require(listingIds.length <= 20, "Too many purchases"); // 限制批量购买数量

        uint256 totalRequired = 0;
        uint256[] memory prices = new uint256[](listingIds.length);

        // 第一遍：验证所有挂单并计算总价
        for (uint256 i = 0; i < listingIds.length; i++) {
            Listing storage listing = listings[listingIds[i]];

            // 检查挂单状态
            require(listing.active, "Listing not active");
            require(msg.sender != listing.seller, "Cannot buy your own NFT");

            prices[i] = listing.price;
            totalRequired += listing.price;
        }

        require(msg.value >= totalRequired, "Insufficient payment");

        // 第二遍：执行购买
        uint256 remainingValue = msg.value;
        for (uint256 i = 0; i < listingIds.length; i++) {
            uint256 listingId = listingIds[i];
            Listing storage listing = listings[listingId];

            // 再次检查状态（防止重入）
            if (!listing.active || listing.seller == msg.sender) {
                continue; // 跳过已处理的或无效的挂单
            }

            // 先更新状态（CEI原则）
            listing.active = false;

            // 计算手续费
            uint256 fee = (listing.price * platformFee) / 10000;

            // 获取版税信息
            (address royaltyReceiver, uint256 royaltyAmount) = _getRoyaltyInfo(
                listing.nftContract,
                listing.tokenId,
                listing.price
            );

            // 计算卖家收益
            uint256 sellerAmount = listing.price - fee - royaltyAmount;

            // 转移NFT
            IERC721(listing.nftContract).safeTransferFrom(
                listing.seller,
                msg.sender,
                listing.tokenId
            );

            // 资金分配：版税 -> 平台手续费 -> 卖家收益
            if (royaltyAmount > 0 && royaltyReceiver != address(0)) {
                (bool successRoyalty, ) = royaltyReceiver.call{value: royaltyAmount}("");
                require(successRoyalty, "Royalty transfer failed");
            }

            (bool successSeller, ) = listing.seller.call{value: sellerAmount}("");
            require(successSeller, "Transfer to seller failed");

            (bool successFee, ) = feeRecipient.call{value: fee}("");
            require(successFee, "Transfer fee failed");

            emit NFTSold(listingId, msg.sender, listing.seller, listing.price);
        }

        // 退还多余资金
        if (remainingValue > totalRequired) {
            (bool successRefund, ) = msg.sender.call{
                value: remainingValue - totalRequired
            }("");
            require(successRefund, "Refund failed");
        }
    }

    /**
     * @dev 更新挂单白名单
     * @param listingId 挂单ID
     * @param whitelist 新的白名单地址数组
     * @notice 只有卖家可以更新白名单
     */
    function updateWhitelist(uint256 listingId, address[] calldata whitelist) external {
        Listing storage listing = listings[listingId];

        require(listing.active, "Listing not active");
        require(listing.seller == msg.sender, "Not the seller");

        listing.whitelist = whitelist;
    }

    /**
     * @dev 添加白名单买家
     * @param listingId 挂单ID
     * @param buyer 要添加的买家地址
     * @notice 只有卖家可以添加白名单买家
     */
    function addToWhitelist(uint256 listingId, address buyer) external {
        Listing storage listing = listings[listingId];

        require(listing.active, "Listing not active");
        require(listing.seller == msg.sender, "Not the seller");
        require(buyer != address(0), "Invalid buyer address");
        require(!_isWhitelisted(listing.whitelist, buyer), "Already whitelisted");

        listing.whitelist.push(buyer);
    }

    /**
     * @dev 从白名单移除买家
     * @param listingId 挂单ID
     * @param buyer 要移除的买家地址
     * @notice 只有卖家可以移除白名单买家
     */
    function removeFromWhitelist(uint256 listingId, address buyer) external {
        Listing storage listing = listings[listingId];

        require(listing.active, "Listing not active");
        require(listing.seller == msg.sender, "Not the seller");

        for (uint256 i = 0; i < listing.whitelist.length; i++) {
            if (listing.whitelist[i] == buyer) {
                // 将最后一个元素移到当前位置，然后删除最后一个
                listing.whitelist[i] = listing.whitelist[listing.whitelist.length - 1];
                listing.whitelist.pop();
                break;
            }
        }
    }

    /**
     * @dev 检查地址是否在白名单中
     * @param whitelist 白名单数组
     * @param buyer 要检查的地址
     * @return 是否在白名单中
     */
    function _isWhitelisted(address[] memory whitelist, address buyer) internal pure returns (bool) {
        for (uint256 i = 0; i < whitelist.length; i++) {
            if (whitelist[i] == buyer) {
                return true;
            }
        }
        return false;
    }

    /**
     * @dev 获取挂单的白名单
     * @param listingId 挂单ID
     * @return 白名单地址数组
     */
    function getWhitelist(uint256 listingId) external view returns (address[] memory) {
        return listings[listingId].whitelist;
    }

    /**
     * @dev 检查买家是否可以购买挂单
     * @param listingId 挂单ID
     * @param buyer 买家地址
     * @return 是否可以购买
     */
    function canBuy(uint256 listingId, address buyer) external view returns (bool) {
        Listing memory listing = listings[listingId];

        if (!listing.active) return false;
        if (buyer == listing.seller) return false;
        if (listing.isPrivate && !_isWhitelisted(listing.whitelist, buyer)) return false;

        return true;
    }
}
