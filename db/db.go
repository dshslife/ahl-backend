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
	createUsers := "CREATE TABLE IF NOT EXISTS `users` (`id` INT(11) NOT NULL AUTO_INCREMENT, `google_id` VARCHAR(255) NOT NULL, `name` VARCHAR(255) NOT NULL, `email` VARCHAR(255) NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"
	createTeachers := "CREATE TABLE IF NOT EXISTS `teachers` (`id` INT(11) NOT NULL AUTO_INCREMENT, `google_id` VARCHAR(255) NOT NULL, `name` VARCHAR(255) NOT NULL, `email` VARCHAR(255) NOT NULL, `access` VARCHAR(255) NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"
	createAdmins := "CREATE TABLE IF NOT EXISTS `admins` (`id` INT(11) NOT NULL AUTO_INCREMENT, `google_id` VARCHAR(255) NOT NULL, `name` VARCHAR(255) NOT NULL, `email` VARCHAR(255) NOT NULL, `access_to_all` VARCHAR(255) NOT NULL, `admin_access` VARCHAR(255) NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"
	createStudents := "CREATE TABLE IF NOT EXISTS `students` (`id` INT(11) NOT NULL AUTO_INCREMENT, `name` VARCHAR(255) NOT NULL, `email` VARCHAR(255) NOT NULL, `grade` INT(11) NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"
	createLessons := "CREATE TABLE IF NOT EXISTS `lessons` (`id` INT(11) NOT NULL AUTO_INCREMENT, `user_id` INT(11) NOT NULL, `name` VARCHAR(255) NOT NULL, `teacher` VARCHAR(255) NOT NULL, `location` VARCHAR(255) NOT NULL, `period` TIME NOT NULL, `day` INT(11) NOT NULL, PRIMARY KEY (`id`), FOREIGN KEY (`user_id`) REFERENCES `users`(`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"
	createCafeteria := "CREATE TABLE IF NOT EXISTS `cafeteria_menus` ( `id` INT(11) NOT NULL AUTO_INCREMENT, `date` DATE NOT NULL, `meal` VARCHAR(255) NOT NULL, `items` TEXT NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"
	createChecklists := "CREATE TABLE IF NOT EXISTS `checklists` (`id` INT(11) NOT NULL AUTO_INCREMENT, `title` TEXT NOT NULL, `UserID` VARCHAR(255) NOT NULL, `items` TEXT NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"
	createEvents := "CREATE TABLE IF NOT EXISTS `schoolevents` (`id` INT(11) NOT NULL AUTO_INCREMENT, `month` INT(11) NOT NULL, `school` TEXT NOT NULL, `events` VARCHAR(255) NOT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;"

	// Execute query
	queries := [8]string{createUsers,
		createTeachers,
		createAdmins,
		createStudents,
		createLessons,
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
func VerifyToken(token string) (string, error) {
	// Verify the JWT token
	SecretKey := os.Getenv("SECRET_KEY")
	VerifiedToken, err := utils.VerifyJWT(token, SecretKey)
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

/*// Get all students
func GetStudents() (*models.Student, error) {
	students, err := GetAllStudents()
	if err != nil {
		return nil, err
	}

	return students, err
}*/

// GetAllTeachers returns all teachers
func GetAllTeachers() (*models.Teacher, error) {
	// Prepare query
	query := "SELECT * FROM teachers"

	// Execute query
	row := db.QueryRow(query)

	// Scan row into teacher object
	var teacher models.Teacher
	err := row.Scan(&teacher.ID, &teacher.Name, &teacher.Email, &teacher.Access)
	if err != nil {
		return nil, err
	}

	return &teacher, nil
}

// GetTeacherByID returns a teacher by ID
func GetTeacherByID(id string) (*models.Teacher, error) {
	// Prepare query
	query := "SELECT * FROM teachers WHERE id = ?"

	// Execute query
	row := db.QueryRow(query, id)

	// Scan row into teacher object
	var teacher models.Teacher
	err := row.Scan(&teacher.ID, &teacher.Name, &teacher.Email, &teacher.Access)
	if err != nil {
		return nil, err
	}

	return &teacher, nil
}

// CreateTeacher creates a new teacher
func CreateTeacher(teacher *models.Teacher) (int64, error) {
	// Prepare query
	query := "INSERT INTO teachers (id, name, email, access) VALUES (?, ?, ?, ?)"

	// Execute query
	result, err := db.Exec(query, teacher.ID, teacher.Name, teacher.Email, teacher.Access)
	if err != nil {
		return 0, err
	}

	// Get the ID of the newly created student
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	teacher.ID = int(id)

	return id, nil
}

// UpdateTeacher updates a teacher
func UpdateTeacher(teacher *models.Teacher) error {
	// Prepare query
	query := "UPDATE teachers SET name = ?, email = ?, access = ? WHERE id = ?"

	// Execute query
	_, err := db.Exec(query, teacher.Name, teacher.Email, teacher.Access, teacher.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteTeacher deletes a teacher by ID
func DeleteTeacher(id int) error {
	// Prepare query
	query := "DELETE FROM teachers WHERE id = ?"

	// Execute query
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}

// GetAllAdmins returns all admins
func GetAllAdmins() (*models.Admin, error) {
	// Prepare query
	query := "SELECT * FROM admins"

	// Execute query
	row := db.QueryRow(query)

	// Scan row into admin object
	var admin models.Admin
	err := row.Scan(&admin.ID, &admin.Name, &admin.Email, &admin.AccessToAll, &admin.AdminAccess)
	if err != nil {
		return nil, err
	}

	return &admin, nil
}

// GetAdminByID returns an admin by ID
func GetAdminByID(id string) (*models.Admin, error) {
	// Prepare query
	query := "SELECT * FROM admins WHERE id = ?"

	// Execute query
	row := db.QueryRow(query, id)

	// Scan row into admin object
	var admin models.Admin
	err := row.Scan(&admin.ID, &admin.Name, &admin.Email, &admin.AccessToAll, &admin.AdminAccess)
	if err != nil {
		return nil, err
	}

	return &admin, nil
}

// CreateAdmin creates a new admin
func CreateAdmin(admin *models.Admin) (int64, error) {
	// Prepare query
	query := "INSERT INTO admins (id, name, email, access_to_all, admin_access) VALUES (?, ?, ?, ?, ?)"

	// Execute query
	result, err := db.Exec(query, admin.ID, admin.Name, admin.Email, admin.AccessToAll, admin.AdminAccess)
	if err != nil {
		return 0, err
	}

	// Get the ID of the newly created student
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	admin.ID = int(id)

	return id, nil
}

// UpdateAdmin updates an admin
func UpdateAdmin(admin *models.Admin) error {
	// Prepare query
	query := "UPDATE admins SET name = ?, email = ?, access_to_all = ?, admin_access = ? WHERE id = ?"

	// Execute query
	_, err := db.Exec(query, admin.Name, admin.Email, admin.AccessToAll, admin.AdminAccess, admin.ID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteAdmin deletes an admin by ID
func DeleteAdmin(id int) error {
	// Prepare query
	query := "DELETE FROM admins WHERE id = ?"

	// Execute query
	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
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
		err = rows.Scan(&timetable.ID, &timetable.StudentID, &timetable.Day, &timetable.Period, &timetable.Subject, &timetable.IsPublic)
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
	err = rows.Scan(&timetable.ID, &timetable.StudentID, &timetable.Day, &timetable.Period, &timetable.Subject, &timetable.IsPublic)
	if err != nil {
		return nil, err
	}

	return &timetable, nil
}

// CreateTimetable creates a new timetable
func CreateTimetable(timetable *models.Timetable) (int64, error) {
	// Prepare query
	query := "INSERT INTO timetables (student_id, day, period, subject, IsPublic) VALUES (?, ?, ?, ?, ?)"

	// Execute query
	result, err := db.Exec(query, timetable.StudentID, timetable.Day, timetable.Period, timetable.Subject, timetable.IsPublic)
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
	query := "UPDATE timetables SET student_id = ?, day = ?, period = ?, subject = ?, IsPublic = ? WHERE id = ?"

	// Execute query
	_, err := db.Exec(query, timetable.StudentID, timetable.Day, timetable.Period, timetable.Subject, timetable.IsPublic, timetable.ID)
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

// GetMenuByID returns the cafeteria menu for a specific date by ID
func GetMenuByID(id int) (*models.CafeteriaMenu, error) {
	// Prepare query
	query := "SELECT * FROM cafeteria_menu WHERE id = ?"

	// Execute query
	row := db.QueryRow(query, id)

	// Scan row into cafeteria menu object
	var menu models.CafeteriaMenu
	err := row.Scan(&menu.Date, &menu.Meal, &menu.Items)
	if err != nil {
		return nil, err
	}

	return &menu, nil
}

// GetMenu returns the cafeteria menu for a specific date
func GetMenu(date time.Time) (*models.CafeteriaMenu, error) {
	// Prepare query
	query := "SELECT * FROM cafeteria_menu WHERE date = ?"

	// Execute query
	row := db.QueryRow(query, date)

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

// GetCafeteriaMenus returns the cafeteria menus
func GetCafeteriaMenus() (*models.CafeteriaMenu, error) {
	// Prepare query
	query := "SELECT * FROM cafeteria_menu"

	// Execute query
	row := db.QueryRow(query)

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

// CreateCafeteriaMenu creates a new cafeteria menu
func CreateCafeteriaMenu(menu *models.CafeteriaMenu) (int, error) {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return -1, err
	}

	// Prepare query to insert menu
	menuQuery := "INSERT INTO cafeteria_menu (date, meal) VALUES (?, ?)"

	// Execute query to insert menu
	menuResult, err := tx.Exec(menuQuery, menu.Date, menu.Meal)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	// Get the ID of the newly created menu
	menuID, err := menuResult.LastInsertId()
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	// Prepare query to insert menu items
	itemQuery := "INSERT INTO menu_items (menu_id, name, allergy, vegetari) VALUES (?, ?, ?, ?)"

	// Execute query to insert menu items
	for _, item := range menu.Items {
		_, err = tx.Exec(itemQuery, menuID, item.Name, item.Allergy, item.Vegetari)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	return menu.ID, err
}

// UpdateCafeteriaMenu updates a cafeteria menu
func UpdateCafeteriaMenu(menu *models.CafeteriaMenu) error {
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

// DeleteCafeteriaMenu deletes a cafeteria menu by ID
func DeleteCafeteriaMenu(id int) error {
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

// GetChecklists returns a checklist
func GetChecklists(studentID string) (*models.Checklist, error) {
	// Prepare query
	query := "SELECT * FROM checklists WHERE UserID = ?"

	// Execute query
	row := db.QueryRow(query, studentID)

	// Scan row into checklist object
	var checklist models.Checklist
	err := row.Scan(&checklist.ID, &checklist.Title, &checklist.Items, &checklist.UserID)
	if err != nil {
		return nil, err
	}

	return &checklist, nil
}

// CreateChecklist creates a new checklist
func CreateChecklist(checklist *models.Checklist) (int, error) {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return -1, err
	}

	// Prepare query to insert checklist
	checklistQuery := "INSERT INTO checklists (id, title, UserID) VALUES (?, ?, ?)"

	// Execute query to insert checklist
	checklistResult, err := tx.Exec(checklistQuery, checklist.ID, checklist.Title, checklist.UserID)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	// Get the ID of the newly created checklist
	checklistID, err := checklistResult.LastInsertId()
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	// Prepare query to insert checklist items
	itemQuery := "INSERT INTO checklist_items (id, item_id, Text, Complete, IsPublic, SharedWith) VALUES (?, ?, ?, ?, ?, ?)"

	// Execute query to insert checklist items
	for _, item := range checklist.Items {
		_, err = tx.Exec(itemQuery, checklistID, item.ID, item.Text, item.Complete, item.IsPublic, item.SharedWith)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	return checklist.ID, err
}

// UpdateChecklist updates a checklist
func UpdateChecklist(checklist *models.Checklist) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare query to update checklist
	checklistQuery := "UPDATE checklist SET id = ?, Title = ? WHERE UserID = ?"

	// Execute query to update checklist
	_, err = tx.Exec(checklistQuery, checklist.ID, checklist.Title, checklist.UserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Prepare query to delete existing checklist items
	deleteItemsQuery := "DELETE FROM checklist_items WHERE id = ?"

	// Execute query to delete existing checklist items
	_, err = tx.Exec(deleteItemsQuery, checklist.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Prepare query to insert new checklist items
	insertItemsQuery := "INSERT INTO checklist_items (id, item_id, Text, Complete, IsPublic, SharedWith) VALUES (?, ?, ?, ?, ?, ?)"

	// Execute query to insert new checklist items
	for _, item := range checklist.Items {
		_, err = tx.Exec(insertItemsQuery, checklist.ID, item.ID, item.Text, item.Complete, item.IsPublic, item.SharedWith)
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

// DeleteChecklist deletes a checklist by ID
func DeleteChecklist(id int) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare query to delete checklist items
	deleteItemsQuery := "DELETE FROM checklist_items WHERE id = ?"

	// Execute query to delete checklist items
	_, err = tx.Exec(deleteItemsQuery, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Prepare query to delete checklist
	deleteChecklistQuery := "DELETE FROM checklist WHERE id = ?"

	// Execute query to delete checklist
	_, err = tx.Exec(deleteChecklistQuery, id)
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

// GetChecklistItemByID returns the checklist item by ID
func GetChecklistItemByID(id int) (*models.Checklist, error) {
	// Prepare query
	query := "SELECT * FROM checklist WHERE id = ?"

	// Execute query
	row := db.QueryRow(query, id)

	// Scan row into checklist object
	var checklist models.Checklist
	err := row.Scan(&checklist.ID, &checklist.UserID, &checklist.Title)
	if err != nil {
		return nil, err
	}

	// Prepare query to get checklist items
	query = "SELECT * FROM checklist_items WHERE id = ?"

	// Execute query
	rows, err := db.Query(query, checklist.ID)
	if err != nil {
		return nil, err
	}

	// Scan rows into item objects and add them to the checklist
	for rows.Next() {
		var item models.Items
		err = rows.Scan(&item.ID, &item.Text, &item.Complete, &item.IsPublic, &item.SharedWith)
		if err != nil {
			return nil, err
		}

		checklist.Items = append(checklist.Items, item)
	}

	return &checklist, nil
}

// GetChecklistItems returns the checklist items
func GetChecklistItems(id int) (*models.Checklist, error) {
	// Prepare query
	query := "SELECT * FROM checklist WHERE id = ?"

	// Execute query
	row := db.QueryRow(query, id)

	// Scan row into checklist object
	var checklist models.Checklist
	err := row.Scan(&checklist.ID, &checklist.UserID, &checklist.Title)
	if err != nil {
		return nil, err
	}

	// Prepare query to get checklist items
	query = "SELECT * FROM checklist_items WHERE id = ?"

	// Execute query
	rows, err := db.Query(query, checklist.ID)
	if err != nil {
		return nil, err
	}

	// Scan rows into item objects and add them to the checklist
	for rows.Next() {
		var item models.Items
		err = rows.Scan(&item.ID, &item.Text, &item.Complete, &item.IsPublic, &item.SharedWith)
		if err != nil {
			return nil, err
		}

		checklist.Items = append(checklist.Items, item)
	}

	return &checklist, nil
}

// CreateChecklistItem creates a new checklist item
func CreateChecklistItem(checklist *models.Checklist) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare query to insert checklist
	checklistQuery := "INSERT INTO checklists (id, title, UserID) VALUES (?, ?, ?)"

	// Execute query to insert checklist
	checklistResult, err := tx.Exec(checklistQuery, checklist.ID, checklist.Title, checklist.UserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Get the ID of the newly created checklist
	checklistID, err := checklistResult.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	// Prepare query to insert checklist items
	itemQuery := "INSERT INTO checklist_items (id, item_id, Text, Complete, IsPublic, SharedWith) VALUES (?, ?, ?, ?, ?, ?)"

	// Execute query to insert checklist items
	for _, item := range checklist.Items {
		_, err = tx.Exec(itemQuery, checklistID, item.ID, item.Text, item.Complete, item.IsPublic, item.SharedWith)
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

	return err
}

// UpdateChecklistItem updates a checklist item
func UpdateChecklistItem(checklist *models.Checklist) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare query to update checklist
	checklistQuery := "UPDATE checklist SET id = ?, Title = ? WHERE UserID = ?"

	// Execute query to update checklist
	_, err = tx.Exec(checklistQuery, checklist.ID, checklist.Title, checklist.UserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Prepare query to delete existing checklist items
	deleteItemsQuery := "DELETE FROM checklist_items WHERE id = ?"

	// Execute query to delete existing checklist items
	_, err = tx.Exec(deleteItemsQuery, checklist.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Prepare query to insert new checklist items
	insertItemsQuery := "INSERT INTO checklist_items (id, item_id, Text, Complete, IsPublic, SharedWith) VALUES (?, ?, ?, ?, ?, ?)"

	// Execute query to insert new checklist items
	for _, item := range checklist.Items {
		_, err = tx.Exec(insertItemsQuery, checklist.ID, item.ID, item.Text, item.Complete, item.IsPublic, item.SharedWith)
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

// DeleteChecklistItem deletes a checklist item by ID
func DeleteChecklistItem(id int) error {
	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Prepare query to delete checklist items
	deleteItemsQuery := "DELETE FROM checklist_items WHERE id = ?"

	// Execute query to delete checklist items
	_, err = tx.Exec(deleteItemsQuery, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Prepare query to delete checklist
	deleteChecklistQuery := "DELETE FROM checklist WHERE id = ?"

	// Execute query to delete checklist
	_, err = tx.Exec(deleteChecklistQuery, id)
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

// GetAllEvents returns all events
func GetAllEvents() (*models.Events, error) {
	// Prepare query
	query := "SELECT * FROM schoolevents"

	// Execute query
	row := db.QueryRow(query)

	// Scan row into student object
	var events models.Events
	err := row.Scan(&events.ID, &events.Month, &events.School, &events.Event, &events.Exists)
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
	err := row.Scan(&events.ID, &events.Month, &events.School, &events.Event, &events.Exists)
	if err != nil {
		return nil, err
	}

	return &events, nil
}

// CreateEvents creates a new event
func CreateEvents(events *models.Events) (int64, error) {
	// Prepare query
	query := "INSERT INTO schoolevents (month, school, events, exists) VALUES (?, ?, ?, ?)"

	// Execute query
	result, err := db.Exec(query, events.Month, events.School, events.Event, &events.Exists)
	if err != nil {
		return 0, err
	}

	// Get the ID of the newly created student
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	events.ID = int(id)

	return id, nil
}
