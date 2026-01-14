package repository

import (
	"go_blog/internal/app/model"
	"go_blog/internal/pkg/database"

	"gorm.io/gorm"
)

// PostRepository 文章数据访问层
type PostRepository struct {
	db *gorm.DB
}

// NewPostRepository 创建文章仓库实例
func NewPostRepository() *PostRepository {
	return &PostRepository{db: database.DB}
}

// Create 创建文章
func (r *PostRepository) Create(post *model.Post) error {
	return r.db.Create(post).Error
}

// GetByID 根据ID获取文章
func (r *PostRepository) GetByID(id uint) (*model.Post, error) {
	var post model.Post
	err := r.db.Preload("User").First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetList 获取文章列表（分页）
func (r *PostRepository) GetList(page, pageSize int) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var total int64

	// 获取总数
	err := r.db.Model(&model.Post{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = r.db.Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// GetByUserID 获取用户的文章列表
func (r *PostRepository) GetByUserID(userID uint, page, pageSize int) ([]*model.Post, int64, error) {
	var posts []*model.Post
	var total int64

	// 获取总数
	err := r.db.Model(&model.Post{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	err = r.db.Preload("User").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// Update 更新文章
func (r *PostRepository) Update(post *model.Post) error {
	return r.db.Save(post).Error
}

// Delete 删除文章
func (r *PostRepository) Delete(id uint) error {
	return r.db.Delete(&model.Post{}, id).Error
}

// Exists 检查文章是否存在
func (r *PostRepository) Exists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Post{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}
