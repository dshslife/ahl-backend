package models

// Admin struct represents an admin
type Admin struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	AccessToAll bool   `json:"access_to_all"`
	AdminAccess bool   `json:"admin_access"`
}
