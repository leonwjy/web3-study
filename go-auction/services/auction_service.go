package services

import (
	"errors"
	"log/slog"

	"go-auction/dto"
	"go-auction/models"
	"go-auction/repositories"
	"go-auction/vo"

	"gorm.io/gorm"
)

// AuctionService 拍卖服务
type AuctionService struct {
	repo    *repositories.AuctionRepository
	bidRepo *repositories.BidRepository
	nftRepo *repositories.NFTRepository
}

// NewAuctionService 创建拍卖服务实例
func NewAuctionService() *AuctionService {
	return &AuctionService{
		repo:    repositories.NewAuctionRepository(),
		bidRepo: repositories.NewBidRepository(),
		nftRepo: repositories.NewNFTRepository(),
	}
}

// GetByID 根据ID获取拍卖详情
func (s *AuctionService) GetByID(id uint64) (*vo.AuctionVO, error) {
	auction, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAuctionNotFound
		}
		slog.Error("查询拍卖失败", "error", err, "id", id)
		return nil, err
	}

	auctionVO := vo.ToAuctionVO(auction)

	// 可选：加载关联的 NFT 信息
	nft, err := s.nftRepo.GetByContractAndTokenID(auction.NFTContract, auction.TokenID)
	if err == nil && nft != nil {
		auctionVO.NFT = vo.ToNFTVO(nft)
	}

	return auctionVO, nil
}

// GetList 获取拍卖列表（支持状态、卖家、NFT合约过滤）
func (s *AuctionService) GetList(req *dto.AuctionListRequest) (*vo.AuctionListVO, error) {
	// 默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 20 // 默认每页20条
	}
	if pageSize > 100 {
		pageSize = 100 // 最大100条
	}

	auctions, total, err := s.repo.GetList(page, pageSize, req.Status, req.Seller, req.NFTContract)
	if err != nil {
		slog.Error("获取拍卖列表失败", "error", err)
		return nil, err
	}

	return &vo.AuctionListVO{
		Total: total,
		List:  vo.ToAuctionVOList(auctions),
	}, nil
}

// GetActiveAuctions 获取活跃拍卖列表
func (s *AuctionService) GetActiveAuctions(req *dto.AuctionListRequest) (*vo.AuctionListVO, error) {
	// 设置状态为 active
	req.Status = models.AuctionStatusActive
	return s.GetList(req)
}

// GetEndedAuctions 获取已结束拍卖列表
func (s *AuctionService) GetEndedAuctions(req *dto.AuctionListRequest) (*vo.AuctionListVO, error) {
	// 设置状态为 ended
	req.Status = models.AuctionStatusEnded
	return s.GetList(req)
}

// GetBySeller 获取卖家的拍卖列表
func (s *AuctionService) GetBySeller(seller string, req *dto.AuctionListRequest) (*vo.AuctionListVO, error) {
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

	auctions, total, err := s.repo.GetBySeller(seller, page, pageSize)
	if err != nil {
		slog.Error("获取卖家拍卖列表失败", "error", err, "seller", seller)
		return nil, err
	}

	return &vo.AuctionListVO{
		Total: total,
		List:  vo.ToAuctionVOList(auctions),
	}, nil
}
