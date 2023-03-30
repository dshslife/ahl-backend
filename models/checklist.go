package models

// Checklist struct represents a to-do list
type Checklist struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Items  []Items `json:"items"`
	UserID string  `json:"user_id"`
}

// Item struct represents an item in a to-do list
type Items struct {
	ID         int       `json:"id"`
	Text       string    `json:"text"`
	Complete   bool      `json:"complete"`
	IsPublic   bool      `json:"isPublic"`
	SharedWith []Friends `json:"shared_with"`
}

// Friends struct represents an person who can share to-do list
type Friends struct {
	ID        int    `json:"id"`
	studentID string `json:"studentID"`
	School    string `json:"school"`
	Grade     int    `json:"grade"`
	Class     int    `json:"class"`
	Num       int    `json:"num"`
	Name      string `json:"name"`
}
