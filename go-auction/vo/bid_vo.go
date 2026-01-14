package vo

import (
	"go-auction/models"
	"time"
)

// BidVO 出价记录响应数据
type BidVO struct {
	ID             uint64    `json:"id"`
	AuctionID      uint64    `json:"auction_id"`
	Bidder         string    `json:"bidder"`
	BidAmountUSD   string    `json:"bid_amount_usd"`
	OriginalAmount string    `json:"original_amount"`
	PaymentToken   string    `json:"payment_token"`
	BlockNumber    uint64    `json:"block_number"`
	TxHash         string    `json:"tx_hash"`
	LogIndex       uint      `json:"log_index"`
	CreatedAt      time.Time `json:"created_at"`
}

// BidListVO 出价列表响应
type BidListVO struct {
	Total int64    `json:"total"`
	List  []*BidVO `json:"list"`
}

// ToBidVO 将 Bid 模型转换为 BidVO
func ToBidVO(bid *models.Bid) *BidVO {
	if bid == nil {
		return nil
	}
	return &BidVO{
		ID:             bid.ID,
		AuctionID:      bid.AuctionID,
		Bidder:         bid.Bidder,
		BidAmountUSD:   bid.BidAmountUSD,
		OriginalAmount: bid.OriginalAmount,
		PaymentToken:   bid.PaymentToken,
		BlockNumber:    bid.BlockNumber,
		TxHash:         bid.TxHash,
		LogIndex:       bid.LogIndex,
		CreatedAt:      bid.CreatedAt,
	}
}

// ToBidVOList 将 Bid 模型列表转换为 BidVO 列表
func ToBidVOList(bids []*models.Bid) []*BidVO {
	list := make([]*BidVO, 0, len(bids))
	for _, bid := range bids {
		list = append(list, ToBidVO(bid))
	}
	return list
}
