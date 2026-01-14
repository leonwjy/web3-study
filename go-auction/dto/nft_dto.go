package dto

// NFTListRequest NFT列表请求
type NFTListRequest struct {
	Owner    string `form:"owner" binding:"omitempty"`
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
}

// NFTDetailRequest NFT详情请求
type NFTDetailRequest struct {
	ID uint64 `uri:"id" binding:"required,min=1"`
}

// NFTByContractRequest 根据合约和TokenID获取NFT请求
type NFTByContractRequest struct {
	ContractAddress string `form:"contract_address" binding:"required"`
	TokenID         uint64 `form:"token_id" binding:"required,min=0"`
}
