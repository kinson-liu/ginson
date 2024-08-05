package model

import (
	"ginson/config"
	"ginson/core/const/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	DefaultFields `bson:",inline"`
	Username      string `json:"username" bson:"username,omitempty"`
	Password      string `json:"-" bson:"password,omitempty"`
	Nickname      string `json:"nickname" bson:"nickname,omitempty"`
	Avatar        string `json:"avatar" bson:"avatar,omitempty"`
	Phone         string `json:"phone" bson:"phone,omitempty"`
	Email         string `json:"email" bson:"email,omitempty"`
	Locked        bool   `json:"locked" bson:"locked,omitempty"`
}

func (u *User) CollectionName() string {
	return "user"
}

func (u *User) DB() *mongo.Collection {
	return db.Mongo.Database(config.Config.DB.Schema).Collection(u.CollectionName())
}
