package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Item struct {
	ID primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ItemName string `json:"item_name,omitempty" bson:"item_name,omitempty"`
	ItemDescription string `json:"item_description,omitempty" bson:"item_description,omitempty"`
	ItemType string `json:"item_type,omitempty" bson:"item_type,omitempty"`
	ItemImg string `json:"item_img,omitempty" bson:"item_img,omitempty"`
	ItemPrice float32 `json:"item_price,omitempty" bson:"item_price,omitempty"`
	EstimatePrepareTime  int64 `json:"estimate_prepare_time,omitempty" bson:"estimate_prepare_time,omitempty"`
	CreatedAt int64 `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt int64 `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt int64 `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}
