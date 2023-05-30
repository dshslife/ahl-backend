package utils

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/username/schoolapp/models"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"time"
)

// ValidateNewAccount ValidateNewAccount와 동일하나 ID 검사는 안함
func ValidateNewAccount(account *models.Account) error {
	if account == nil {
		return errors.New("account is nil")
	}

	if account.Name == "" {
		return errors.New("account name is empty")
	}
	if account.Email == "" {
		return errors.New("account email is empty")
	}
	if !isEmailValid(account.Email) {
		return errors.New("account email is invalid")
	}
	if len(account.Password) == 0 {
		return errors.New("account password is empty")
	}
	// PermissionInfo 필드 값 확인
	switch account.GetLevel() {
	case models.STUDENT:
		studentInfo, ok := account.PermissionInfo.(models.StudentInfo)
		if !ok {
			return errors.New("invalid permission info for student")
		}
		if studentInfo.SchoolId == "" {
			return errors.New("school id is empty for student")
		}
		if studentInfo.Grade == 0 {
			return errors.New("grade is empty for student")
		}
		if studentInfo.Class == 0 {
			return errors.New("class is empty for student")
		}
		if studentInfo.Number == 0 {
			return errors.New("number is empty for student")
		}
	case models.TEACHER:
		teacherInfo, ok := account.PermissionInfo.(models.TeacherInfo)
		if !ok {
			return errors.New("invalid permission info for teacher")
		}
		if teacherInfo.SchoolId == "" {
			return errors.New("school id is empty for teacher")
		}
	case models.ADMIN:
		_, ok := account.PermissionInfo.(models.AdminInfo)
		if !ok {
			return errors.New("invalid permission info for admin")
		}
	default:
		return errors.New("invalid permission level")
	}
	return nil
}

// Validate account depending on their type
func ValidateAccount(account *models.Account) error {
	if account == nil {
		return errors.New("account is nil")
	}

	// Account 필드 값 확인
	if account.UserId == uuid.Nil {
		return errors.New("account user id is empty")
	}
	if account.Name == "" {
		return errors.New("account name is empty")
	}
	if account.Email == "" {
		return errors.New("account email is empty")
	}
	if !isEmailValid(account.Email) {
		return errors.New("account email is invalid")
	}
	if len(account.Password) == 0 {
		return errors.New("account password is empty")
	}
	// PermissionInfo 필드 값 확인
	switch account.GetLevel() {
	case models.STUDENT:
		studentInfo, ok := account.PermissionInfo.(models.StudentInfo)
		if !ok {
			return errors.New("invalid permission info for student")
		}
		if studentInfo.SchoolId == "" {
			return errors.New("school id is empty for student")
		}
		if studentInfo.Grade == 0 {
			return errors.New("grade is empty for student")
		}
		if studentInfo.Class == 0 {
			return errors.New("class is empty for student")
		}
		if studentInfo.Number == 0 {
			return errors.New("number is empty for student")
		}
	case models.TEACHER:
		teacherInfo, ok := account.PermissionInfo.(models.TeacherInfo)
		if !ok {
			return errors.New("invalid permission info for teacher")
		}
		if teacherInfo.SchoolId == "" {
			return errors.New("school id is empty for teacher")
		}
	case models.ADMIN:
		_, ok := account.PermissionInfo.(models.AdminInfo)
		if !ok {
			return errors.New("invalid permission info for admin")
		}
	default:
		return errors.New("invalid permission level")
	}
	return nil
}

// ValidateTimetable Validate a timetable object
func ValidateTimetable(tt *models.Timetable) error {
	if tt == nil {
		return errors.New("timetable is nil")
	}

	if len(tt.Entries) == 0 {
		return errors.New("timetable has no entries")
	}

	return nil
}

func ValidateTimeTableEntry(entry *models.TimetableEntry) error {
	if entry.TeacherId == uuid.Nil {
		return errors.New("teacher id field is empty")
	}
	if entry.Location == "" {
		return errors.New("location field is empty")
	}
	if entry.Day == "" {
		return errors.New("day field is empty")
	}
	if entry.Period == "" {
		return errors.New("period field is empty")
	}
	if entry.Subject == "" {
		return errors.New("subject field is empty")
	}
	return nil
}

// Validate a checklist object
func ValidateChecklist(checklist *models.Checklist) error {
	if checklist == nil {
		return errors.New("checklist is nil")
	}
	if checklist.StudentId == uuid.Nil {
		return errors.New("student id is required")
	}
	if checklist.Title == "" {
		return errors.New("title is required")
	}

	for _, item := range checklist.Items {
		if item.Text == "" {
			return errors.New("item text is required")
		}
	}

	return nil
}

// Validate a cafeteria menu item
func ValidateCafeteriaMenu(menu *models.CafeteriaMenu) error {
	if menu.SchoolId == "" {
		return errors.New("missing school ID")
	}

	if menu.MealName == "" {
		return errors.New("missing meal name")
	}

	if menu.Date == "" {
		return errors.New("missing date")
	}

	if len(menu.Items) == 0 {
		return errors.New("menu has no items")
	}

	for _, item := range menu.Items {
		if item.Name == "" {
			return errors.New("missing item name in menu")
		}

		if len(item.Allergies) > 0 {
			for _, allergy := range item.Allergies {
				if !(1 <= allergy && allergy <= 13) {
					return fmt.Errorf("invalid allergy type: %s", allergy)
				}
			}
		}

		if item.Contents == "" {
			return errors.New("missing contents in menu item")
		}
	}

	return nil
}

func ValidateEvents(events *models.Events) error {
	if events.SchoolId == "" {
		return fmt.Errorf("school ID is required")
	}
	if events.Month < 1 || events.Month > 12 {
		return fmt.Errorf("month must be between 1 and 12")
	}
	for _, event := range events.Events {
		// Date validation
		if _, err := time.Parse("2006-01-02", event.Date); err != nil {
			return fmt.Errorf("invalid date format: %s", event.Date)
		}
		// AttendanceType validation
		if !isValidAttendanceType(event.FirstGradeAttends) {
			return fmt.Errorf("invalid first grade attendance type: %d", event.FirstGradeAttends)
		}
		if !isValidAttendanceType(event.SecondGradeAttends) {
			return fmt.Errorf("invalid second grade attendance type: %d", event.SecondGradeAttends)
		}
		if !isValidAttendanceType(event.ThirdGradeAttends) {
			return fmt.Errorf("invalid third grade attendance type: %d", event.ThirdGradeAttends)
		}
	}
	return nil
}

func ValidateSchool(events *models.School) error {
	if events.SchoolId == "" {
		return fmt.Errorf("school ID is required")
	}
	if events.RegionId == "" {
		return fmt.Errorf("region ID is required")
	}
	if events.SchoolName == "" {
		return fmt.Errorf("school name ID is required")
	}
	if events.RegionName == "" {
		return fmt.Errorf("region name is required")
	}
	return nil
}

func isValidAttendanceType(attendanceType models.AttendanceType) bool {
	return attendanceType == models.YES || attendanceType == models.IGNORED || attendanceType == models.NO
}

// Check if email address is valid
func isEmailValid(email string) bool {
	// Regular expression for email validation
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")

	return emailRegex.MatchString(email)
}

func HashPassword(password []byte) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func VerifyPassword(hashedPassword, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	return err == nil
}
