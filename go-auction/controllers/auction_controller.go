package controllers

import (
	"errors"
	"strconv"

	"go-auction/dto"
	"go-auction/services"
	"go-auction/utils/response"

	"github.com/gin-gonic/gin"
)

// AuctionController 拍卖控制器
type AuctionController struct {
	service *services.AuctionService
}

// NewAuctionController 创建拍卖控制器实例
func NewAuctionController() *AuctionController {
	return &AuctionController{
		service: services.NewAuctionService(),
	}
}

// GetByID 获取拍卖详情
// @Summary 获取拍卖详情
// @Description 根据ID获取拍卖详细信息
// @Tags 拍卖
// @Accept json
// @Produce json
// @Param id path int true "拍卖ID"
// @Success 200 {object} response.Response{data=vo.AuctionVO} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "拍卖不存在"
// @Router /api/v1/auctions/{id} [get]
func (c *AuctionController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "无效的拍卖ID")
		return
	}

	auction, err := c.service.GetByID(id)
	if err != nil {
		if errors.Is(err, services.ErrAuctionNotFound) {
			response.NotFound(ctx, err.Error())
		} else {
			response.InternalError(ctx, "获取拍卖失败")
		}
		return
	}

	response.Success(ctx, auction)
}

// GetList 获取拍卖列表
// @Summary 获取拍卖列表
// @Description 分页获取拍卖列表，支持状态、卖家、NFT合约过滤
// @Tags 拍卖
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param status query string false "拍卖状态" Enums(active, ended, cancelled)
// @Param seller query string false "卖家地址"
// @Param nft_contract query string false "NFT合约地址"
// @Success 200 {object} response.Response{data=vo.AuctionListVO} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/auctions [get]
func (c *AuctionController) GetList(ctx *gin.Context) {
	var req dto.AuctionListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "参数验证失败: "+err.Error())
		return
	}

	listVO, err := c.service.GetList(&req)
	if err != nil {
		response.InternalError(ctx, "获取拍卖列表失败")
		return
	}

	response.Success(ctx, listVO)
}

// GetActiveAuctions 获取活跃拍卖列表
// @Summary 获取活跃拍卖列表
// @Description 获取所有活跃状态的拍卖列表
// @Tags 拍卖
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=vo.AuctionListVO} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/auctions/active [get]
func (c *AuctionController) GetActiveAuctions(ctx *gin.Context) {
	var req dto.AuctionListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "参数验证失败: "+err.Error())
		return
	}

	listVO, err := c.service.GetActiveAuctions(&req)
	if err != nil {
		response.InternalError(ctx, "获取活跃拍卖列表失败")
		return
	}

	response.Success(ctx, listVO)
}

// GetEndedAuctions 获取已结束拍卖列表
// @Summary 获取已结束拍卖列表
// @Description 获取所有已结束状态的拍卖列表
// @Tags 拍卖
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=vo.AuctionListVO} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/auctions/ended [get]
func (c *AuctionController) GetEndedAuctions(ctx *gin.Context) {
	var req dto.AuctionListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "参数验证失败: "+err.Error())
		return
	}

	listVO, err := c.service.GetEndedAuctions(&req)
	if err != nil {
		response.InternalError(ctx, "获取已结束拍卖列表失败")
		return
	}

	response.Success(ctx, listVO)
}

// GetBySeller 获取卖家的拍卖列表
// @Summary 获取卖家的拍卖列表
// @Description 获取指定卖家的所有拍卖列表
// @Tags 拍卖
// @Accept json
// @Produce json
// @Param seller path string true "卖家地址"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=vo.AuctionListVO} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/auctions/seller/{seller} [get]
func (c *AuctionController) GetBySeller(ctx *gin.Context) {
	seller := ctx.Param("seller")
	if seller == "" {
		response.BadRequest(ctx, "卖家地址不能为空")
		return
	}

	var req dto.AuctionListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "参数验证失败: "+err.Error())
		return
	}

	listVO, err := c.service.GetBySeller(seller, &req)
	if err != nil {
		response.InternalError(ctx, "获取卖家拍卖列表失败")
		return
	}

	response.Success(ctx, listVO)
}
