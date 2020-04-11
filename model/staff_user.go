package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type StaffUser struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName string `json:"first_name,omitempty" bson:"first_name,omitempty"`
	LastName string `json:"last_name,omitempty" bson:"last_name,omitempty"`
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
	Token  int64 `json:"token,omitempty" bson:"token,omitempty"`
	IsAdmin bool `json:"is_admin, omitempty" bson:"is_admin,omitempty"`
	CreatedAt int64 `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt int64 `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt int64 `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}
