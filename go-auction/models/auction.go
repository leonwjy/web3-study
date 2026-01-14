package models

import (
	"time"
)

// Auction 拍卖信息模型
type Auction struct {
	ID                uint64    `gorm:"primaryKey;comment:拍卖ID"`
	NFTContract       string    `gorm:"type:varchar(42);not null;index:idx_nft_contract_token;comment:NFT合约地址"`
	TokenID           uint64    `gorm:"not null;index:idx_nft_contract_token;comment:NFT TokenID"`
	Seller            string    `gorm:"type:varchar(42);not null;index;comment:卖家地址"`
	StartingPrice     string    `gorm:"type:decimal(36,18);not null;comment:起拍价USD"`
	CurrentHighestBid string    `gorm:"type:decimal(36,18);default:0;comment:当前最高出价USD"`
	HighestBidder     string    `gorm:"type:varchar(42);default:'';comment:当前最高出价者"`
	StartTime         uint64    `gorm:"not null;comment:拍卖开始时间"`
	EndTime           uint64    `gorm:"not null;index;comment:拍卖结束时间"`
	PaymentToken      string    `gorm:"type:varchar(42);not null;comment:支付代币地址"`
	Status            string    `gorm:"type:varchar(20);not null;default:'active';index;comment:拍卖状态"`
	CreatedAt         time.Time `gorm:"comment:创建时间"`
	UpdatedAt         time.Time `gorm:"comment:更新时间"`
}

// TableName 指定表名
func (Auction) TableName() string {
	return "auctions"
}

// AuctionStatus 拍卖状态常量
const (
	AuctionStatusActive    = "active"
	AuctionStatusEnded     = "ended"
	AuctionStatusCancelled = "cancelled"
)
