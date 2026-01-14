package dto

// RegisterRequest 用户注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32" example:"johndoe"`
	Password string `json:"password" binding:"required,min=6,max=32" example:"password123"`
	Email    string `json:"email" binding:"required,email" example:"john@example.com"`
}

// LoginRequest 用户登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"johndoe"`
	Password string `json:"password" binding:"required" example:"password123"`
}
