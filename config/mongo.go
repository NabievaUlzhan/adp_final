package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"final4/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var DB *mongo.Database

func ConnectMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB ping failed:", err)
	}

	fmt.Println("✅ MongoDB connected successfully!")
	DB = client.Database("food_store")

	_ = DB.CreateCollection(ctx, "users")
	_ = DB.CreateCollection(ctx, "products")
	_ = DB.CreateCollection(ctx, "orders")
	_ = DB.CreateCollection(ctx, "carts")

	fmt.Println("📦 Collections ready: users, products, orders, carts")
}

func SeedDemoUsers() {
	collection := DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existing models.User
	err := collection.FindOne(ctx, bson.M{"email": "admin@store.com"}).Decode(&existing)
	if err == nil {
		fmt.Println("Demo users already exist")
		return
	}

	hashedAdminPass, _ := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)
	hashedUserPass, _ := bcrypt.GenerateFromPassword([]byte("1234"), bcrypt.DefaultCost)

	admin := models.User{
		Name:     "Admin",
		Email:    "admin@store.com",
		Password: string(hashedAdminPass),
		Role:     "admin",
	}

	customer := models.User{
		Name:     "User",
		Email:    "user@store.com",
		Password: string(hashedUserPass),
		Role:     "customer",
	}

	_, _ = collection.InsertOne(ctx, admin)
	_, _ = collection.InsertOne(ctx, customer)

	fmt.Println("✅ Demo users created: admin@store.com / 1234, user@store.com / 1234")
}
