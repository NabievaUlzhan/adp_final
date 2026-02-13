package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"final4/config"
	"final4/models"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	user.Role = "customer"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", 500)
		return
	}

	user.Password = string(hashedPassword)

	collection := config.DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, "Error creating user", 500)
		return
	}

	w.Write([]byte("User registered successfully"))
}

func Login(w http.ResponseWriter, r *http.Request) {
	var input models.User
	json.NewDecoder(r.Body).Decode(&input)

	collection := config.DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err := collection.FindOne(ctx, bson.M{
		"email": input.Email,
	}).Decode(&user)

	if err != nil {
		http.Error(w, "Invalid credentials", 401)
		return
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(input.Password),
	)

	if err != nil {
		http.Error(w, "Invalid credentials", 401)
		return
	}

	json.NewEncoder(w).Encode(user)
}
