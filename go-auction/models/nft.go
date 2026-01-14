package models

import (
	"time"
)

// NFT NFT信息模型
type NFT struct {
	ID              uint64    `gorm:"primaryKey;comment:NFT ID"`
	ContractAddress string    `gorm:"type:varchar(42);not null;uniqueIndex:idx_nft_unique;comment:NFT合约地址"`
	TokenID         uint64    `gorm:"not null;uniqueIndex:idx_nft_unique;comment:Token ID"`
	Name            string    `gorm:"type:varchar(255);comment:NFT名称"`
	ImageURL        string    `gorm:"type:varchar(512);comment:图片URL"`
	Description     string    `gorm:"type:text;comment:描述"`
	Owner           string    `gorm:"type:varchar(42);index;comment:当前所有者"`
	CreatedAt       time.Time `gorm:"comment:创建时间"`
	UpdatedAt       time.Time `gorm:"comment:更新时间"`
}

// TableName 指定表名
func (NFT) TableName() string {
	return "nfts"
}
