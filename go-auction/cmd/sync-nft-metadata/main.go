package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	"go-auction/blockchain"
	"go-auction/config"
	"go-auction/models"
	"go-auction/repositories"

	"github.com/ethereum/go-ethereum/common"
)

func main() {
	// 解析命令行参数
	var (
		batchSize   = flag.Int("batch-size", 10, "每批处理的NFT数量")
		concurrency = flag.Int("concurrency", 5, "并发数")
		dryRun      = flag.Bool("dry-run", false, "仅显示将要处理的NFT，不实际更新")
	)
	flag.Parse()

	// 初始化日志
	initLogger()

	slog.Info("开始同步 NFT 元数据", "batch_size", *batchSize, "concurrency", *concurrency, "dry_run", *dryRun)

	// 加载配置
	cfg := config.Load()
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "local"
	}

	// 初始化数据库
	config.InitDatabase(cfg)
	defer config.Close()

	// 初始化Redis（可选）
	config.InitRedis(cfg)
	defer config.CloseRedis()

	// 初始化区块链客户端
	if cfg.Blockchain.ContractAddress == "" || cfg.Blockchain.ContractAddress == "0x..." {
		slog.Error("未配置合约地址，无法同步元数据")
		os.Exit(1)
	}

	bcClient, err := blockchain.NewClient(&cfg.Blockchain)
	if err != nil {
		slog.Error("区块链客户端初始化失败", "error", err)
		os.Exit(1)
	}
	defer bcClient.Close()

	// 创建 Repository
	nftRepo := repositories.NewNFTRepository()

	// 创建元数据服务
	metadataService := blockchain.NewNFTMetadataService(bcClient.GetRPCClient(), &cfg.Blockchain.NFTMetadata)

	// 创建适配器
	nftRepoAdapter := &nftRepositoryAdapter{repo: *nftRepo}

	// 获取所有需要同步的NFT（image_url 或 description 为空）
	nfts, err := nftRepo.GetNFTsWithoutMetadata()
	if err != nil {
		slog.Error("获取NFT列表失败", "error", err)
		os.Exit(1)
	}

	total := len(nfts)
	if total == 0 {
		slog.Info("没有需要同步的NFT")
		return
	}

	slog.Info("找到需要同步的NFT", "total", total)

	if *dryRun {
		slog.Info("Dry run 模式：仅显示将要处理的NFT")
		for _, nft := range nfts {
			slog.Info("将处理NFT",
				"id", nft.ID,
				"contract", nft.ContractAddress,
				"token_id", nft.TokenID,
				"name", nft.Name,
				"has_image", nft.ImageURL != "",
				"has_description", nft.Description != "")
		}
		return
	}

	// 使用 Worker Pool 模式并发处理
	ctx := context.Background()
	jobs := make(chan *models.NFT, *batchSize)
	var wg sync.WaitGroup

	// 启动 Worker
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go worker(ctx, i, jobs, metadataService, nftRepoAdapter, &wg)
	}

	// 发送任务
	go func() {
		defer close(jobs)
		for _, nft := range nfts {
			jobs <- nft
		}
	}()

	// 等待所有 Worker 完成
	wg.Wait()

	slog.Info("NFT 元数据同步完成", "total", total)
}

// worker 工作协程
func worker(ctx context.Context, id int, jobs <-chan *models.NFT, metadataService *blockchain.NFTMetadataService, nftRepo blockchain.NFTRepository, wg *sync.WaitGroup) {
	defer wg.Done()

	processed := 0
	for nft := range jobs {
		processed++
		slog.Info("Worker处理NFT",
			"worker_id", id,
			"nft_id", nft.ID,
			"contract", nft.ContractAddress,
			"token_id", nft.TokenID,
			"processed", processed)

		// 调用同步方法（同步版本）
		if err := syncNFTMetadata(ctx, metadataService, nftRepo, nft.ContractAddress, nft.TokenID); err != nil {
			// 检查是否是 token 不存在的错误
			if strings.Contains(err.Error(), "token does not exist") {
				slog.Info("跳过NFT（链上不存在）",
					"worker_id", id,
					"contract", nft.ContractAddress,
					"token_id", nft.TokenID,
					"note", "NFT可能在拍卖合约中或已被销毁")
			} else {
				slog.Warn("同步NFT元数据失败",
					"worker_id", id,
					"contract", nft.ContractAddress,
					"token_id", nft.TokenID,
					"error", err)
			}
		} else {
			slog.Info("NFT元数据同步成功",
				"worker_id", id,
				"contract", nft.ContractAddress,
				"token_id", nft.TokenID)
		}

		// 避免请求过快
		time.Sleep(100 * time.Millisecond)
	}

	slog.Info("Worker完成", "worker_id", id, "processed", processed)
}

// syncNFTMetadata 同步NFT元数据（同步版本）
func syncNFTMetadata(ctx context.Context, metadataService *blockchain.NFTMetadataService, nftRepo blockchain.NFTRepository, contractAddress string, tokenID uint64) error {
	contractAddr := common.HexToAddress(contractAddress)
	tokenIDBig := big.NewInt(int64(tokenID))

	// 获取元数据
	metadata, err := metadataService.GetMetadata(ctx, contractAddr, tokenIDBig)
	if err != nil {
		return fmt.Errorf("failed to get metadata: %w", err)
	}

	// 获取NFT记录
	nft, err := nftRepo.GetByContractAndTokenID(contractAddress, tokenID)
	if err != nil {
		return fmt.Errorf("failed to get NFT: %w", err)
	}

	if nft == nil {
		return fmt.Errorf("NFT not found")
	}

	// 更新元数据字段
	updated := false
	if metadata.Name != "" && nft.Name == "" {
		nft.Name = metadata.Name
		updated = true
	}
	if metadata.Image != "" && nft.ImageURL == "" {
		nft.ImageURL = metadata.Image
		updated = true
	}
	if metadata.Description != "" && nft.Description == "" {
		nft.Description = metadata.Description
		updated = true
	}

	if updated {
		if err := nftRepo.Update(nft); err != nil {
			return fmt.Errorf("failed to update NFT: %w", err)
		}
	}

	return nil
}

// nftRepositoryAdapter 适配器
type nftRepositoryAdapter struct {
	repo repositories.NFTRepository
}

func (a *nftRepositoryAdapter) GetByContractAndTokenID(contractAddress string, tokenID uint64) (*blockchain.NFTModel, error) {
	nft, err := a.repo.GetByContractAndTokenID(contractAddress, tokenID)
	if err != nil {
		return nil, err
	}
	if nft == nil {
		return nil, nil
	}
	return &blockchain.NFTModel{
		ID:              nft.ID,
		ContractAddress: nft.ContractAddress,
		TokenID:         nft.TokenID,
		Name:            nft.Name,
		ImageURL:        nft.ImageURL,
		Description:     nft.Description,
		Owner:           nft.Owner,
	}, nil
}

func (a *nftRepositoryAdapter) Update(nft *blockchain.NFTModel) error {
	model := &models.NFT{
		ID:              nft.ID,
		ContractAddress: nft.ContractAddress,
		TokenID:         nft.TokenID,
		Name:            nft.Name,
		ImageURL:        nft.ImageURL,
		Description:     nft.Description,
		Owner:           nft.Owner,
	}
	return a.repo.Update(model)
}

// initLogger 初始化日志配置
func initLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	var handler slog.Handler
	if os.Getenv("GO_ENV") == "prd" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
}
