package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateToken() (string, error) {
	bytes := make([]byte, 8) // 8 byte = 16 karakter hex
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	fmt.Println("Generated token:", hex.EncodeToString(bytes))
	return hex.EncodeToString(bytes), nil
}

// func GenerateUserAttendanceToken()
