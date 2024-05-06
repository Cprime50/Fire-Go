package main

import (
	"log"
	"log/slog"
	"os"

	"firebase.google.com/go/v4/auth"

	"github.com/cprime50/fire-go/db"
	"github.com/cprime50/fire-go/role"

	"github.com/cprime50/fire-go/middleware"
	"github.com/cprime50/fire-go/profile"
	"github.com/cprime50/fire-go/quote"

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

func RegisterRoutes(r *gin.Engine, client *auth.Client) {
	s := profile.ProfileServiceImpl{} // Corrected instantiation

	profileRoutes := r.Group("/profile")
	profileRoutes.Use(middleware.Auth(client))
	{
		profileRoutes.POST("/create", func(c *gin.Context) {
			profile.CreateProfileHandler(c, s)
		})
		profileRoutes.PUT("/update", func(c *gin.Context) {
			profile.UpdateProfileHandler(c, s)
		})
		profileRoutes.DELETE("/delete/:id", func(c *gin.Context) {
			profile.DeleteProfileHandler(c, s)
		})
		profileRoutes.GET("/:id", func(c *gin.Context) {
			profile.GetProfileHandler(c, s)
		})
	}

	quoteService := &quote.QuoteServiceImpl{}

	quoteRoutes := r.Group("/quote")
	quoteRoutes.Use(middleware.Auth(client))
	{
		quoteRoutes.POST("/create", func(c *gin.Context) {
			quote.CreateQuoteHandler(c, quoteService)
		})
		quoteRoutes.PUT("/update", func(c *gin.Context) {
			quote.UpdateQuoteHandler(c, quoteService)
		})
		quoteRoutes.DELETE("/delete/:id", func(c *gin.Context) {
			quote.DeleteQuoteHandler(c, quoteService)
		})
		quoteRoutes.GET("/", func(c *gin.Context) {
			quote.GetQuotesHandler(c, quoteService)
		})
		quoteRoutes.GET("/quotes/:profile-id", func(c *gin.Context) {
			quote.GetQuotesByUserIdHandler(c, quoteService)
		})
		quoteRoutes.PUT("/approve/:id", func(c *gin.Context) {
			quote.ApproveQuoteHandler(c, quoteService)
		})
		quoteRoutes.GET("/unapproved", func(c *gin.Context) {
			quote.GetUnapprovedQuotesHandler(c, quoteService)
		})
	}
}

// Admin routes
func RegisterAdminRoutes(r *gin.Engine, client *auth.Client) {
	profileService := profile.ProfileServiceImpl{}
	quoteService := &quote.QuoteServiceImpl{}
	adminService := role.NewAdminService(client)

	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.Auth(client), middleware.RoleAuth("admin"))
	{
		adminRoutes.GET("/profiles", func(c *gin.Context) {
			profile.GetAllProfilesHandler(c, profileService)
		})
		// Update this line to use the separated handler for approving quotes
		adminRoutes.POST("/quote/approve/:id", func(c *gin.Context) {
			quote.ApproveQuoteHandler(c, quoteService)
		})
		// Update this line to use the separated handler for getting unapproved quotes
		adminRoutes.GET("/quote/unapproved", func(c *gin.Context) {
			quote.GetUnapprovedQuotesHandler(c, quoteService)
		})
		adminRoutes.POST("/make", func(ctx *gin.Context) {
			role.MakeAdminHandler(ctx, adminService)
		})
		adminRoutes.DELETE("/remove", func(ctx *gin.Context) {
			role.RemoveAdminHandler(ctx, adminService)
		})
	}
}
