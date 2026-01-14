package router

import (
	v1 "go_blog/api/v1"
	"go_blog/internal/app/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter 配置路由
func SetupRouter(mode string) *gin.Engine {
	// 设置 Gin 模式
	gin.SetMode(mode)

	r := gin.New()

	// 全局中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 初始化 API
	userAPI := v1.NewUserAPI()
	postAPI := v1.NewPostAPI()
	commentAPI := v1.NewCommentAPI()

	// Swagger 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 路由组
	apiV1 := r.Group("/api/v1")
	{
		// 公开路由（无需认证）
		apiV1.POST("/register", userAPI.Register)
		apiV1.POST("/login", userAPI.Login)

		// 文章公开路由
		apiV1.GET("/posts", postAPI.GetList)
		apiV1.GET("/posts/:id", postAPI.GetByID)

		// 评论公开路由
		apiV1.GET("/posts/:id/comments", commentAPI.GetByPostID)

		// 需要认证的路由
		auth := apiV1.Group("")
		auth.Use(middleware.JWTAuth())
		{
			// 文章管理
			auth.POST("/posts", postAPI.Create)
			auth.PUT("/posts/:id", postAPI.Update)
			auth.DELETE("/posts/:id", postAPI.Delete)

			// 评论管理
			auth.POST("/posts/:id/comments", commentAPI.Create)
		}
	}

	return r
}
