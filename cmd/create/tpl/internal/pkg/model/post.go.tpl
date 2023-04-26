package model

import (
	"time"

	"gorm.io/gorm"

	"{[.RootPackage]}/pkg/util/id"
)

// PostM 是数据库中 post 记录 struct 格式的映射.
type PostM struct {
	ID        int64     `gorm:"column:id;primary_key"`
	Username  string    `gorm:"column:username;not null"`
	PostID    string    `gorm:"column:post_id;not null"`
	Title     string    `gorm:"column:title;not null"`
	Content   string    `gorm:"column:content"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// TableName 用来指定映射的 MySQL 表名.
func (p *PostM) TableName() string {
	return "post"
}

// BeforeCreate 在创建数据库记录之前生成 postID.
func (p *PostM) BeforeCreate(tx *gorm.DB) error {
	p.PostID = "post-" + id.GenShortID()

	return nil
}
