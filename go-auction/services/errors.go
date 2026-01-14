package services

import "errors"

// 共享错误定义
var (
	ErrAuctionNotFound = errors.New("拍卖不存在")
	ErrBidNotFound     = errors.New("出价记录不存在")
	ErrNFTNotFound     = errors.New("NFT不存在")
)
