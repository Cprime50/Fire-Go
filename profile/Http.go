package profile

import (
	"errors"
	"net/http"

	"github.com/cprime50/fire-go/middleware"
	"github.com/gin-gonic/gin"
)

func CreateProfileHandler(c *gin.Context, s ProfileServiceImpl) {
	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	response, err := s.CreateProfile(user.UserID, user.Email)
	if err != nil {
		if errors.Is(err, ErrProfileAlreadyExists) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": response.Message, "profile": response.Profile})
}

func UpdateProfileHandler(c *gin.Context, s ProfileServiceImpl) {
	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if err := c.ShouldBindJSON(&UpdateProfileReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	response, err := s.UpdateProfile(user.UserID, UpdateProfileReq.Bio, UpdateProfileReq.Username)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": response.Message, "profile": response.Profile})
}

func DeleteProfileHandler(c *gin.Context, service ProfileServiceImpl) {
	profileId := c.Param("id")

	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	err := service.DeleteProfile(profileId, user.Role, user.UserID)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile deleted successfully"})
}

func GetProfileHandler(c *gin.Context, service ProfileServiceImpl) {
	userID := c.Param("id")

	profile, err := service.GetProfile(userID)
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": profile})
}

func GetAllProfilesHandler(c *gin.Context, service ProfileServiceImpl) {
	profiles, err := service.GetAllProfiles()
	if err != nil {
		if errors.Is(err, ErrProfileNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
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
