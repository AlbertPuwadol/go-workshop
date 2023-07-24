package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hashtag struct {
	ID      primitive.ObjectID `bson:"_id"`
	Keyword string             `bson:"keyword"`
}
