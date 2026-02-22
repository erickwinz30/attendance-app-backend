package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// LogToFile menulis log ke file seperti laravel.log
func LogToFile(message string) error {
	logDir := "logs"

	// Buat folder logs jika belum ada
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("gagal membuat folder logs: %w", err)
	}

	// Nama file log dengan format: app-YYYY-MM-DD.log
	logFileName := fmt.Sprintf("app-%s.log", time.Now().Format("2006-01-02"))
	logFilePath := filepath.Join(logDir, logFileName)

	// Buka atau buat file
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("gagal membuka file log: %w", err)
	}
	defer file.Close()

	// Format log dengan timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %s\n", timestamp, message)

	// Tulis ke file
	if _, err := file.WriteString(logEntry); err != nil {
		return fmt.Errorf("gagal menulis ke file log: %w", err)
	}

	return nil
}

// LogEditUser menulis log khusus untuk edit user
func LogEditUser(userID int, data map[string]interface{}) error {
	message := fmt.Sprintf("[EDIT USER] User ID: %d, Data: %+v", userID, data)
	return LogToFile(message)
}
