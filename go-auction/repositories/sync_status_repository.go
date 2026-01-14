package repositories

import (
	"go-auction/config"
	"go-auction/models"

	"gorm.io/gorm"
)

// SyncStatusRepository 区块同步状态数据访问层
type SyncStatusRepository struct {
	db *gorm.DB
}

// NewSyncStatusRepository 创建同步状态仓库实例
func NewSyncStatusRepository() *SyncStatusRepository {
	return &SyncStatusRepository{db: config.GetDB()}
}

// Create 创建同步状态记录
func (r *SyncStatusRepository) Create(syncStatus *models.SyncStatus) error {
	return r.db.Create(syncStatus).Error
}

// GetByContractAddress 根据合约地址获取同步状态
func (r *SyncStatusRepository) GetByContractAddress(contractAddress string) (*models.SyncStatus, error) {
	var syncStatus models.SyncStatus
	err := r.db.Where("contract_address = ?", contractAddress).
		First(&syncStatus).Error
	if err != nil {
		return nil, err
	}
	return &syncStatus, nil
}

// GetOrCreate 获取或创建同步状态（如果不存在则创建，LastSyncedBlock 默认为 0）
func (r *SyncStatusRepository) GetOrCreate(contractAddress string) (*models.SyncStatus, error) {
	var syncStatus models.SyncStatus
	err := r.db.FirstOrCreate(&syncStatus, models.SyncStatus{
		ContractAddress: contractAddress,
		LastSyncedBlock: 0,
	}).Error
	if err != nil {
		return nil, err
	}
	return &syncStatus, nil
}

// UpdateLastSyncedBlock 更新最后同步的区块号
func (r *SyncStatusRepository) UpdateLastSyncedBlock(contractAddress string, blockNumber uint64) error {
	return r.db.Model(&models.SyncStatus{}).
		Where("contract_address = ?", contractAddress).
		Update("last_synced_block", blockNumber).Error
}
