package quote

import (
	"errors"
	"log"
	"net/http"

	"github.com/cprime50/fire-go/middleware"
	"github.com/gin-gonic/gin"
)

type QuoteRequest struct {
	Quote string `json:"quote"`
}

type QuoteUpdateRequest struct {
	Id    string `json:"id"`
	Quote string `json:"quote"`
}

// @Summary Create a new quote
// @Description This endpoint allows a user to submit a new quote.
// @Tags Quote
// @Accept json
// @Produce json
// @Param Authorization header string true "ID token"
// @Param quote body QuoteRequest true "The quote to be submitted."
// @Success 201 {object} map[string]string "{\"message\":\"Quote created successfully\", \"quote\":\"Example quote text\"}"
// @Failure 401 {object} map[string]string "{\"error\":\"Unauthorized\"}"
// @Failure 400 {object} map[string]string "{\"error\":\"Invalid request body\"}"
// @Failure 500 {object} map[string]string "{\"error\":\"Failed to create quote\"}"
// @Router /quote/create [post]
func CreateQuote(c *gin.Context) {
	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var quoteRequest QuoteRequest
	userData := user.(*middleware.User)
	if err := c.ShouldBindJSON(&quoteRequest); err != nil {
		log.Println("Error binding incoming json data", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	err := createQuote(&Quote{
		UserId:   userData.UserID,
		Quote:    quoteRequest.Quote,
		Approved: false,
	})
	if err != nil {
		log.Println("CreatePost: Error failed to create post", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create quote"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Quote created successfully", "quote": quoteRequest.Quote})
}

// @Summary Update an existing quote
// @Description This endpoint allows a user to update an existing quote. Admins can update any quote.
// @Tags Quote
// @Accept json
// @Produce json
// @Param Authorization header string true "ID token"
// @Param id path string true "The ID of the quote to be updated."
// @Param quote body QuoteUpdateRequest true "The updated quote text."
// @Success 200 {object} map[string]string "{\"message\":\"Quote updated successfully\", \"quote\":\"Updated quote text\"}"
// @Failure 401 {object} map[string]string "{\"error\":\"Unauthorized access\"}"
// @Failure 400 {object} map[string]string "{\"error\":\"Invalid request body\"}"
// @Failure 404 {object} map[string]string "{\"error\":\"Quote not found\"}"
// @Failure 403 {object} map[string]string "{\"error\":\"not authorized\"}"
// @Failure 500 {object} map[string]string "{\"error\":\"Failed to update quote\"}"
// @Router /quote/update [put]
func UpdateQuote(c *gin.Context) {
	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}
	uid := user.(*middleware.User).UserID
	role := user.(*middleware.User).Role

	var requestBody QuoteUpdateRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println("Error binding incoming json data", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	quote, err := getQuoteById(requestBody.Id)
	if err != nil {
		if errors.Is(err, ErrQuoteNotFound) {
			log.Printf("getQuote: Quote not found with id %s", requestBody.Id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Quote not found"})
		} else {
			log.Printf("getQuote: Error getting quote with id %s: %v", requestBody.Id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update quote"})
		}
		return
	}

	if role != "admin" {
		if uid != quote.UserId {
			log.Printf("Unauthorized access to delete quote for user %s with role %s :", user.(*middleware.User).Email, role)
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
	}

	err = updateQuote(&Quote{
		Id:       requestBody.Id,
		Quote:    requestBody.Quote,
		Approved: false,
	})
	if err != nil {
		log.Println("Update Quote, error updating quote", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update quote"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Quote created successfully", "quote": requestBody.Quote})
}

// @Summary Delete an existing quote
// @Description This endpoint allows a user to delete an existing quote. Admins can delete any quote.
// @Tags Quote
// @Accept json
// @Produce json
// @Param Authorization header string true "ID token"
// @Param id path string true "The ID of the quote to be deleted."
// @Success 200 {object} map[string]string "{\"message\":\"Quote deleted successfully\"}"
// @Failure 401 {object} map[string]string "{\"error\":\"Unauthorized access\"}"
// @Failure 404 {object} map[string]string "{\"error\":\"Quote not found\"}"
// @Failure 403 {object} map[string]string "{\"error\":\"not authorized\"}"
// @Failure 400 {object} map[string]string "{\"error\":\"Failed to delete quote\"}"
// @Router /quote/delete/{id} [delete]
func DeleteQuote(c *gin.Context) {
	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}
	uid := user.(*middleware.User).UserID
	role := user.(*middleware.User).Role

	id := c.Param("id")
	quote, err := getQuoteById(id)
	if err != nil {
		if errors.Is(err, ErrQuoteNotFound) {
			log.Printf("getQuote: Quote not found with id %s", id)
			c.JSON(http.StatusNotFound, gin.H{"error": "Quote not found"})
		} else {
			log.Printf("getQuote: Error getting quote with id %s: %v", id, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update quote"})
		}
		return
	}
	if role != "admin" {
		if uid != quote.UserId {
			log.Printf("Unauthorized access to delete quote for user %s with role %s :", user.(*middleware.User).Email, role)
			c.JSON(http.StatusForbidden, gin.H{"error": "not authorized"})
			return
		}
	}

	err = deleteQuote(id)
	if err != nil {
		log.Println("Error deleting quote:", err)
		c.JSON(http.StatusBadRequest, gin.H{"Failed to delete quote ": err})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Quote deleted successfully"})
}

// @Summary Get quotes
// @Description This endpoint retrieves quotes. Admins can retrieve all quotes, while other users can only retrieve approved quotes.
// @Tags Quote
// @Accept json
// @Produce json
// @Param Authorization header string true "ID token"
// @Success 200 {object} map[string][]Quote "A list of quotes."
// @Failure 401 {object} map[string]string "{\"error\":\"Unauthorized\"}"
// @Failure 404 {object} map[string]string "{\"error\":\"Quote not found\"}"
// @Failure 500 {object} map[string]string "{\"error\":\"Failed to get quote\"}"
// @Router /quote [get]
func GetQuotes(c *gin.Context) {
	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	role := user.(*middleware.User).Role

	var quotes []*Quote
	var err error

	if role == "admin" {
		quotes, err = getAllQuotes()
	} else {
		quotes, err = getAllApprovedQuotes()
	}

	if err != nil {
		if errors.Is(err, ErrQuoteNotFound) {
			log.Print("getQuote: Quotes not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Quote not found"})
		} else {
			log.Printf("getQuote: Error getting quotes %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get quote"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"quotes": quotes})
}

// @Summary Get quotes by user ID
// @Description Retrieve quotes by user ID
// @Tags quote
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT token"
// @Param profile-id path string true "User ID"
// @Success 200 {array} Quote "List of quotes"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "Quote not found"
// @Failure 500 {string} string "Failed to get quote"
// @Router /quote/{profile-id} [get]
func GetQuotesByUserId(c *gin.Context) {
	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var quotes []*Quote
	var err error

	userID := user.(*middleware.User).UserID
	role := user.(*middleware.User).Role

	if role == "admin" || userID == c.Param("profile-id") {
		quotes, err = getQuotesByUserId(userID)
	} else {
		quotes, err = getApprovedQuotesByUserId(userID)
	}

	if err != nil {
		if errors.Is(err, ErrQuoteNotFound) {
			log.Print("getQuote: Quotes not found")
			c.JSON(http.StatusNotFound, gin.H{"error": "Quote not found"})
		} else {
			log.Printf("getQuote: Error getting quotes %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get quote"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"quotes": quotes})
}

// @Summary Approve quote
// @Description Approve a quote by ID
// @Tags admin
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT token"
// @Param id path string true "Quote ID"
// @Success 200 {string} string "Quote approved successfully"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Unauthorized"
// @Failure 404 {string} string "Quote not found"
// @Failure 500 {string} string "Failed to approve quote"
// @Router /admin/quote/approve/{id} [post]
func ApproveQuote(c *gin.Context) {
	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	role := user.(*middleware.User).Role
	if role != "admin" {
		log.Printf("Unauthorized access for user %s with role %s :", user.(*middleware.User).Email, user.(*middleware.User).Role)
		c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized"})
		return
	}

	quoteID := c.Param("id")
	err := approveQuote(quoteID)
	if err != nil {
		if errors.Is(err, ErrQuoteNotFound) {
			log.Println("ApproveQuote: Error quote not found", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "Quote not found"})
		} else {
			log.Println("ApproveQuote: Error approving quote", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to approve quote"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Quote approved successfully"})
}

// @Summary Get unapproved quotes
// @Description Get all unapproved quotes
// @Tags admin
// @Accept json
// @Produce json
// @Param Authorization header string true "JWT token"
// @Success 200 {array} Quote "List of unapproved quotes"
// @Failure 401 {string} string "Unauthorized"
// @Failure 404 {string} string "No unapproved quotes found"
// @Failure 500 {string} string "Failed to get unapproved quotes"
// @Router /admin/quote/unapproved [get]
func GetUnapprovedQuotes(c *gin.Context) {
	unapprovedQuotes, err := getUnapprovedQuotes()
	if err != nil {
		if err == ErrQuoteNotFound {
			log.Println("GetUnapprovedQuotes, No unapproved quotes found", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "No unapproved quotes found"})
			return
		}
		log.Println("GetUnapprovedQuotes, Failed to get unapproved quotes", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get unapproved quotes"})
		return
	}
	c.JSON(http.StatusOK, unapprovedQuotes)
}

func getUserFromCtx(ctx *gin.Context) (any, bool) {
	user, exists := ctx.Get("user")
	if !exists {
		return "", false
	}
	return user, true
}
