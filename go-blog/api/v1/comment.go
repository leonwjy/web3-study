package v1

import (
	"errors"
	"strconv"

	"go_blog/internal/app/dto"
	"go_blog/internal/app/middleware"
	"go_blog/internal/app/service"
	"go_blog/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// CommentAPI 评论接口
type CommentAPI struct {
	service *service.CommentService
}

// NewCommentAPI 创建评论接口实例
func NewCommentAPI() *CommentAPI {
	return &CommentAPI{
		service: service.NewCommentService(),
	}
}

// Create 创建评论
// @Summary 创建评论
// @Description 对指定文章发表评论（需要登录）
// @Tags 评论
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "文章ID"
// @Param request body dto.CreateCommentRequest true "评论内容"
// @Success 200 {object} response.Response{data=vo.CommentVO} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 404 {object} response.Response "文章不存在"
// @Router /api/v1/posts/{id}/comments [post]
func (a *CommentAPI) Create(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的文章ID")
		return
	}

	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数验证失败: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	comment, err := a.service.Create(userID, uint(postID), &req)
	if err != nil {
		if errors.Is(err, service.ErrPostNotFound) {
			response.NotFound(c, err.Error())
		} else {
			response.InternalError(c, "创建评论失败")
		}
		return
	}

	response.SuccessWithMsg(c, "评论成功", comment)
}

// GetByPostID 获取文章评论列表
// @Summary 获取文章评论列表
// @Description 分页获取指定文章的评论列表
// @Tags 评论
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=vo.CommentListVO} "获取成功"
// @Failure 404 {object} response.Response "文章不存在"
// @Router /api/v1/posts/{id}/comments [get]
func (a *CommentAPI) GetByPostID(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的文章ID")
		return
	}

	var req dto.CommentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数验证失败: "+err.Error())
		return
	}

	comments, err := a.service.GetByPostID(uint(postID), &req)
	if err != nil {
		if errors.Is(err, service.ErrPostNotFound) {
			response.NotFound(c, err.Error())
		} else {
			response.InternalError(c, "获取评论列表失败")
		}
		return
	}

	response.Success(c, comments)
}
