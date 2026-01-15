package blockchain

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"sync"
	"time"

	"go-auction/repositories"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"go-auction/blockchain/contract"
)

// EventListener 事件监听器
type EventListener struct {
	client       *Client
	contract     *contract.Auction
	syncRepo     *repositories.SyncStatusRepository
	auctionRepo  *repositories.AuctionRepository
	bidRepo      *repositories.BidRepository
	nftRepo      *repositories.NFTRepository
	handler      *EventHandler
	processor    *EventProcessor
	startBlock   uint64
	currentBlock uint64
	isRunning    bool
	stopChan     chan struct{}
	mu           sync.RWMutex
	batchSize    uint64        // 每次同步的区块数量
	useWS        bool          // 是否使用WebSocket
	pollInterval time.Duration // HTTP轮询间隔
}

// NewEventListener 创建事件监听器
func NewEventListener(
	client *Client,
	syncRepo *repositories.SyncStatusRepository,
	auctionRepo *repositories.AuctionRepository,
	bidRepo *repositories.BidRepository,
	nftRepo *repositories.NFTRepository,
	startBlock uint64,
) (*EventListener, error) {
	// 创建合约实例
	contractInstance, err := contract.NewAuction(client.GetContractAddress(), client.GetRPCClient())
	if err != nil {
		return nil, fmt.Errorf("failed to create contract instance: %w", err)
	}

	// 创建事件处理器
	handler := NewEventHandler(auctionRepo, bidRepo, nftRepo, syncRepo, client.GetContractAddress())

	// 创建Worker Pool事件处理器
	processor := NewEventProcessor(handler, contractInstance, 5) // 默认5个worker

	return &EventListener{
		client:       client,
		contract:     contractInstance,
		syncRepo:     syncRepo,
		auctionRepo:  auctionRepo,
		bidRepo:      bidRepo,
		nftRepo:      nftRepo,
		handler:      handler,
		processor:    processor,
		startBlock:   startBlock,
		currentBlock: startBlock,
		isRunning:    false,
		stopChan:     make(chan struct{}),
		batchSize:    1000, // 每次同步1000个区块
		useWS:        client.GetWSClient() != nil,
		pollInterval: 5 * time.Second, // HTTP轮询间隔5秒
	}, nil
}

// Start 启动事件监听
func (l *EventListener) Start(ctx context.Context) error {
	l.mu.Lock()
	if l.isRunning {
		l.mu.Unlock()
		return fmt.Errorf("event listener is already running")
	}
	l.isRunning = true
	l.mu.Unlock()

	slog.Info("启动事件监听器", "start_block", l.startBlock, "contract", l.client.GetContractAddress().Hex())

	// 启动Worker Pool
	l.processor.Start()

	// 获取或创建同步状态
	syncStatus, err := l.syncRepo.GetOrCreate(l.client.GetContractAddress().Hex())
	if err != nil {
		l.processor.Stop()
		return fmt.Errorf("failed to get sync status: %w", err)
	}

	// 如果数据库中有同步状态，使用数据库中的区块号
	if syncStatus.LastSyncedBlock > 0 {
		l.currentBlock = syncStatus.LastSyncedBlock
		slog.Info("从数据库恢复同步状态", "last_synced_block", l.currentBlock)
	} else {
		l.currentBlock = l.startBlock
	}

	// 先同步历史事件
	go l.syncHistoricalEvents(ctx)

	// 然后订阅新事件（WebSocket或HTTP轮询）
	if l.useWS {
		go l.subscribeNewEvents(ctx)
	} else {
		go l.startPollingFallback(ctx)
	}

	return nil
}

// Stop 停止监听
func (l *EventListener) Stop() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.isRunning {
		return
	}

	l.isRunning = false
	close(l.stopChan)

	// 停止Worker Pool
	if l.processor != nil {
		l.processor.Stop()
	}

	slog.Info("事件监听器已停止")
}

// syncHistoricalEvents 同步历史事件
func (l *EventListener) syncHistoricalEvents(ctx context.Context) {
	slog.Info("开始同步历史事件", "from_block", l.currentBlock)

	for {
		select {
		case <-l.stopChan:
			return
		case <-ctx.Done():
			return
		default:
		}

		// 获取最新区块号
		latestBlock, err := l.client.GetLatestBlock(ctx)
		if err != nil {
			slog.Error("获取最新区块号失败", "error", err)
			time.Sleep(5 * time.Second)
			continue
		}

		// 如果已经同步到最新区块，等待新区块
		if l.currentBlock >= latestBlock {
			slog.Debug("已同步到最新区块", "current_block", l.currentBlock, "latest_block", latestBlock)
			time.Sleep(5 * time.Second)
			continue
		}

		// 计算本次同步的结束区块
		toBlock := l.currentBlock + l.batchSize
		if toBlock > latestBlock {
			toBlock = latestBlock
		}

		// 同步事件
		if err := l.syncEventsInRange(ctx, l.currentBlock, toBlock); err != nil {
			slog.Error("同步事件失败", "error", err, "from_block", l.currentBlock, "to_block", toBlock)
			time.Sleep(5 * time.Second)
			continue
		}

		// 更新同步状态
		l.currentBlock = toBlock + 1
		if err := l.syncRepo.UpdateLastSyncedBlock(l.client.GetContractAddress().Hex(), l.currentBlock-1); err != nil {
			slog.Error("更新同步状态失败", "error", err)
		}

		slog.Info("同步进度", "current_block", l.currentBlock-1, "latest_block", latestBlock)
	}
}

// syncEventsInRange 同步指定范围内的事件（批量查询优化）
func (l *EventListener) syncEventsInRange(ctx context.Context, fromBlock, toBlock uint64) error {
	// 如果范围太大，分批查询（每批最多1000个区块）
	const maxBatchSize = 1000
	if toBlock-fromBlock > maxBatchSize {
		// 分批处理
		for start := fromBlock; start <= toBlock; start += maxBatchSize {
			end := start + maxBatchSize - 1
			if end > toBlock {
				end = toBlock
			}
			if err := l.syncEventsInRange(ctx, start, end); err != nil {
				return err
			}
		}
		return nil
	}

	// 构建查询过滤器
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(fromBlock)),
		ToBlock:   big.NewInt(int64(toBlock)),
		Addresses: []common.Address{l.client.GetContractAddress()},
	}

	// 查询日志（带重试）
	var logs []types.Log
	err := RetryWithBackoffDefault(ctx, func() error {
		var err error
		logs, err = l.client.GetRPCClient().FilterLogs(ctx, query)
		return err
	})
	if err != nil {
		return fmt.Errorf("failed to filter logs: %w", err)
	}

	if len(logs) == 0 {
		return nil
	}

	slog.Info("找到事件日志", "count", len(logs), "from_block", fromBlock, "to_block", toBlock)

	// 使用Worker Pool并发处理日志
	for _, log := range logs {
		l.processor.ProcessLog(ctx, log)
	}

	return nil
}

// startPollingFallback HTTP轮询备用机制
func (l *EventListener) startPollingFallback(ctx context.Context) {
	slog.Info("启动HTTP轮询备用模式", "poll_interval", l.pollInterval)

	ticker := time.NewTicker(l.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-l.stopChan:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			// 获取最新区块号
			latestBlock, err := l.client.GetLatestBlock(ctx)
			if err != nil {
				slog.Error("获取最新区块号失败", "error", err)
				continue
			}

			// 如果当前区块落后，同步新区块
			if latestBlock > l.currentBlock {
				toBlock := l.currentBlock + l.batchSize
				if toBlock > latestBlock {
					toBlock = latestBlock
				}

				if err := l.syncEventsInRange(ctx, l.currentBlock, toBlock); err != nil {
					slog.Error("轮询同步事件失败", "error", err, "from_block", l.currentBlock, "to_block", toBlock)
					continue
				}

				// 更新同步状态
				l.currentBlock = toBlock + 1
				if err := l.syncRepo.UpdateLastSyncedBlock(l.client.GetContractAddress().Hex(), l.currentBlock-1); err != nil {
					slog.Error("更新同步状态失败", "error", err)
				}
			}
		}
	}
}

// subscribeNewEvents 订阅新事件（WebSocket模式）
func (l *EventListener) subscribeNewEvents(ctx context.Context) {
	// 等待历史事件同步完成
	time.Sleep(10 * time.Second)

	wsClient := l.client.GetWSClient()
	if wsClient == nil {
		slog.Warn("WebSocket 客户端不可用，切换到HTTP轮询模式")
		go l.startPollingFallback(ctx)
		return
	}

	slog.Info("开始订阅新事件（WebSocket模式）")

	for {
		select {
		case <-l.stopChan:
			return
		case <-ctx.Done():
			return
		default:
		}

		// 构建订阅过滤器
		query := ethereum.FilterQuery{
			Addresses: []common.Address{l.client.GetContractAddress()},
		}

		// 订阅日志
		logsChan := make(chan types.Log)
		sub, err := wsClient.SubscribeFilterLogs(ctx, query, logsChan)
		if err != nil {
			slog.Error("订阅事件失败", "error", err)
			time.Sleep(5 * time.Second)
			continue
		}

		slog.Info("事件订阅成功")

		// 处理订阅的日志
		for {
			select {
			case <-l.stopChan:
				sub.Unsubscribe()
				return
			case <-ctx.Done():
				sub.Unsubscribe()
				return
			case err := <-sub.Err():
				slog.Error("订阅错误", "error", err)
				sub.Unsubscribe()
				// 尝试重连
				if err := l.client.Reconnect(); err != nil {
					slog.Error("重连失败，切换到HTTP轮询模式", "error", err)
					// WebSocket失败，切换到HTTP轮询
					go l.startPollingFallback(ctx)
					return
				}
				// 重连成功，继续订阅
				break
			case log := <-logsChan:
				l.processor.ProcessLog(ctx, log)
			}
		}
	}
}
