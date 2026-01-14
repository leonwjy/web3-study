package dto

// AuctionListRequest 拍卖列表请求
type AuctionListRequest struct {
	Page        int    `form:"page" binding:"omitempty,min=1"`
	PageSize    int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Status      string `form:"status" binding:"omitempty,oneof=active ended cancelled"`
	Seller      string `form:"seller" binding:"omitempty"`
	NFTContract string `form:"nft_contract" binding:"omitempty"`
}

// AuctionDetailRequest 拍卖详情请求
type AuctionDetailRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1"`
}
