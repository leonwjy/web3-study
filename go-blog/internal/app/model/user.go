package model

import (
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Username string `gorm:"type:varchar(32);uniqueIndex;not null;comment:用户名"`
	Password string `gorm:"type:varchar(128);not null;comment:密码"`
	Salt     string `gorm:"type:varchar(32);not null;comment:密码盐"`
	Email    string `gorm:"type:varchar(128);uniqueIndex;not null;comment:邮箱"`

	// 关联
	Posts    []Post    `gorm:"foreignKey:UserID"`
	Comments []Comment `gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
