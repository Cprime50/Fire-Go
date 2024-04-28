package admin

import (
	"fmt"
	"log"
	"net/http"

	"firebase.google.com/go/v4/auth"
	"github.com/cprime50/fire-go/middleware"
	"github.com/gin-gonic/gin"
)

type EmailInput struct {
	Email string `json:"email"`
}

func MakeAdmin(ctx *gin.Context, client *auth.Client) {
	var input EmailInput
	if err := ctx.BindJSON(&input); err != nil {
		log.Printf("Error binding JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	email, err := validateInput(ctx, input)
	if err != nil {
		log.Printf("Error validating email: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := middleware.AssignRole(ctx.Request.Context(), client, email, "admin"); err != nil {
		log.Printf("Error assigning admin role: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %s is now an admin", input.Email)})
}

func RemoveAdmin(ctx *gin.Context, client *auth.Client) {
	var input EmailInput
	if err := ctx.BindJSON(&input); err != nil {
		log.Printf("Error binding JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	email, err := validateInput(ctx, input)
	if err != nil {
		log.Printf("Error validating email: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := middleware.AssignRole(ctx.Request.Context(), client, email, "user"); err != nil {
		log.Printf("Error assigning user role: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %s admin rights have been revoked", input.Email)})
}
