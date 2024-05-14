package service

import (
	"context"
	"encoding/json"
	"fmt"
	"ginchat/common"
	"ginchat/define"
	"ginchat/helper"
	"ginchat/middlewares"
	"ginchat/models"
	"ginchat/utils"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {

	var params define.UserReg

	if err := c.BindJSON(&params); err != nil {
		common.Response(c, 401, err.Error(), nil)
		return
	}
	if params.PassWord != params.ResetPassWord {
		common.Response(c, -1, "两次输入的密码不一致", nil)
		return
	}
	isMobile := helper.CheckMobile(params.Phone)
	if !isMobile {
		common.Response(c, -1, "手机号格式有误!", nil)
		return
	}
	//根据手机号判断是否已经注册
	user := models.FindUserByMobile(params.Phone)

	if user.Name != "" {
		common.Response(c, -1, "手机号已被注册!", gin.H{})
		return
	}
	salt := fmt.Sprintf("%06d", rand.Int31())
	params.PassWord = helper.GetMd5(params.PassWord)
	params.Identity = salt

	params.LoginTime = time.Now()
	params.LoginOutTime = time.Now()
	params.HeartbeatTime = time.Now()
	err, result := models.Create(params)
	if err != nil {
		common.Response(c, 401, "系统错误，注册失败!", params)
		return
	}

	common.Response(c, 200, "success", gin.H{
		"id": result.ID,
	})
}

func Login(c *gin.Context) {
	var params define.UserLogin

	if err := c.BindJSON(&params); err != nil {
		common.Response(c, 401, err.Error(), nil)
		return
	}
	if params.Phone == "" {
		common.Response(c, -1, "手机号码不能为空!", nil)
		return
	}
	if params.PassWord == "" {
		common.Response(c, -1, "密码不能为空!", nil)
		return
	}

	user, err := models.FindUserByMobilePwd(params.Phone, helper.GetMd5(params.PassWord))
	if err != nil {
		common.Response(c, -1, "手机号或密码有误!", nil)
		return
	}
	//颁发token
	token, err := helper.GenerateToken(user.Phone, int(user.ID))
	if err != nil {
		common.Response(c, -1, "签发token失败!", nil)
		return
	}
	//在redis或者mysql中存储一份，保证登出操作

	common.Response(c, 200, "success", gin.H{
		"userInfo": user,
		"token":    token,
	})
}

func GetUserInfo(c *gin.Context) {
	_, err := CheckAuth(c)
	if err != nil {
		common.Response(c, -1, "获取失败，token校验失败!", nil)
		return
	}
	userId, _ := strconv.Atoi(c.DefaultQuery("user_id", "0"))
	userInfo, err := models.FindUserId(userId)
	if err != nil {
		common.Response(c, -1, "用户不存在!", nil)
		return
	}
	common.Response(c, http.StatusOK, "获取用户详情信息成功", userInfo)
}

func GetUserDetails(c *gin.Context) {
	user, err := CheckAuth(c)
	if err != nil {
		common.Response(c, -1, "获取失败，token校验失败!", nil)
		return
	}

	userInfo, err := models.FindUserId(user.Id)
	if err != nil {
		common.Response(c, -1, "用户不存在!", nil)
		return
	}
	common.Response(c, http.StatusOK, "获取用户详情信息成功", userInfo)
}

func UpdateUser(c *gin.Context) {
	var params define.UserSave

	if err := c.BindJSON(&params); err != nil {
		common.Response(c, 401, err.Error(), nil)
		return
	}

	isMobile := helper.CheckMobile(params.Phone)
	if !isMobile {
		common.Response(c, -1, "手机号格式有误!", nil)
		return
	}

	user, err := CheckAuth(c)
	if err != nil {
		common.Response(c, -1, "获取失败，token校验失败!", nil)
		return
	}

	userInfo := models.FindUserByMobile(params.Phone)

	if userInfo.Phone == params.Phone && userInfo.ID != uint(user.Id) {
		common.Response(c, -1, "手机号已被使用!", gin.H{})
		return
	}

	_, err = models.FindUserId(user.Id)
	if err != nil {
		common.Response(c, -1, "用户不存在!", nil)
		return
	}
	u := models.UserBasic{
		Name:     params.Name,
		Phone:    params.Phone,
		PassWord: helper.GetMd5(params.PassWord),
	}
	err = models.UpdateUser(user.Id, &u)
	if err != nil {
		common.Response(c, -1, "系统出小差了，修改失败!", nil)
		return
	}
	common.Response(c, 200, "success", nil)
}

func DeleteUser(c *gin.Context) {
	user, err := CheckAuth(c)
	if err != nil {
		common.Response(c, -1, "获取失败，token校验失败!", nil)
		return
	}
	_, err = models.FindUserId(user.Id)
	if err != nil {
		common.Response(c, -1, "用户不存在!", nil)
		return
	}
	err = models.DeleteUser(user.Id)
	if err != nil {
		common.Response(c, -1, "系统出错，删除失败!", nil)
		return
	}
	common.Response(c, http.StatusOK, "删除成功!", nil)
}

func SendCode(c *gin.Context) {
	json := make(map[string]interface{})
	c.BindJSON(&json)

	email, ok := json["email"]
	if !ok {
		common.Response(c, 403, "请传入email", nil)
		return
	}
	//判断邮箱是否被注册
	id := models.GetUserBasicByEmail(email.(string))

	if id > 0 {
		msg := fmt.Sprintf(email.(string) + "已被注册!")
		common.Response(c, 403, msg, nil)
		return
	}
	code := helper.GetRandCode()
	err := helper.SendCode(email.(string), code)
	if err != nil {
		log.Printf("[Send Code 发送失败]:%v\n", err)
		common.Response(c, 403, "发送失败", nil)
		return
	}

	//发送成功存入redis中，有效期为30分钟
	err = models.CacheSet(define.CACHE_PREFIX+email.(string), code, time.Second*1800)
	if err != nil {
		log.Printf("[Cache Set Code ERROR]:%v\n", err)
	}

	common.Response(c, 200, "success", nil)
}

// 获取用户列表
func GetUserList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	user, err := CheckAuth(c)
	if err != nil {
		common.Response(c, -1, "获取失败，token校验失败!", nil)
		return
	}
	page = (page - 1) * pageSize
	userList, err := models.FindUserList(user.Id, page, pageSize)
	if err != nil {
		common.Response(c, -1, "获取用户列表失败!", nil)
		return
	}
	//统计用户总数
	count := models.CountUserTotal()
	common.Response(c, http.StatusOK, "success", gin.H{
		"userList": userList,
		"total":    count,
	})
}

// 校验用户token
func CheckAuth(c *gin.Context) (*helper.UserClaims, error) {
	var err error
	u, _ := c.Get(middlewares.USER_CLAIMS)
	if u == nil {
		log.Printf("[GET user_claims ERROR]:%v\n")
		return nil, err
	}
	user := u.(*helper.UserClaims)
	return user, err
}

// 单文件上传
func FileUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		common.Response(c, -1, "请选择文件上传!", nil)
		return
	}
	extName := path.Ext(file.Filename)

	allowExt := map[string]bool{
		".jpg":  true,
		".png":  true,
		".gif":  true,
		".jpeg": true,
		".docx": true,
		".xlsx": true,
		".pptx": true,
	}
	if _, ok := allowExt[extName]; !ok {
		common.Response(c, -1, "仅支持png,jpg,jpeg,docx,xlsx,gig,pptx类型文件!", nil)
		return
	}
	day := time.Now().Format("2006-01-02")
	dir := "./static/file/" + day
	e := os.MkdirAll(dir, 0755)
	if e != nil {
		fmt.Println(e, "+++++")
		common.Response(c, -1, "创建文件目录出错，上传文件失败!", nil)
		return
	}

	fileMaxSize := 10 << 20 //10M
	if int(file.Size) > fileMaxSize {
		common.Response(c, -1, "文件不允许大于4M!", nil)
		return
	}
	fileName := strconv.FormatInt(time.Now().Unix(), 10) + extName
	dst := path.Join(dir, fileName)

	c.SaveUploadedFile(file, dst)
	response := gin.H{
		"file": dst,
	}
	common.Response(c, 200, "上传成功!", response)
}

// 获取消息列表
func MessageList(c *gin.Context) {

	user, err := CheckAuth(c)
	if err != nil {
		common.Response(c, -1, "获取失败，token校验失败!", nil)
		return
	}
	userId := user.Id
	start, _ := strconv.Atoi(c.DefaultQuery("page", "0"))
	end, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	target_id, _ := strconv.Atoi(c.DefaultQuery("target_id", "0"))
	if target_id == 0 {
		common.Response(c, -1, "target_id empty!", nil)
		return
	}
	userInfo, err := models.FindUserId(userId)
	if err != nil {
		common.Response(c, -1, "用户不存在!", nil)
		return
	}
	targetUserInfo, err := models.FindUserId(target_id)
	if err != nil {
		common.Response(c, -1, "用户不存在!", nil)
		return
	}
	var key string
	targetIdStr := strconv.Itoa(target_id)
	userIdStr := strconv.Itoa(userId)
	if target_id > userId {
		key = "chat_" + userIdStr + "_" + targetIdStr
	} else {
		key = "chat_" + targetIdStr + "_" + userIdStr
	}
	// list, err := utils.Redis.ZRevRange(context.Background(), key, int64(start), int64(end)).Result()
	list, err := utils.Redis.ZRange(context.Background(), key, int64(start), int64(end)).Result()
	if err != nil {
		common.Response(c, -1, "获取聊天记录失败!", nil)
		return
	}
	type messageInfo struct {
		Content     interface{} `json:"content"`
		Group       string      `json:"group"`
		Uindex      int         `json:"uindex"`
		Uname       string      `json:"uname"`
		ContentType string      `json:"contentType"`
		Uface       string      `json:"uface"`
		Date        string      `json:"date"`
		TargetId    int         `json:"target_id"`
	}
	msgList := make([]messageInfo, 0)

	for _, value := range list {
		message := new(define.SendMsg)
		err := json.Unmarshal([]byte(value), &message)
		if err != nil {
			break
		}
		var info = messageInfo{}
		if message.TargetId == target_id {
			info = messageInfo{
				Group:       "chat",
				Content:     message.Content,
				Uindex:      userId,
				Uname:       userInfo.Name,
				ContentType: "txt",
				Uface:       userInfo.Avatar,
				Date:        message.Date,
				TargetId:    target_id,
			}
		} else {
			info = messageInfo{
				Group:       "chat",
				Content:     message.Content,
				Uindex:      target_id,
				Uname:       targetUserInfo.Name,
				ContentType: "txt",
				Uface:       targetUserInfo.Avatar,
				Date:        message.Date,
				TargetId:    userId,
			}
		}
		msgList = append(msgList, info)
	}
	common.Response(c, 200, "success", msgList)
	return
}
