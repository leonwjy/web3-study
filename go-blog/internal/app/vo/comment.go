package vo

import (
	"time"

	"go_blog/internal/app/model"
)

// CommentVO 评论响应数据
type CommentVO struct {
	ID        uint      `json:"id" example:"1"`
	Content   string    `json:"content" example:"这是一条评论"`
	UserID    uint      `json:"user_id" example:"1"`
	Author    string    `json:"author" example:"johndoe"`
	PostID    uint      `json:"post_id" example:"1"`
	CreatedAt time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
}

// CommentListVO 评论列表响应数据
type CommentListVO struct {
	Total int64        `json:"total" example:"50"`
	List  []*CommentVO `json:"list"`
}

// ToCommentVO 将 Comment 模型转换为 CommentVO
func ToCommentVO(comment *model.Comment) *CommentVO {
	if comment == nil {
		return nil
	}
	vo := &CommentVO{
		ID:        comment.ID,
		Content:   comment.Content,
		UserID:    comment.UserID,
		PostID:    comment.PostID,
		CreatedAt: comment.CreatedAt,
	}
	// 如果关联了用户，填充作者名
	if comment.User.ID != 0 {
		vo.Author = comment.User.Username
	}
	return vo
}

// ToCommentVOList 将 Comment 模型列表转换为 CommentVO 列表
func ToCommentVOList(comments []*model.Comment) []*CommentVO {
	list := make([]*CommentVO, 0, len(comments))
	for _, comment := range comments {
		list = append(list, ToCommentVO(comment))
	}
	return list
}
