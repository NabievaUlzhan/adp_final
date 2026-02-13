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

func GetProducts(w http.ResponseWriter, r *http.Request) {
	collection := config.DB.Collection("products")

	filter := bson.M{}

	category := r.URL.Query().Get("category")
	search := r.URL.Query().Get("search")
	sortBy := r.URL.Query().Get("sort")

	if category != "" {
		filter["category"] = category
	}

	if search != "" {
		filter["name"] = bson.M{"$regex": search, "$options": "i"}
	}

	opts := options.Find()

	if sortBy == "price_asc" {
		opts.SetSort(bson.M{"price": 1})
	}
	if sortBy == "price_desc" {
		opts.SetSort(bson.M{"price": -1})
	}
	if sortBy == "newest" {
		opts.SetSort(bson.M{"created_at": -1})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, _ := collection.Find(ctx, filter, opts)

	var products []models.Product
	cursor.All(ctx, &products)

	json.NewEncoder(w).Encode(products)
}

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	json.NewDecoder(r.Body).Decode(&product)

	product.CreatedAt = time.Now()

	collection := config.DB.Collection("products")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection.InsertOne(ctx, product)

	json.NewEncoder(w).Encode(product)
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	objID, _ := primitive.ObjectIDFromHex(id)

	var updated models.Product
	json.NewDecoder(r.Body).Decode(&updated)

	collection := config.DB.Collection("products")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	set := bson.M{
		"name":      updated.Name,
		"category":  updated.Category,
		"image_url": updated.ImageURL,
		"price":     updated.Price,
		"stock":     updated.Stock,
	}

	collection.UpdateOne(ctx,
		bson.M{"_id": objID},
		bson.M{"$set": set},
	)

	w.Write([]byte("Product updated"))
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	objID, _ := primitive.ObjectIDFromHex(id)

	collection := config.DB.Collection("products")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	collection.DeleteOne(ctx, bson.M{"_id": objID})

	w.Write([]byte("Product deleted"))
}
