package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"

	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

type User struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
}

var (
	tokenCache = make(map[string]*auth.Token)
	cacheMutex sync.Mutex
)

func Auth(client *auth.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()

		header := ctx.Request.Header.Get("Authorization")
		if header == "" {
			log.Println("Missing Authorization header")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, Invalid Token"})
			return
		}
		idToken := strings.Split(header, "Bearer ")
		if len(idToken) != 2 {
			log.Println("Invalid Authorization header")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, Invalid Token"})
			return
		}
		tokenID := idToken[1]

		cacheMutex.Lock()
		defer cacheMutex.Unlock()

		// Check if token is in cache
		if token, ok := tokenCache[tokenID]; ok {
			processToken(ctx, client, token)
			log.Println("Auth time:", time.Since(startTime))
			return
		}

		// Token not in cache, verify it
		token, err := client.VerifyIDToken(context.Background(), tokenID)
		if err != nil {
			log.Printf("Error verifying token. Error: %v\n", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, Invalid Token"})
			return
		}

		// Cache verified token
		tokenCache[tokenID] = token

		processToken(ctx, client, token)
		log.Println("Auth time:", time.Since(startTime))
	}
}

func processToken(ctx *gin.Context, client *auth.Client, token *auth.Token) {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	email, ok := token.Claims["email"].(string)
	if !ok {
		log.Println("Email claim not found in token")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, Invalid Token"})
		return
	}
	log.Println("auth email is ", email)

	role, ok := token.Claims["role"].(string)
	if email != adminEmail {
		log.Println("adminEmail", adminEmail)
		log.Println("Email", email)
	}
	if email == adminEmail && role == "user" || !ok {
		if err := MakeAdmin(ctx, client, adminEmail); err != nil {
			log.Printf("Error making adminEmail admin: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something wrong happened"})
			return
		}
		role = "admin"
	}
	if !ok {
		if err := MakeUser(ctx, client, token.UID); err != nil {
			log.Printf("Error making user regular user: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Something wrong happened"})
			return
		}
		role = "user"
	}

	user := &User{
		UserID: token.UID,
		Email:  email,
		Role:   role,
	}

	ctx.Set("user", user)

	log.Println("Successfully authenticated")
	log.Printf("Email: %v\n", user.Email)
	log.Printf("Role: %v\n", user.Role)

	ctx.Next()
}

func InitAuth() (*auth.Client, error) {
	var firebaseCredFile = os.Getenv("FIREBASE_KEY")
	opt := option.WithCredentialsFile(firebaseCredFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v", err)
		return nil, err
	}

	client, errAuth := app.Auth(context.Background())
	if errAuth != nil {
		log.Fatalf("error initializing firebase auth: %v", errAuth)
		return nil, errAuth
	}

	return client, nil
}
