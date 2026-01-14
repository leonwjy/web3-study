package model

import (
	"gorm.io/gorm"
)

// Comment 评论模型
type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null;comment:评论内容"`
	UserID  uint   `gorm:"index;not null;comment:评论者ID"`
	PostID  uint   `gorm:"index;not null;comment:文章ID"`

	// 关联
	User User `gorm:"foreignKey:UserID"`
	Post Post `gorm:"foreignKey:PostID"`
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}
