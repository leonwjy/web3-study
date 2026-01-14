package services

import (
	"errors"
	"log/slog"

	"go-auction/dto"
	"go-auction/repositories"
	"go-auction/vo"

	"gorm.io/gorm"
)

// BidService 出价服务
type BidService struct {
	repo        *repositories.BidRepository
	auctionRepo *repositories.AuctionRepository
}

// NewBidService 创建出价服务实例
func NewBidService() *BidService {
	return &BidService{
		repo:        repositories.NewBidRepository(),
		auctionRepo: repositories.NewAuctionRepository(),
	}
}

// GetByAuctionID 获取拍卖的出价列表
func (s *BidService) GetByAuctionID(auctionID uint64, req *dto.BidListRequest) (*vo.BidListVO, error) {
	// 验证拍卖是否存在
	exists, err := s.auctionRepo.Exists(auctionID)
	if err != nil {
		slog.Error("检查拍卖失败", "error", err, "auction_id", auctionID)
		return nil, err
	}
	if !exists {
		return nil, ErrAuctionNotFound
	}

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

	bids, total, err := s.repo.GetByAuctionID(auctionID, page, pageSize)
	if err != nil {
		slog.Error("获取出价列表失败", "error", err, "auction_id", auctionID)
		return nil, err
	}

	return &vo.BidListVO{
		Total: total,
		List:  vo.ToBidVOList(bids),
	}, nil
}

// GetByBidder 获取出价者的出价历史
func (s *BidService) GetByBidder(bidder string, req *dto.BidHistoryRequest) (*vo.BidListVO, error) {
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

	bids, total, err := s.repo.GetByBidder(bidder, page, pageSize)
	if err != nil {
		slog.Error("获取出价历史失败", "error", err, "bidder", bidder)
		return nil, err
	}

	return &vo.BidListVO{
		Total: total,
		List:  vo.ToBidVOList(bids),
	}, nil
}

// GetHighestBid 获取拍卖的最高出价
func (s *BidService) GetHighestBid(auctionID uint64) (*vo.BidVO, error) {
	// 验证拍卖是否存在
	exists, err := s.auctionRepo.Exists(auctionID)
	if err != nil {
		slog.Error("检查拍卖失败", "error", err, "auction_id", auctionID)
		return nil, err
	}
	if !exists {
		return nil, ErrAuctionNotFound
	}

	bid, err := s.repo.GetHighestBid(auctionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrBidNotFound
		}
		slog.Error("获取最高出价失败", "error", err, "auction_id", auctionID)
		return nil, err
	}

	return vo.ToBidVO(bid), nil
}
