package models

import (
	"time"
)

// Bid 出价记录模型
type Bid struct {
	ID             uint64    `gorm:"primaryKey;comment:出价记录ID"`
	AuctionID      uint64    `gorm:"not null;index;comment:拍卖ID"`
	Bidder         string    `gorm:"type:varchar(42);not null;index;comment:出价者地址"`
	BidAmountUSD   string    `gorm:"type:decimal(36,18);not null;comment:出价金额USD"`
	OriginalAmount string    `gorm:"type:decimal(36,18);not null;comment:原始出价金额"`
	PaymentToken   string    `gorm:"type:varchar(42);not null;comment:支付代币地址"`
	BlockNumber    uint64    `gorm:"not null;uniqueIndex:idx_event_unique;comment:区块号"`
	TxHash         string    `gorm:"type:varchar(66);not null;uniqueIndex:idx_event_unique;comment:交易哈希"`
	LogIndex       uint      `gorm:"not null;uniqueIndex:idx_event_unique;comment:日志索引"`
	CreatedAt      time.Time `gorm:"comment:创建时间"`
}

// TableName 指定表名
func (Bid) TableName() string {
	return "bids"
}
