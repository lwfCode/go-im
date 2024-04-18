package models

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Mongo = InitMongoDB("im", "mongodb://localhost:27017", "admin", "admin")

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
