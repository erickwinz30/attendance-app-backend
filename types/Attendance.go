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

type TodayAttendance struct {
	UserID         int       `json:"user_id"`
	UserName       string    `json:"user_name"`
	UserEmail      string    `json:"user_email"`
	DepartmentName string    `json:"department_name"`
	Position       string    `json:"position"`
	CheckInTime    time.Time `json:"check_in_time"`
	Token          string    `json:"token"`
	IsUsed         bool      `json:"is_used"`
	Status         string    `json:"status"` // "on-time" or "late"
}

type AbsentUser struct {
	UserID         int    `json:"user_id"`
	UserName       string `json:"user_name"`
	UserEmail      string `json:"user_email"`
	DepartmentName string `json:"department_name"`
	Position       string `json:"position"`
}

type TodayAttendanceListResponse struct {
	Date        string            `json:"date"`
	TotalAttend int               `json:"total_attend"`
	TotalLate   int               `json:"total_late"`
	TotalAbsent int               `json:"total_absent"`
	Attendances []TodayAttendance `json:"attendances"`
	AbsentUsers []AbsentUser      `json:"absent_users"`
}

type MonthlyAttendanceListResponse struct {
	Month       string            `json:"month"`
	Year        string            `json:"year"`
	TotalAttend int               `json:"total_attend"`
	TotalLate   int               `json:"total_late"`
	TotalAbsent int               `json:"total_absent"`
	Attendances []TodayAttendance `json:"attendances"`
	AbsentUsers []AbsentUser      `json:"absent_users"`
}
