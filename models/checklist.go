package models

// Checklist struct represents a to-do list
type Checklist struct {
	ID        DbId            `json:"id"`
	StudentId UserId          `json:"student_id"`
	Title     string          `json:"title"`
	Items     []ChecklistItem `json:"items"`
}

// ChecklistItem struct represents an item in a to-do list
type ChecklistItem struct {
	Text     string `json:"text"`
	Complete bool   `json:"complete"`
	IsPublic bool   `json:"is_public"`
	// 아래 필드는 무조건 StudentInfo#Friends에 등록된 친구만 포함할 것!
	SharedWith []UserId `json:"shared_with"`
}
