package services

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"go-auction/dto"
	"go-auction/repositories"
	"go-auction/utils"
	"go-auction/vo"

	"gorm.io/gorm"
)

// NFTService NFT服务
type NFTService struct {
	repo *repositories.NFTRepository
}

// NewNFTService 创建NFT服务实例
func NewNFTService() *NFTService {
	return &NFTService{
		repo: repositories.NewNFTRepository(),
	}
}

// GetByID 根据ID获取NFT
func (s *NFTService) GetByID(id uint64) (*vo.NFTVO, error) {
	ctx := context.Background()
	cacheKey := utils.CacheKey("nft", id)

	// 尝试从缓存获取
	var nftVO vo.NFTVO
	if err := utils.CacheGet(ctx, cacheKey, &nftVO); err == nil {
		slog.Debug("从缓存获取NFT详情", "id", id)
		return &nftVO, nil
	}

	// 缓存未命中，从数据库查询
	nft, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNFTNotFound
		}
		slog.Error("查询NFT失败", "error", err, "id", id)
		return nil, err
	}

	nftVO = *vo.ToNFTVO(nft)

	// 存入缓存（TTL: 10分钟）
	if err := utils.CacheSet(ctx, cacheKey, nftVO, 10*time.Minute); err != nil {
		slog.Warn("缓存NFT详情失败", "error", err, "id", id)
	}

	return &nftVO, nil
}

// GetByContractAndTokenID 根据合约地址和TokenID获取NFT
func (s *NFTService) GetByContractAndTokenID(contractAddress string, tokenID uint64) (*vo.NFTVO, error) {
	ctx := context.Background()
	cacheKey := utils.CacheKey("nft", contractAddress, tokenID)

	// 尝试从缓存获取
	var nftVO vo.NFTVO
	if err := utils.CacheGet(ctx, cacheKey, &nftVO); err == nil {
		slog.Debug("从缓存获取NFT详情", "contract", contractAddress, "token_id", tokenID)
		return &nftVO, nil
	}

	// 缓存未命中，从数据库查询
	nft, err := s.repo.GetByContractAndTokenID(contractAddress, tokenID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNFTNotFound
		}
		slog.Error("查询NFT失败", "error", err, "contract", contractAddress, "token_id", tokenID)
		return nil, err
	}

	nftVO = *vo.ToNFTVO(nft)

	// 存入缓存（TTL: 10分钟）
	if err := utils.CacheSet(ctx, cacheKey, nftVO, 10*time.Minute); err != nil {
		slog.Warn("缓存NFT详情失败", "error", err, "contract", contractAddress, "token_id", tokenID)
	}

	return &nftVO, nil
}

// GetByOwner 获取所有者的NFT列表
func (s *NFTService) GetByOwner(owner string, req *dto.NFTListRequest) (*vo.NFTListVO, error) {
	// 默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	nfts, total, err := s.repo.GetByOwner(owner, page, pageSize)
	if err != nil {
		slog.Error("获取NFT列表失败", "error", err, "owner", owner)
		return nil, err
	}

	return &vo.NFTListVO{
		Total: total,
		List:  vo.ToNFTVOList(nfts),
	}, nil
}
