package vo

import (
	"go-auction/models"
	"time"
)

// AuctionVO 拍卖响应数据
type AuctionVO struct {
	ID                uint64    `json:"id"`
	NFTContract       string    `json:"nft_contract"`
	TokenID           uint64    `json:"token_id"`
	Seller            string    `json:"seller"`
	StartingPrice     string    `json:"starting_price"`
	CurrentHighestBid string    `json:"current_highest_bid"`
	HighestBidder     string    `json:"highest_bidder"`
	StartTime         uint64    `json:"start_time"`
	EndTime           uint64    `json:"end_time"`
	PaymentToken      string    `json:"payment_token"`
	Status            string    `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	// 关联数据（可选）
	NFT *NFTVO `json:"nft,omitempty"`
}

// AuctionListVO 拍卖列表响应
type AuctionListVO struct {
	Total int64        `json:"total"`
	List  []*AuctionVO `json:"list"`
}

// ToAuctionVO 将 Auction 模型转换为 AuctionVO
func ToAuctionVO(auction *models.Auction) *AuctionVO {
	if auction == nil {
		return nil
	}
	return &AuctionVO{
		ID:                auction.ID,
		NFTContract:       auction.NFTContract,
		TokenID:           auction.TokenID,
		Seller:            auction.Seller,
		StartingPrice:     auction.StartingPrice,
		CurrentHighestBid: auction.CurrentHighestBid,
		HighestBidder:     auction.HighestBidder,
		StartTime:         auction.StartTime,
		EndTime:           auction.EndTime,
		PaymentToken:      auction.PaymentToken,
		Status:            auction.Status,
		CreatedAt:         auction.CreatedAt,
		UpdatedAt:         auction.UpdatedAt,
	}
}

// ToAuctionVOList 将 Auction 模型列表转换为 AuctionVO 列表
func ToAuctionVOList(auctions []*models.Auction) []*AuctionVO {
	list := make([]*AuctionVO, 0, len(auctions))
	for _, auction := range auctions {
		list = append(list, ToAuctionVO(auction))
	}
	return list
}
