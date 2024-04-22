package models

import "context"

type MessageBasic struct {
	Identity     string `bson:"identity"`
	UserIdentity string `bson:"user_identity"`
	RoomIdentity string `bson:"room_identity"`
	Data         string `bson:"data"`
	CreatedAt    string `bson:"created_at"`
	UpdatedAt    string `bson:"updated_at"`
}

func (u MessageBasic) CollectionName() string {
	return "message_basic"
}

func InsertOneMessageBasic(msg *MessageBasic) error {
	_, err := Mongo.Collection(MessageBasic{}.CollectionName()).InsertOne(context.Background(), msg)
	return err
}
