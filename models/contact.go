package models

// 人员关系
type Contact struct {
	UserId    int //发送者
	TargetId  int //接收者
	Type      int //1-好友2-拉黑3-删除4-陌生人5-其他
	CreatedAt string
	UpdatedAt string
}

func (msg *Contact) TableName() string {
	return "contact"
}
