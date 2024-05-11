package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var myKey = []byte("testToken")

type UserClaims struct {
	Mobile string
	Id     int
	jwt.RegisteredClaims
}

func TestJwtToken(t *testing.T) {

	mobile := "13612900319"
	id := 10086
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
		t.Fatalf(err.Error())
	}
	fmt.Println(tokenString, "++++生成token")
}

func TestAnalyseToken(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJNb2JpbGUiOiIxMzYxMjkwMDMxOSIsIklkIjoxMDA4Niwic3ViIjoiVG9rZW4iLCJleHAiOjE3MTQyOTY0NjYsImlhdCI6MTcxNDI4NjM4Nn0.M8DGeNDtfH5LhcPbqc-Tq4QKwfNB4GURt06Fsp_dsyc"

	userClaim := new(UserClaims)
	claims, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		t.Fatalf(err.Error())
	}
	if !claims.Valid {
		fmt.Errorf("analyse Token Error:%v", err)
	}
	fmt.Println(userClaim, "+++++token解析结果")
}
