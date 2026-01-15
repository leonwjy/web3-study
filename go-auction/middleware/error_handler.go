package middleware

import (
	"log/slog"
	"net/http"

	"go-auction/utils"

	"github.com/gin-gonic/gin"
)

// ErrorHandler 全局错误处理中间件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("Panic recovered", "error", err, "path", c.Request.URL.Path, "method", c.Request.Method)

				// 返回统一错误响应
				utils.InternalError(c, "服务器内部错误")
				c.Abort()
			}
		}()

		c.Next()

		// 检查是否有错误（通过状态码判断）
		if c.Writer.Status() >= http.StatusBadRequest {
			// 如果还没有设置响应体，设置默认错误响应
			if !c.Writer.Written() {
				switch c.Writer.Status() {
				case http.StatusNotFound:
					utils.NotFound(c, "资源不存在")
				case http.StatusUnauthorized:
					utils.Unauthorized(c, "未授权")
				case http.StatusForbidden:
					utils.Forbidden(c, "禁止访问")
				case http.StatusBadRequest:
					utils.BadRequest(c, "请求参数错误")
				default:
					utils.InternalError(c, "服务器内部错误")
				}
			}
		}
	}
}
