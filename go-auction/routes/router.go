package routes

import (
	"go-auction/controllers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter 配置路由
func SetupRouter() *gin.Engine {
	r := gin.New()

	// 全局中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Swagger 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 路由组
	apiV1 := r.Group("/api/v1")
	{
		// 初始化控制器
		auctionCtrl := controllers.NewAuctionController()
		bidCtrl := controllers.NewBidController()
		nftCtrl := controllers.NewNFTController()

		// 拍卖路由（注意：具体路由需要在参数路由之前定义）
		apiV1.GET("/auctions/active", auctionCtrl.GetActiveAuctions)
		apiV1.GET("/auctions/ended", auctionCtrl.GetEndedAuctions)
		apiV1.GET("/auctions/seller/:seller", auctionCtrl.GetBySeller)
		apiV1.GET("/auctions/:id", auctionCtrl.GetByID)
		apiV1.GET("/auctions", auctionCtrl.GetList)

		// 出价路由
		apiV1.GET("/bids/auction/:auction_id", bidCtrl.GetByAuctionID)
		apiV1.GET("/bids/bidder/:bidder", bidCtrl.GetByBidder)
		apiV1.GET("/bids/highest/:auction_id", bidCtrl.GetHighestBid)

		// NFT路由（注意：具体路由需要在参数路由之前定义）
		apiV1.GET("/nfts/by-contract", nftCtrl.GetByContractAndTokenID)
		apiV1.GET("/nfts/:id", nftCtrl.GetByID)
		apiV1.GET("/nfts", nftCtrl.GetByOwner)
	}

	return r
}
