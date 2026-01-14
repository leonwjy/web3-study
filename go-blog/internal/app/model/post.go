package model

import (
	"gorm.io/gorm"
)

// Post 文章模型
type Post struct {
	gorm.Model
	Title   string `gorm:"type:varchar(200);not null;comment:文章标题"`
	Content string `gorm:"type:longtext;not null;comment:文章内容"`
	UserID  uint   `gorm:"index;not null;comment:作者ID"`

	// 关联
	User     User      `gorm:"foreignKey:UserID"`
	Comments []Comment `gorm:"foreignKey:PostID"`
}

// TableName 指定表名
func (Post) TableName() string {
	return "posts"
}
