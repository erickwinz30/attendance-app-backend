package types

import "time"

type User struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Position       string    `json:"position"`
	DepartmentID   int       `json:"department_id"`
	DepartmentName string    `json:"department_name"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

type CreateUserRequest struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Position     string `json:"position"`
	DepartmentID int    `json:"department_id"`
	Status       string `json:"status"`
}

type CheckUserTemp struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	DepartmentID int    `json:"department_id"`
}
