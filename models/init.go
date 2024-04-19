package models

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Mongo = InitMongoDB("im", "mongodb://localhost:27017", "admin", "admin")
var RedisCli = InitRedis("192.168.6.66:6380", "skyinfor2016")

func InitMongoDB(dbName, address, userName, pasword string) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().SetAuth(options.Credential{
		Username:    userName,
		Password:    pasword,
		PasswordSet: false,
	}).ApplyURI(address))
	if err != nil {
		log.Println("connect MongoDB Error:", err)
		return nil
	}
	return client.Database(dbName)
}

func InitRedis(address, password string) *redis.Client {
	config := &redis.Options{
		Addr:         address,
		Password:     password,
		DB:           0, // 使用默认DB
		PoolSize:     15,
		MinIdleConns: 10, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。
		//超时
		//DialTimeout:  5 * time.Second, //连接建立超时时间，默认5秒。
		//ReadTimeout:  3 * time.Second, //读超时，默认3秒， -1表示取消读超时
		//WriteTimeout: 3 * time.Second, //写超时，默认等于读超时
		//PoolTimeout:  4 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。
	}
	redisClient := redis.NewClient(config)
	return redisClient
}
