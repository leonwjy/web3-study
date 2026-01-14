package repositories

import (
	"go-auction/config"
	"go-auction/models"

	"gorm.io/gorm"
)

// AuctionRepository 拍卖数据访问层
type AuctionRepository struct {
	db *gorm.DB
}

// NewAuctionRepository 创建拍卖仓库实例
func NewAuctionRepository() *AuctionRepository {
	return &AuctionRepository{db: config.GetDB()}
}

// Create 创建拍卖
func (r *AuctionRepository) Create(auction *models.Auction) error {
	return r.db.Create(auction).Error
}

// GetByID 根据ID获取拍卖
func (r *AuctionRepository) GetByID(id uint64) (*models.Auction, error) {
	var auction models.Auction
	err := r.db.First(&auction, id).Error
	if err != nil {
		return nil, err
	}
	return &auction, nil
}

// Update 更新拍卖
func (r *AuctionRepository) Update(auction *models.Auction) error {
	return r.db.Save(auction).Error
}

// Delete 删除拍卖
func (r *AuctionRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Auction{}, id).Error
}

// Exists 检查拍卖是否存在
func (r *AuctionRepository) Exists(id uint64) (bool, error) {
	var count int64
	err := r.db.Model(&models.Auction{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// GetList 获取拍卖列表（分页，支持过滤）
func (r *AuctionRepository) GetList(page, pageSize int, status, seller, nftContract string) ([]*models.Auction, int64, error) {
	var auctions []*models.Auction
	var total int64

	query := r.db.Model(&models.Auction{})

	// 应用过滤条件
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if seller != "" {
		query = query.Where("seller = ?", seller)
	}
	if nftContract != "" {
		query = query.Where("nft_contract = ?", nftContract)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = query.Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&auctions).Error
	if err != nil {
		return nil, 0, err
	}

	return auctions, total, nil
}

// GetBySeller 获取卖家的拍卖列表
func (r *AuctionRepository) GetBySeller(seller string, page, pageSize int) ([]*models.Auction, int64, error) {
	var auctions []*models.Auction
	var total int64

	// 获取总数
	err := r.db.Model(&models.Auction{}).Where("seller = ?", seller).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = r.db.Where("seller = ?", seller).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&auctions).Error
	if err != nil {
		return nil, 0, err
	}

	return auctions, total, nil
}

// GetByNFT 根据NFT获取拍卖（唯一查询）
func (r *AuctionRepository) GetByNFT(nftContract string, tokenID uint64) (*models.Auction, error) {
	var auction models.Auction
	err := r.db.Where("nft_contract = ? AND token_id = ?", nftContract, tokenID).
		First(&auction).Error
	if err != nil {
		return nil, err
	}
	return &auction, nil
}

// GetActiveAuctions 获取活跃拍卖列表
func (r *AuctionRepository) GetActiveAuctions(page, pageSize int) ([]*models.Auction, int64, error) {
	return r.GetList(page, pageSize, models.AuctionStatusActive, "", "")
}

// GetEndedAuctions 获取已结束拍卖列表
func (r *AuctionRepository) GetEndedAuctions(page, pageSize int) ([]*models.Auction, int64, error) {
	return r.GetList(page, pageSize, models.AuctionStatusEnded, "", "")
}

// UpdateStatus 更新拍卖状态
func (r *AuctionRepository) UpdateStatus(id uint64, status string) error {
	return r.db.Model(&models.Auction{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// UpdateHighestBid 更新最高出价信息
func (r *AuctionRepository) UpdateHighestBid(id uint64, bidder string, amount string) error {
	return r.db.Model(&models.Auction{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"highest_bidder":      bidder,
			"current_highest_bid": amount,
		}).Error
}
