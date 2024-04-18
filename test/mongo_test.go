package test

import (
	"context"
	"fmt"
	"im/models"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestFindOne(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().SetAuth(options.Credential{
		Username:    "admin",
		Password:    "admin",
		PasswordSet: false,
	}).ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Error(err)
	}
	db := client.Database("im")
	user := new(models.UserBasic)
	err = db.Collection("user_basic").FindOne(context.Background(), bson.D{}).Decode(user)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("user =>", user)
}
