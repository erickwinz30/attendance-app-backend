package controllers

import (
	"backend/database"
	"backend/types"
	"database/sql"
	"fmt"
)

func GetAllUsers() ([]types.User, error) {
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

	var users []types.User
	for rows.Next() {
		var user types.User
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

func GetUser(userID int) (*types.User, error) {
	var user types.User
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

func CreateUser(req types.CreateUserRequest) (types.CreateUserRequest, error) {
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

func SearchUsers(query string) ([]types.User, error) {
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

	var users []types.User
	for rows.Next() {
		var user types.User
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
