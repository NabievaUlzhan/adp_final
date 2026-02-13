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
)

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("UserID")
	objID, _ := primitive.ObjectIDFromHex(userID)

	cartCollection := config.DB.Collection("carts")
	orderCollection := config.DB.Collection("orders")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cart models.Cart
	err := cartCollection.FindOne(ctx, bson.M{"user_id": objID}).Decode(&cart)
	if err != nil {
		http.Error(w, "Cart not found", http.StatusBadRequest)
		return
	}

	var orderItems []models.OrderItem

	for _, item := range cart.Items {
		orderItems = append(orderItems, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	order := models.Order{
		UserID:    objID,
		Items:     orderItems,
		CreatedAt: time.Now(),
	}

	_, err = orderCollection.InsertOne(ctx, order)
	if err != nil {
		http.Error(w, "Failed to create order", http.StatusInternalServerError)
		return
	}

	_, err = cartCollection.UpdateOne(
		ctx,
		bson.M{"user_id": objID},
		bson.M{"$set": bson.M{"items": []models.CartItem{}}},
	)
	if err != nil {
		http.Error(w, "Failed to clear cart", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Order created successfully"))
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	role := r.Header.Get("Role")
	userID := r.Header.Get("UserID")

	collection := config.DB.Collection("orders")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var orders []models.Order

	if role == "admin" {
		cursor, _ := collection.Find(ctx, bson.M{})
		cursor.All(ctx, &orders)
	} else {
		userObj, _ := primitive.ObjectIDFromHex(userID)
		cursor, _ := collection.Find(ctx, bson.M{"user_id": userObj})
		cursor.All(ctx, &orders)
	}

	json.NewEncoder(w).Encode(orders)
}

func CancelOrder(w http.ResponseWriter, r *http.Request) {
	orderID := mux.Vars(r)["id"]
	objID, _ := primitive.ObjectIDFromHex(orderID)
	userID := r.Header.Get("UserID")
	userObj, _ := primitive.ObjectIDFromHex(userID)

	collection := config.DB.Collection("orders")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var order models.Order
	err := collection.FindOne(ctx, bson.M{"_id": objID, "user_id": userObj}).Decode(&order)
	if err != nil {
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	if time.Since(order.CreatedAt) > 24*time.Hour {
		http.Error(w, "Too late to cancel (24h limit)", http.StatusBadRequest)
		return
	}

	collection.DeleteOne(ctx, bson.M{"_id": objID, "user_id": userObj})

	w.Write([]byte("Order cancelled"))
}
