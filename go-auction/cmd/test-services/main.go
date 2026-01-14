package main

import (
	"log/slog"
	"os"
	"time"

	"go-auction/config"
	"go-auction/dto"
	"go-auction/models"
	"go-auction/repositories"
	"go-auction/services"
)

func main() {
	// 初始化日志
	initLogger()

	slog.Info("开始测试 Service 层...")

	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	config.InitDatabase(cfg)
	defer config.Close()

	// 准备测试数据
	setupTestData()

	// 测试 AuctionService
	slog.Info("=== 测试 AuctionService ===")
	testAuctionService()

	// 测试 BidService
	slog.Info("=== 测试 BidService ===")
	testBidService()

	// 测试 NFTService
	slog.Info("=== 测试 NFTService ===")
	testNFTService()

	slog.Info("Service 层测试完成！")
}

func setupTestData() {
	// 创建测试 NFT
	nftRepo := repositories.NewNFTRepository()
	nft := &models.NFT{
		ContractAddress: "0xTESTNFT123456789012345678901234567890",
		TokenID:         uint64(time.Now().Unix()),
		Name:            "Test NFT for Service",
		ImageURL:        "https://example.com/nft.png",
		Description:     "Test NFT",
		Owner:           "0xSELLER123456789012345678901234567890",
	}
	nftRepo.Create(nft)

	// 创建测试拍卖
	auctionRepo := repositories.NewAuctionRepository()
	auction := &models.Auction{
		NFTContract:       nft.ContractAddress,
		TokenID:           nft.TokenID,
		Seller:            "0xSELLER123456789012345678901234567890",
		StartingPrice:     "1000.000000000000000000",
		CurrentHighestBid: "0",
		HighestBidder:     "",
		StartTime:         uint64(time.Now().Unix()),
		EndTime:           uint64(time.Now().Add(24 * time.Hour).Unix()),
		PaymentToken:      "0x0000000000000000000000000000000000000000",
		Status:            models.AuctionStatusActive,
	}
	if err := auctionRepo.Create(auction); err != nil {
		slog.Error("创建测试拍卖失败", "error", err)
		return
	}

	// 创建测试出价（使用创建后的 auction.ID）
	bidRepo := repositories.NewBidRepository()
	// 使用正确的交易哈希格式（66字符：0x + 64个十六进制字符）
	txHash := "0x" + "1234567890123456789012345678901234567890123456789012345678901234"
	bid := &models.Bid{
		AuctionID:      auction.ID,
		Bidder:         "0xBIDDER123456789012345678901234567890",
		BidAmountUSD:   "1200.000000000000000000",
		OriginalAmount: "1.000000000000000000",
		PaymentToken:   "0x0000000000000000000000000000000000000000",
		BlockNumber:    uint64(time.Now().Unix()),
		TxHash:         txHash,
		LogIndex:       0,
	}
	if err := bidRepo.Create(bid); err != nil {
		slog.Warn("创建测试出价失败（可能是唯一索引冲突）", "error", err)
	}
}

func testAuctionService() {
	service := services.NewAuctionService()

	// 测试 GetList
	req := &dto.AuctionListRequest{
		Page:     1,
		PageSize: 10,
	}
	listVO, err := service.GetList(req)
	if err != nil {
		slog.Error("GetList失败", "error", err)
	} else {
		slog.Info("GetList成功", "total", listVO.Total, "count", len(listVO.List))
		if len(listVO.List) > 0 {
			slog.Info("第一条拍卖", "id", listVO.List[0].ID, "status", listVO.List[0].Status)
		}
	}

	// 测试 GetActiveAuctions
	activeList, err := service.GetActiveAuctions(req)
	if err != nil {
		slog.Error("GetActiveAuctions失败", "error", err)
	} else {
		slog.Info("GetActiveAuctions成功", "total", activeList.Total, "count", len(activeList.List))
	}

	// 测试 GetByID（使用第一个拍卖的ID）
	if listVO != nil && len(listVO.List) > 0 {
		auctionID := listVO.List[0].ID
		auctionVO, err := service.GetByID(auctionID)
		if err != nil {
			slog.Error("GetByID失败", "error", err)
		} else {
			slog.Info("GetByID成功", "id", auctionVO.ID, "seller", auctionVO.Seller)
		}
	}
}

func testBidService() {
	service := services.NewBidService()

	// 先获取一个拍卖ID
	auctionRepo := repositories.NewAuctionRepository()
	auctions, _, _ := auctionRepo.GetList(1, 1, "", "", "")
	if len(auctions) == 0 {
		slog.Warn("没有测试数据，跳过 BidService 测试")
		return
	}
	auctionID := auctions[0].ID

	// 测试 GetByAuctionID
	req := &dto.BidListRequest{
		AuctionID: auctionID,
		Page:      1,
		PageSize:  10,
	}
	listVO, err := service.GetByAuctionID(auctionID, req)
	if err != nil {
		slog.Error("GetByAuctionID失败", "error", err)
	} else {
		slog.Info("GetByAuctionID成功", "total", listVO.Total, "count", len(listVO.List))
	}

	// 测试 GetHighestBid
	highestBid, err := service.GetHighestBid(auctionID)
	if err != nil {
		slog.Error("GetHighestBid失败", "error", err)
	} else {
		slog.Info("GetHighestBid成功", "bidder", highestBid.Bidder, "amount", highestBid.BidAmountUSD)
	}
}

func testNFTService() {
	service := services.NewNFTService()

	// 先获取一个NFT
	nftRepo := repositories.NewNFTRepository()
	nfts, _, _ := nftRepo.GetByOwner("0xSELLER123456789012345678901234567890", 1, 1)
	if len(nfts) == 0 {
		slog.Warn("没有测试数据，跳过 NFTService 测试")
		return
	}
	nft := nfts[0]

	// 测试 GetByID
	nftVO, err := service.GetByID(nft.ID)
	if err != nil {
		slog.Error("GetByID失败", "error", err)
	} else {
		slog.Info("GetByID成功", "id", nftVO.ID, "name", nftVO.Name)
	}

	// 测试 GetByContractAndTokenID
	nftVO2, err := service.GetByContractAndTokenID(nft.ContractAddress, nft.TokenID)
	if err != nil {
		slog.Error("GetByContractAndTokenID失败", "error", err)
	} else {
		slog.Info("GetByContractAndTokenID成功", "id", nftVO2.ID, "name", nftVO2.Name)
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
