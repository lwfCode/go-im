package service

import (
	"im/common"
	"im/helper"
	"im/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
	Sex      int64  `json:"sex"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

// 用户登陆
func Login(c *gin.Context) {
	var params User

	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "304", "error": err.Error()})
		return
	}

	//查询mongoDB中用户账号密码
	user, err := models.GetUserBasicByAccountPassword(params.Account, helper.GetMd5(params.Password))

	if err != nil {
		common.Response(c, http.StatusOK, "用户不存在", nil)
		return
	}
	token, err := helper.GenerateToken(user.Identity, user.Email)
	if err != nil {
		common.Response(c, http.StatusBadRequest, "系统有误", nil)
		return
	}

	common.Response(c, http.StatusOK, "登陆成功", gin.H{
		"token": token,
	})
}

// 用户注册
func Register(c *gin.Context) {
	var params User

	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "304", "error": err.Error()})
		return
	}

	//判断数据库中是否已存在用户
	id, err := models.GetUserBasicByAccount(params.Account)
	if err != nil {
		common.Response(c, http.StatusBadRequest, "系统出错!", nil)
		return
	}
	if id > 0 {
		common.Response(c, http.StatusBadRequest, "该账号已经注册，请勿重复操作!", nil)
		return
	}

	result, err := models.InsertUserInfo(
		params.Account, params.Password, params.Nickname, params.Email, params.Avatar, params.Sex,
	)
	if err != nil {
		common.Response(c, http.StatusBadRequest, "注册失败", nil)
		return
	}

	common.Response(c, http.StatusOK, "注册成功", result)
}

// 获取用户详情信息
func GetUserDetails(c *gin.Context) {
	u, _ := c.Get("user_claims")
	if u == nil {
		log.Printf("[GET user_claims ERROR]:%v\n")
	}
	user := u.(*helper.UserClaims)
	userInfo, err := models.GetUserBasicById(user.Identity)
	if err != nil {
		log.Printf("[DB ERROR]:%v\n", err)
		common.Response(c, http.StatusBadRequest, "系统查询有误!", nil)
		return
	}
	common.Response(c, http.StatusOK, "获取用户详情信息成功", userInfo)
}
