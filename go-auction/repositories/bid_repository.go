package repositories

import (
	"go-auction/config"
	"go-auction/models"

	"gorm.io/gorm"
)

// BidRepository 出价记录数据访问层
type BidRepository struct {
	db *gorm.DB
}

// NewBidRepository 创建出价记录仓库实例
func NewBidRepository() *BidRepository {
	return &BidRepository{db: config.GetDB()}
}

// Create 创建出价记录
func (r *BidRepository) Create(bid *models.Bid) error {
	return r.db.Create(bid).Error
}

// GetByID 根据ID获取出价记录
func (r *BidRepository) GetByID(id uint64) (*models.Bid, error) {
	var bid models.Bid
	err := r.db.First(&bid, id).Error
	if err != nil {
		return nil, err
	}
	return &bid, nil
}

// GetByAuctionID 获取拍卖的出价列表（分页，按时间倒序）
func (r *BidRepository) GetByAuctionID(auctionID uint64, page, pageSize int) ([]*models.Bid, int64, error) {
	var bids []*models.Bid
	var total int64

	// 获取总数
	err := r.db.Model(&models.Bid{}).Where("auction_id = ?", auctionID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = r.db.Where("auction_id = ?", auctionID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&bids).Error
	if err != nil {
		return nil, 0, err
	}

	return bids, total, nil
}

// GetByBidder 获取出价者的出价列表
func (r *BidRepository) GetByBidder(bidder string, page, pageSize int) ([]*models.Bid, int64, error) {
	var bids []*models.Bid
	var total int64

	// 获取总数
	err := r.db.Model(&models.Bid{}).Where("bidder = ?", bidder).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = r.db.Where("bidder = ?", bidder).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&bids).Error
	if err != nil {
		return nil, 0, err
	}

	return bids, total, nil
}

// GetHighestBid 获取拍卖的最高出价
func (r *BidRepository) GetHighestBid(auctionID uint64) (*models.Bid, error) {
	var bid models.Bid
	err := r.db.Where("auction_id = ?", auctionID).
		Order("bid_amount_usd DESC").
		First(&bid).Error
	if err != nil {
		return nil, err
	}
	return &bid, nil
}

// ExistsByEvent 检查事件是否已处理（用于去重）
func (r *BidRepository) ExistsByEvent(blockNumber uint64, txHash string, logIndex uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Bid{}).
		Where("block_number = ? AND tx_hash = ? AND log_index = ?", blockNumber, txHash, logIndex).
		Count(&count).Error
	return count > 0, err
}
