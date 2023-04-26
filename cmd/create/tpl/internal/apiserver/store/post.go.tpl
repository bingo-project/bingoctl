package store

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"{[.RootPackage]}/internal/pkg/model"
)

// PostStore 定义了 user 模块在 store 层所实现的方法.
type PostStore interface {
	Create(ctx context.Context, post *model.PostM) error
	Get(ctx context.Context, username, postID string) (*model.PostM, error)
	Update(ctx context.Context, post *model.PostM) error
	List(ctx context.Context, username string, offset, limit int) (int64, []*model.PostM, error)
	Delete(ctx context.Context, username string, postIDs []string) error
}

// PostStore 接口的实现.
type posts struct {
	db *gorm.DB
}

// 确保 posts 实现了 PostStore 接口.
var _ PostStore = (*posts)(nil)

func newPosts(db *gorm.DB) *posts {
	return &posts{db: db}
}

// Create 插入一条 post 记录.
func (u *posts) Create(ctx context.Context, post *model.PostM) error {
	return u.db.Create(&post).Error
}

// Get 根据用户名查询指定 user 的数据库记录.
func (u *posts) Get(ctx context.Context, username, postID string) (post *model.PostM, err error) {
	err = u.db.Where("username = ?", username).
		Where("post_id = ?", postID).
		First(&post).Error

	return
}

// Update 更新一条 user 数据库记录.
func (u *posts) Update(ctx context.Context, post *model.PostM) error {
	return u.db.Save(&post).Error
}

// List 根据 offset 和 limit 返回 user 列表.
func (u *posts) List(ctx context.Context, username string, offset, limit int) (count int64, ret []*model.PostM, err error) {
	err = u.db.Where("username = ?", username).
		Offset(offset).Limit(defaultLimit(limit)).Order("id desc").Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error

	return
}

// Delete 根据 username 删除数据库 user 记录.
func (u *posts) Delete(ctx context.Context, username string, postIDs []string) error {
	err := u.db.Where("username = ?", username).
		Where("post_id IN (?)", postIDs).
		Delete(&model.PostM{}).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return nil
}
