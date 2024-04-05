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

// @Summary Make user admin
// @Description Make a user an admin by email
// @Tags admin
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT token"
// @Param email body EmailInput true "User email to make admin"
// @Success 200 {string} string "User is now an admin"
// @Failure 400 {string} string "Invalid request body"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /admin/make [post]
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

// @Summary Remove admin rights
// @Description Remove admin rights from a user by email
// @Tags admin
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT token"
// @Param email body EmailInput true "User email to remove admin rights"
// @Success 200 {string} string "User admin rights have been revoked"
// @Failure 400 {string} string "Invalid request body"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal Server Error"
// @Router /admin/remove [delete]
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
