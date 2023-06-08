package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/username/schoolapp/models"
	"github.com/username/schoolapp/utils"
	"log"
	"os"
	"time"
)

var db *sql.DB

func createTables() {
	// prepare query
	createSchools := "CREATE TABLE IF NOT EXISTS `schools` (id INT(11) NOT NULL AUTO_INCREMENT PRIMARY KEY, `school_id` VARCHAR(255) NOT NULL, `region_id` VARCHAR(255) NOT NULL, `school_name` VARCHAR(255) NOT NULL, `region_name` VARCHAR(255) NOT NULL, `school_email_only` BOOL NOT NULL, `school_email` VARCHAR(255) NOT NULL)  ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"
	createAccounts := "CREATE TABLE IF NOT EXISTS `accounts` (`id` INT(11) NOT NULL AUTO_INCREMENT PRIMARY KEY, `user_id` TINYBLOB NOT NULL, `name` VARCHAR(255) NOT NULL, `email` VARCHAR(255) NOT NULL, `password` TINYBLOB NOT NULL, `permission_level` TINYINT NOT NULL, `school_id` VARCHAR(255), `timetable_list` LONGBLOB, `timetable_is_public` BOOL,`grade` TINYINT, `class` TINYINT, `number` TINYINT, `checklist_id` INT(11), `friends` LONGBLOB) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"
	createTimeTables := "CREATE TABLE IF NOT EXISTS `timetables` (`id` INT(11) NOT NULL AUTO_INCREMENT, `teacher_id` TINYBLOB NOT NULL, `location` VARCHAR(255) NOT NULL, `day` INT(11) NOT NULL, `period` TIME NOT NULL, `subject` VARCHAR(255) NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"
	createCafeteria := "CREATE TABLE IF NOT EXISTS `cafeteria_menus` ( `id` INT(11) NOT NULL AUTO_INCREMENT, `school_id` VARCHAR(255) NOT NULL, `meal_name` VARCHAR(255) NOT NULL, `date` DATE NOT NULL, `contents` TEXT NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"
	createChecklists := "CREATE TABLE IF NOT EXISTS `checklists` (`id` INT(11) NOT NULL AUTO_INCREMENT, `student_id` TINYBLOB NOT NULL, `title` TEXT NOT NULL, `items` TEXT NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"
	createEvents := "CREATE TABLE IF NOT EXISTS `schoolevents` (`id` INT(11) NOT NULL AUTO_INCREMENT, `school_id` VARCHAR(255) NOT NULL, `month` INT(11) NOT NULL, `events` TEXT NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"

	// Execute query
	queries := []string{createSchools,
		createAccounts,
		createTimeTables,
		createCafeteria,
		createChecklists,
		createEvents}
	for i := range queries {
		_, err := db.Exec(queries[i])
		if err != nil {
			log.Fatalf("Error creating database with query %s: %s", queries[i], err.Error())
		}
	}
}

// Connect connects to the database
func Connect() {
	var err error

	// Get database configuration from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Create data source name (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Connect to database
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err.Error())
	}

	// Set database connection parameters
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	// Check if database is alive
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error pinging database: %s", err.Error())
	}

	// Create tables if they don't exist
	createTables()
}

// Close closes the database connection
func Close() {
	err := db.Close()
	if err != nil {
		log.Fatalf("Error closing database connection: %s", err.Error())
	}
}

// GetAccountByEmail returns a user by Email
func GetAccountByEmail(Email *string) (*models.Account, error) {
	// Prepare query
	query := "SELECT * FROM accounts WHERE email = ?"

	// Execute query
	row := db.QueryRow(query, Email)

	// Scan row into flataccount object
	var flataccount models.FlatAccount
	err := row.Scan(&flataccount.DbId, &flataccount.UserId, &flataccount.Name, &flataccount.Email, &flataccount.Password, &flataccount.PermissionLevel, &flataccount.SchoolId, &flataccount.TimeTableEntries, &flataccount.TimeTableIsPublic, &flataccount.Grade, &flataccount.Class, &flataccount.Number, &flataccount.ChecklistId, &flataccount.Friends)
	if err != nil {
		return nil, err
	}
	user, err := flataccount.Restore()
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetAccountById returns a student by ID
func GetAccountById(id *uuid.UUID) (*models.Account, error) {
	// Prepare query
	query := "SELECT * FROM accounts WHERE user_id = ?"

	// Execute query
	row := db.QueryRow(query, id)

	// Scan row into user object
	var flataccount models.FlatAccount
	err := row.Scan(&flataccount.DbId, &flataccount.UserId, &flataccount.Name, &flataccount.Email, &flataccount.Password, &flataccount.PermissionLevel, &flataccount.SchoolId, &flataccount.TimeTableEntries, &flataccount.TimeTableIsPublic, &flataccount.Grade, &flataccount.Class, &flataccount.Number, &flataccount.ChecklistId, &flataccount.Friends)
	if err != nil {
		return nil, err
	}
	user, err := flataccount.Restore()
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateAccount creates a new student
func CreateAccount(account *models.Account) (models.DbId, error) {
	// Prepare query
	query := "INSERT INTO accounts (user_id, name, email, password, permission_level, school_id, timetable_list, timetable_is_public, grade, class, number, checklist_id, friends) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	var result sql.Result
	var err error

	flataccount, err := account.ToSql()
	if err != nil {
		return models.DbId(0), err
	}

	result, err = db.Exec(query, flataccount.UserId, flataccount.Name, flataccount.Email, flataccount.Password, flataccount.PermissionLevel, flataccount.SchoolId, flataccount.TimeTableEntries, flataccount.TimeTableIsPublic, flataccount.Grade, flataccount.Class, flataccount.Number, flataccount.ChecklistId, flataccount.Friends)

	if err != nil {
		return 0, err
	}

	// Get the ID of the newly created account
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	account.DbId = models.DbId(id)

	return account.DbId, nil
}

// UpdateAccount updates a student
func UpdateAccount(account *models.Account) error {
	// Prepare query
	query := "UPDATE accounts SET user_id = ?, name = ?, email = ?, password = ?, permission_level = ?, school_id = ?, timetable_list = ?, timetable_is_public = ?, grade = ?, class = ?, number = ?, checklist_id = ?, friends = ? WHERE id = ?"

	flataccount, err := account.ToSql()
	if err != nil {
		return err
	}
	// Execute query
	_, err = db.Exec(query, flataccount.UserId, flataccount.Name, flataccount.Email, flataccount.Password, flataccount.PermissionLevel, flataccount.SchoolId, flataccount.TimeTableEntries, flataccount.TimeTableIsPublic, flataccount.Grade, flataccount.Class, flataccount.Number, flataccount.ChecklistId, flataccount.ChecklistId, flataccount.Friends, flataccount.DbId)
	if err != nil {
		return err
	}

	return nil
}

// DeleteAccount deletes a student by ID
func DeleteAccount(id models.DbId) error {
	// Prepare query
	query := "DELETE FROM accounts WHERE id = ?"

	// Execute query
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

// GetTimeTableEntry returns a list of timetables for a student
func GetTimeTableEntry(id models.DbId) (*models.TimetableEntry, error) {
	// Prepare query
	query := "SELECT * FROM timetables WHERE id = ?"

	// Execute query
	rows, err := db.Query(query, id)
	if err != nil {
		return nil, err
	}

	// Scan rows into entry objects
	var entry models.TimetableEntry
	err = rows.Scan(&entry.ID, &entry.TeacherId, &entry.Location, &entry.Day, &entry.Period, &entry.Subject)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

// CreateTimetable creates a new timetable
func CreateTimetable(entry *models.TimetableEntry) (models.DbId, error) {
	err := utils.ValidateTimeTableEntry(entry)
	if err != nil {
		return 0, err
	}

	// Prepare query
	query := "INSERT INTO timetables (teacher_id, location, day, period, subject) VALUES (?, ?, ?, ?, ?)"

	// Execute query
	result, err := db.Exec(query, entry.TeacherId, entry.Location, entry.Day, entry.Period, entry.Subject)
	if err != nil {
		return 0, err
	}

	// Get the ID of the newly created entry
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	entry.ID = models.DbId(id)

	return entry.ID, nil
}

// UpdateTimetable updates a timetable
func UpdateTimetable(entry *models.TimetableEntry) error {
	err := utils.ValidateTimeTableEntry(entry)
	if err != nil {
		return err
	}
	// Prepare query
	query := "UPDATE timetables SET id = ?, teacher_id = ?, location = ?, day = ?, period = ?, subject = ? WHERE id = ?"

	// Execute query
	_, err = db.Exec(query, entry.ID, entry.TeacherId, entry.Location, entry.Day, entry.Period, entry.Subject, entry.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTimetable deletes a timetable by ID
func DeleteTimetable(id models.DbId) error {
	// Prepare query
	query := "DELETE FROM timetables WHERE id = ?"

	// Execute query
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

// GetMenuByID returns the cafeteria menu for a specific date by ID
func GetMenuByID(id models.DbId) (*models.CafeteriaMenu, error) {
	// Prepare query
	query := "SELECT * FROM cafeteria_menus WHERE id = ?"

	// Execute query
	row := db.QueryRow(query, id)

	// Scan row into cafeteria menu object
	var menu models.CafeteriaMenu
	err := row.Scan(&menu.ID, &menu.SchoolId, &menu.MealName, &menu.Date, &menu.Contents)
	if err != nil {
		return nil, err
	}

	return &menu, nil
}

// GetMenu returns the cafeteria menu for a specific date
func GetMenu(date time.Time) (*models.CafeteriaMenu, error) {
	// Prepare query
	query := "SELECT * FROM cafeteria_menus WHERE date = ?"

	// Execute query
	row := db.QueryRow(query, date)

	// Scan row into cafeteria menu object
	var menu models.CafeteriaMenu
	err := row.Scan(&menu.ID, &menu.SchoolId, &menu.MealName, &menu.Date, &menu.Contents)
	if err != nil {
		return nil, err
	}

	return &menu, nil
}

// CreateMenu creates a new cafeteria menu
func CreateMenu(menu *models.CafeteriaMenu) (models.DbId, error) {
	err := utils.ValidateCafeteriaMenu(menu)
	if err != nil {
		return 0, err
	}

	// Prepare query to insert menu
	menuQuery := "INSERT INTO cafeteria_menus (school_id, meal_name, date, contents) VALUES (?, ?, ?, ?)"

	// Execute query to insert menu
	menuResult, err := db.Exec(menuQuery, menu.SchoolId, menu.MealName, menu.Date, menu.Contents)
	if err != nil {
		return 0, err
	}

	// Get the ID of the newly created menu
	menuID, err := menuResult.LastInsertId()
	if err != nil {
		return 0, err
	}
	menu.ID = models.DbId(menuID)

	return models.DbId(menuID), nil
}

// UpdateMenu updates a cafeteria menu
func UpdateMenu(menu *models.CafeteriaMenu) error {
	err := utils.ValidateCafeteriaMenu(menu)
	if err != nil {
		return err
	}

	// Prepare query to insert menu
	menuQuery := "UPDATE cafeteria_menus SET school_id = ?, meal_name = ?, date = ?, contents = ? WHERE id = ?"

	// Execute query to insert menu
	_, err = db.Exec(menuQuery, menu.SchoolId, menu.MealName, menu.Date, menu.Contents, menu.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteMenu deletes a cafeteria menu by ID
func DeleteMenu(id models.DbId) error {
	// Prepare query to delete menu
	deleteMenuQuery := "DELETE FROM cafeteria_menus WHERE id = ?"

	// Execute query to delete menu
	_, err := db.Exec(deleteMenuQuery, id)
	if err != nil {
		return err
	}

	return nil
}

// GetChecklistsOfStudent returns a checklist
func GetChecklistsOfStudent(studentID *uuid.UUID) (*models.Checklist, error) {
	// Prepare query
	query := "SELECT * FROM checklists WHERE student_id = ?"

	// Execute query
	row := db.QueryRow(query, studentID)

	// Scan row into checklist object
	var checklist models.Checklist
	err := row.Scan(&checklist.ID, &checklist.Title, &checklist.Items, &checklist.Items)
	if err != nil {
		return nil, err
	}

	return &checklist, nil
}

func GetChecklistsById(id models.DbId) (*models.Checklist, error) {
	// Prepare query
	query := "SELECT * FROM checklists WHERE id = ?"

	// Execute query
	row := db.QueryRow(query, id)

	// Scan row into checklist object
	var checklist models.Checklist
	err := row.Scan(&checklist.ID, &checklist.Title, &checklist.Items, &checklist.Items)
	if err != nil {
		return nil, err
	}

	return &checklist, nil
}

// CreateChecklist creates a new checklist
func CreateChecklist(checklist *models.Checklist) (models.DbId, error) {
	err := utils.ValidateChecklist(checklist)
	if err != nil {
		return 0, err
	}

	// Prepare query to insert menu
	createQuery := "INSERT INTO checklists (student_id, title, items) VALUES (?, ?, ?)"

	// Execute query to insert menu
	result, err := db.Exec(createQuery, checklist.StudentId, checklist.Title, checklist.Items)
	if err != nil {
		return 0, err
	}

	// Get the ID of the newly created checklist
	listId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	checklist.ID = models.DbId(listId)

	return models.DbId(listId), nil
}

// UpdateChecklist updates a checklist
func UpdateChecklist(checklist *models.Checklist) error {
	err := utils.ValidateChecklist(checklist)
	if err != nil {
		return err
	}

	// Prepare query to insert checklist
	update := "UPDATE checklists SET student_id = ?, title = ?, items = ? WHERE id = ?"

	// Execute query to insert checklist
	_, err = db.Exec(update, checklist.StudentId, checklist.Title, checklist.Items, checklist.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteChecklist deletes a checklist by ID
func DeleteChecklist(id models.DbId) error {
	// Prepare query to deleteQuery checklist items
	deleteQuery := "DELETE FROM checklists WHERE id = ?"

	// Execute query to deleteQuery checklist items
	_, err := db.Exec(deleteQuery, id)
	if err != nil {
		return err
	}
	return nil
}

// GetAllEvents returns all events
func GetAllEvents() (*models.Events, error) {
	// Prepare query
	query := "SELECT * FROM schoolevents"

	// Execute query
	row := db.QueryRow(query)

	// Scan row into student object
	var events models.Events
	err := row.Scan(&events.ID, &events.SchoolId, &events.Month, &events.Events)
	if err != nil {
		return nil, err
	}

	return &events, nil
}

// GetEventsByMonth returns an event by Month
func GetEventsByMonth(month int) (*models.Events, error) {
	// Prepare query
	query := "SELECT * FROM schoolevents WHERE month = ?"

	// Execute query
	row := db.QueryRow(query, month)

	// Scan row into student object
	var events models.Events
	err := row.Scan(&events.ID, &events.SchoolId, &events.Month, &events.Events)
	if err != nil {
		return nil, err
	}

	return &events, nil
}

// CreateEvents creates a new event
func CreateEvents(events *models.Events) (models.DbId, error) {
	err := utils.ValidateEvents(events)
	if err != nil {
		return 0, err
	}
	// Prepare query
	query := "INSERT INTO schoolevents (school_id, month, events) VALUES (?, ?, ?)"

	// Execute query
	result, err := db.Exec(query, events.SchoolId, events.Month, &events.Events)
	if err != nil {
		return 0, err
	}

	// Get the ID of the newly created student
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	events.ID = models.DbId(id)

	return events.ID, nil
}

func GetAllAccounts() ([]models.Account, error) {
	// Prepare query
	query := "SELECT * FROM accounts"

	// Execute query
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	// Iterate through rows and create Account objects
	var accounts []models.Account
	for rows.Next() {
		var account models.Account
		err := rows.Scan(&account.DbId, &account.UserId, &account.Name, &account.Email, &account.Password, &account.PermissionInfo)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

func CreateSchool(school *models.School) (models.DbId, error) {
	err := utils.ValidateSchool(school)
	if err != nil {
		return 0, err
	}
	// Prepare query
	query := "INSERT INTO schools (school_id, region_id, school_name, region_name, school_email_only, school_email) VALUES (?, ?, ?, ?, ?, ?)"

	// Execute query
	result, err := db.Exec(query, school.SchoolId, school.RegionId, school.SchoolName, school.RegionName, school.SchoolEmailOnly, school.SchoolEmail)
	if err != nil {
		return 0, err
	}

	// Get the ID of the newly created school
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	school.ID = models.DbId(id)

	return school.ID, nil
}

func GetSchool(id models.SchoolId) (*models.School, error) {
	// Prepare query
	query := "SELECT * FROM schools WHERE school_id = ?"

	// Execute query
	row := db.QueryRow(query, id)

	// Scan row into student object
	var school models.School
	err := row.Scan(&school.ID, &school.SchoolId, &school.RegionId, &school.SchoolName, &school.RegionName, &school.SchoolEmailOnly)
	if err != nil {
		return nil, err
	}

	return &school, nil
}

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
	if !utils.IsEmailValid(account.Email) {
		return errors.New("account email is invalid")
	}
	if len(account.Password) == 0 {
		return errors.New("account password is empty")
	}
	if account.PermissionInfo == nil {
		return errors.New("permission info is empty")
	}
	// PermissionInfo 필드 값 확인
	switch account.PermissionInfo.GetLevel() {
	case models.STUDENT:
		studentInfo, ok := account.PermissionInfo.(models.StudentInfo)
		if !ok {
			return errors.New("invalid permission info for student")
		}
		if studentInfo.SchoolId == "" {
			return errors.New("school id is empty for student")
		}
		school, _ := GetSchool(studentInfo.SchoolId)
		if school == nil {
			return errors.New("unregistered school id")
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
