package handlers

import (
	"context"
	"encoding/json"
	"final4/config"
	"final4/models"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/genai"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CallGeminiAI(history string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("API Key is missing!")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{APIKey: apiKey})
	if err != nil {
		return "", err
	}

	prompt := fmt.Sprintf(`
You are a recommendation system for a bakery food store which sells products as pastry, cake, bread, dessert etc.
Based on this users order history, recommend 3 products user might buy next or buy again(products which user buy usually). 
Write names of these products and why you decided that without symbold like * # etc.
User data:
%s
`, history)

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-3-flash-preview",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		return "", err
	}

	return result.Text(), nil
}

func GetRecommendations(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("UserID")

	orderCollection := config.DB.Collection("orders")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userObj, _ := primitive.ObjectIDFromHex(userID)

	cursor, _ := orderCollection.Find(ctx, bson.M{"user_id": userObj})

	var orders []models.Order
	cursor.All(ctx, &orders)

	if len(orders) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"mode":  "default",
			"items": []string{"Strawberry Cake", "Macarons Box", "Donuts"},
		})
		return
	}

	history := ""
	for _, order := range orders {
		for _, item := range order.Items {
			history += item.ProductID.Hex() + " "
		}
	}

	result, err := CallGeminiAI(history)
	if err != nil {
		http.Error(w, "AI error", 500)
		return
	}

	w.Write([]byte(result))
}
