package vo

import (
	"time"

	"go_blog/internal/app/model"
)

// UserVO 用户响应数据
type UserVO struct {
	ID        uint      `json:"id" example:"1"`
	Username  string    `json:"username" example:"johndoe"`
	Email     string    `json:"email" example:"john@example.com"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

// LoginVO 登录响应数据
type LoginVO struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  UserVO `json:"user"`
}

// ToUserVO 将 User 模型转换为 UserVO
func ToUserVO(user *model.User) *UserVO {
	if user == nil {
		return nil
	}
	return &UserVO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
}
