package dto

// CreatePostRequest 创建文章请求
type CreatePostRequest struct {
	Title   string `json:"title" binding:"required,min=1,max=200" example:"我的第一篇博客"`
	Content string `json:"content" binding:"required,min=1" example:"这是文章内容..."`
}

// UpdatePostRequest 更新文章请求
type UpdatePostRequest struct {
	Title   string `json:"title" binding:"omitempty,min=1,max=200" example:"更新后的标题"`
	Content string `json:"content" binding:"omitempty,min=1" example:"更新后的内容..."`
}

// PostListRequest 文章列表请求
type PostListRequest struct {
	Page     int `form:"page" binding:"omitempty,min=1" example:"1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=100" example:"10"`
}
