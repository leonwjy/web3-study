package blockchain

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"sync"
	"time"

	"go-auction/config"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client 区块链客户端
type Client struct {
	rpcClient    *ethclient.Client
	wsClient     *ethclient.Client
	rpcURL       string
	wsURL        string
	chainID      *big.Int
	contractAddr common.Address
	mu           sync.RWMutex
}

// NewClient 创建区块链客户端
func NewClient(cfg *config.BlockchainConfig) (*Client, error) {
	// 解析合约地址
	contractAddr := common.HexToAddress(cfg.ContractAddress)
	if contractAddr == (common.Address{}) {
		return nil, fmt.Errorf("invalid contract address: %s", cfg.ContractAddress)
	}

	// 创建 RPC 客户端
	rpcClient, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %w", err)
	}

	// 创建 WebSocket 客户端（可选，如果配置了）
	var wsClient *ethclient.Client
	if cfg.WSURL != "" {
		wsClient, err = ethclient.Dial(cfg.WSURL)
		if err != nil {
			slog.Warn("failed to connect to WebSocket, will use RPC only", "error", err)
			wsClient = nil
		}
	}

	// 获取链 ID
	chainID := big.NewInt(int64(cfg.ChainID))

	client := &Client{
		rpcClient:    rpcClient,
		wsClient:     wsClient,
		rpcURL:       cfg.RPCURL,
		wsURL:        cfg.WSURL,
		chainID:      chainID,
		contractAddr: contractAddr,
	}

	// 验证连接
	if err := client.verifyConnection(); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to verify connection: %w", err)
	}

	slog.Info("区块链客户端初始化成功", "contract", contractAddr.Hex(), "chain_id", chainID.String())
	return client, nil
}

// verifyConnection 验证连接是否正常
func (c *Client) verifyConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := c.rpcClient.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get block number: %w", err)
	}

	return nil
}

// GetRPCClient 获取 RPC 客户端
func (c *Client) GetRPCClient() *ethclient.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.rpcClient
}

// GetWSClient 获取 WebSocket 客户端（如果可用）
func (c *Client) GetWSClient() *ethclient.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.wsClient
}

// GetContractAddress 获取合约地址
func (c *Client) GetContractAddress() common.Address {
	return c.contractAddr
}

// GetChainID 获取链 ID
func (c *Client) GetChainID() *big.Int {
	return c.chainID
}

// GetLatestBlock 获取最新区块号
func (c *Client) GetLatestBlock(ctx context.Context) (uint64, error) {
	blockNumber, err := c.rpcClient.BlockNumber(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get latest block: %w", err)
	}
	return blockNumber, nil
}

// Reconnect 重新连接（用于连接断开时，带指数退避）
func (c *Client) Reconnect() error {
	return RetryWithBackoffDefault(context.Background(), func() error {
		return c.reconnectInternal()
	})
}

// reconnectInternal 内部重连逻辑
func (c *Client) reconnectInternal() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 关闭旧连接
	if c.rpcClient != nil {
		c.rpcClient.Close()
	}
	if c.wsClient != nil {
		c.wsClient.Close()
	}

	// 重新连接 RPC
	rpcClient, err := ethclient.Dial(c.rpcURL)
	if err != nil {
		return fmt.Errorf("failed to reconnect RPC: %w", err)
	}
	c.rpcClient = rpcClient

	// 重新连接 WebSocket（如果配置了）
	if c.wsURL != "" {
		wsClient, err := ethclient.Dial(c.wsURL)
		if err != nil {
			slog.Warn("failed to reconnect WebSocket", "error", err)
			c.wsClient = nil
		} else {
			c.wsClient = wsClient
		}
	}

	// 验证连接
	if err := c.verifyConnection(); err != nil {
		return fmt.Errorf("failed to verify reconnection: %w", err)
	}

	slog.Info("区块链客户端重连成功")
	return nil
}

// Close 关闭所有连接
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var errs []error
	if c.rpcClient != nil {
		c.rpcClient.Close()
		c.rpcClient = nil
	}
	if c.wsClient != nil {
		c.wsClient.Close()
		c.wsClient = nil
	}

	if len(errs) > 0 {
		return fmt.Errorf("errors closing clients: %v", errs)
	}

	slog.Info("区块链客户端已关闭")
	return nil
}
