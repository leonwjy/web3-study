package controllers

import (
	"errors"
	"strconv"

	"go-auction/dto"
	"go-auction/services"
	"go-auction/utils"

	"github.com/gin-gonic/gin"
)

// NFTController NFT控制器
type NFTController struct {
	service *services.NFTService
}

// NewNFTController 创建NFT控制器实例
func NewNFTController() *NFTController {
	return &NFTController{
		service: services.NewNFTService(),
	}
}

// GetByID 获取NFT详情
// @Summary 获取NFT详情
// @Description 根据ID获取NFT详细信息
// @Tags NFT
// @Accept json
// @Produce json
// @Param id path int true "NFT ID"
// @Success 200 {object} response.Response{data=vo.NFTVO} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "NFT不存在"
// @Router /api/v1/nfts/{id} [get]
func (c *NFTController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(ctx, "无效的NFT ID")
		return
	}

	nftVO, err := c.service.GetByID(id)
	if err != nil {
		if errors.Is(err, services.ErrNFTNotFound) {
			utils.NotFound(ctx, err.Error())
		} else {
			utils.InternalError(ctx, "获取NFT失败")
		}
		return
	}

	utils.Success(ctx, nftVO)
}

// GetByContractAndTokenID 根据合约和TokenID获取NFT
// @Summary 根据合约和TokenID获取NFT
// @Description 根据合约地址和TokenID获取NFT信息
// @Tags NFT
// @Accept json
// @Produce json
// @Param contract_address query string true "合约地址"
// @Param token_id query int true "Token ID"
// @Success 200 {object} response.Response{data=vo.NFTVO} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 404 {object} response.Response "NFT不存在"
// @Router /api/v1/nfts/by-contract [get]
func (c *NFTController) GetByContractAndTokenID(ctx *gin.Context) {
	var req dto.NFTByContractRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(ctx, "参数验证失败: "+err.Error())
		return
	}

	nftVO, err := c.service.GetByContractAndTokenID(req.ContractAddress, req.TokenID)
	if err != nil {
		if errors.Is(err, services.ErrNFTNotFound) {
			utils.NotFound(ctx, err.Error())
		} else {
			utils.InternalError(ctx, "获取NFT失败")
		}
		return
	}

	utils.Success(ctx, nftVO)
}

// GetByOwner 获取所有者的NFT列表
// @Summary 获取所有者的NFT列表
// @Description 获取指定所有者的所有NFT列表
// @Tags NFT
// @Accept json
// @Produce json
// @Param owner query string true "所有者地址"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=vo.NFTListVO} "获取成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/nfts [get]
func (c *NFTController) GetByOwner(ctx *gin.Context) {
	var req dto.NFTListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(ctx, "参数验证失败: "+err.Error())
		return
	}

	if req.Owner == "" {
		utils.BadRequest(ctx, "所有者地址不能为空")
		return
	}

	listVO, err := c.service.GetByOwner(req.Owner, &req)
	if err != nil {
		utils.InternalError(ctx, "获取NFT列表失败")
		return
	}

	utils.Success(ctx, listVO)
}
