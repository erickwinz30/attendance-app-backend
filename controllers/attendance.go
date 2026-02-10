package controllers

import (
	"backend/database"
	"backend/types"
	"crypto/rand"
	"encoding/hex"
	"fmt"
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
