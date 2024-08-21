package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blog struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Title     string             `json:"title" validate:"required,min=2,max=30"`
	Content   string             `json:"content" validate:"required"`
	Category  string             `json:"category"`
	Tags      []string           `json:"tags"`
	CreatedAt time.Time          `json:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt"`
}
