package service

import (
	"errors"
	"log/slog"

	"go_blog/internal/app/dto"
	"go_blog/internal/app/model"
	"go_blog/internal/app/repository"
	"go_blog/internal/app/vo"

	"gorm.io/gorm"
)

var (
	ErrPostNotFound = errors.New("文章不存在")
	ErrNotAuthor    = errors.New("无权操作此文章")
)

// PostService 文章服务
type PostService struct {
	repo        *repository.PostRepository
	commentRepo *repository.CommentRepository
}

// NewPostService 创建文章服务实例
func NewPostService() *PostService {
	return &PostService{
		repo:        repository.NewPostRepository(),
		commentRepo: repository.NewCommentRepository(),
	}
}

// Create 创建文章
func (s *PostService) Create(userID uint, req *dto.CreatePostRequest) (*vo.PostVO, error) {
	post := &model.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	if err := s.repo.Create(post); err != nil {
		slog.Error("创建文章失败", "error", err)
		return nil, err
	}

	// 重新获取文章以加载关联数据
	post, err := s.repo.GetByID(post.ID)
	if err != nil {
		return nil, err
	}

	slog.Info("文章创建成功", "post_id", post.ID, "user_id", userID)
	return vo.ToPostVO(post), nil
}

// GetByID 根据ID获取文章
func (s *PostService) GetByID(id uint) (*vo.PostVO, error) {
	post, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}
	return vo.ToPostVO(post), nil
}

// GetList 获取文章列表
func (s *PostService) GetList(req *dto.PostListRequest) (*vo.PostListVO, error) {
	// 默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	posts, total, err := s.repo.GetList(page, pageSize)
	if err != nil {
		slog.Error("获取文章列表失败", "error", err)
		return nil, err
	}

	return &vo.PostListVO{
		Total: total,
		List:  vo.ToPostVOList(posts),
	}, nil
}

// Update 更新文章
func (s *PostService) Update(id, userID uint, req *dto.UpdatePostRequest) (*vo.PostVO, error) {
	// 获取文章
	post, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	// 检查是否为作者
	if post.UserID != userID {
		return nil, ErrNotAuthor
	}

	// 更新字段
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}

	if err := s.repo.Update(post); err != nil {
		slog.Error("更新文章失败", "error", err)
		return nil, err
	}

	slog.Info("文章更新成功", "post_id", id, "user_id", userID)
	return vo.ToPostVO(post), nil
}

// Delete 删除文章
func (s *PostService) Delete(id, userID uint) error {
	// 获取文章
	post, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPostNotFound
		}
		return err
	}

	// 检查是否为作者
	if post.UserID != userID {
		return ErrNotAuthor
	}

	// 删除文章的所有评论
	if err := s.commentRepo.DeleteByPostID(id); err != nil {
		slog.Error("删除文章评论失败", "error", err)
		return err
	}

	// 删除文章
	if err := s.repo.Delete(id); err != nil {
		slog.Error("删除文章失败", "error", err)
		return err
	}

	slog.Info("文章删除成功", "post_id", id, "user_id", userID)
	return nil
}

// Exists 检查文章是否存在
func (s *PostService) Exists(id uint) (bool, error) {
	return s.repo.Exists(id)
}
