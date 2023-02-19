package models

// Checklist struct represents a to-do list
type Checklist struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Items  []Item `json:"items"`
	UserID string `json:"user_id"`
}

// Item struct represents an item in a to-do list
type Item struct {
	ID       int    `json:"id"`
	Text     string `json:"text"`
	Complete bool   `json:"complete"`
}
