package controllers

import (
	"errors"
	"strconv"

	"go-auction/dto"
	"go-auction/services"
	"go-auction/utils"

	"github.com/gin-gonic/gin"
)

// BidController 出价控制器
type BidController struct {
	service *services.BidService
}

// NewBidController 创建出价控制器实例
func NewBidController() *BidController {
	return &BidController{
		service: services.NewBidService(),
	}
}

// GetByAuctionID 获取拍卖的出价列表
// @Summary 获取拍卖的出价列表
// @Description 获取指定拍卖的所有出价记录
// @Tags 出价
// @Accept json
// @Produce json
// @Param auction_id path int true "拍卖ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=vo.BidListVO} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "拍卖不存在"
// @Router /api/v1/bids/auction/{auction_id} [get]
func (c *BidController) GetByAuctionID(ctx *gin.Context) {
	auctionID, err := strconv.ParseUint(ctx.Param("auction_id"), 10, 64)
	if err != nil {
		utils.BadRequest(ctx, "无效的拍卖ID")
		return
	}

	var req dto.BidListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(ctx, "参数验证失败: "+err.Error())
		return
	}

	req.AuctionID = auctionID
	listVO, err := c.service.GetByAuctionID(auctionID, &req)
	if err != nil {
		if errors.Is(err, services.ErrAuctionNotFound) {
			utils.NotFound(ctx, err.Error())
		} else {
			utils.InternalError(ctx, "获取出价列表失败")
		}
		return
	}

	utils.Success(ctx, listVO)
}

// GetByBidder 获取出价者的出价历史
// @Summary 获取出价者的出价历史
// @Description 获取指定出价者的所有出价记录
// @Tags 出价
// @Accept json
// @Produce json
// @Param bidder path string true "出价者地址"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=vo.BidListVO} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/bids/bidder/{bidder} [get]
func (c *BidController) GetByBidder(ctx *gin.Context) {
	bidder := ctx.Param("bidder")
	if bidder == "" {
		utils.BadRequest(ctx, "出价者地址不能为空")
		return
	}

	var req dto.BidHistoryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(ctx, "参数验证失败: "+err.Error())
		return
	}

	req.Bidder = bidder
	listVO, err := c.service.GetByBidder(bidder, &req)
	if err != nil {
		utils.InternalError(ctx, "获取出价历史失败")
		return
	}

	utils.Success(ctx, listVO)
}

// GetHighestBid 获取拍卖的最高出价
// @Summary 获取拍卖的最高出价
// @Description 获取指定拍卖的最高出价记录
// @Tags 出价
// @Accept json
// @Produce json
// @Param auction_id path int true "拍卖ID"
// @Success 200 {object} response.Response{data=vo.BidVO} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "拍卖不存在或没有出价记录"
// @Router /api/v1/bids/highest/{auction_id} [get]
func (c *BidController) GetHighestBid(ctx *gin.Context) {
	auctionID, err := strconv.ParseUint(ctx.Param("auction_id"), 10, 64)
	if err != nil {
		utils.BadRequest(ctx, "无效的拍卖ID")
		return
	}

	bidVO, err := c.service.GetHighestBid(auctionID)
	if err != nil {
		if errors.Is(err, services.ErrAuctionNotFound) {
			utils.NotFound(ctx, err.Error())
		} else if errors.Is(err, services.ErrBidNotFound) {
			utils.NotFound(ctx, err.Error())
		} else {
			utils.InternalError(ctx, "获取最高出价失败")
		}
		return
	}

	utils.Success(ctx, bidVO)
}
