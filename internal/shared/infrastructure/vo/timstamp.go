package vo

import "time"

type Timestamp struct {
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

func (t Timestamp) GetCreatedAt() time.Time {
	return t.CreatedAt
}

func (t Timestamp) GetUpdatedAt() time.Time {
	return t.UpdatedAt
}
