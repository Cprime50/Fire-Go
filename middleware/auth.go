package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
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
	adminEmail = os.Getenv("ADMIN_EMAIL")
)

func Auth(client *auth.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()

		header := ctx.Request.Header.Get("Authorization")
		if header == "" {
			log.Println("Missing Authorization header")
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		idToken := strings.Split(header, "Bearer ")
		if len(idToken) != 2 {
			log.Println("Invalid Authorization header")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		token, err := client.VerifyIDToken(context.Background(), idToken[1])
		if err != nil {
			log.Printf("Error verifying token. Error: %v\n", err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		email, ok := token.Claims["email"].(string)
		if !ok {
			log.Println("Email claim not found in token")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		role, ok := token.Claims["role"].(string)
		if !ok {
			if email == adminEmail {
				if err := MakeAdmin(ctx, client, adminEmail); err != nil {
					log.Printf("Error making adminEmail %v admin: %v\n", adminEmail, err)
					ctx.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				role = "admin"
			} else {
				if err := MakeUser(ctx, client, token.UID); err != nil {
					log.Printf("Error making user: %v regular user role: %v\n", email, err)
					ctx.AbortWithStatus(http.StatusInternalServerError)
					return
				}
				role = "user"
			}
		}

		user := &User{
			UserID: token.UID,
			Email:  email,
			Role:   role,
		}

		log.Println("Auth time:", time.Since(startTime))
		ctx.Set("user", user)

		log.Println("Successfully authenticated")
		log.Printf("Email: %v\n", user.Email)
		log.Printf("Role: %v\n", user.Role)

		ctx.Next()
	}
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
