package db

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/username/schoolapp/models"
	"github.com/username/schoolapp/utils"
	"log"
	"os"
	"time"
)

var db *sql.DB

func createTables() {
	// prepare query
	createSchools := "CREATE TABLE IF NOT EXISTS `schools` (id INT(11) NOT NULL AUTO_INCREMENT PRIMARY KEY, `school_id` VARCHAR(255) NOT NULL, `region_id` VARCHAR(255) NOT NULL, `school_name` VARCHAR(255) NOT NULL, `region_name` VARCHAR(255) NOT NULL)  ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"
	createAccounts := "CREATE TABLE IF NOT EXISTS `accounts` (`id` INT(11) NOT NULL AUTO_INCREMENT PRIMARY KEY, `user_id` VARCHAR(255) NOT NULL, `name` VARCHAR(255) NOT NULL, `email` VARCHAR(255) NOT NULL, `password` VARCHAR(255) NOT NULL, `permission_level` TINYINT NOT NULL, `school_id` VARCHAR(255), `timetable` TEXT NOT NULL, `grade` INT, `class` INT, `number` INT, `checklist_id` VARCHAR(255) NOT NULL, `friends` TEXT) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"
	createTimeTables := "CREATE TABLE IF NOT EXISTS `timetables` (`id` INT(11) NOT NULL AUTO_INCREMENT, `teacher_id` INT(11) NOT NULL, `location` VARCHAR(255) NOT NULL, `day` INT(11) NOT NULL, `period` TIME NOT NULL, `subject` VARCHAR(255) NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"
	createCafeteria := "CREATE TABLE IF NOT EXISTS `cafeteria_menus` ( `id` INT(11) NOT NULL AUTO_INCREMENT, `school_id` VARCHAR(255) NOT NULL, `meal_name` VARCHAR(255) NOT NULL, `date` DATE NOT NULL, `items` TEXT NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"
	createChecklists := "CREATE TABLE IF NOT EXISTS `checklists` (`id` INT(11) NOT NULL AUTO_INCREMENT, `student_id` VARCHAR(255) NOT NULL, `title` TEXT NOT NULL, `items` TEXT NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"
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

// VerifyToken verifies a JWT token and returns the user ID
func VerifyToken(token *string) (string, error) {
	// Verify the JWT token
	SecretKey := os.Getenv("SECRET_KEY")
	VerifiedToken, err := utils.VerifyJWT(*token, SecretKey)
	if err != nil {
		return "", err
	}

	claims := VerifiedToken.Claims.(jwt.MapClaims)

	// Extract the user ID from the token claims
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("error: extracting user ID from token")
	}

	return userID, nil
}

// GetAccountByEmail returns a user by Email
func GetAccountByEmail(Email *string) (*models.Account, error) {
	// Prepare query
	query := "SELECT * FROM accounts WHERE email = ?"

	// Execute query
	row := db.QueryRow(query, Email)

	// Scan row into user object
	var user models.Account
	err := row.Scan(&user.DbId, &user.UserId, &user.Name, &user.Email, &user.Password, &user.PermissionInfo /*그냥 &PermissionInfo만 적어도 되나..?*/)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetAccountById returns a student by ID
func GetAccountById(id *models.UserId) (*models.Account, error) {
	// Prepare query
	query := "SELECT * FROM accounts WHERE user_id = ?"

	// Execute query
	row := db.QueryRow(query, id)

	// Scan row into user object
	var user models.Account
	err := row.Scan(&user.DbId, &user.UserId, &user.Name, &user.Email, &user.Password, &user.PermissionInfo /*그냥 &PermissionInfo만 적어도 되나..?*/)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateAccount creates a new student
func CreateAccount(account *models.Account) (models.DbId, error) {
	// Prepare query
	query := "INSERT INTO accounts (user_id, name, email, password, permission_level, school_id, timetable, grade, class, number, checklist_id, friends) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

	var result sql.Result
	var err error
	switch account.GetLevel() {
	case models.UNKNOWN:
		info := account.PermissionInfo.(models.Unknown)
		result, err = db.Exec(query, account.UserId, account.Name, account.Email, account.Password, info.GetLevel(), "", "", 0, 0)
	case models.STUDENT:
		info := account.PermissionInfo.(models.StudentInfo)
		result, err = db.Exec(query, account.UserId, account.Name, account.Email, account.Password, info.GetLevel(), info.SchoolId, info.Timetable, info.Grade, info.Class)
	case models.TEACHER:
		info := account.PermissionInfo.(models.TeacherInfo)
		result, err = db.Exec(query, account.UserId, account.Name, account.Email, account.Password, info.GetLevel(), info.SchoolId, "", 0, 0)
	case models.ADMIN:
		info := account.PermissionInfo.(models.AdminInfo)
		result, err = db.Exec(query, account.UserId, account.Name, account.Email, account.Password, info.GetLevel(), "", "", 0, 0)
	}

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
	// TODO 아래 쿼리문에 인자 더 추가하기, PermissionInfo가 기술할 수 있는 모든 종류의 인자가 있어야 함
	query := "UPDATE accounts SET user_id = ?, name = ?, email = ?, password = ?, school_id = ?, timetable = ?, grade = ?, class = ?, number = ?, checklist_id = ?, friends = ? WHERE id = ?"

	info := account.PermissionInfo.(models.StudentInfo)
	// Execute query
	_, err := db.Exec(query, account.UserId, account.Name, account.Email, account.Password, info.SchoolId, info.Timetable, info.Grade, info.Class, info.Number, info.ChecklistId, info.ChecklistId, info.Friends, account.DbId)
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
	query := "INSERT INTO timetables (id, teacher_id, location, day, period, subject) VALUES (?, ?, ?, ?, ?, ?)"

	// Execute query
	result, err := db.Exec(query, entry.ID, entry.TeacherId, entry.Location, entry.Day, entry.Period, entry.Subject)
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
	err := row.Scan(&menu.ID, &menu.SchoolId, &menu.MealName, &menu.Date, &menu.Items)
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
	err := row.Scan(&menu.ID, &menu.SchoolId, &menu.MealName, &menu.Date, &menu.Items)
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
	menuQuery := "INSERT INTO cafeteria_menus (id, school_id, meal_name, date, items) VALUES (?, ?, ?, ?, ?)"

	// Execute query to insert menu
	menuResult, err := db.Exec(menuQuery, menu.ID, menu.SchoolId, menu.MealName, menu.Date, menu.Items)
	if err != nil {
		return 0, err
	}

	// Get the ID of the newly created menu
	menuID, err := menuResult.LastInsertId()
	if err != nil {
		return 0, err
	}

	return models.DbId(menuID), nil
}

// UpdateMenu updates a cafeteria menu
func UpdateMenu(menu *models.CafeteriaMenu) error {
	err := utils.ValidateCafeteriaMenu(menu)
	if err != nil {
		return err
	}

	// Prepare query to insert menu
	menuQuery := "UPDATE cafeteria_menus SET school_id = ?, meal_name = ?, date = ?, items = ? WHERE id = ?"

	// Execute query to insert menu
	_, err = db.Exec(menuQuery, menu.SchoolId, menu.MealName, menu.Date, menu.Items, menu.ID)
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
func GetChecklistsOfStudent(studentID *models.UserId) (*models.Checklist, error) {
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
	createQuery := "INSERT INTO checklists (id, student_id, title, items) VALUES (?, ?, ?, ?)"

	// Execute query to insert menu
	result, err := db.Exec(createQuery, checklist.ID, checklist.StudentId, checklist.Title, checklist.Items)
	if err != nil {
		return 0, err
	}

	// Get the ID of the newly created checklist
	listId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

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
	query := "INSERT INTO schoolevents (id, school_id, month, events) VALUES (?, ?, ?, ?)"

	// Execute query
	result, err := db.Exec(query, events.ID, events.SchoolId, events.Month, &events.Events)
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
