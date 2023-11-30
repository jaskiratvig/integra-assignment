package models

type User struct {
    ID         int    `json:"id"`
    UserName   string `json:"user_name"`
    FirstName  string `json:"first_name"`
    LastName   string `json:"last_name"`
    Email      string `json:"email"`
    UserStatus string `json:"user_status"`
    Department string `json:"department"`
}