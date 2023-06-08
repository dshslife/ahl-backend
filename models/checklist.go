package models

import (
	"encoding/json"
	"github.com/google/uuid"
)

// Checklist struct represents a to-do list
type Checklist struct {
	ID        DbId            `json:"id"`
	StudentId uuid.UUID       `json:"student_id"`
	Title     string          `json:"title"`
	Items     []ChecklistItem `json:"items"`
}

func (checklist Checklist) Flatten() (FlatCheckList, error) {
	data, err := json.Marshal(checklist.Items)
	if err != nil {
		return FlatCheckList{}, err
	}
	return FlatCheckList{
		ID:        checklist.ID,
		StudentId: checklist.StudentId,
		Title:     checklist.Title,
		Items:     string(data),
	}, nil
}

func (flatten FlatCheckList) Restore() (Checklist, error) {
	var result []ChecklistItem
	err := json.Unmarshal([]byte(flatten.Items), &result)
	if err != nil {
		return Checklist{}, err
	}

	return Checklist{
		ID:        flatten.ID,
		StudentId: flatten.StudentId,
		Title:     flatten.Title,
		Items:     result,
	}, nil
}

type FlatCheckList struct {
	ID        DbId
	StudentId uuid.UUID
	Title     string
	Items     string
}

// ChecklistItem struct represents an item in a to-do list
type ChecklistItem struct {
	Text     string `json:"text"`
	Complete bool   `json:"complete"`
	IsPublic bool   `json:"is_public"`
	// 아래 필드는 무조건 StudentInfo#Friends에 등록된 친구만 포함할 것!
	SharedWith []uuid.UUID `json:"shared_with"`
}
