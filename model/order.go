package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	TableID int32 `json:"table_id,omitempty" bson:"table_id,omitempty"`
	StaffUserID string `json:"staff_user_id,omitempty" bson:"staff_user_id,omitempty"`
	GuestUserID string `json:"guest_user_id,omitempty" bson:"guest_user_id,omitempty"`
	Items []Item `json:"items,omitempty" bson:"items,omitempty"`
	TotalPrice float32 `json:"total_price,omitempty" bson:"total_price,omitempty"`
	CreatedAt int64 `json:"created_at,omitempty" bson:"created_at,omitempty"`
	ReadyAt int64 `json:"ready_at,omitempty" bson:"ready_at,omitempty"`
	DeliveredAt int64 `json:"delivered_at,omitempty" bson:"delivered_at,omitempty"`
	PayedAt int64 `json:"payed_at,omitempty" bson:"payed_at,omitempty"`
	UpdatedAt int64 `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt int64 `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}
