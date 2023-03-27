package models

// Timetable struct represents a student's timetable
type Timetable struct {
	ID        int    `json:"id"`
	StudentID string `json:"student_id"`
	Teacher   string `json:"teacher"`
	Location  string `json:"location"`
	Day       string `json:"day"`
	Period    string `json:"period"`
	Subject   string `json:"subject"`
	IsPublic  bool   `json:"isPublic"`
}

// Timetables is a slice of Timetable objects
type Timetables []Timetable
