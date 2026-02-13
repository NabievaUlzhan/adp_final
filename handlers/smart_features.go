package handlers

import (
	"context"
	"encoding/json"
	"final4/config"
	"final4/models"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func WeeklySmartBasket(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("UserID")
	userObj, _ := primitive.ObjectIDFromHex(userID)

	orderCollection := config.DB.Collection("orders")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, _ := orderCollection.Find(ctx, bson.M{"user_id": userObj})

	var orders []models.Order
	cursor.All(ctx, &orders)

	frequency := make(map[string]int)

	for _, order := range orders {
		for _, item := range order.Items {
			frequency[item.ProductID.Hex()]++
		}
	}

	type Recommendation struct {
		ProductID string `json:"product_id"`
		Score     int    `json:"score"`
	}

	var recs []Recommendation

	for id, count := range frequency {
		if count >= 2 {
			recs = append(recs, Recommendation{id, count})
		}
	}

	json.NewEncoder(w).Encode(recs)
}

func PredictRestock(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("UserID")
	userObj, _ := primitive.ObjectIDFromHex(userID)

	orderCollection := config.DB.Collection("orders")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, _ := orderCollection.Find(ctx, bson.M{"user_id": userObj})

	var orders []models.Order
	cursor.All(ctx, &orders)

	now := time.Now()
	var reminders []string

	for _, order := range orders {
		days := int(now.Sub(order.CreatedAt).Hours() / 24)
		if days > 7 {
			reminders = append(reminders, "You may need to restock items from order "+order.ID.Hex())
		}
	}

	json.NewEncoder(w).Encode(reminders)
}

func AntiWasteMode(w http.ResponseWriter, r *http.Request) {
	productCollection := config.DB.Collection("products")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, _ := productCollection.Find(ctx, bson.M{"shelf_life": bson.M{"$lt": 5}})

	var products []models.Product
	cursor.All(ctx, &products)

	json.NewEncoder(w).Encode(products)
}
