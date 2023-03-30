package models

type Events struct {
	Month  int               `json:"month"`
	School string            `json:"school"`
	Event  []EventCollection `json:"event"`
	ID     int               `json:"id"`
	Exists bool              `json:"exists"`
}

type EventCollection struct {
	ClassType     string `json:"class_type"`
	Date          string `json:"date"`
	EventName     string `json:"event_name"`
	EventContents string `json:"event_contents"`
	FirstGrade    string `json:"first_grade"`
	SecondGrade   string `json:"second_grade"`
	ThirdGrade    string `json:"third_grade"`
	ModifiedDate  string `json:"modified_date"`
}
