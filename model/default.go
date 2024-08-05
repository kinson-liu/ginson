package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type DefaultFields struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt,omitempty"`
}

// DefaultUpdateAt changes the default updateAt field
func (df *DefaultFields) DefaultUpdateAt() {
	df.UpdatedAt = time.Now().Local()
}

// DefaultCreateAt changes the default createAt field
func (df *DefaultFields) DefaultCreateAt() {
	if df.CreatedAt.IsZero() {
		df.CreatedAt = time.Now().Local()
	}
}

// DefaultId changes the default _id field
func (df *DefaultFields) DefaultId() {
	if df.Id.IsZero() {
		df.Id = primitive.NewObjectID()
	}
}

func (df *DefaultFields) Defaults() {
	df.DefaultId()
	df.DefaultCreateAt()
	df.DefaultUpdateAt()
}
