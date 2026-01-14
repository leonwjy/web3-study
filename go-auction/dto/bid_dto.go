package dto

// BidListRequest 出价列表请求
type BidListRequest struct {
	AuctionID uint64 `form:"auction_id" binding:"required,min=1"`
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PageSize  int    `form:"page_size" binding:"omitempty,min=1,max=100"`
}

// BidHistoryRequest 出价历史请求
type BidHistoryRequest struct {
	Bidder   string `form:"bidder" binding:"required"`
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
}

// HighestBidRequest 最高出价请求
type HighestBidRequest struct {
	AuctionID uint64 `uri:"auction_id" binding:"required,min=1"`
}
