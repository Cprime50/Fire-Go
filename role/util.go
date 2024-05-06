package role

import (
	"fmt"
	"regexp"

	"github.com/gin-gonic/gin"
)

func validateInput(ctx *gin.Context, input EmailInput) (string, error) {
	if err := ctx.BindJSON(&input); err != nil {
		return "", fmt.Errorf("invalid JSON format")
	}
	emailOk := ValidateEmail(input.Email)
	if !emailOk {
		return "", fmt.Errorf("invalid email format")
	}
	return input.Email, nil
}

// Regex for email validation
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
