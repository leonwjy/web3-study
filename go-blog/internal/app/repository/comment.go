package repository

import (
	"go_blog/internal/app/model"
	"go_blog/internal/pkg/database"

	"gorm.io/gorm"
)

// CommentRepository 评论数据访问层
type CommentRepository struct {
	db *gorm.DB
}

// NewCommentRepository 创建评论仓库实例
func NewCommentRepository() *CommentRepository {
	return &CommentRepository{db: database.DB}
}

// Create 创建评论
func (r *CommentRepository) Create(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

// GetByID 根据ID获取评论
func (r *CommentRepository) GetByID(id uint) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.Preload("User").Preload("Post").First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetByPostID 获取文章的评论列表（分页）
func (r *CommentRepository) GetByPostID(postID uint, page, pageSize int) ([]*model.Comment, int64, error) {
	var comments []*model.Comment
	var total int64

	// 获取总数
	err := r.db.Model(&model.Comment{}).Where("post_id = ?", postID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = r.db.Preload("User").
		Where("post_id = ?", postID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// GetByUserID 获取用户的评论列表
func (r *CommentRepository) GetByUserID(userID uint, page, pageSize int) ([]*model.Comment, int64, error) {
	var comments []*model.Comment
	var total int64

	// 获取总数
	err := r.db.Model(&model.Comment{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = r.db.Preload("User").Preload("Post").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}

	return comments, total, nil
}

// Delete 删除评论
func (r *CommentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Comment{}, id).Error
}

// DeleteByPostID 删除文章的所有评论
func (r *CommentRepository) DeleteByPostID(postID uint) error {
	return r.db.Where("post_id = ?", postID).Delete(&model.Comment{}).Error
}
