package main

import (
	"database/sql"
	"integraAssignment/controllers"
	"log"

	_ "integraAssignment/docs"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func main() {

	// PostgreSQL connection string
	dbConnString := "postgres://myuser:mypassword@localhost/mydatabase?sslmode=disable"

	// Open database connection
	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Check if the database is connected
	if err = db.Ping(); err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Assign the database connection to the controllers
	controllers.SetDatabase(db)

	e := echo.New()

	// User routes
	e.GET("/users", controllers.GetAllUsers)
	e.POST("/users", controllers.CreateUser)
	e.GET("/users/:id", controllers.GetUser)
	e.PUT("/users/:id", controllers.UpdateUser)
	e.DELETE("/users/:id", controllers.DeleteUser)

	// Serve static files
	// Navigate to http://localhost:8080/ to see the table of users
	e.Static("/", "public")

	// Swagger route
	// Navigate to http://localhost:8080/swagger/index.html to view the swagger file
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start the server
	e.Logger.Fatal(e.Start(":8080"))
}
