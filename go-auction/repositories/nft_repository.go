package repositories

import (
	"go-auction/config"
	"go-auction/models"

	"errors"
	"gorm.io/gorm"
)

// NFTRepository NFT信息数据访问层
type NFTRepository struct {
	db *gorm.DB
}

// NewNFTRepository 创建NFT仓库实例
func NewNFTRepository() *NFTRepository {
	return &NFTRepository{db: config.GetDB()}
}

// Create 创建NFT记录
func (r *NFTRepository) Create(nft *models.NFT) error {
	return r.db.Create(nft).Error
}

// GetByID 根据ID获取NFT
func (r *NFTRepository) GetByID(id uint64) (*models.NFT, error) {
	var nft models.NFT
	err := r.db.First(&nft, id).Error
	if err != nil {
		return nil, err
	}
	return &nft, nil
}

// Update 更新NFT
func (r *NFTRepository) Update(nft *models.NFT) error {
	return r.db.Save(nft).Error
}

// Delete 删除NFT
func (r *NFTRepository) Delete(id uint64) error {
	return r.db.Delete(&models.NFT{}, id).Error
}

// GetByContractAndTokenID 根据合约地址和TokenID获取NFT（唯一查询）
// 如果记录不存在，返回 (nil, nil) 而不是错误
func (r *NFTRepository) GetByContractAndTokenID(contractAddress string, tokenID uint64) (*models.NFT, error) {
	var nft models.NFT
	err := r.db.Where("contract_address = ? AND token_id = ?", contractAddress, tokenID).
		First(&nft).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 记录不存在是正常情况，返回 nil 而不是错误
			return nil, nil
		}
		return nil, err
	}
	return &nft, nil
}

// GetByOwner 获取所有者的NFT列表
func (r *NFTRepository) GetByOwner(owner string, page, pageSize int) ([]*models.NFT, int64, error) {
	var nfts []*models.NFT
	var total int64

	// 获取总数
	err := r.db.Model(&models.NFT{}).Where("owner = ?", owner).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = r.db.Where("owner = ?", owner).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&nfts).Error
	if err != nil {
		return nil, 0, err
	}

	return nfts, total, nil
}

// UpdateOwner 更新NFT所有者
func (r *NFTRepository) UpdateOwner(contractAddress string, tokenID uint64, owner string) error {
	return r.db.Model(&models.NFT{}).
		Where("contract_address = ? AND token_id = ?", contractAddress, tokenID).
		Update("owner", owner).Error
}

// Exists 检查NFT是否存在
func (r *NFTRepository) Exists(contractAddress string, tokenID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&models.NFT{}).
		Where("contract_address = ? AND token_id = ?", contractAddress, tokenID).
		Count(&count).Error
	return count > 0, err
}

// GetNFTsWithoutMetadata 获取缺少元数据的NFT列表（image_url 或 description 为空）
func (r *NFTRepository) GetNFTsWithoutMetadata() ([]*models.NFT, error) {
	var nfts []*models.NFT
	err := r.db.Where("image_url = '' OR description = '' OR name = '' OR name LIKE 'NFT #%'").
		Find(&nfts).Error
	if err != nil {
		return nil, err
	}
	return nfts, nil
}
