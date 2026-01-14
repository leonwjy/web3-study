package dto

// CreateCommentRequest 创建评论请求
type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1,max=1000" example:"这是一条评论"`
}

// CommentListRequest 评论列表请求
type CommentListRequest struct {
	Page     int `form:"page" binding:"omitempty,min=1" example:"1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100" example:"10"`
}
