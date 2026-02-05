package controllers

import (
	"backend/database"
	"backend/types"
	"database/sql"
	"log"
)

func CheckAuthentication(userID int) (types.AuthCheckResponse, error) {
	var hrdID int
	var tempUser types.CheckUserTemp

	err := database.DB.QueryRow(`
		SELECT id 
		FROM departments 
		WHERE name = 'HR'
	`).Scan(&hrdID)

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
	isHRD := false
	if tempUser.DepartmentID != 0 && tempUser.DepartmentID == hrdID {
		isHRD = true
	}

	log.Printf("User authenticated successfully: ID=%d, Name=%s, IsHRD=%t", tempUser.ID, tempUser.Name, isHRD)

	userAuthInfo := types.UserAuthInfo{
		ID:    tempUser.ID,
		Name:  tempUser.Name,
		Email: tempUser.Email,
		IsHRD: isHRD,
	}

	return types.AuthCheckResponse{
		Authenticated: true,
		User:          &userAuthInfo,
	}, nil
}
