package vo

import (
	"go-auction/models"
	"time"
)

// NFTVO NFT响应数据
type NFTVO struct {
	ID              uint64    `json:"id"`
	ContractAddress string    `json:"contract_address"`
	TokenID         uint64    `json:"token_id"`
	Name            string    `json:"name"`
	ImageURL        string    `json:"image_url"`
	Description     string    `json:"description"`
	Owner           string    `json:"owner"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// NFTListVO NFT列表响应
type NFTListVO struct {
	Total int64    `json:"total"`
	List  []*NFTVO `json:"list"`
}

// ToNFTVO 将 NFT 模型转换为 NFTVO
func ToNFTVO(nft *models.NFT) *NFTVO {
	if nft == nil {
		return nil
	}
	return &NFTVO{
		ID:              nft.ID,
		ContractAddress: nft.ContractAddress,
		TokenID:         nft.TokenID,
		Name:            nft.Name,
		ImageURL:        nft.ImageURL,
		Description:     nft.Description,
		Owner:           nft.Owner,
		CreatedAt:       nft.CreatedAt,
		UpdatedAt:       nft.UpdatedAt,
	}
}

// ToNFTVOList 将 NFT 模型列表转换为 NFTVO 列表
func ToNFTVOList(nfts []*models.NFT) []*NFTVO {
	list := make([]*NFTVO, 0, len(nfts))
	for _, nft := range nfts {
		list = append(list, ToNFTVO(nft))
	}
	return list
}
