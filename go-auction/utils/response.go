package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// 业务状态码定义
const (
	CodeSuccess       = 200
	CodeBadRequest    = 400
	CodeUnauthorized  = 401
	CodeForbidden     = 403
	CodeNotFound      = 404
	CodeInternalError = 500
)

// 状态码对应的消息
var codeMsg = map[int]string{
	CodeSuccess:       "success",
	CodeBadRequest:    "请求参数错误",
	CodeUnauthorized:  "未授权",
	CodeForbidden:     "禁止访问",
	CodeNotFound:      "资源不存在",
	CodeInternalError: "服务器内部错误",
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  codeMsg[CodeSuccess],
		Data: data,
	})
}

// SuccessWithMsg 成功响应（自定义消息）
func SuccessWithMsg(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  msg,
		Data: data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
	})
}

// FailWithCode 失败响应（使用预定义状态码）
func FailWithCode(c *gin.Context, code int) {
	msg, ok := codeMsg[code]
	if !ok {
		msg = "未知错误"
	}
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
	})
}

// BadRequest 请求参数错误
func BadRequest(c *gin.Context, msg string) {
	if msg == "" {
		msg = codeMsg[CodeBadRequest]
	}
	Fail(c, CodeBadRequest, msg)
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context, msg string) {
	if msg == "" {
		msg = codeMsg[CodeUnauthorized]
	}
	c.JSON(http.StatusUnauthorized, Response{
		Code: CodeUnauthorized,
		Msg:  msg,
	})
}

// Forbidden 禁止访问
func Forbidden(c *gin.Context, msg string) {
	if msg == "" {
		msg = codeMsg[CodeForbidden]
	}
	c.JSON(http.StatusForbidden, Response{
		Code: CodeForbidden,
		Msg:  msg,
	})
}

// NotFound 资源不存在
func NotFound(c *gin.Context, msg string) {
	if msg == "" {
		msg = codeMsg[CodeNotFound]
	}
	Fail(c, CodeNotFound, msg)
}

// InternalError 服务器内部错误
func InternalError(c *gin.Context, msg string) {
	if msg == "" {
		msg = codeMsg[CodeInternalError]
	}
	Fail(c, CodeInternalError, msg)
}
