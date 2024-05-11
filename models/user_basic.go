package models

import (
	"context"
	"ginchat/define"
	"ginchat/utils"
	"time"

	"gorm.io/gorm"
)

type UserBasic struct {
	gorm.Model
	Name          string
	PassWord      string
	Phone         string
	Email         string
	Identity      string
	ClientIp      string
	ClientPort    string
	LoginTime     time.Time
	HeartbeatTime time.Time
	LogOutTime    time.Time `gorm:"column:login_out_time" json:"login_out_time"`
	IsLogout      bool
	DeviceInfo    string
	Avatar        string `json:"avatar"`
}

func (u *UserBasic) TableName() string {
	return "user_basic"
}

// 新增用户
func Create(u define.UserReg) (error, UserBasic) {
	user := UserBasic{
		Name:      u.Name,
		PassWord:  u.PassWord,
		Phone:     u.Phone,
		Email:     u.Email,
		Identity:  u.Identity,
		LoginTime: u.LoginTime,
	}
	err := utils.DB.Create(&user).Error
	return err, user
}

func GetUserBasicByEmail(email string) int64 {
	var count int64
	utils.DB.Where("email = ?", email).Count(&count)
	return count
}

/** 根据手机号查询用户*/
func FindUserByMobile(mobile string) *UserBasic {
	user := new(UserBasic)
	utils.DB.Where("phone = ?", mobile).First(&user)
	return user
}

/** 根据用户id查询*/
func FindUserId(id int) (*UserBasic, error) {
	user := new(UserBasic)
	err := utils.DB.Where("id = ?", id).First(&user).Error
	return user, err
}

func FindUserByMobilePwd(mobile, password string) (*UserBasic, error) {
	user := new(UserBasic)
	err := utils.DB.Where("phone = ? and pass_word = ?", mobile, password).First(&user).Error
	return user, err
}

func UpdateUser(userId int, u *UserBasic) error {
	err := utils.DB.Where("id = ?", userId).Updates(u).Error
	return err
}

func DeleteUser(userId int) error {
	user := new(UserBasic)
	//Unscoped() 硬删除
	err := utils.DB.Unscoped().Where("id = ?", userId).Delete(&user).Error
	return err
}

func FindUserList(UserId int, page, pageSize int) ([]*UserBasic, error) {
	users := make([]*UserBasic, pageSize)
	err := utils.DB.Where("id != ?", UserId).
		Order("id desc").Limit(pageSize).Offset(page).Find(&users).Error
	return users, err
}

func CountUserTotal() int64 {
	var count int64
	utils.DB.Model(&UserBasic{}).Count(&count)
	return count
}

func CacheSet(key string, code interface{}, t time.Duration) error {
	return utils.Redis.Set(context.Background(), key, code, time.Second*1800).Err()
}
