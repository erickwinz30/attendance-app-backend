package controllers

import (
	"backend/database"
	"backend/types"
	"database/sql"
	"log"
)

func CheckAuthentication(userID int) (types.AuthCheckResponse, error) {
	var hrdID int
	var adminID int
	var tempUser types.CheckUserTemp

	err := database.DB.QueryRow(`
		SELECT id 
		FROM departments 
		WHERE name = 'HR'
	`).Scan(&hrdID)

	err = database.DB.QueryRow(`
		SELECT id 
		FROM departments 
		WHERE name = 'Administrator'
	`).Scan(&adminID)

	if err != nil {
		log.Printf("Error fetching HRD department: %v", err)
		return types.AuthCheckResponse{Authenticated: false}, err
	}

	err = database.DB.QueryRow(`
		SELECT id, name, email, department_id
		FROM users
		WHERE id = $1
	`, userID).Scan(&tempUser.ID, &tempUser.Name, &tempUser.Email, &tempUser.DepartmentID)

	if err == sql.ErrNoRows {
		log.Printf("User not found with ID: %d", userID)
		return types.AuthCheckResponse{Authenticated: false}, nil
	}

	if err != nil {
		log.Printf("Error fetching user with ID %d: %v", userID, err)
		return types.AuthCheckResponse{Authenticated: false}, err
	}

	//check if user is from hrd department
	var role string

	if tempUser.DepartmentID != 0 && tempUser.DepartmentID == adminID {
		role = "Admin"
	} else if tempUser.DepartmentID != 0 && tempUser.DepartmentID == hrdID {
		role = "HR"
	} else {
		role = "Employee"
	}

	log.Printf("User authenticated successfully: ID=%d, Name=%s, Role=%s", tempUser.ID, tempUser.Name, role)

	userAuthInfo := types.UserAuthInfo{
		ID:    tempUser.ID,
		Name:  tempUser.Name,
		Email: tempUser.Email,
		Role:  role,
	}

	// check if user is attend today
	rows, err := database.DB.Query(`
    SELECT user_id, is_used 
    FROM attendance_tokens
    WHERE user_id = $1 AND DATE(created_at) = CURRENT_DATE
`, userID)

	if err != nil {
		log.Printf("Error fetching attendance for user %d: %v", userID, err)
		return types.AuthCheckResponse{Authenticated: false}, err
	}
	defer rows.Close()

	// Check if any token has been used today
	isAttend := false
	for rows.Next() {
		var userIDFromDB int
		var isUsed bool

		if err := rows.Scan(&userIDFromDB, &isUsed); err != nil {
			log.Printf("Error scanning attendance row: %v", err)
			continue
		}

		if isUsed {
			isAttend = true
			break
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating attendance rows: %v", err)
	}

	return types.AuthCheckResponse{
		Authenticated: true,
		User:          &userAuthInfo,
		IsAttended:    isAttend,
	}, nil
}
