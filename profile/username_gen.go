package profile

import (
	"fmt"
	"strings"
)

func generateUsername(email string) (string, error) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid email format")
	}

	usernamePart := parts[0]
	names := strings.Split(usernamePart, ".")
	if len(names) == 0 {
		return "", fmt.Errorf("no names found in email")
	}

	username := "gopher" + strings.Join(names, "")
	return username, nil
}
