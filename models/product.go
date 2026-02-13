package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Category  string             `bson:"category"`
	ImageURL  string             `bson:"image_url"`
	Price     float64            `bson:"price"`
	Stock     int                `bson:"stock"`
	ShelfLife int                `bson:"shelf_life"` // days (optional)
	Tags      []string           `bson:"tags"`       // optional: healthy, family, budget...
	CreatedAt time.Time          `bson:"created_at"`
}
