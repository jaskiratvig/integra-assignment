package controllers

import (
	"database/sql"
	"integraAssignment/models"
	"log"
	"net/http"

	"github.com/Masterminds/squirrel"
	"github.com/labstack/echo/v4"
)

// db is the database connection object
var db *sql.DB

// SetDatabase assigns a database connection to the controllers
func SetDatabase(database *sql.DB) {
	db = database
}

// CreateUser godoc
// @Summary Create a new user
// @Description Add a new user to the system
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.User true "Create User"
// @Success 201 {object} models.ApiResponse
// @Failure 400 {object} models.ApiResponse
// @Failure 500 {object} models.ApiResponse
// @Router /users [post]
func CreateUser(c echo.Context) error {
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, models.ApiResponse{Success: false, Message: "Invalid request", Data: nil})
	}

	sb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Check if username exists
	query, args, err := sb.Select("user_name").From("users").Where(squirrel.Eq{"user_name": u.UserName}).ToSql()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ApiResponse{Success: false, Message: "Error building query", Data: nil})
	}
	err = db.QueryRow(query, args...).Scan(&u.UserName)
	if err == nil {
		return c.JSON(http.StatusBadRequest, models.ApiResponse{Success: false, Message: "Username already exists", Data: nil})
	}

	// Create user
	insertQuery, insertArgs, err := sb.Insert("users").Columns("user_name", "first_name", "last_name", "email", "user_status", "department").
		Values(u.UserName, u.FirstName, u.LastName, u.Email, u.UserStatus, u.Department).ToSql()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ApiResponse{Success: false, Message: "Error building insert query", Data: nil})
	}
	_, err = db.Exec(insertQuery, insertArgs...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ApiResponse{Success: false, Message: "Error inserting user", Data: nil})
	}

	return c.JSON(http.StatusCreated, models.ApiResponse{Success: true, Message: "User created successfully", Data: u})
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieve a list of all users
// @Tags users
// @Produce  json
// @Success 200 {object} []models.User
// @Failure 500 {object} models.ApiResponse
// @Router /users [get]
func GetAllUsers(c echo.Context) error {
	var users []models.User

	sb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Building query using Squirrel
	query, args, err := sb.Select("*").From("users").ToSql()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ApiResponse{Success: false, Message: "Error building query", Data: nil})
	}

	// Executing query
	log.Println("Executing SQL query:", query, "Args:", args)
	rows, err := db.Query(query, args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return c.JSON(http.StatusInternalServerError, models.ApiResponse{Success: false, Message: "Error executing query", Data: nil})
	}
	defer rows.Close()

	// Scanning rows
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.UserName, &u.FirstName, &u.LastName, &u.Email, &u.UserStatus, &u.Department); err != nil {
			return c.JSON(http.StatusInternalServerError, models.ApiResponse{Success: false, Message: "Error scanning user data", Data: nil})
		}
		users = append(users, u)
	}

	// Checking for errors during row iteration
	if err := rows.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, models.ApiResponse{Success: false, Message: "Error completing user retrieval", Data: nil})
	}

	return c.JSON(http.StatusOK, models.ApiResponse{Success: true, Message: "Users retrieved successfully", Data: users})
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Retrieve a user by their ID
// @Tags users
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} models.ApiResponse
// @Failure 500 {object} models.ApiResponse
// @Router /users/{id} [get]
func GetUser(c echo.Context) error {
	id := c.Param("id")
	u := models.User{}

	sb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Building query using Squirrel
	query, args, err := sb.Select("*").From("users").Where(squirrel.Eq{"user_id": id}).ToSql()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ApiResponse{Success: false, Message: "Error building query", Data: nil})
	}

	// Executing query
	err = db.QueryRow(query, args...).Scan(&u.ID, &u.UserName, &u.FirstName, &u.LastName, &u.Email, &u.UserStatus, &u.Department)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, models.ApiResponse{Success: false, Message: "User not found", Data: nil})
		}
		return c.JSON(http.StatusInternalServerError, models.ApiResponse{Success: false, Message: "Error retrieving user", Data: nil})
	}

	return c.JSON(http.StatusOK, models.ApiResponse{Success: true, Message: "User retrieved successfully", Data: u})
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update an existing user's details
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Param user body models.User true "Update User"
// @Success 200 {object} models.ApiResponse
// @Failure 400 {object} models.ApiResponse
// @Failure 500 {object} models.ApiResponse
// @Router /users/{id} [put]
func UpdateUser(c echo.Context) error {
	id := c.Param("id")
	u := new(models.User)

	sb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, models.ApiResponse{Success: false, Message: "Invalid request", Data: nil})
	}

	// Building query using Squirrel
	query, args, err := sb.Update("users").SetMap(map[string]interface{}{
		"user_name": u.UserName, "first_name": u.FirstName, "last_name": u.LastName,
		"email": u.Email, "user_status": u.UserStatus, "department": u.Department,
	}).Where(squirrel.Eq{"user_id": id}).ToSql()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ApiResponse{Success: false, Message: "Error building query", Data: nil})
	}

	// Executing query
	_, err = db.Exec(query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ApiResponse{Success: false, Message: "Error updating user", Data: nil})
	}

	return c.JSON(http.StatusOK, models.ApiResponse{Success: true, Message: "User updated successfully", Data: u})
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Remove a user from the system by their ID
// @Tags users
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} models.ApiResponse
// @Failure 404 {object} models.ApiResponse
// @Failure 500 {object} models.ApiResponse
// @Router /users/{id} [delete]
func DeleteUser(c echo.Context) error {
	id := c.Param("id")

	sb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	// Building query using Squirrel
	query, args, err := sb.Delete("users").Where(squirrel.Eq{"user_id": id}).ToSql()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ApiResponse{Success: false, Message: "Error building query", Data: nil})
	}

	// Executing query
	result, err := db.Exec(query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ApiResponse{Success: false, Message: "Error deleting user", Data: nil})
	}

	// Checking if the user was actually deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ApiResponse{Success: false, Message: "Error checking deletion result", Data: nil})
	}

	if rowsAffected == 0 {
		return c.JSON(http.StatusNotFound, models.ApiResponse{Success: false, Message: "User not found", Data: nil})
	}

	return c.JSON(http.StatusOK, models.ApiResponse{Success: true, Message: "User deleted successfully", Data: nil})
}
