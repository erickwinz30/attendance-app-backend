package types

import "time"

type AttendanceToken struct {
	UserID    int       `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiredAt time.Time `json:"expired_at" db:"expired_at"`
	IsUsed    bool      `json:"is_used" db:"is_used"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type UserReceivedAttendanceToken struct {
	UserID    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expired_at"`
}

type CheckAttendanceToken struct {
	UserID int    `json:"user_id" db:"user_id"`
	Token  string `json:"token" db:"token"`
}

type CheckAttendanceTokenResponse struct {
	Valid      bool       `json:"valid"`
	Is_Used    *bool      `json:"is_used"`
	Expired_At *time.Time `json:"expired_at"`
	Message    string     `json:"message,omitempty"`
}

type SubmitAttendanceResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	UserID  int    `json:"user_id"`
}
