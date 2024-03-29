package main

import (
	"log"
	"log/slog"
	"os"

	"firebase.google.com/go/v4/auth"

	"github.com/cprime50/fire-go/db"
	"github.com/cprime50/fire-go/middleware"
	"github.com/cprime50/fire-go/src"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv()

	// Initialize Firebase authentication middleware
	client, err := middleware.InitAuth()
	if err != nil {
		log.Fatalf("Error initializing Firebase auth: %v", err)
	}

	//Connect db
	Db, err := db.Connect()
	if err != nil {
		slog.Error("Error opening database", "db.Connect", err)
		log.Fatal("Error connecting to Db", err)
	}
	log.Println("Database connected successfully")

	// migrations
	log.Printf("Migrations Started")
	err = db.Migrate(Db)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Db.Close()

	r := gin.Default()
	r.Use(cors.Default())

	// Register routes
	RegisterRoutes(r, client)
	RegisterAdminRoutes(r, client)

	// Set port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Gin server is running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start Gin server: %v", err)
	}
}

func loadEnv() {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	log.Println(".env file loaded successfully")
}

// Auth routes
func RegisterRoutes(r *gin.Engine, client *auth.Client) {
	profileRoutes := r.Group("/profile")
	profileRoutes.Use(middleware.Auth(client))
	{

		profileRoutes.POST("/create", src.CreateProfile)
		profileRoutes.PUT("/update", src.UpdateProfile)
		profileRoutes.DELETE("/delete/:id", src.DeleteProfile)
		profileRoutes.GET("/:id", src.GetProfile)
	}
	quoteRoutes := r.Group("/quote")
	quoteRoutes.Use(middleware.Auth(client))
	{
		quoteRoutes.GET("/", src.GetQuotes)
		quoteRoutes.POST("/create", src.CreateQuote)
		quoteRoutes.PUT("/update", src.UpdateQuote)
		quoteRoutes.DELETE("/delete/:id", src.DeleteQuote)
		quoteRoutes.GET("/:profile-id", src.GetQuotesByUserId)

	}
}

// Admin routes
func RegisterAdminRoutes(r *gin.Engine, client *auth.Client) {

	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.Auth(client), middleware.RoleAuth("admin"))
	{
		adminRoutes.GET("/profiles", src.GetAllProfiles)
		adminRoutes.POST("/quote/approve/:id", src.ApproveQuote)
		adminRoutes.GET("/quotes/unapproved", src.GetUnapprovedQuotes)
		adminRoutes.POST("/make", func(ctx *gin.Context) {
			src.MakeAdmin(ctx, client)
		})
		adminRoutes.DELETE("/remove", func(ctx *gin.Context) {
			src.RemoveAdmin(ctx, client)
		})

	}

}
