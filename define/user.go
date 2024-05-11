package define

import (
	"time"
)

type UserReg struct {
	Name          string `json:"name" binding:"required"`
	PassWord      string `json:"password" binding:"required"`
	ResetPassWord string `json:"reset_pass_word" binding:"required"`
	Phone         string `json:"phone" binding:"required"`
	Email         string `json:"email"`
	Identity      string `json:"identity"`
	LoginTime     time.Time
	HeartbeatTime time.Time
	LoginOutTime  time.Time
}

type UserLogin struct {
	PassWord string `json:"password" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
}

type UserSave struct {
	UserLogin
	Name string `json:"name" binding:"required"`
	// Name string `json:"name" validate:"required||string=1,5||unique"`
}
