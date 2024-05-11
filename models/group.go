package models

// 群
type GroupBasic struct {
	Name      string
	UserId    int    //创建群用户id
	Type      int    //1-100人 2-500 3-1000人+
	Icon      string //群头像
	Desc      string //群描述
	CreatedAt string
	UpdatedAt string
}

func (g *GroupBasic) TableName() string {
	return "group_basic"
}
