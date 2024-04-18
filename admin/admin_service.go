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

func RemoveAdmin(ctx *gin.Context, client *auth.Client) {
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
