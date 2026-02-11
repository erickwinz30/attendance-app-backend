package controllers

import (
	"backend/database"
	"backend/types"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"time"
)

func GenerateToken() string {
	bytes := make([]byte, 8) // 8 byte = 16 karakter hex
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}

	fmt.Println("Generated token:", hex.EncodeToString(bytes))
	return hex.EncodeToString(bytes)
}

func GenerateExpirationTime() time.Time {
	// Set token expiration to 5 minutes from now
	return time.Now().Add(5 * time.Minute)
}

func GenerateUserAttendanceToken(userID int) (types.AttendanceToken, error) {
	generatedToken := GenerateToken()
	expirationTime := GenerateExpirationTime()

	attendanceToken := types.AttendanceToken{
		UserID:    userID,
		Token:     generatedToken,
		ExpiredAt: expirationTime,
		IsUsed:    false,
		CreatedAt: time.Now(),
	}

	_, err := database.DB.Exec(`
		INSERT INTO attendance_tokens (user_id, token, expired_at, is_used, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, attendanceToken.UserID, attendanceToken.Token, attendanceToken.ExpiredAt, attendanceToken.IsUsed, attendanceToken.CreatedAt)

	if err != nil {
		return attendanceToken, fmt.Errorf("gagal insert attendance token: %w", err)
	}

	return attendanceToken, nil
}

func CheckAttendanceToken(data types.CheckAttendanceToken) (types.CheckAttendanceTokenResponse, error) {
	var expired_at time.Time
	var is_used bool

	err := database.DB.QueryRow(`
		SELECT expired_at, is_used
		FROM attendance_tokens
		WHERE user_id = $1 AND token = $2
	`, data.UserID, data.Token).Scan(&expired_at, &is_used)

	if err == sql.ErrNoRows {
		log.Printf("User not found with ID: %d", data.UserID)
		return types.CheckAttendanceTokenResponse{
			Valid:      false,
			Is_Used:    nil,
			Expired_At: nil,
			Message:    "Token not found",
		}, nil
	}

	if err != nil {
		log.Printf("Error fetching attendance token for user ID %d: %v", data.UserID, err)
		return types.CheckAttendanceTokenResponse{
			Valid:      false,
			Is_Used:    nil,
			Expired_At: nil,
			Message:    "Token not found 2",
		}, err
	}

	if is_used {
		return types.CheckAttendanceTokenResponse{
			Valid:      false,
			Is_Used:    &is_used,
			Expired_At: nil,
			Message:    "Token already used",
		}, nil
	}

	if time.Now().After(expired_at) {
		return types.CheckAttendanceTokenResponse{
			Valid:      false,
			Is_Used:    &is_used,
			Expired_At: &expired_at,
			Message:    "Token expired",
		}, nil
	}

	// jika valid
	return types.CheckAttendanceTokenResponse{
		Valid:      true,
		Is_Used:    &is_used,
		Expired_At: &expired_at,
		Message:    "Token is valid",
	}, nil
}
