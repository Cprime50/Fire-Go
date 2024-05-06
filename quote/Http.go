package quote

import (
	"log"
	"net/http"

	"github.com/cprime50/fire-go/middleware"
	"github.com/gin-gonic/gin"
)

func CreateQuoteHandler(c *gin.Context, s QuoteService) {
	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var quoteRequest QuoteRequest
	if err := c.ShouldBindJSON(&quoteRequest); err != nil {
		log.Println("Error binding incoming json data", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidRequestBody.Error()})
		return
	}

	err := s.CreateQuote(user.UserID, quoteRequest.Quote)
	if err != nil {
		log.Println("CreateQuoteHandler: Error failed to create quote", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": ErrCreateQuote.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Quote created successfully", "quote": quoteRequest.Quote})
}

func UpdateQuoteHandler(c *gin.Context, service QuoteService) {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidRequestBody.Error()})
		return
	}

	err := service.UpdateQuote(uid, role, requestBody.Id, requestBody.Quote)
	if err != nil {
		switch err {
		case ErrQuoteNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case ErrNotAuthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Quote updated successfully", "quote": requestBody.Quote})
}

func DeleteQuoteHandler(c *gin.Context, service QuoteService) {
	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}
	uid := user.UserID
	role := user.Role

	id := c.Param("id")

	err := service.DeleteQuote(uid, role, id)
	if err != nil {
		switch err {
		case ErrQuoteNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case ErrNotAuthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Quote deleted successfully"})
}

func GetQuotesHandler(c *gin.Context, service QuoteService) {
	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	role := user.Role

	quotes, err := service.GetQuotes(role)
	if err != nil {
		switch err {
		case ErrQuoteNotFound:
			log.Print("getQuote: Quotes not found")
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case ErrNotAuthorized:
			log.Printf("Unauthorized access for user %s with role %s", user.Email, role)
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			log.Printf("getQuote: Error getting quotes %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"quotes": quotes})
}

func GetQuotesByUserIdHandler(c *gin.Context, service QuoteService) {
	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}
	uid := user.UserID
	role := user.Role
	requestedUserId := c.Param("profile-id")

	quotes, err := service.GetQuotesByUserId(uid, role, requestedUserId)
	if err != nil {
		switch err {
		case ErrQuoteNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case ErrNotAuthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"quotes": quotes})
}

func ApproveQuoteHandler(c *gin.Context, service QuoteService) {
	user, ok := getUserFromCtx(c)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized access"})
		return
	}
	uid := user.UserID
	role := user.Role
	quoteId := c.Param("id")

	err := service.ApproveQuote(uid, role, quoteId)
	if err != nil {
		switch err {
		case ErrQuoteNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case ErrNotAuthorized:
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Quote approved successfully"})
}

func GetUnapprovedQuotesHandler(c *gin.Context, service QuoteService) {
	unapprovedQuotes, err := service.GetUnapprovedQuotes()
	if err != nil {
		switch err {
		case ErrQuoteNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "No unapproved quotes found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get unapproved quotes"})
		}
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
