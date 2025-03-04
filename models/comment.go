package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// Comment 评论模型
type Comment struct {
	ID        string    `json:"id" gorm:"primary_key"`
	PollID    string    `json:"poll_id" gorm:"not null"`
	UserID    string    `json:"user_id" gorm:"not null"`
	Content   string    `json:"content" gorm:"not null"`
	ParentID  string    `json:"parent_id"` // 父评论ID，用于回复功能
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `json:"user,omitempty" gorm:"foreignkey:UserID"`
	Replies   []Comment `json:"replies,omitempty" gorm:"-"` // 不存储在数据库中，用于API响应
}

// BeforeCreate 在创建记录前生成UUID
func (comment *Comment) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New().String())
} 