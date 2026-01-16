package blockchain

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"go-auction/config"
	"go-auction/models"
	"go-auction/repositories"
	"go-auction/utils"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"go-auction/blockchain/contract"
)

// EventHandler 事件处理器
type EventHandler struct {
	auctionRepo    *repositories.AuctionRepository
	bidRepo        *repositories.BidRepository
	nftRepo        *repositories.NFTRepository
	syncRepo       *repositories.SyncStatusRepository
	metadataService *NFTMetadataService
	contractAddr   common.Address
	retryConfig    RetryConfig
}

// NewEventHandler 创建事件处理器
func NewEventHandler(
	auctionRepo *repositories.AuctionRepository,
	bidRepo *repositories.BidRepository,
	nftRepo *repositories.NFTRepository,
	syncRepo *repositories.SyncStatusRepository,
	client *Client,
	contractAddr common.Address,
	cfg *config.NFTMetadataConfig,
) *EventHandler {
	metadataService := NewNFTMetadataService(client.GetRPCClient(), cfg)
	return &EventHandler{
		auctionRepo:     auctionRepo,
		bidRepo:         bidRepo,
		nftRepo:         nftRepo,
		syncRepo:        syncRepo,
		metadataService: metadataService,
		contractAddr:    contractAddr,
		retryConfig:     DefaultRetryConfig,
	}
}

// HandleLog 处理单个日志（统一入口，支持重试）
// 注意：这个方法需要合约实例来解析事件，所以应该在EventListener中调用
// 这里保留一个简化版本，实际解析在EventListener中完成
func (h *EventHandler) HandleLog(ctx context.Context, log types.Log) error {
	return RetryWithBackoff(ctx, h.retryConfig, func() error {
		// 事件解析在EventListener中完成，这里只处理已解析的事件
		// 实际的事件处理逻辑在各个Handle方法中
		return nil
	})
}

// HandleAuctionCreated 处理拍卖创建事件
func (h *EventHandler) HandleAuctionCreated(event *contract.AuctionAuctionCreated, log types.Log) error {
	// 检查事件是否已处理（去重）
	auctionID := uint64(event.AuctionId.Uint64())
	exists, err := h.auctionRepo.Exists(auctionID)
	if err != nil {
		return fmt.Errorf("failed to check auction existence: %w", err)
	}
	if exists {
		slog.Debug("AuctionCreated event already processed", "auction_id", auctionID)
		return nil
	}

	// 创建或更新 NFT 记录
	nft := &models.NFT{
		ContractAddress: event.NftContract.Hex(),
		TokenID:         event.TokenId.Uint64(),
		Owner:           event.Seller.Hex(),
		Name:            fmt.Sprintf("NFT #%d", event.TokenId.Uint64()),
		ImageURL:        "",
		Description:     "",
	}

	existingNFT, err := h.nftRepo.GetByContractAndTokenID(nft.ContractAddress, nft.TokenID)
	if err != nil {
		// 数据库查询错误（非 record not found）
		slog.Warn("failed to query NFT", "error", err, "contract", nft.ContractAddress, "token_id", nft.TokenID)
	} else if existingNFT != nil {
		// NFT 已存在，更新所有者
		existingNFT.Owner = nft.Owner
		if err := h.nftRepo.Update(existingNFT); err != nil {
			slog.Warn("failed to update NFT", "error", err, "nft_id", existingNFT.ID)
		}
		// 如果元数据为空，异步获取并更新
		if existingNFT.ImageURL == "" || existingNFT.Description == "" {
			nftRepoAdapter := &nftRepositoryAdapter{repo: *h.nftRepo}
			h.metadataService.UpdateNFTMetadata(context.Background(), nftRepoAdapter, nft.ContractAddress, nft.TokenID)
		}
	} else {
		// NFT 不存在，创建新记录
		// 注意：ImageURL 和 Description 为空是正常的，因为链上事件不包含元数据
		// 异步获取元数据并更新
		if err := h.nftRepo.Create(nft); err != nil {
			slog.Warn("failed to create NFT", "error", err)
		} else {
			// 异步获取并更新 NFT 元数据
			nftRepoAdapter := &nftRepositoryAdapter{repo: *h.nftRepo}
			h.metadataService.UpdateNFTMetadata(context.Background(), nftRepoAdapter, nft.ContractAddress, nft.TokenID)
		}
	}

	// 创建拍卖记录
	auction := &models.Auction{
		ID:                auctionID,
		NFTContract:       event.NftContract.Hex(),
		TokenID:           event.TokenId.Uint64(),
		Seller:            event.Seller.Hex(),
		StartingPrice:     formatUSD(event.StartingPrice),
		CurrentHighestBid: "0",
		HighestBidder:     "",
		StartTime:         uint64(time.Now().Unix()), // 使用当前时间作为开始时间
		EndTime:           event.EndTime.Uint64(),
		PaymentToken:      event.PaymentToken.Hex(),
		Status:            models.AuctionStatusActive,
	}

	if err := h.auctionRepo.Create(auction); err != nil {
		return fmt.Errorf("failed to create auction: %w", err)
	}

	// 清除相关缓存
	ctx := context.Background()
	utils.CacheDelete(ctx, utils.CacheKey("auction", auction.ID))
	utils.CacheDeletePattern(ctx, "auctions:list:*")

	slog.Info("AuctionCreated event processed", "auction_id", auction.ID, "seller", auction.Seller)
	return nil
}

// HandleBidPlaced 处理出价事件
func (h *EventHandler) HandleBidPlaced(event *contract.AuctionBidPlaced, log types.Log) error {
	// 检查事件是否已处理（去重）
	exists, err := h.bidRepo.ExistsByEvent(log.BlockNumber, log.TxHash.Hex(), log.Index)
	if err != nil {
		return fmt.Errorf("failed to check bid existence: %w", err)
	}
	if exists {
		slog.Debug("BidPlaced event already processed", "auction_id", event.AuctionId.Uint64(), "block", log.BlockNumber)
		return nil
	}

	// 创建出价记录
	bid := &models.Bid{
		AuctionID:      uint64(event.AuctionId.Uint64()),
		Bidder:         event.Bidder.Hex(),
		BidAmountUSD:   formatUSD(event.BidAmountUSD),
		OriginalAmount: formatUSD(event.OriginalAmount),
		PaymentToken:   event.PaymentToken.Hex(),
		BlockNumber:    log.BlockNumber,
		TxHash:         log.TxHash.Hex(),
		LogIndex:       log.Index,
	}

	if err := h.bidRepo.Create(bid); err != nil {
		return fmt.Errorf("failed to create bid: %w", err)
	}

	// 更新拍卖的最高出价信息
	auction, err := h.auctionRepo.GetByID(bid.AuctionID)
	if err != nil {
		slog.Warn("failed to get auction for updating highest bid", "error", err, "auction_id", bid.AuctionID)
	} else {
		// 比较当前出价和最高出价
		currentHighestBid := parseUSD(auction.CurrentHighestBid)
		newBidAmount := event.BidAmountUSD

		if newBidAmount.Cmp(currentHighestBid) > 0 {
			auction.CurrentHighestBid = formatUSD(newBidAmount)
			auction.HighestBidder = bid.Bidder
			if err := h.auctionRepo.Update(auction); err != nil {
				slog.Warn("failed to update auction highest bid", "error", err, "auction_id", auction.ID)
			}
		}
	}

	// 清除相关缓存
	ctx := context.Background()
	utils.CacheDelete(ctx,
		utils.CacheKey("auction", bid.AuctionID),
		utils.CacheKey("bid", "highest", bid.AuctionID),
	)
	utils.CacheDeletePattern(ctx, fmt.Sprintf("bids:auction:%d:*", bid.AuctionID))
	utils.CacheDeletePattern(ctx, "auctions:list:*")

	slog.Info("BidPlaced event processed", "auction_id", bid.AuctionID, "bidder", bid.Bidder, "amount_usd", bid.BidAmountUSD)
	return nil
}

// HandleAuctionEnded 处理拍卖结束事件
func (h *EventHandler) HandleAuctionEnded(event *contract.AuctionAuctionEnded, log types.Log) error {
	auctionID := uint64(event.AuctionId.Uint64())

	// 获取拍卖记录
	auction, err := h.auctionRepo.GetByID(auctionID)
	if err != nil {
		return fmt.Errorf("failed to get auction: %w", err)
	}

	// 更新拍卖状态
	auction.Status = models.AuctionStatusEnded
	auction.CurrentHighestBid = formatUSD(event.FinalPriceUSD)
	if event.Winner != (common.Address{}) {
		auction.HighestBidder = event.Winner.Hex()
	}

	if err := h.auctionRepo.Update(auction); err != nil {
		return fmt.Errorf("failed to update auction: %w", err)
	}

	// 清除相关缓存
	ctx := context.Background()
	utils.CacheDelete(ctx, utils.CacheKey("auction", auctionID))
	utils.CacheDeletePattern(ctx, "auctions:list:*")

	slog.Info("AuctionEnded event processed", "auction_id", auctionID, "winner", auction.HighestBidder)
	return nil
}

// HandleBidWithdrawn 处理出价撤回事件
func (h *EventHandler) HandleBidWithdrawn(event *contract.AuctionBidWithdrawn, log types.Log) error {
	// BidWithdrawn 事件主要用于记录，不需要创建新的 Bid 记录
	// 如果需要记录撤回历史，可以在这里添加逻辑
	slog.Info("BidWithdrawn event processed", "auction_id", event.AuctionId.Uint64(), "bidder", event.Bidder.Hex(), "amount", event.Amount.Uint64())
	return nil
}

// formatUSD 将 big.Int 格式化为 USD 字符串（18位小数）
func formatUSD(amount *big.Int) string {
	if amount == nil {
		return "0"
	}
	// 转换为字符串，保留18位小数
	amountStr := amount.String()
	if len(amountStr) <= 18 {
		// 补零
		zeros := 18 - len(amountStr)
		amountStr = "0." + fmt.Sprintf("%0*s", zeros, "") + amountStr
	} else {
		// 插入小数点
		amountStr = amountStr[:len(amountStr)-18] + "." + amountStr[len(amountStr)-18:]
	}
	return amountStr
}

// parseUSD 将 USD 字符串解析为 big.Int（18位小数）
func parseUSD(amountStr string) *big.Int {
	// 解析带小数点的字符串
	// 例如 "1234.567890123456789012" -> 1234567890123456789012
	amount := new(big.Int)

	// 移除小数点，转换为整数
	amountStr = strings.ReplaceAll(amountStr, ".", "")
	if amountStr == "" {
		return big.NewInt(0)
	}

	amount.SetString(amountStr, 10)
	return amount
}
