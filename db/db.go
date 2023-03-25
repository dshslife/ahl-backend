package db

import (
	"database/sql"
	"fmt"
	"github.com/username/schoolapp/models"
	"github.com/username/schoolapp/utils"
	"log"
	"os"
	"time"
)

var db *sql.DB

func createTables() {
	// prepare query
	query := "CREATE TABLE IF NOT EXISTS `users` (`id` INT(11) NOT NULL AUTO_INCREMENT, `google_id` VARCHAR(255) NOT NULL, `name` VARCHAR(255) NOT NULL, `email` VARCHAR(255) NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci; CREATE TABLE IF NOT EXISTS `teachers` (`id` INT(11) NOT NULL AUTO_INCREMENT, `google_id` VARCHAR(255) NOT NULL, `name` VARCHAR(255) NOT NULL, `email` VARCHAR(255) NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci; CREATE TABLE IF NOT EXISTS `admins` {`id` INT(11) NOT NULL AUTO_INCREMENT, `google_id` VARCHAR(255) NOT NULL, `name` VARCHAR(255) NOT NULL, `email` VARCHAR(255) NOT NULL, PRIMARY KEY (`id`)} ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;CREATE TABLE IF NOT EXISTS `students` (`id` INT(11) NOT NULL AUTO_INCREMENT, `name` VARCHAR(255) NOT NULL, `grade` INT(11) NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci; CREATE TABLE IF NOT EXISTS `lessons` (`id` INT(11) NOT NULL AUTO_INCREMENT, `user_id` INT(11) NOT NULL, `name` VARCHAR(255) NOT NULL, `teacher` VARCHAR(255) NOT NULL, `location` VARCHAR(255) NOT NULL, `period` TIME NOT NULL, `day` INT(11) NOT NULL, PRIMARY KEY (`id`), FOREIGN KEY (`user_id`) REFERENCES `users`(`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci; CREATE TABLE IF NOT EXISTS `cafeteria_menus` ( `id` INT(11) NOT NULL AUTO_INCREMENT, `date` DATE NOT NULL, `meal` VARCHAR(255) NOT NULL, `items` TEXT NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"

	// Execute query
	db.QueryRow(query)
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
// 제발 SecretKey 생성 좀 만들어줘요
func VerifyToken(token string) (string, error) {
	// Verify the JWT token
	claims, err := utils.VerifyJWT(token, secretKey)
	if err != nil {
		return "", err
	}

	// Extract the user ID from the token claims
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", fmt.Errorf("Error extracting user ID from token")
	}

	return userID, nil
}

// GetUserByEmail returns a user by Email
func GetUserByEmail(Email string) (*models.User, error) {
	// Prepare query
	query := "SELECT * FROM users WHERE email = ?"

	// Execute query
	row := db.QueryRow(query, Email)

	// Scan row into user object
	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetAllStudents returns all student
func GetAllStudents() (*models.Student, error) {
	// Prepare query
	query := "SELECT * FROM students"

	// Execute query
	row := db.QueryRow(query)

	// Scan row into student object
	var student models.Student
	err := row.Scan(&student.ID, &student.Name, &student.Email, &student.Grade)
	if err != nil {
		return nil, err
	}

	return &student, nil
}

// GetStudentByID returns a student by ID
func GetStudentByID(id string) (*models.Student, error) {
	// Prepare query
	query := "SELECT * FROM students WHERE id = ?"

	// Execute query
	row := db.QueryRow(query, id)

	// Scan row into student object
	var student models.Student
	err := row.Scan(&student.ID, &student.Name, &student.Email, &student.Grade)
	if err != nil {
		return nil, err
	}

	return &student, nil
}

// CreateStudent creates a new student
func CreateStudent(student *models.Student) (int64, error) {
	// Prepare query
	query := "INSERT INTO students (id, name, email, grade) VALUES (?, ?, ?, ?)"

	// Execute query
	result, err := db.Exec(query, student.ID, student.Name, student.Email, student.Grade)
	if err != nil {
		return 0, err
	}

	// Get the ID of the newly created student
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	student.ID = int(id)

	return id, nil
}

// UpdateStudent updates a student
func UpdateStudent(student *models.Student) error {
	// Prepare query
	query := "UPDATE students SET name = ?, email = ?, grade = ? WHERE id = ?"

	// Execute query
	_, err := db.Exec(query, student.Name, student.Email, student.Grade, student.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteStudent deletes a student by ID
func DeleteStudent(id int) error {
	// Prepare query
	query := "DELETE FROM students WHERE id = ?"

	// Execute query
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

// Get all students
func GetStudents() (*models.Student, error) {
	students, err := GetAllStudents()
	if err != nil {
		return nil, err
	}

	return students, err
}

// GetTimetable returns a list of timetables for a student
func GetTimetable(studentID string) (models.Timetables, error) {
	// Prepare query
	query := "SELECT * FROM timetables WHERE student_id = ?"

	// Execute query
	rows, err := db.Query(query, studentID)
	if err != nil {
		return nil, err
	}

	// Scan rows into timetable objects
	var timetables models.Timetables
	for rows.Next() {
		var timetable models.Timetable
		err = rows.Scan(&timetable.ID, &timetable.StudentID, &timetable.Day, &timetable.Period, &timetable.Subject)
		if err != nil {
			return nil, err
		}

		timetables = append(timetables, timetable)
	}

	return timetables, nil
}

// GetTimetableByID returns a list of timetables for a student
func GetTimetableByID(studentID int) (*models.Timetable, error) {
	// Prepare query
	query := "SELECT * FROM timetables WHERE student_id = ?"

	// Execute query
	rows, err := db.Query(query, studentID)
	if err != nil {
		return nil, err
	}

	// Scan rows into timetable objects
	var timetable models.Timetable
	err = rows.Scan(&timetable.ID, &timetable.StudentID, &timetable.Day, &timetable.Period, &timetable.Subject)
	if err != nil {
		return nil, err
	}

	return &timetable, nil
}

// CreateTimetable creates a new timetable
func CreateTimetable(timetable *models.Timetable) (int64, error) {
	// Prepare query
	query := "INSERT INTO timetables (student_id, day, period, subject) VALUES (?, ?, ?, ?)"

	// Execute query
	result, err := db.Exec(query, timetable.StudentID, timetable.Day, timetable.Period, timetable.Subject)
	if err != nil {
		return 0, err
	}

	// Get the ID of the newly created timetable
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	timetable.ID = int(id)

	return id, nil
}

// UpdateTimetable updates a timetable
func UpdateTimetable(timetable *models.Timetable) error {
	// Prepare query
	query := "UPDATE timetables SET student_id = ?, day = ?, period = ?, subject = ? WHERE id = ?"

	// Execute query
	_, err := db.Exec(query, timetable.StudentID, timetable.Day, timetable.Period, timetable.Subject, timetable.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTimetable deletes a timetable by ID
func DeleteTimetable(id int) error {
	// Prepare query
	query := "DELETE FROM timetables WHERE id = ?"

	// Execute query
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

// GetMenu returns the cafeteria menu for a specific date and meal
func GetMenu(date string, meal string) (*models.CafeteriaMenu, error) {
	// Prepare query
	query := "SELECT * FROM cafeteria_menu WHERE date = ? AND meal = ?"

	// Execute query
	row := db.QueryRow(query, date, meal)

	// Scan row into cafeteria menu object
	var menu models.CafeteriaMenu
	err := row.Scan(&menu.ID, &menu.Date, &menu.Meal)
	if err != nil {
		return nil, err
	}

	// Prepare query to get menu items
	query = "SELECT * FROM menu_items WHERE menu_id = ?"

	// Execute query
	rows, err := db.Query(query, menu.ID)
	if err != nil {
		return nil, err
	}

	// Scan rows into item objects and add them to the menu
	for rows.Next() {
		var item models.Item
		err = rows.Scan(&item.ID, &item.Name, &item.Allergy, &item.Vegetari)
		if err != nil {
			return nil, err
		}

		menu.Items = append(menu.Items, item)
	}

	return &menu, nil
}

// CreateMenu creates a new cafeteria menu
func CreateMenu(menu *models.CafeteriaMenu) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare query to insert menu
	menuQuery := "INSERT INTO cafeteria_menu (date, meal) VALUES (?, ?)"

	// Execute query to insert menu
	menuResult, err := tx.Exec(menuQuery, menu.Date, menu.Meal)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Get the ID of the newly created menu
	menuID, err := menuResult.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Prepare query to insert menu items
	itemQuery := "INSERT INTO menu_items (menu_id, name, allergy, vegetari) VALUES (?, ?, ?, ?)"

	// Execute query to insert menu items
	for _, item := range menu.Items {
		_, err = tx.Exec(itemQuery, menuID, item.Name, item.Allergy, item.Vegetari)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// UpdateMenu updates a cafeteria menu
func UpdateMenu(menu *models.CafeteriaMenu) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare query to update menu
	menuQuery := "UPDATE cafeteria_menu SET date = ?, meal = ? WHERE id = ?"

	// Execute query to update menu
	_, err = tx.Exec(menuQuery, menu.Date, menu.Meal, menu.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Prepare query to delete existing menu items
	deleteItemsQuery := "DELETE FROM menu_items WHERE menu_id = ?"

	// Execute query to delete existing menu items
	_, err = tx.Exec(deleteItemsQuery, menu.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Prepare query to insert new menu items
	insertItemsQuery := "INSERT INTO menu_items (menu_id, name, allergy, vegetari) VALUES (?, ?, ?, ?)"

	// Execute query to insert new menu items
	for _, item := range menu.Items {
		_, err = tx.Exec(insertItemsQuery, menu.ID, item.Name, item.Allergy, item.Vegetari)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

// DeleteMenu deletes a cafeteria menu by ID
func DeleteMenu(id int) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare query to delete menu items
	deleteItemsQuery := "DELETE FROM menu_items WHERE menu_id = ?"

	// Execute query to delete menu items
	_, err = tx.Exec(deleteItemsQuery, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Prepare query to delete menu
	deleteMenuQuery := "DELETE FROM cafeteria_menu WHERE id = ?"

	// Execute query to delete menu
	_, err = tx.Exec(deleteMenuQuery, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
