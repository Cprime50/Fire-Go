package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

func RoleAuth(requiredRole string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userValue, exists := ctx.Get("user")
		if !exists {
			log.Println("User not found in context")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		user, ok := userValue.(*User)
		if !ok || user == nil {
			log.Println("Invalid user data in context")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if user.Role == "" {
			log.Println("User role not set")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if user.Role != requiredRole {
			log.Printf("User with email %s and role %s tried to access a route that was for the %s role only",
				user.Email, user.Role, requiredRole)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		log.Printf("User with email %s and role %s authorized", user.Email, user.Role)
		ctx.Next()
	}
}

func MakeAdmin(ctx context.Context, client *auth.Client, email string) error {
	user, err := client.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("MakeAdmin Error: User with email %s not found", email)
	}

	currentCustomClaims := user.CustomClaims
	if currentCustomClaims == nil {
		currentCustomClaims = map[string]interface{}{}
	}
	currentCustomClaims["role"] = "admin"

	if _, found := currentCustomClaims["admin"]; !found {
		currentCustomClaims["admin"] = true
	}

	if err := client.SetCustomUserClaims(ctx, user.UID, currentCustomClaims); err != nil {
		return fmt.Errorf("MakeAdmin Error: Error setting custom claims %w", err)
	}

	return nil
}

func MakeUser(ctx context.Context, client *auth.Client, userID string) error {
	user, err := client.GetUser(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("MakeUser Error: User with ID %s not found", userID)
	}
	currentCustomClaims := user.CustomClaims
	if currentCustomClaims == nil {
		currentCustomClaims = map[string]interface{}{}
	}
	currentCustomClaims["role"] = "user"

	if err := client.SetCustomUserClaims(ctx, user.UID, currentCustomClaims); err != nil {
		return fmt.Errorf("MakeUser Error: Error setting custom claims: %w", err)
	}

	return nil
}

func RemoveAdmin(ctx context.Context, client *auth.Client, email string) error {
	user, err := client.GetUserByEmail(ctx, email)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("RemoveAdmin Error: User with email %s not found", email)
	}

	currentCustomClaims := user.CustomClaims
	if currentCustomClaims == nil {
		currentCustomClaims = map[string]interface{}{}
	}
	currentCustomClaims["role"] = "user"

	delete(currentCustomClaims, "admin")

	if err := client.SetCustomUserClaims(ctx, user.UID, currentCustomClaims); err != nil {
		return fmt.Errorf("RemoveAdmin Error: Error setting custom claims %w", err)
	}

	return nil
}
