package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"firebase.google.com/go/v4/auth"

	"github.com/cprime50/fire-go/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	// Serve static html file to test firebase auth in the browser
	http.Handle("/", http.FileServer(http.Dir(".")))

	loadEnv()

	go func() {
		// Serve Gin server
		r := gin.Default()
		r.Use(cors.Default())
		client, err := middleware.InitAuth()
		if err != nil {
			log.Println(err)
			return
		}

		RegisterRoutes(r, client)
		RegisterAdminRoutes(r, client)

		// Set port
		port := os.Getenv("PORT")
		if port == "" {
			port = "localhost:8080" // Default port
		}

		// Start server
		log.Printf("Gin server is running on port %s", port)
		if err := r.Run(port); err != nil {
			log.Fatalf("Failed to start Gin server: %v", err)
		}
	}()

	// Start static html server
	log.Println("Static file server is running on port 8000")
	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatalf("Failed to start static file server: %v", err)
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
	routes := r.Group("/profile")
	routes.Use(middleware.Auth(client))
	{
		routes.GET("/", GetProfile)
		routes.GET("/time", func(c *gin.Context) {
			currentTime := time.Now()
			c.JSON(http.StatusOK, gin.H{"time is": currentTime.Format(time.RFC3339)})
		})
	}
}

func GetProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		log.Println("Error Authenticated user should exist in context \n User does not exist in ctx")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	profile := user.(*middleware.User)

	c.JSON(http.StatusOK, profile)
}

// Admin routes
func RegisterAdminRoutes(r *gin.Engine, client *auth.Client) {

	adminRoutes := r.Group("/admin")
	adminRoutes.Use(middleware.Auth(client), middleware.RoleAuth("admin"))
	{
		adminRoutes.POST("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "welcome admin"})
		})
		adminRoutes.POST("/make", func(ctx *gin.Context) {
			makeAdmin(ctx, client)
		})
		adminRoutes.DELETE("/remove", func(ctx *gin.Context) {
			removeAdmin(ctx, client)
		})
	}

}

type EmailInput struct {
	Email string `json:"email"`
}

func makeAdmin(ctx *gin.Context, client *auth.Client) {
	var input EmailInput

	email, err := validateInput(ctx, input)
	if err != nil {
		log.Print("error validating email: invalid format", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = middleware.MakeAdmin(ctx.Request.Context(), client, email)
	if err != nil {
		log.Print("error making admin:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %s is now an admin", input.Email)})
}

func removeAdmin(ctx *gin.Context, client *auth.Client) {
	var input EmailInput

	email, err := validateInput(ctx, input)
	if err != nil {
		log.Print("error validating email: invalid format", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = middleware.RemoveAdmin(ctx.Request.Context(), client, email)
	if err != nil {
		log.Print("error making admin: Invalid email format")
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %s admin rights have been revoked", input.Email)})
}

func validateInput(ctx *gin.Context, input EmailInput) (string, error) {
	if err := ctx.BindJSON(&input); err != nil {
		return "", fmt.Errorf("invalid JSON format")
	}
	emailOk := ValidateEmail(input.Email)
	if !emailOk {
		return "", fmt.Errorf("invalid email format")
	}
	return input.Email, nil
}

// Regex for email validation
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
