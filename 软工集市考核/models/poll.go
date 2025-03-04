package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// 投票类型
const (
	PollTypeBinary = "binary"  // 二分选项（是/否）
	PollTypeSingle = "single"  // 单选
	PollTypeMulti  = "multi"   // 多选
)

// Poll 投票模型
type Poll struct {
	ID          string    `json:"id" gorm:"primary_key"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	Type        string    `json:"type" gorm:"not null"` // binary, single, multi
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	EndTime     time.Time `json:"end_time"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	Options     []Option  `json:"options" gorm:"foreignkey:PollID"`
}

// Option 选项模型
type Option struct {
	ID        string    `json:"id" gorm:"primary_key"`
	PollID    string    `json:"poll_id" gorm:"not null"`
	Text      string    `json:"text" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Votes     []Vote    `json:"votes,omitempty" gorm:"foreignkey:OptionID"`
}

// Vote 投票记录模型
type Vote struct {
	ID        string    `json:"id" gorm:"primary_key"`
	PollID    string    `json:"poll_id" gorm:"not null"`
	OptionID  string    `json:"option_id" gorm:"not null"`
	UserID    string    `json:"user_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}

// User 用户模型
type User struct {
	ID        string    `json:"id" gorm:"primary_key"`
	Username  string    `json:"username" gorm:"unique;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate 在创建记录前生成UUID
func (poll *Poll) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New().String())
}

// BeforeCreate 在创建记录前生成UUID
func (option *Option) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New().String())
}

// BeforeCreate 在创建记录前生成UUID
func (vote *Vote) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New().String())
}

// BeforeCreate 在创建记录前生成UUID
func (user *User) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("ID", uuid.New().String())
} 