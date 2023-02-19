package models

// CafeteriaMenu struct represents a cafeteria menu
type CafeteriaMenu struct {
	ID    int    `json:"id"`
	Date  string `json:"date"`
	Meal  string `json:"meal"`
	Items []Item `json:"items"`
}

// Item struct represents an item in a cafeteria menu
type Item struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Allergy  string `json:"allergy,omitempty"`
	Vegetari bool   `json:"vegetari"`
}
