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
	ErrCommentNotFound = errors.New("评论不存在")
)

// CommentService 评论服务
type CommentService struct {
	repo     *repository.CommentRepository
	postRepo *repository.PostRepository
}

// NewCommentService 创建评论服务实例
func NewCommentService() *CommentService {
	return &CommentService{
		repo:     repository.NewCommentRepository(),
		postRepo: repository.NewPostRepository(),
	}
}

// Create 创建评论
func (s *CommentService) Create(userID, postID uint, req *dto.CreateCommentRequest) (*vo.CommentVO, error) {
	// 检查文章是否存在
	exists, err := s.postRepo.Exists(postID)
	if err != nil {
		slog.Error("检查文章失败", "error", err)
		return nil, err
	}
	if !exists {
		return nil, ErrPostNotFound
	}

	comment := &model.Comment{
		Content: req.Content,
		UserID:  userID,
		PostID:  postID,
	}

	if err := s.repo.Create(comment); err != nil {
		slog.Error("创建评论失败", "error", err)
		return nil, err
	}

	// 重新获取评论以加载关联数据
	comment, err = s.repo.GetByID(comment.ID)
	if err != nil {
		return nil, err
	}

	slog.Info("评论创建成功", "comment_id", comment.ID, "post_id", postID, "user_id", userID)
	return vo.ToCommentVO(comment), nil
}

// GetByPostID 获取文章的评论列表
func (s *CommentService) GetByPostID(postID uint, req *dto.CommentListRequest) (*vo.CommentListVO, error) {
	// 检查文章是否存在
	exists, err := s.postRepo.Exists(postID)
	if err != nil {
		slog.Error("检查文章失败", "error", err)
		return nil, err
	}
	if !exists {
		return nil, ErrPostNotFound
	}

	// 默认分页参数
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	comments, total, err := s.repo.GetByPostID(postID, page, pageSize)
	if err != nil {
		slog.Error("获取评论列表失败", "error", err)
		return nil, err
	}

	return &vo.CommentListVO{
		Total: total,
		List:  vo.ToCommentVOList(comments),
	}, nil
}

// GetByID 根据ID获取评论
func (s *CommentService) GetByID(id uint) (*vo.CommentVO, error) {
	comment, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCommentNotFound
		}
		return nil, err
	}
	return vo.ToCommentVO(comment), nil
}

// Delete 删除评论（只有评论作者可以删除）
func (s *CommentService) Delete(id, userID uint) error {
	comment, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCommentNotFound
		}
		return err
	}

	if comment.UserID != userID {
		return ErrNotAuthor
	}

	if err := s.repo.Delete(id); err != nil {
		slog.Error("删除评论失败", "error", err)
		return err
	}

	slog.Info("评论删除成功", "comment_id", id, "user_id", userID)
	return nil
}
