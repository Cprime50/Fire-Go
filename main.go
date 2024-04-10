package main

import (
	"log"
	"log/slog"
	"os"

	"firebase.google.com/go/v4/auth"

	"github.com/cprime50/fire-go/admin"
	"github.com/cprime50/fire-go/db"
	docs "github.com/cprime50/fire-go/docs"
	"github.com/cprime50/fire-go/middleware"
	"github.com/cprime50/fire-go/profile"
	"github.com/cprime50/fire-go/quote"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title FireGo
// @description A server for a simple Go application

// @host cprime50.github.io
// @BasePath /api
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

	// Configure SwaggerInfo
	docs.SwaggerInfo.BasePath = "/api"

	r := gin.Default()
	r.Use(cors.Default())

	// Register routes
	RegisterRoutes(r, client)
	RegisterAdminRoutes(r, client)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

func RegisterRoutes(r *gin.Engine, client *auth.Client) {
	profileRoutes := r.Group("/profile")
	profileRoutes.Use(middleware.Auth(client))
	{

		profileRoutes.POST("/create", profile.CreateProfile)
		profileRoutes.PUT("/update", profile.UpdateProfile)
		profileRoutes.DELETE("/delete/:id", profile.DeleteProfile)
		profileRoutes.GET("/:id", profile.GetProfile)
	}

	quoteRoutes := r.Group("/quote")
	quoteRoutes.Use(middleware.Auth(client))
	{
		quoteRoutes.GET("/", quote.GetQuotes)
		quoteRoutes.POST("/create", quote.CreateQuote)
		quoteRoutes.PUT("/update", quote.UpdateQuote)
		quoteRoutes.DELETE("/delete/:id", quote.DeleteQuote)
		quoteRoutes.GET("/:profile-id", quote.GetQuotesByUserId)

	}
}

// Admin routes
func RegisterAdminRoutes(r *gin.Engine, client *auth.Client) {

	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.Auth(client), middleware.RoleAuth("admin"))
	{
		adminRoutes.GET("/profiles", profile.GetAllProfiles)
		adminRoutes.POST("/quote/approve/:id", quote.ApproveQuote)
		adminRoutes.GET("/quote/unapproved", quote.GetUnapprovedQuotes)
		adminRoutes.POST("/make", func(ctx *gin.Context) {
			admin.MakeAdmin(ctx, client)
		})
		adminRoutes.DELETE("/remove", func(ctx *gin.Context) {
			admin.RemoveAdmin(ctx, client)
		})

	}

}
