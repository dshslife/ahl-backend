package utils

import "regexp"

// Check if email address is valid
func IsEmailValid(email string) bool {
	// Regular expression for email validation
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")

	return emailRegex.MatchString(email)
}
