package models

// TimetableEntry struct represents a lesson
// Multiple Timetable may share same TimetableEntry
type TimetableEntry struct {
	ID        DbId   `json:"id"`
	TeacherId UserId `json:"teacher"`
	Location  string `json:"location"`
	Day       string `json:"day"`
	Period    string `json:"period"`
	Subject   string `json:"subject"`
}

// Timetable is a holder of TimetableEntry objects and its visibility
// This struct is per-user, each user has their very own TimeTable
type Timetable struct {
	Entries  []DbId `json:"entries"`
	IsPublic bool   `json:"isPublic"`
}
