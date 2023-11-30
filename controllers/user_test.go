package controllers_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"integraAssignment/controllers"
	"integraAssignment/models"
	"net/http"
	"net/http/httptest"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("User Controller Tests", func() {
	var (
		e    *echo.Echo
		req  *http.Request
		rec  *httptest.ResponseRecorder
		mock sqlmock.Sqlmock
		db   *sql.DB
		err  error
	)

	setupMockAndRequest := func(method, path string, body []byte) {
		e = echo.New()
		rec = httptest.NewRecorder()
		db, mock, err = sqlmock.New()
		Expect(err).ShouldNot(HaveOccurred())
		controllers.SetDatabase(db)
		req = httptest.NewRequest(method, path, bytes.NewBuffer(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}

	AfterEach(func() {
		db.Close()
		mock.ExpectClose()
	})

	Describe("CreateUser", func() {
		Context("when the request is valid", func() {
			It("creates a new user", func() {
				user := models.User{
					UserName:   "testuser",
					FirstName:  "Test",
					LastName:   "User",
					Email:      "testuser@example.com",
					UserStatus: "A",
					Department: "Development",
				}
				requestBody, _ := json.Marshal(user)
				setupMockAndRequest(http.MethodPost, "/users", requestBody)

				mock.ExpectQuery(`SELECT user_name FROM users WHERE user_name = \$1`).
					WithArgs(user.UserName).
					WillReturnRows(sqlmock.NewRows([]string{"user_name"}))

				mock.ExpectExec(`INSERT INTO users`).
					WithArgs(user.UserName, user.FirstName, user.LastName, user.Email, user.UserStatus, user.Department).
					WillReturnResult(sqlmock.NewResult(1, 1))

				c := e.NewContext(req, rec)
				err := controllers.CreateUser(c)
				Expect(err).ToNot(HaveOccurred())
				Expect(rec.Code).To(Equal(http.StatusCreated))
			})
		})
	})

	Describe("GetAllUsers", func() {
		Context("when the query executes successfully", func() {
			It("retrieves all users", func() {
				setupMockAndRequest(http.MethodGet, "/users", nil)

				mock.ExpectQuery(`SELECT \* FROM users`).
					WillReturnRows(sqlmock.NewRows([]string{"id", "user_name", "first_name", "last_name", "email", "user_status", "department"}).
						AddRow(1, "testuser", "Test", "User", "testuser@example.com", "active", "Development").
						AddRow(2, "anotheruser", "Another", "User", "anotheruser@example.com", "inactive", "HR"))

				c := e.NewContext(req, rec)
				err := controllers.GetAllUsers(c)
				Expect(err).ToNot(HaveOccurred())
				Expect(rec.Code).To(Equal(http.StatusOK))
			})
		})
	})
})
