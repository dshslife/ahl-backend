package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"

	"github.com/username/schoolapp/models"
)

// Validate a student object
func ValidateStudent(student *models.Student) error {
	if student.Name == "" {
		return errors.New("Name is required")
	}

	if student.Email == "" {
		return errors.New("Email is required")
	}

	if !isEmailValid(student.Email) {
		return errors.New("Invalid email address")
	}

	if student.Grade < 1 || student.Grade > 3 {
		return errors.New("Grade should be between 1 and 3")
	}

	return nil
}

// Validate a timetable object
func ValidateTimetable(timetable *models.Timetable) error {
	if timetable.StudentID == "" {
		return errors.New("Student ID is required")
	}

	if timetable.Day == "" {
		return errors.New("Day is required")
	}

	if timetable.Period == "" {
		return errors.New("Period is required")
	}

	if timetable.Subject == "" {
		return errors.New("Subject is required")
	}

	return nil
}

// Validate a checklist object
func ValidateChecklist(checklist *models.Checklist) error {
	if checklist.Title == "" {
		return errors.New("Title is required")
	}

	return nil
}

// Validate a cafeteria menu item
func ValidateCafeteriaMenu(menu *models.CafeteriaMenu) error {
	if menu.Date == "" {
		return errors.New("Date is required")
	}

	if menu.Meal == "" {
		return errors.New("Meal is required")
	}

	return nil
}

// Check if email address is valid
func isEmailValid(email string) bool {
	// Regular expression for email validation
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")

	return emailRegex.MatchString(email)
}

func HashPassword(password []byte) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func VerifyPassword(hashedPassword, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	return err == nil
}

func GenerateJWT(userId int, secretKey string) (string, error) {
	// Define the expiration time for the token
	expirationTime := time.Now().Add(24 * time.Hour).Unix()

	// Create a new claims instance
	claims := jwt.MapClaims{}
	claims["user_id"] = userId
	claims["exp"] = expirationTime

	// Create a new token instance using the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
