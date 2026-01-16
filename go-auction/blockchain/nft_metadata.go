package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"math/big"
	"net/http"
	"strings"
	"time"

	"go-auction/config"
	"go-auction/models"
	"go-auction/repositories"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ERC721Metadata NFT 元数据结构（符合 ERC721 Metadata 标准）
type ERC721Metadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
	// 其他可选字段
	Attributes []map[string]interface{} `json:"attributes,omitempty"`
	ExternalURL string                  `json:"external_url,omitempty"`
}

// NFTMetadataService NFT 元数据获取服务
type NFTMetadataService struct {
	client        *ethclient.Client
	httpClient    *http.Client
	ipfsGateway   string
	arweaveGateway string
	userAgent     string
	updateTimeout time.Duration
}

// NewNFTMetadataService 创建 NFT 元数据服务
func NewNFTMetadataService(client *ethclient.Client, cfg *config.NFTMetadataConfig) *NFTMetadataService {
	return &NFTMetadataService{
		client:         client,
		ipfsGateway:    cfg.IPFSGateway,
		arweaveGateway: cfg.ArweaveGateway,
		userAgent:      cfg.UserAgent,
		updateTimeout:  cfg.UpdateTimeout,
		httpClient: &http.Client{
			Timeout: cfg.HTTPTimeout,
		},
	}
}

// GetMetadata 获取 NFT 元数据
// 流程：1. 调用合约的 tokenURI() 方法 2. 从 URI 获取 JSON 3. 解析元数据
func (s *NFTMetadataService) GetMetadata(ctx context.Context, contractAddress common.Address, tokenID *big.Int) (*ERC721Metadata, error) {
	// 1. 调用合约的 tokenURI() 方法
	tokenURI, err := s.getTokenURI(ctx, contractAddress, tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tokenURI: %w", err)
	}

	// tokenURI 为空表示 URI 未设置，这是正常情况（不是错误）
	if tokenURI == "" {
		return nil, fmt.Errorf("tokenURI is empty (URI not set for this token)")
	}

	// 2. 从 URI 获取 JSON 元数据
	metadataJSON, contentType, err := s.fetchMetadataFromURI(ctx, tokenURI)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata from URI: %w", err)
	}

	// 3. 验证 Content-Type，必须是 JSON
	if !strings.Contains(contentType, "application/json") && !strings.Contains(contentType, "text/json") {
		slog.Warn("tokenURI returned non-JSON content",
			"contract", contractAddress.Hex(),
			"token_id", tokenID.String(),
			"token_uri", tokenURI,
			"content_type", contentType)
		return nil, fmt.Errorf("tokenURI returned non-JSON content (Content-Type: %s)", contentType)
	}

	// 4. 解析 JSON
	var metadata ERC721Metadata
	if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
		slog.Warn("failed to parse metadata JSON",
			"contract", contractAddress.Hex(),
			"token_id", tokenID.String(),
			"token_uri", tokenURI,
			"error", err)
		return nil, fmt.Errorf("failed to parse metadata JSON: %w", err)
	}

	return &metadata, nil
}

// getTokenOwner 调用 ERC721 合约的 ownerOf() 方法检查 token 是否存在
func (s *NFTMetadataService) getTokenOwner(ctx context.Context, contractAddress common.Address, tokenID *big.Int) (common.Address, error) {
	// ERC721 ownerOf 方法签名: ownerOf(uint256) returns (address)
	// Method ID: 0x6352211e
	methodID := common.Hex2Bytes("6352211e")

	// 编码参数 (uint256 tokenId)
	paramType, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to create param type: %w", err)
	}

	arguments := abi.Arguments{
		{Type: paramType},
	}

	encoded, err := arguments.Pack(tokenID)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to pack arguments: %w", err)
	}

	// 构造调用数据
	data := append(methodID, encoded...)

	// 调用合约
	callMsg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}

	result, err := s.client.CallContract(ctx, callMsg, nil)
	if err != nil {
		// 如果 revert，说明 token 不存在
		if strings.Contains(err.Error(), "revert") || strings.Contains(err.Error(), "execution reverted") {
			return common.Address{}, nil // 返回零地址表示不存在
		}
		return common.Address{}, fmt.Errorf("failed to call ownerOf: %w", err)
	}

	if len(result) == 0 {
		return common.Address{}, fmt.Errorf("empty result from ownerOf")
	}

	// 解码返回值 (address)
	returnType, err := abi.NewType("address", "", nil)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to create return type: %w", err)
	}

	returnArgs := abi.Arguments{
		{Type: returnType},
	}

	decoded, err := returnArgs.Unpack(result)
	if err != nil {
		return common.Address{}, fmt.Errorf("failed to unpack result: %w", err)
	}

	if len(decoded) == 0 {
		return common.Address{}, fmt.Errorf("empty decoded result")
	}

	// 尝试多种类型断言，因为 abi 解码可能返回不同类型的值
	var owner common.Address
	switch v := decoded[0].(type) {
	case common.Address:
		owner = v
	case [20]byte:
		owner = common.BytesToAddress(v[:])
	case []byte:
		if len(v) >= 20 {
			owner = common.BytesToAddress(v[len(v)-20:])
		} else {
			return common.Address{}, fmt.Errorf("invalid address length: %d", len(v))
		}
	default:
		return common.Address{}, fmt.Errorf("invalid return type: %T, expected common.Address", decoded[0])
	}

	return owner, nil
}

// getTokenURI 调用 ERC721 合约的 tokenURI() 方法
func (s *NFTMetadataService) getTokenURI(ctx context.Context, contractAddress common.Address, tokenID *big.Int) (string, error) {
	// 注意：不再预先检查 ownerOf，因为：
	// 1. 某些 RPC 节点可能对 ownerOf 调用返回错误（即使 token 存在）
	// 2. 直接调用 tokenURI 更简单，如果 token 不存在会 revert，我们捕获错误即可
	// 3. NFT 可能在拍卖合约中，ownerOf 应该返回拍卖合约地址，但某些实现可能有问题

	// ERC721 tokenURI 方法签名: tokenURI(uint256) returns (string)
	// Method ID: 0xc87b56dd
	methodID := common.Hex2Bytes("c87b56dd")

	// 编码参数 (uint256 tokenId)
	paramType, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create param type: %w", err)
	}

	arguments := abi.Arguments{
		{Type: paramType},
	}

	encoded, err := arguments.Pack(tokenID)
	if err != nil {
		return "", fmt.Errorf("failed to pack arguments: %w", err)
	}

	// 构造调用数据
	data := append(methodID, encoded...)

	// 调用合约
	callMsg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}

	result, err := s.client.CallContract(ctx, callMsg, nil)
	if err != nil {
		// 检查是否是 revert 错误
		if strings.Contains(err.Error(), "revert") || strings.Contains(err.Error(), "execution reverted") {
			return "", nil // 返回空字符串，表示 URI 未设置或 token 不存在
		}
		return "", fmt.Errorf("failed to call contract: %w", err)
	}

	// 解码返回值 (string)
	returnType, err := abi.NewType("string", "", nil)
	if err != nil {
		return "", fmt.Errorf("failed to create return type: %w", err)
	}

	returnArgs := abi.Arguments{
		{Type: returnType},
	}

	decoded, err := returnArgs.Unpack(result)
	if err != nil {
		return "", fmt.Errorf("failed to unpack result: %w", err)
	}

	if len(decoded) == 0 {
		return "", fmt.Errorf("empty result")
	}

	tokenURI, ok := decoded[0].(string)
	if !ok {
		return "", fmt.Errorf("invalid return type")
	}

	return tokenURI, nil
}

// fetchMetadataFromURI 从 URI 获取元数据 JSON
// 支持 IPFS (ipfs://) 和 HTTP/HTTPS
// 返回：响应体、Content-Type、错误
func (s *NFTMetadataService) fetchMetadataFromURI(ctx context.Context, uri string) ([]byte, string, error) {
	// 处理 IPFS URI
	if strings.HasPrefix(uri, "ipfs://") {
		ipfsHash := strings.TrimPrefix(uri, "ipfs://")
		uri = fmt.Sprintf("%s/%s", strings.TrimSuffix(s.ipfsGateway, "/"), ipfsHash)
	} else if strings.HasPrefix(uri, "ar://") {
		// 处理 Arweave URI
		arHash := strings.TrimPrefix(uri, "ar://")
		uri = fmt.Sprintf("%s/%s", strings.TrimSuffix(s.arweaveGateway, "/"), arHash)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", s.userAgent)
	req.Header.Set("Accept", "application/json")

	// 发送请求
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to fetch metadata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Warn("unexpected HTTP status",
			"uri", uri,
			"status_code", resp.StatusCode)
		return nil, "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read response: %w", err)
	}

	return body, contentType, nil
}

// UpdateNFTMetadata 更新 NFT 元数据（异步）
// 在事件处理时调用，异步获取并更新元数据
func (s *NFTMetadataService) UpdateNFTMetadata(ctx context.Context, nftRepo NFTRepository, contractAddress string, tokenID uint64) {
	go func() {
		// 使用独立的 context，避免被取消
		updateCtx, cancel := context.WithTimeout(context.Background(), s.updateTimeout)
		defer cancel()

		contractAddr := common.HexToAddress(contractAddress)
		tokenIDBig := big.NewInt(int64(tokenID))

		// 获取元数据
		metadata, err := s.GetMetadata(updateCtx, contractAddr, tokenIDBig)
		if err != nil {
			slog.Warn("failed to fetch NFT metadata",
				"contract", contractAddress,
				"token_id", tokenID,
				"error", err)
			return
		}

		// 更新数据库
		nft, err := nftRepo.GetByContractAndTokenID(contractAddress, tokenID)
		if err != nil {
			slog.Warn("failed to get NFT for metadata update",
				"contract", contractAddress,
				"token_id", tokenID,
				"error", err)
			return
		}

		if nft == nil {
			slog.Warn("NFT not found for metadata update",
				"contract", contractAddress,
				"token_id", tokenID)
			return
		}

		// 更新元数据字段
		updated := false
		// 更新名称（如果为空或者是默认格式 "NFT #X"）
		if metadata.Name != "" && (nft.Name == "" || strings.HasPrefix(nft.Name, "NFT #")) {
			nft.Name = metadata.Name
			updated = true
		}
		// 更新图片URL（如果为空）
		if metadata.Image != "" && nft.ImageURL == "" {
			nft.ImageURL = metadata.Image
			updated = true
		}
		// 更新描述（如果为空）
		if metadata.Description != "" && nft.Description == "" {
			nft.Description = metadata.Description
			updated = true
		}

		if updated {
			if err := nftRepo.Update(nft); err != nil {
				slog.Warn("failed to update NFT metadata",
					"contract", contractAddress,
					"token_id", tokenID,
					"error", err)
				return
			}

			slog.Info("NFT metadata updated",
				"contract", contractAddress,
				"token_id", tokenID,
				"name", metadata.Name,
				"has_image", metadata.Image != "")
		}
	}()
}

// NFTRepository 接口，用于解耦
type NFTRepository interface {
	GetByContractAndTokenID(contractAddress string, tokenID uint64) (*NFTModel, error)
	Update(nft *NFTModel) error
}

// NFTModel NFT 模型接口
type NFTModel struct {
	ID              uint64
	ContractAddress string
	TokenID         uint64
	Name            string
	ImageURL        string
	Description     string
	Owner           string
}

// nftRepositoryAdapter 适配器，将 repositories.NFTRepository 适配为 NFTRepository 接口
type nftRepositoryAdapter struct {
	repo repositories.NFTRepository
}

func (a *nftRepositoryAdapter) GetByContractAndTokenID(contractAddress string, tokenID uint64) (*NFTModel, error) {
	nft, err := a.repo.GetByContractAndTokenID(contractAddress, tokenID)
	if err != nil {
		return nil, err
	}
	if nft == nil {
		return nil, nil
	}
	return &NFTModel{
		ID:              nft.ID,
		ContractAddress: nft.ContractAddress,
		TokenID:         nft.TokenID,
		Name:            nft.Name,
		ImageURL:        nft.ImageURL,
		Description:     nft.Description,
		Owner:           nft.Owner,
	}, nil
}

func (a *nftRepositoryAdapter) Update(nft *NFTModel) error {
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
