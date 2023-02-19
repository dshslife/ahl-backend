package models

// Student struct represents a student
type Student struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Grade int    `json:"grade"`
}
