package vo

import (
	"time"

	"go_blog/internal/app/model"
)

// PostVO 文章响应数据
type PostVO struct {
	ID        uint      `json:"id" example:"1"`
	Title     string    `json:"title" example:"我的第一篇博客"`
	Content   string    `json:"content" example:"这是文章内容..."`
	UserID    uint      `json:"user_id" example:"1"`
	Author    string    `json:"author" example:"johndoe"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

// PostListVO 文章列表响应数据
type PostListVO struct {
	Total int64     `json:"total" example:"100"`
	List  []*PostVO `json:"list"`
}

// ToPostVO 将 Post 模型转换为 PostVO
func ToPostVO(post *model.Post) *PostVO {
	if post == nil {
		return nil
	}
	vo := &PostVO{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		UserID:    post.UserID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}
	// 如果关联了用户，填充作者名
	if post.User.ID != 0 {
		vo.Author = post.User.Username
	}
	return vo
}

// ToPostVOList 将 Post 模型列表转换为 PostVO 列表
func ToPostVOList(posts []*model.Post) []*PostVO {
	list := make([]*PostVO, 0, len(posts))
	for _, post := range posts {
		list = append(list, ToPostVO(post))
	}
	return list
}
