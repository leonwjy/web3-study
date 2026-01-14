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

// PostAPI 文章接口
type PostAPI struct {
	service *service.PostService
}

// NewPostAPI 创建文章接口实例
func NewPostAPI() *PostAPI {
	return &PostAPI{
		service: service.NewPostService(),
	}
}

// Create 创建文章
// @Summary 创建文章
// @Description 创建一篇新的博客文章（需要登录）
// @Tags 文章
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token"
// @Param request body dto.CreatePostRequest true "文章信息"
// @Success 200 {object} response.Response{data=vo.PostVO} "创建成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Router /api/v1/posts [post]
func (a *PostAPI) Create(c *gin.Context) {
	var req dto.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数验证失败: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	post, err := a.service.Create(userID, &req)
	if err != nil {
		response.InternalError(c, "创建文章失败")
		return
	}

	response.SuccessWithMsg(c, "创建成功", post)
}

// GetByID 获取文章详情
// @Summary 获取文章详情
// @Description 根据ID获取文章详细信息
// @Tags 文章
// @Accept json
// @Produce json
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response{data=vo.PostVO} "获取成功"
// @Failure 404 {object} response.Response "文章不存在"
// @Router /api/v1/posts/{id} [get]
func (a *PostAPI) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的文章ID")
		return
	}

	post, err := a.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, service.ErrPostNotFound) {
			response.NotFound(c, err.Error())
		} else {
			response.InternalError(c, "获取文章失败")
		}
		return
	}

	response.Success(c, post)
}

// GetList 获取文章列表
// @Summary 获取文章列表
// @Description 分页获取所有文章列表
// @Tags 文章
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.Response{data=vo.PostListVO} "获取成功"
// @Router /api/v1/posts [get]
func (a *PostAPI) GetList(c *gin.Context) {
	var req dto.PostListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "参数验证失败: "+err.Error())
		return
	}

	posts, err := a.service.GetList(&req)
	if err != nil {
		response.InternalError(c, "获取文章列表失败")
		return
	}

	response.Success(c, posts)
}

// Update 更新文章
// @Summary 更新文章
// @Description 更新指定文章（只有作者可以更新）
// @Tags 文章
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "文章ID"
// @Param request body dto.UpdatePostRequest true "更新内容"
// @Success 200 {object} response.Response{data=vo.PostVO} "更新成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权操作"
// @Failure 404 {object} response.Response "文章不存在"
// @Router /api/v1/posts/{id} [put]
func (a *PostAPI) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的文章ID")
		return
	}

	var req dto.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数验证失败: "+err.Error())
		return
	}

	userID := middleware.GetUserID(c)
	post, err := a.service.Update(uint(id), userID, &req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPostNotFound):
			response.NotFound(c, err.Error())
		case errors.Is(err, service.ErrNotAuthor):
			response.Forbidden(c, err.Error())
		default:
			response.InternalError(c, "更新文章失败")
		}
		return
	}

	response.SuccessWithMsg(c, "更新成功", post)
}

// Delete 删除文章
// @Summary 删除文章
// @Description 删除指定文章（只有作者可以删除）
// @Tags 文章
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param Authorization header string true "Bearer Token"
// @Param id path int true "文章ID"
// @Success 200 {object} response.Response "删除成功"
// @Failure 401 {object} response.Response "未授权"
// @Failure 403 {object} response.Response "无权操作"
// @Failure 404 {object} response.Response "文章不存在"
// @Router /api/v1/posts/{id} [delete]
func (a *PostAPI) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的文章ID")
		return
	}

	userID := middleware.GetUserID(c)
	err = a.service.Delete(uint(id), userID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPostNotFound):
			response.NotFound(c, err.Error())
		case errors.Is(err, service.ErrNotAuthor):
			response.Forbidden(c, err.Error())
		default:
			response.InternalError(c, "删除文章失败")
		}
		return
	}

	response.SuccessWithMsg(c, "删除成功", nil)
}
