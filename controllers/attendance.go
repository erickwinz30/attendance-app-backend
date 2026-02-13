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

func GenerateUserAttendanceToken(userID int) (types.UserReceivedAttendanceToken, error) {
	generatedToken := GenerateToken()
	expirationTime := GenerateExpirationTime()

	attendanceToken := types.AttendanceToken{
		UserID:    userID,
		Token:     generatedToken,
		ExpiredAt: expirationTime,
		IsUsed:    false,
		CreatedAt: time.Now(),
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return types.UserReceivedAttendanceToken{}, fmt.Errorf("gagal memulai transaction: %w", err)
	}

	// Defer rollback untuk memastikan rollback jika terjadi error
	defer func() {
		if err != nil {
			tx.Rollback()
			log.Printf("Transaction rolled back for user ID %d", userID)
		}
	}()

	_, err = tx.Exec(`
		INSERT INTO attendance_tokens (user_id, token, expired_at, is_used, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, attendanceToken.UserID, attendanceToken.Token, attendanceToken.ExpiredAt, attendanceToken.IsUsed, attendanceToken.CreatedAt)

	if err != nil {
		return types.UserReceivedAttendanceToken{}, fmt.Errorf("gagal insert attendance token: %w", err)
	}

	// Commit transaction jika semua berhasil
	err = tx.Commit()
	if err != nil {
		return types.UserReceivedAttendanceToken{}, fmt.Errorf("gagal commit transaction: %w", err)
	}

	// send generated token to user
	userReceivedToken := types.UserReceivedAttendanceToken{
		UserID:    attendanceToken.UserID,
		Token:     attendanceToken.Token,
		ExpiredAt: attendanceToken.ExpiredAt,
	}

	log.Printf("Generated attendance token for user ID %d: %s", userReceivedToken.UserID, userReceivedToken.Token)

	return userReceivedToken, nil
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

func SubmitAttendance(submitReq types.UserReceivedAttendanceToken) (types.SubmitAttendanceResponse, error) {
	// cek terlebih dahulu apakah token user expired dan apakah sudah terpakai
	var expired_at time.Time
	var is_used bool

	err := database.DB.QueryRow(`
		SELECT expired_at, is_used
		FROM attendance_tokens
		WHERE user_id = $1 AND token = $2
	`, submitReq.UserID, submitReq.Token).Scan(&expired_at, &is_used)

	if err != nil {
		log.Printf("Error fetching attendance token for user ID %d: %v", submitReq.UserID, err)
		return types.SubmitAttendanceResponse{}, err
	}

	if is_used {
		return types.SubmitAttendanceResponse{
			Success: false,
			Message: "Token already used",
			UserID:  submitReq.UserID,
		}, nil
	}

	if time.Now().After(expired_at) {
		return types.SubmitAttendanceResponse{
			Success: false,
			Message: "Token expired",
			UserID:  submitReq.UserID,
		}, nil
	}

	// jika valid, maka mulai update is_used menjadi true

	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction for submitting attendance: %v", err)
		return types.SubmitAttendanceResponse{}, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			log.Printf("Transaction rolled back for submitting attendance: %v", err)
		}
	}()

	_, err = tx.Exec(`
		UPDATE attendance_tokens
		SET is_used = true
		WHERE user_id = $1 AND token = $2
	`, submitReq.UserID, submitReq.Token)

	if err != nil {
		log.Printf("Failed to update attendance token as used for user ID %d: %v", submitReq.UserID, err)
		return types.SubmitAttendanceResponse{}, err
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("Failed to commit transaction for submitting attendance for user ID %d: %v", submitReq.UserID, err)
		return types.SubmitAttendanceResponse{}, err
	}

	log.Printf("Attendance submitted successfully for user ID %d with token %s", submitReq.UserID, submitReq.Token)
	return types.SubmitAttendanceResponse{
		Success: true,
		Message: fmt.Sprintf("User with ID %d attendance submitted successfully", submitReq.UserID),
		UserID:  submitReq.UserID,
	}, nil
}
