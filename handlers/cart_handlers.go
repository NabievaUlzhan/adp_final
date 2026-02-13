package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"final4/config"
	"final4/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("UserID")
	userObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid UserID", http.StatusBadRequest)
		return
	}

	collection := config.DB.Collection("carts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cart models.Cart
	err = collection.FindOne(ctx, bson.M{"user_id": userObj}).Decode(&cart)
	if err != nil {
		json.NewEncoder(w).Encode(models.Cart{Items: []models.CartItem{}})
		return
	}

	json.NewEncoder(w).Encode(cart)
}

func AddToCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("UserID")
	userObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid UserID", http.StatusBadRequest)
		return
	}

	var body struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	if body.ProductID == "" || body.Quantity < 1 {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	productObj, err := primitive.ObjectIDFromHex(body.ProductID)
	if err != nil {
		http.Error(w, "Invalid product_id", http.StatusBadRequest)
		return
	}

	item := models.CartItem{
		ProductID: productObj,
		Quantity:  body.Quantity,
	}

	collection := config.DB.Collection("carts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Update().SetUpsert(true)

	_, err = collection.UpdateOne(
		ctx,
		bson.M{"user_id": userObj},
		bson.M{"$push": bson.M{"items": item}},
		opts,
	)
	if err != nil {
		http.Error(w, "Failed to add to cart", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Added"))
}

func UpdateCartQuantity(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("UserID")
	productID := mux.Vars(r)["id"]

	userObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid UserID", http.StatusBadRequest)
		return
	}
	productObj, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		http.Error(w, "Invalid product id", http.StatusBadRequest)
		return
	}

	var body struct {
		Quantity int `json:"quantity"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)
	if body.Quantity < 1 {
		http.Error(w, "Quantity must be >= 1", http.StatusBadRequest)
		return
	}

	collection := config.DB.Collection("carts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, _ = collection.UpdateOne(ctx,
		bson.M{
			"user_id":          userObj,
			"items.product_id": productObj,
		},
		bson.M{
			"$set": bson.M{
				"items.$.quantity": body.Quantity,
			},
		},
	)

	w.Write([]byte("Updated"))
}

func RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("UserID")
	productID := mux.Vars(r)["id"]

	userObj, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		http.Error(w, "Invalid UserID", http.StatusBadRequest)
		return
	}
	productObj, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		http.Error(w, "Invalid product id", http.StatusBadRequest)
		return
	}

	collection := config.DB.Collection("carts")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, _ = collection.UpdateOne(ctx,
		bson.M{"user_id": userObj},
		bson.M{
			"$pull": bson.M{
				"items": bson.M{"product_id": productObj},
			},
		},
	)

	w.Write([]byte("Removed"))
}
