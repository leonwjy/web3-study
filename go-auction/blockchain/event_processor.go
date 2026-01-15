package blockchain

import (
	"context"
	"log/slog"
	"sync"

	"go-auction/blockchain/contract"

	"github.com/ethereum/go-ethereum/core/types"
)

// EventProcessor 事件处理器（Worker Pool模式）
type EventProcessor struct {
	workerCount int
	jobQueue    chan types.Log
	handler     *EventHandler
	contract    *contract.Auction
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewEventProcessor 创建事件处理器
func NewEventProcessor(handler *EventHandler, contractInstance *contract.Auction, workerCount int) *EventProcessor {
	if workerCount <= 0 {
		workerCount = 5 // 默认5个worker
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &EventProcessor{
		workerCount: workerCount,
		jobQueue:    make(chan types.Log, workerCount*10), // 缓冲区大小为worker数量的10倍
		handler:     handler,
		contract:    contractInstance,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start 启动事件处理器
func (p *EventProcessor) Start() {
	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
	slog.Info("事件处理器启动", "worker_count", p.workerCount)
}

// Stop 停止事件处理器
func (p *EventProcessor) Stop() {
	p.cancel()
	close(p.jobQueue)
	p.wg.Wait()
	slog.Info("事件处理器已停止")
}

// ProcessLog 处理日志（异步）
func (p *EventProcessor) ProcessLog(ctx context.Context, log types.Log) {
	select {
	case <-ctx.Done():
		return
	case <-p.ctx.Done():
		return
	case p.jobQueue <- log:
		// 成功加入队列
	default:
		// 队列已满，记录警告但继续处理（避免阻塞）
		slog.Warn("事件队列已满，跳过处理", "block", log.BlockNumber, "tx", log.TxHash.Hex())
	}
}

// worker 工作协程
func (p *EventProcessor) worker(id int) {
	defer p.wg.Done()

	slog.Debug("Worker启动", "worker_id", id)

	for {
		select {
		case <-p.ctx.Done():
			slog.Debug("Worker停止", "worker_id", id)
			return
		case log, ok := <-p.jobQueue:
			if !ok {
				slog.Debug("Worker队列已关闭", "worker_id", id)
				return
			}
			p.processLog(p.ctx, log, id)
		}
	}
}

// processLog 处理单个日志
func (p *EventProcessor) processLog(ctx context.Context, log types.Log, workerID int) {
	// 尝试解析不同类型的事件
	var err error

	if event, parseErr := p.contract.ParseAuctionCreated(log); parseErr == nil {
		err = p.handler.HandleAuctionCreated(event, log)
	} else if event, parseErr := p.contract.ParseBidPlaced(log); parseErr == nil {
		err = p.handler.HandleBidPlaced(event, log)
	} else if event, parseErr := p.contract.ParseAuctionEnded(log); parseErr == nil {
		err = p.handler.HandleAuctionEnded(event, log)
	} else if event, parseErr := p.contract.ParseBidWithdrawn(log); parseErr == nil {
		err = p.handler.HandleBidWithdrawn(event, log)
	} else {
		slog.Debug("未识别的事件类型", "worker_id", workerID, "topics", len(log.Topics), "tx", log.TxHash.Hex())
		return
	}

	if err != nil {
		slog.Error("Worker处理事件失败", "worker_id", workerID, "error", err, "block", log.BlockNumber, "tx", log.TxHash.Hex())
	} else {
		slog.Debug("Worker处理事件成功", "worker_id", workerID, "block", log.BlockNumber, "tx", log.TxHash.Hex())
	}
}
