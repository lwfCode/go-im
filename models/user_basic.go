package models

import (
	"context"
	"im/helper"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserBasic struct {
	Identity  string `bson:"_id"`
	Account   string `bson:"account"`
	Password  string `bson:"password"`
	Nickname  string `bson:"nickname"`
	Sex       int    `bson:"sex"`
	Email     string `bson:"email"`
	Avatar    string `bson:"avatar"`
	CreatedAt string `bson:"created_at"`
	UpdateAt  string `bson:"update_at"`
}

func (u UserBasic) CollectionName() string {
	return "user_basic"
}

func GetUserBasicByAccountPassword(account, password string) (*UserBasic, error) {
	ub := new(UserBasic)

	filter := bson.D{
		{"account", account},
		{"password", password},
	}
	err := Mongo.Collection(UserBasic{}.CollectionName()).
		FindOne(context.Background(), filter).Decode(ub)

	return ub, err
}

func GetUserBasicByAccount(account string) (int64, error) {

	filter := bson.D{
		{"account", account},
	}
	return Mongo.Collection(UserBasic{}.CollectionName()).
		CountDocuments(context.Background(), filter)
}

func GetUserBasicById(id primitive.ObjectID) (*UserBasic, error) {
	ub := new(UserBasic)

	filter := bson.D{
		{"_id", id},
	}
	err := Mongo.Collection(UserBasic{}.CollectionName()).
		FindOne(context.Background(), filter).Decode(ub)

	return ub, err
}

func InsertUserInfo(
	account string, pasword string, nickname string, email string, avatar string, sex int64,
) (*mongo.InsertOneResult, error) {

	formattedTime := time.Now().Format("2006-01-02 15:04:05")
	data := bson.D{
		{"account", account}, {"password", helper.GetMd5(pasword)},
		{"nickname", nickname}, {"email", email},
		{"avatar", avatar}, {"sex", sex}, {"created_at", formattedTime}, {"update_at", formattedTime},
	}

	res, err := Mongo.Collection(UserBasic{}.CollectionName()).InsertOne(
		context.Background(), data,
	)
	return res, err
}

func GetUserBasicByEmail(email string) (int64, error) {

	filter := bson.D{
		{"email", email},
	}
	return Mongo.Collection(UserBasic{}.CollectionName()).
		CountDocuments(context.Background(), filter)
}
