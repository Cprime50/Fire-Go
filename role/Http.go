package role

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MakeAdminHandler(ctx *gin.Context, service AdminService) {
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

	if err := service.MakeAdmin(email); err != nil {
		log.Printf("Error assigning admin role: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %s is now an admin", input.Email)})
}

func RemoveAdminHandler(ctx *gin.Context, service AdminService) {
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

	if err := service.RemoveAdmin(email); err != nil {
		log.Printf("Error assigning user role: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User %s admin rights have been revoked", input.Email)})
}
