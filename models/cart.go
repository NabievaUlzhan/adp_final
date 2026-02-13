package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type CartItem struct {
	ProductID primitive.ObjectID `bson:"product_id"`
	Quantity  int                `bson:"quantity"`
}

type Cart struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID primitive.ObjectID `bson:"user_id"`
	Items  []CartItem         `bson:"items"`
}
