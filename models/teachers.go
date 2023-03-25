package models

// Teacher struct represents a teacher
type Teacher struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Access bool   `json:"access"`
}
