package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"go-auction/config"
)

func main() {
	// 初始化日志
	initLogger()

	slog.Info("开始测试 API 端点...")

	// 加载配置
	cfg := config.Load()

	// 初始化数据库和Redis（确保服务可用）
	config.InitDatabase(cfg)
	defer config.Close()
	config.InitRedis(cfg)
	defer config.CloseRedis()

	// 等待服务器启动
	time.Sleep(2 * time.Second)

	baseURL := "http://localhost:8080/api/v1"

	// 测试拍卖相关API
	slog.Info("=== 测试拍卖相关API ===")
	testAuctionAPIs(baseURL)

	// 测试出价相关API
	slog.Info("=== 测试出价相关API ===")
	testBidAPIs(baseURL)

	// 测试NFT相关API
	slog.Info("=== 测试NFT相关API ===")
	testNFTAPIs(baseURL)

	slog.Info("API 测试完成！")
}

func testAuctionAPIs(baseURL string) {
	// 测试获取拍卖列表
	testGET(baseURL+"/auctions", "获取拍卖列表")

	// 测试获取活跃拍卖列表
	testGET(baseURL+"/auctions/active", "获取活跃拍卖列表")

	// 测试获取已结束拍卖列表
	testGET(baseURL+"/auctions/ended", "获取已结束拍卖列表")

	// 测试获取拍卖详情（使用ID=1）
	testGET(baseURL+"/auctions/1", "获取拍卖详情")

	// 测试获取卖家拍卖列表
	testGET(baseURL+"/auctions/seller/0xSELLER123456789012345678901234567890", "获取卖家拍卖列表")
}

func testBidAPIs(baseURL string) {
	// 测试获取拍卖的出价列表
	testGET(baseURL+"/bids/auction/1", "获取拍卖的出价列表")

	// 测试获取出价者的出价历史
	testGET(baseURL+"/bids/bidder/0xBIDDER123456789012345678901234567890", "获取出价者的出价历史")

	// 测试获取最高出价
	testGET(baseURL+"/bids/highest/1", "获取最高出价")
}

func testNFTAPIs(baseURL string) {
	// 测试获取NFT列表
	testGET(baseURL+"/nfts?owner=0xSELLER123456789012345678901234567890", "获取NFT列表")

	// 测试获取NFT详情
	testGET(baseURL+"/nfts/1", "获取NFT详情")

	// 测试根据合约和TokenID获取NFT
	testGET(baseURL+"/nfts/by-contract?contract_address=0xTESTNFT123456789012345678901234567890&token_id=1", "根据合约和TokenID获取NFT")
}

func testGET(url string, description string) {
	slog.Info("测试", "description", description, "url", url)

	resp, err := http.Get(url)
	if err != nil {
		slog.Error("请求失败", "error", err, "url", url)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("读取响应失败", "error", err)
		return
	}

	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, body, "", "  "); err != nil {
		slog.Error("格式化JSON失败", "error", err)
		return
	}

	slog.Info("响应", "status", resp.StatusCode, "body", prettyJSON.String())
	fmt.Println()
}

func initLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(handler))
}
