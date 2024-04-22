package models

type RoomBasic struct {
	Number       string `bson:"number"`
	Name         string `bson:"name"`
	Info         string `bson:"info"`
	UserIdentity string `bson:"user_identity"`
	CreatedAt    string `bson:"created_at"`
	UpdatedAt    string `bson:"updated_at"`
}

func (u RoomBasic) CollectionName() string {
	return "room_basic"
}
