package v1

import (
	"errors"

	"go_blog/internal/app/dto"
	"go_blog/internal/app/service"
	"go_blog/internal/pkg/response"

	"github.com/gin-gonic/gin"
)

// UserAPI 用户接口
type UserAPI struct {
	service *service.UserService
}

// NewUserAPI 创建用户接口实例
func NewUserAPI() *UserAPI {
	return &UserAPI{
		service: service.NewUserService(),
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 创建新用户账号
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "注册信息"
// @Success 200 {object} response.Response{data=vo.UserVO} "注册成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Router /api/v1/register [post]
func (a *UserAPI) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数验证失败: "+err.Error())
		return
	}

	user, err := a.service.Register(&req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUsernameExists):
			response.BadRequest(c, err.Error())
		case errors.Is(err, service.ErrEmailExists):
			response.BadRequest(c, err.Error())
		default:
			response.InternalError(c, "注册失败")
		}
		return
	}

	response.SuccessWithMsg(c, "注册成功", user)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取JWT Token
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=vo.LoginVO} "登录成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 401 {object} response.Response "用户名或密码错误"
// @Router /api/v1/login [post]
func (a *UserAPI) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数验证失败: "+err.Error())
		return
	}

	loginVO, err := a.service.Login(&req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			response.Unauthorized(c, err.Error())
		} else {
			response.InternalError(c, "登录失败")
		}
		return
	}

	response.SuccessWithMsg(c, "登录成功", loginVO)
}
