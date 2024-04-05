package profile

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cprime50/fire-go/middleware"
	"github.com/gin-gonic/gin"
)

// @Summary Create profile
// @Description Create a new user profile
// @Tags profile
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT token"
// @Success 201 {string} string "Profile created successfully"
// @Failure 400 {string} string "Profile already exists"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Failed to create profile"
// @Router /profile/create [post]
func CreateProfile(c *gin.Context) {
	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userData := user.(*middleware.User)
	existingProfile, err := getProfileByUserId(userData.UserID)
	if err == nil && existingProfile != nil {
		log.Printf("Profile already exist for user %s with emial %s", userData.UserID, userData.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Profile already exists"})
		return
	} else if err != nil && !errors.Is(err, ErrProfileNotFound) {
		log.Printf("Error checking profile existence: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile"})
		return
	}

	username, err := generateUsername(userData.Email)
	if err != nil {
		log.Printf("Error generating username: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile"})
		return
	}

	err = createProfile(&Profile{
		UserId:   userData.UserID,
		Email:    userData.Email,
		UserName: username,
		Bio:      "",
	})
	if err != nil {
		log.Printf("Error creating profile: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create profile"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Profile created successfully", "user_email": userData.Email})
}

// @Summary Update profile
// @Description Update the bio and username of the user profile
// @Tags profile
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT token"
// @Param bio body string true "New bio"
// @Param username body string true "New username"
// @Success 200 {string} string "Profile updated successfully"
// @Failure 400 {string} string "Invalid request body"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Profile not found"
// @Failure 500 {string} string "Failed to update profile"
// @Router /profile/update [put]
func UpdateProfile(c *gin.Context) {
	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userData := user.(*middleware.User)

	var requestBody struct {
		Bio      string `json:"bio"`
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Update the profile bio in the database
	err := updateProdile(&Profile{
		UserId:    userData.UserID,
		Bio:       requestBody.Bio,
		UserName:  requestBody.Username,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to update profile: %v", err)})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

// @Summary Delete a user's profile
// @Description Delete a user's profile by ID
// @Tags profile
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param Authorization header string true "JWT token"
// @Success 200 {string} string "Profile deleted successfully"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Not authorized"
// @Failure 404 {string} string "Profile not found"
// @Failure 500 {string} string "Failed to delete profile"
// @Router /profile/delete/{id} [delete]
func DeleteProfile(c *gin.Context) {
	userID := c.Param("id")

	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userData := user.(*middleware.User)
	if userData.Role != "admin" {
		if userData.UserID != userID {
			log.Printf("DeleteProfile:Error User with id %s and role  %s not allowed access to delete user with id %s", userData.UserID, userData.Role, userID)
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

// @Summary Get a user's profile
// @Description Retrieve a user's profile by ID
// @Tags profile
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param Authorization header string true "JWT token"
// @Success 200 {string} string "User profile data"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Profile not found"
// @Failure 500 {string} string "Failed to retrieve profile"
// @Router /profile/{id} [get]
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

// @Summary Get all profiles
// @Description Retrieve all user profiles
// @Tags admin
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT token"
// @Success 200 {array} Profile "List of user profiles"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "No profiles found"
// @Failure 500 {string} string "Failed to retrieve profiles"
// @Router /admin/profile [get]
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

func getUserFromCtx(ctx *gin.Context) (any, bool) {
	user, exists := ctx.Get("user")
	if !exists {
		return "", false
	}
	return user, true
}
