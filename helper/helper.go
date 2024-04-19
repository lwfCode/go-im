package helper

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"im/define"
	"log"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jordan-wright/email"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserClaims struct {
	// Identity string `json:"identity"`
	Identity primitive.ObjectID `json:"identity"`
	Email    string             `json:"email"`
	jwt.RegisteredClaims
}

// GetMd5
// 生成 md5
func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

var myKey = []byte("im")

// GenerateToken
// 生成 token
func GenerateToken(identity, email string) (string, error) {
	objId, err := primitive.ObjectIDFromHex(identity)
	if err != nil {
		return "", err
	}
	UserClaim := &UserClaims{
		Identity: objId,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 30)), //设置过期时间 在当前基础上 添加一个小时后 过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),                       //颁发时间
			Subject:   "Token",                                              //主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// AnalyseToken
// 解析 token
func AnalyseToken(tokenString string) (*UserClaims, error) {
	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !claims.Valid {
		return nil, fmt.Errorf("analyse Token Error:%v", err)
	}
	return userClaim, nil
}

// SendCode
// 发送验证码
func SendCode(toUserEmail, code string) error {
	e := email.NewEmail()
	e.From = "go-im <13262713396@163.com>"
	e.To = []string{toUserEmail}
	e.Subject = "【Go-IM】您好，验证码已发送，请查收"
	e.HTML = []byte("您的验证码：<b>" + code + "</b>")

	MailPwd, err := define.GetMailPwd()
	if err != nil {
		log.Printf("auth error")
		return err
	}
	return e.SendWithTLS("smtp.163.com:465",
		smtp.PlainAuth("", "13262713396@163.com", MailPwd, "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
}

// GetCode
// 生成验证码
func GetCode() string {
	rand.Seed(time.Now().UnixNano())
	res := ""
	for i := 0; i < 6; i++ {
		res += strconv.Itoa(rand.Intn(10))
	}
	return res
}

// GetUUID
// 生成唯一码
func GetUUID() string {
	u := uuid.New()

	return fmt.Sprintf("%x", u)
}

// 生成随机6位验证码
func GetRandCode() string {
	code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	return code
}
