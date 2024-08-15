package vo

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ID struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
}

func NewID(oid primitive.ObjectID) ID {
	return ID{ID: oid}
}

func (i *ID) Hex() string {
	return i.ID.Hex()
}
