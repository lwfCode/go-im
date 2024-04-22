package models

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

type UserRoom struct {
	UserIdentity string `bson:"user_identity"`
	RoomIdentity string `bson:"room_identity"`
	Message      string `bson:"message"`
	CreatedAt    string `bson:"created_at"`
	UpdatedAt    string `bson:"updated_at"`
}

func (u UserRoom) CollectionName() string {
	return "user_room"
}

func GetUserRoomByUserIdentityRoomIdentity(UserIdentity, RoomIdentity string) (*UserRoom, error) {
	ub := new(UserRoom)

	filter := bson.D{
		{"user_identity", UserIdentity},
		{"room_identity", RoomIdentity},
	}
	err := Mongo.Collection(UserRoom{}.CollectionName()).
		FindOne(context.Background(), filter).Decode(ub)

	return ub, err
}

func GetUserRoomByRoomIdentity(RoomIdentity string) ([]*UserRoom, error) {
	filter := bson.D{
		{"room_identity", RoomIdentity},
	}
	cur, err := Mongo.Collection(UserRoom{}.CollectionName()).Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	urs := make([]*UserRoom, 0)
	for cur.Next(context.Background()) {
		ur := new(UserRoom)
		err := cur.Decode(ur)
		if err != nil {
			return nil, err
		}
		urs = append(urs, ur)
	}
	return urs, nil
}
