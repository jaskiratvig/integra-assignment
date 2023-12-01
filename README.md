# Go Web Application

This project is a simple Go web application developed using Echo framework. It allows users to view, create, update, and delete user records in a database. The application follows the MVC pattern and includes Swagger documentation for API endpoints.

## Technical Requirements

- **Language:** Go 1.20
- **Web Framework:** Echo (https://echo.labstack.com/)
- **Data Access:** Squirrel (https://github.com/Masterminds/squirrel)
- **Testing:**
  - Ginkgo (https://github.com/onsi/ginkgo)
  - Gomega (https://github.com/onsi/gomega)

## Application Features

- **View Users:** Display a list of existing users in a simple grid format.
- **Create User:** Add new users to the system.
- **Update User:** Modify details of existing users.
- **Delete User:** Remove users from the system.
- **Error Handling:** Basic error handling and validation, such as checking for duplicate usernames.
- **Swagger Documentation:** Accessible at `http://localhost:8080/swagger/index.html` for API details.

## Running the Application

1. **Set up the Database:**
   - The application uses PostgreSQL, but any SQL database can be used.
   - Update the PostgreSQL connection string in `main.go` with your database credentials.

2. **Start the Server:**
   - Run `go run main.go` from the project root.
   - The server will start on `http://localhost:8080`.

3. **Viewing the User List:**
   - Navigate to `http://localhost:8080` in your web browser to see the table of users.

4. **API Endpoints:**
   - Use the Swagger UI at `http://localhost:8080/swagger/index.html` to interact with the API.

## Testing

- Run `go test ./...` in the project root to execute the test suite.

## Style Guide

- This project adheres to the Go style guidelines as per https://github.com/uber-go/guide/blob/master/style.md.
