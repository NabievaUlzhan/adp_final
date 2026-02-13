package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"final4/config"
	"final4/models"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
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

	url := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=%s",
		apiKey,
	)

	prompt := fmt.Sprintf(`
You are a recommendation system for a food store.
Based on this purchase history, recommend 3 products user might buy next.
User data:
%s
`, history)

	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": prompt},
				},
			},
		},
	}

	jsonBody, _ := json.Marshal(body)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	return fmt.Sprintf("%v", result), nil
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
			"items": []string{},
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
