package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type GuestUser struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string `json:"username,omitempty" bson:"username,omitempty"`
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
	Token  int64 `json:"token,omitempty" bson:"token,omitempty"`
	LoggedAt  int64 `json:"logged_at,omitempty" bson:"logged_at,omitempty"`
	CreatedAt int64 `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt int64 `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt int64 `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}
