package helper

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"ginchat/define"
	"log"
	"math/rand"
	"net"
	"net/smtp"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jordan-wright/email"
)

type UserClaims struct {
	Mobile string `json:"mobile"`
	Id     int    `json:"id"`
	jwt.RegisteredClaims
}

var myKey = []byte("ginChat")

func GetMd5(s string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

// GenerateToken
// 生成 token
func GenerateToken(mobile string, id int) (string, error) {
	UserClaim := &UserClaims{
		Mobile: mobile,
		Id:     id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 10080)), //设置过期时间，7天过期
			IssuedAt:  jwt.NewNumericDate(time.Now()),                          //颁发时间
			Subject:   "Token",                                                 //主题
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
	e.Subject = "【Go-IM】您好，验证码已发送，30分钟内有效。"
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

// mobile verify
func CheckMobile(mobile string) bool {
	regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"

	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobile)
}

// email verify
func CheckEmail(email string) bool {
	//pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// 获取服务器Ip
func GetServerIp() (ip string) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return ""
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				ip = ipNet.IP.String()
			}
		}
	}

	return
}
