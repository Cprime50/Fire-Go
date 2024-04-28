package profile

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/cprime50/fire-go/middleware"
	"github.com/gin-gonic/gin"
)

func CreateProfileHandler(c *gin.Context, s ProfileService) {
	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	response, err := s.CreateProfile(user.UserID, user.Email)
	if err != nil {
		if errors.Is(err, ErrProfileAlreadyExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": ErrProfileAlreadyExists})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": ErrCreateProfile})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": response.Message, "profile": response.Profile})
}

func UpdateProfileHandler(c *gin.Context, s ProfileService) {
	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err := c.ShouldBindJSON(&UpdateProfileReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidRequestBody})
		return
	}
	response, err := s.UpdateProfile(user.UserID, UpdateProfileReq.Bio, UpdateProfileReq.Username)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": ErrProfileNotFound})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": ErrUpdateProfile})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": response.Message, "profile": response.Profile})
}

func DeleteProfile(c *gin.Context) {
	userID := c.Param("id")

	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if user.Role != "admin" {
		if user.UserID != userID {
			log.Printf("DeleteProfile:Error User with id %s and role  %s not allowed access to delete user with id %s", user.UserID, user.Role, userID)
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
	}

	err := deleteProfile(userID)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			log.Printf("DeleteProfile: Profile not found for userID %s", userID)
			c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		} else {
			log.Printf("DeleteProfile: Error deleting profile for userID %s: %v", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete profile"})
		}
		return
	}
	log.Printf("DeleteProfile: Profile deleted successfully for userID %s", userID)
	c.JSON(http.StatusOK, gin.H{"message": "Profile deleted successfully"})
}

func GetProfile(c *gin.Context) {
	userID := c.Param("id")
	profile, err := getProfileByUserId(userID)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			log.Printf("GetProfile: Profile not found for userID %s", userID)
			c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		} else {
			log.Printf("GetProfile: Error retrieving profile for userID %s: %v", userID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to retrieve profile: %v", err)})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": profile})
}

func GetAllProfiles(c *gin.Context) {
	profiles, err := getAllProfiles()
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			log.Print("GetProfile: Profile not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "No profiles found"})
		} else {
			log.Printf("GetProfile: Error retrieving profile for userID %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve profiles"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"profiles": profiles})
}

func getUserFromCtx(ctx *gin.Context) (*middleware.User, bool) {
	user, exists := ctx.Get("user")
	if !exists {
		return nil, false
	}
	return user.(*middleware.User), true
}
