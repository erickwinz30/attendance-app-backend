package controllers

import (
	"backend/database"
	"database/sql"
	"fmt"
	"time"
)

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

func GetAllUsers() ([]User, error) {
	rows, err := database.DB.Query(`
        SELECT 
            u.id, u.name, u.email, u.phone, u.position, 
            u.department_id, d.name as department_name,
            u.status, u.created_at 
        FROM users u
        INNER JOIN departments d ON u.department_id = d.id
        ORDER BY u.created_at DESC
    `)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.Position,
			&user.DepartmentID,
			&user.DepartmentName,
			&user.Status,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func GetUser(userID int) (*User, error) {
	var user User
	err := database.DB.QueryRow(`
        SELECT 
            u.id, u.name, u.email, u.phone, u.position, 
            u.department_id, d.name as department_name,
            u.status, u.created_at 
        FROM users u
        INNER JOIN departments d ON u.department_id = d.id
        WHERE u.id = $1
    `, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Position,
		&user.DepartmentID,
		&user.DepartmentName,
		&user.Status,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateUser(req CreateUserRequest) (CreateUserRequest, error) {
	fmt.Println("Creating user with data (controller):", req)

	// cek apakah email sudah ada
	var existingUserID int
	err := database.DB.QueryRow(`
				SELECT id FROM users WHERE email = $1
		`, req.Email).Scan(&existingUserID)

	if err != nil && err != sql.ErrNoRows {
		// error unexpected (bukan "tidak ditemukan")
		return req, fmt.Errorf("gagal memeriksa email: %w", err)
	}
	if err == nil {
		return req, fmt.Errorf("email sudah terdaftar")
	}

	// cek apakah nomor telepon sudah ada
	err = database.DB.QueryRow(`
				SELECT id FROM users WHERE phone = $1
		`, req.Phone).Scan(&existingUserID)

	if err != nil && err != sql.ErrNoRows {
		// error unexpected (bukan "tidak ditemukan")
		return req, fmt.Errorf("gagal memeriksa nomor telepon: %w", err)
	}
	if err == nil {
		return req, fmt.Errorf("nomor telepon sudah terdaftar")
	}

	_, err = database.DB.Exec(`
		INSERT INTO users (name, email, phone, position, department_id, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`, req.Name, req.Email, req.Phone, req.Position, req.DepartmentID, req.Status)

	if err != nil {
		return req, fmt.Errorf("gagal insert user: %w", err)
	}

	return req, nil
}

func SearchUsers(query string) ([]User, error) {
	rows, err := database.DB.Query(`
        SELECT 
            u.id, u.name, u.email, u.phone, u.position, 
            u.department_id, d.name as department_name,
            u.status, u.created_at 
        FROM users u
        INNER JOIN departments d ON u.department_id = d.id
        WHERE u.name ILIKE $1 
          OR u.email ILIKE $1 
          OR u.phone ILIKE $1
          OR u.position ILIKE $1
        ORDER BY u.created_at DESC
    `, "%"+query+"%")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.Position,
			&user.DepartmentID,
			&user.DepartmentName,
			&user.Status,
			&user.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
