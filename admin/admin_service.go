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
// @Description This endpoint allows an admin to promote a user to an admin role.
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "ID token"
// @Param body body EmailInput true "The email of the user to be promoted to admin."
// @Success 200 {object} map[string]string "{\"message\":\"User example@example.com is now an admin\"}"
// @Failure 400 {object} map[string]string "{\"error\":\"error validating email: invalid format\"}"
// @Failure 500 {object} map[string]string "{\"error\":\"error making admin:\"}"
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

// @Summary Remove user admin rights
// @Description This endpoint allows an admin to revoke admin rights from a user.
// @Tags Admin
// @Accept json
// @Produce json
// @Param Authorization header string true "ID token"
// @Param body body EmailInput true "The email of the user to have admin rights revoked."
// @Success 200 {object} map[string]string "{\"message\":\"User example@example.com admin rights have been revoked\"}"
// @Failure 400 {object} map[string]string "{\"error\":\"error validating email: invalid format\"}"
// @Failure 500 {object} map[string]string "{\"error\":\"error removing admin rights:\"}"
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
