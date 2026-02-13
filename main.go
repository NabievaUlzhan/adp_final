package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"final4/config"
	"final4/handlers"
	"final4/middleware"
	"final4/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func serveTemplate(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/"+name)
	}
}

func seedProductsIfEmpty() {
	col := config.DB.Collection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := col.CountDocuments(ctx, bson.M{})
	if err != nil || count > 0 {
		return
	}

	now := time.Now()
	products := []interface{}{
		models.Product{Name: "Butter Croissant", Category: "pastry", ImageURL: "/static/images/bread.svg", Price: 2.50, Stock: 80, CreatedAt: now},
		models.Product{Name: "French Baguette", Category: "bread", ImageURL: "/static/images/bread.svg", Price: 1.80, Stock: 120, CreatedAt: now},
		models.Product{Name: "Chocolate Donut", Category: "dessert", ImageURL: "/static/images/veggie.svg", Price: 1.60, Stock: 90, CreatedAt: now},
		models.Product{Name: "Cheesecake Slice", Category: "cake", ImageURL: "/static/images/apple.svg", Price: 4.90, Stock: 40, CreatedAt: now},
		models.Product{Name: "Cinnamon Roll", Category: "pastry", ImageURL: "/static/images/rice.svg", Price: 2.90, Stock: 70, CreatedAt: now},
		models.Product{Name: "Latte", Category: "drink", ImageURL: "/static/images/milk.svg", Price: 3.20, Stock: 200, CreatedAt: now},
	}

	_, _ = col.InsertMany(ctx, products)
	fmt.Println("✅ Demo products seeded")
}

func main() {
	// 1) Mongo
	config.ConnectMongo()
	config.SeedDemoUsers() // admin@store.com / 1234, user@store.com / 1234
	seedProductsIfEmpty()  // demo products if empty

	// 2) Router
	r := mux.NewRouter()

	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Pages (HTML)
	r.HandleFunc("/", serveTemplate("index.html")).Methods("GET")
	r.HandleFunc("/login", serveTemplate("login.html")).Methods("GET")
	r.HandleFunc("/register", serveTemplate("register.html")).Methods("GET")
	r.HandleFunc("/cart", serveTemplate("cart.html")).Methods("GET")
	r.HandleFunc("/orders", serveTemplate("orders.html")).Methods("GET")
	r.HandleFunc("/create-product", serveTemplate("create_product.html")).Methods("GET")
	r.HandleFunc("/update-product", serveTemplate("update_product.html")).Methods("GET")

	// API (JSON)
	// AUTH
	r.HandleFunc("/api/login", handlers.Login).Methods("POST")
	r.HandleFunc("/api/register", handlers.Register).Methods("POST")

	// PRODUCTS
	r.HandleFunc("/api/products", handlers.GetProducts).Methods("GET")
	r.HandleFunc("/api/products", middleware.AdminOnly(handlers.CreateProduct)).Methods("POST")
	r.HandleFunc("/api/products/{id}", middleware.AdminOnly(handlers.UpdateProduct)).Methods("PUT")
	r.HandleFunc("/api/products/{id}", middleware.AdminOnly(handlers.DeleteProduct)).Methods("DELETE")

	// CART
	r.HandleFunc("/api/cart", middleware.CustomerOnly(handlers.GetCart)).Methods("GET")
	r.HandleFunc("/api/cart/add", middleware.CustomerOnly(handlers.AddToCart)).Methods("POST")
	r.HandleFunc("/api/cart/update/{id}", middleware.CustomerOnly(handlers.UpdateCartQuantity)).Methods("PUT")
	r.HandleFunc("/api/cart/remove/{id}", middleware.CustomerOnly(handlers.RemoveFromCart)).Methods("DELETE")

	// ORDERS
	r.HandleFunc("/api/orders", middleware.CustomerOnly(handlers.CreateOrder)).Methods("POST")
	r.HandleFunc("/api/orders", handlers.GetOrders).Methods("GET")
	r.HandleFunc("/api/orders/{id}", middleware.CustomerOnly(handlers.CancelOrder)).Methods("DELETE")

	// AI
	r.HandleFunc("/api/recommendations", middleware.CustomerOnly(handlers.GetRecommendations)).Methods("GET")


	// SMART FEATURES
	r.HandleFunc("/api/smart-basket", middleware.CustomerOnly(handlers.WeeklySmartBasket)).Methods("GET")
	r.HandleFunc("/api/predict-restock", middleware.CustomerOnly(handlers.PredictRestock)).Methods("GET")
	r.HandleFunc("/api/anti-waste", handlers.AntiWasteMode).Methods("GET")
	fmt.Println("🚀 Server running on http://localhost:8080")
	_ = http.ListenAndServe(":8080", r)
}
