package models

import (
	"time"
)

// SyncStatus 区块同步状态模型
type SyncStatus struct {
	ID              uint64    `gorm:"primaryKey;comment:记录ID"`
	ContractAddress string    `gorm:"type:varchar(42);not null;uniqueIndex;comment:合约地址"`
	LastSyncedBlock uint64    `gorm:"not null;default:0;comment:最后同步的区块号"`
	UpdatedAt       time.Time `gorm:"comment:更新时间"`
}

// TableName 指定表名
func (SyncStatus) TableName() string {
	return "sync_status"
}
