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

func CreateQuote(c *gin.Context) {
	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	var quoteRequest QuoteRequest
	if err := c.ShouldBindJSON(&quoteRequest); err != nil {
		log.Println("Error binding incoming json data", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	err := createQuote(&Quote{
		UserId:   user.UserID,
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

func UpdateQuote(c *gin.Context) {
	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}
	uid := user.UserID
	role := user.Role

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
			log.Printf("Unauthorized access to delete quote for user %s with role %s :", user.Email, role)
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

func DeleteQuote(c *gin.Context) {
	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}
	uid := user.UserID
	role := user.Role

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
			log.Printf("Unauthorized access to delete quote for user %s with role %s :", user.Email, role)
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

func GetQuotes(c *gin.Context) {
	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	role := user.Role

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

func GetQuotesByUserId(c *gin.Context) {
	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var quotes []*Quote
	var err error

	userID := user.UserID
	role := user.Role

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

func ApproveQuote(c *gin.Context) {
	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	role := user.Role
	if role != "admin" {
		log.Printf("Unauthorized access for user %s with role %s :", user.Email, user.Role)
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

func getUserFromCtx(ctx *gin.Context) (*middleware.User, bool) {
	user, exists := ctx.Get("user")
	if !exists {
		return nil, false
	}
	return user.(*middleware.User), true
}
