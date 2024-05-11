package models

import (
	"ginchat/define"
)

// 消息内容
type Message struct {
	// gorm.Model
	define.SendMsg

	FormId    int //发送者
	CreatedAt string
	UpdatedAt string
	// CreatedAt   *time.Time `gorm:"column:created_at;type:datetime;comment:创建时间" json:"created_at"`
	// UpdatedAt   *time.Time `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updated_at"`
}

func (msg *Message) TableName() string {
	return "message"
}
