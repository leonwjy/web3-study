package main

import (
	"log/slog"
	"os"
	"time"

	"go-auction/config"
	"go-auction/models"
	"go-auction/repositories"
)

func main() {
	// 初始化日志
	initLogger()

	slog.Info("开始测试 Repository 层...")

	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	config.InitDatabase(cfg)
	defer config.Close()

	// 测试 AuctionRepository
	slog.Info("=== 测试 AuctionRepository ===")
	testAuctionRepository()

	// 测试 BidRepository
	slog.Info("=== 测试 BidRepository ===")
	testBidRepository()

	// 测试 NFTRepository
	slog.Info("=== 测试 NFTRepository ===")
	testNFTRepository()

	// 测试 SyncStatusRepository
	slog.Info("=== 测试 SyncStatusRepository ===")
	testSyncStatusRepository()

	slog.Info("Repository 层测试完成！")
}

func testAuctionRepository() {
	repo := repositories.NewAuctionRepository()

	// 创建测试拍卖
	auction := &models.Auction{
		NFTContract:       "0x1234567890123456789012345678901234567890",
		TokenID:           1,
		Seller:            "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
		StartingPrice:     "1000.000000000000000000",
		CurrentHighestBid: "0",
		HighestBidder:     "",
		StartTime:         uint64(time.Now().Unix()),
		EndTime:           uint64(time.Now().Add(24 * time.Hour).Unix()),
		PaymentToken:      "0x0000000000000000000000000000000000000000",
		Status:            models.AuctionStatusActive,
	}

	err := repo.Create(auction)
	if err != nil {
		slog.Error("创建拍卖失败", "error", err)
		return
	}
	slog.Info("创建拍卖成功", "id", auction.ID)

	// 根据ID查询
	found, err := repo.GetByID(auction.ID)
	if err != nil {
		slog.Error("查询拍卖失败", "error", err)
	} else {
		slog.Info("查询拍卖成功", "id", found.ID, "seller", found.Seller)
	}

	// 测试 GetList
	auctions, total, err := repo.GetList(1, 10, "", "", "")
	if err != nil {
		slog.Error("获取拍卖列表失败", "error", err)
	} else {
		slog.Info("获取拍卖列表成功", "total", total, "count", len(auctions))
	}

	// 测试 GetActiveAuctions
	activeAuctions, _, err := repo.GetActiveAuctions(1, 10)
	if err != nil {
		slog.Error("获取活跃拍卖失败", "error", err)
	} else {
		slog.Info("获取活跃拍卖成功", "count", len(activeAuctions))
	}

	// 测试 UpdateStatus
	err = repo.UpdateStatus(auction.ID, models.AuctionStatusEnded)
	if err != nil {
		slog.Error("更新状态失败", "error", err)
	} else {
		slog.Info("更新状态成功")
	}
}

func testBidRepository() {
	repo := repositories.NewBidRepository()

	// 创建测试出价（使用不同的区块号和交易哈希避免唯一索引冲突）
	bid := &models.Bid{
		AuctionID:      1,
		Bidder:         "0x1111111111111111111111111111111111111111",
		BidAmountUSD:   "1200.000000000000000000",
		OriginalAmount: "1.000000000000000000",
		PaymentToken:   "0x0000000000000000000000000000000000000000",
		BlockNumber:    uint64(time.Now().Unix()), // 使用时间戳作为区块号避免冲突
		TxHash:         "0x3333333333333333333333333333333333333333333333333333333333333333",
		LogIndex:       0,
	}

	err := repo.Create(bid)
	if err != nil {
		slog.Error("创建出价失败", "error", err)
		return
	}
	slog.Info("创建出价成功", "id", bid.ID)

	// 测试 ExistsByEvent
	exists, err := repo.ExistsByEvent(bid.BlockNumber, bid.TxHash, bid.LogIndex)
	if err != nil {
		slog.Error("检查事件失败", "error", err)
	} else {
		slog.Info("检查事件成功", "exists", exists)
	}

	// 测试 GetByAuctionID
	bids, total, err := repo.GetByAuctionID(1, 1, 10)
	if err != nil {
		slog.Error("获取拍卖出价列表失败", "error", err)
	} else {
		slog.Info("获取拍卖出价列表成功", "total", total, "count", len(bids))
	}
}

func testNFTRepository() {
	repo := repositories.NewNFTRepository()

	// 创建测试NFT（使用不同的TokenID避免唯一索引冲突）
	nft := &models.NFT{
		ContractAddress: "0x1234567890123456789012345678901234567890",
		TokenID:         uint64(time.Now().Unix()), // 使用时间戳作为TokenID避免冲突
		Name:            "Test NFT",
		ImageURL:        "https://example.com/image.png",
		Description:     "This is a test NFT",
		Owner:           "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
	}

	err := repo.Create(nft)
	if err != nil {
		slog.Error("创建NFT失败", "error", err)
		return
	}
	slog.Info("创建NFT成功", "id", nft.ID)

	// 测试 GetByContractAndTokenID
	found, err := repo.GetByContractAndTokenID(nft.ContractAddress, nft.TokenID)
	if err != nil {
		slog.Error("查询NFT失败", "error", err)
	} else {
		slog.Info("查询NFT成功", "id", found.ID, "name", found.Name)
	}

	// 测试 Exists
	exists, err := repo.Exists(nft.ContractAddress, nft.TokenID)
	if err != nil {
		slog.Error("检查NFT存在失败", "error", err)
	} else {
		slog.Info("检查NFT存在成功", "exists", exists)
	}
}

func testSyncStatusRepository() {
	repo := repositories.NewSyncStatusRepository()

	// 测试 GetOrCreate
	syncStatus, err := repo.GetOrCreate("0x1234567890123456789012345678901234567890")
	if err != nil {
		slog.Error("GetOrCreate失败", "error", err)
		return
	}
	slog.Info("GetOrCreate成功", "contract", syncStatus.ContractAddress, "block", syncStatus.LastSyncedBlock)

	// 测试 UpdateLastSyncedBlock
	err = repo.UpdateLastSyncedBlock("0x1234567890123456789012345678901234567890", 2000)
	if err != nil {
		slog.Error("更新同步区块失败", "error", err)
	} else {
		slog.Info("更新同步区块成功")
	}

	// 验证更新
	updated, err := repo.GetByContractAddress("0x1234567890123456789012345678901234567890")
	if err != nil {
		slog.Error("查询同步状态失败", "error", err)
	} else {
		slog.Info("查询同步状态成功", "block", updated.LastSyncedBlock)
	}
}

// initLogger 初始化日志配置
func initLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(handler))
}
