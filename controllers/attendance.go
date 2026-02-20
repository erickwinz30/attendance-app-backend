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

func GetTodayAttendance() (types.TodayAttendanceListResponse, error) {
	var attendances []types.TodayAttendance

	// Get tolerance time from work_hours
	var toleranceTime string
	err := database.DB.QueryRow(`
		SELECT tolerance_time
		FROM work_hours
		ORDER BY id DESC
		LIMIT 1
	`).Scan(&toleranceTime)

	if err != nil {
		log.Printf("Error fetching work hours: %v", err)
		return types.TodayAttendanceListResponse{}, err
	}

	rows, err := database.DB.Query(`
		SELECT 
			u.id,
			u.name,
			u.email,
			d.name as department_name,
			u.position,
			at.created_at,
			at.token,
			at.is_used
		FROM attendance_tokens at
		JOIN users u ON at.user_id = u.id
		JOIN departments d ON u.department_id = d.id
		WHERE DATE(at.created_at) = CURRENT_DATE AND at.is_used = true
		ORDER BY at.created_at ASC
	`)

	if err != nil {
		log.Printf("Error fetching today's attendance: %v", err)
		return types.TodayAttendanceListResponse{}, err
	}
	defer rows.Close()

	totalLate := 0

	for rows.Next() {
		var attendance types.TodayAttendance
		err := rows.Scan(
			&attendance.UserID,
			&attendance.UserName,
			&attendance.UserEmail,
			&attendance.DepartmentName,
			&attendance.Position,
			&attendance.CheckInTime,
			&attendance.Token,
			&attendance.IsUsed,
		)
		if err != nil {
			log.Printf("Error scanning attendance row: %v", err)
			continue
		}

		// Determine status (on-time or late)
		checkInTime := attendance.CheckInTime.Format("15:04:05")
		if checkInTime > toleranceTime {
			attendance.Status = "late"
			totalLate++
		} else {
			attendance.Status = "on-time"
		}

		attendances = append(attendances, attendance)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating attendance rows: %v", err)
		return types.TodayAttendanceListResponse{}, err
	}

	// Get absent users (active users who haven't attended today)
	var absentUsers []types.AbsentUser
	absentRows, err := database.DB.Query(`
		SELECT 
			u.id,
			u.name,
			u.email,
			d.name as department_name,
			u.position
		FROM users u
		JOIN departments d ON u.department_id = d.id
		WHERE u.status = 'active'
		  AND u.id NOT IN (
			  SELECT DISTINCT user_id 
			  FROM attendance_tokens 
			  WHERE DATE(created_at) = CURRENT_DATE AND is_used = true
		  )
		ORDER BY u.name ASC
	`)

	if err != nil {
		log.Printf("Error fetching absent users: %v", err)
		// Continue even if absent users query fails
	} else {
		defer absentRows.Close()

		for absentRows.Next() {
			var absentUser types.AbsentUser
			err := absentRows.Scan(
				&absentUser.UserID,
				&absentUser.UserName,
				&absentUser.UserEmail,
				&absentUser.DepartmentName,
				&absentUser.Position,
			)
			if err != nil {
				log.Printf("Error scanning absent user row: %v", err)
				continue
			}
			absentUsers = append(absentUsers, absentUser)
		}

		if err = absentRows.Err(); err != nil {
			log.Printf("Error iterating absent user rows: %v", err)
		}
	}

	// Format today's date
	today := time.Now().Format("2006-01-02")

	response := types.TodayAttendanceListResponse{
		Date:        today,
		TotalAttend: len(attendances),
		TotalLate:   totalLate,
		TotalAbsent: len(absentUsers),
		Attendances: attendances,
		AbsentUsers: absentUsers,
	}

	return response, nil
}

func GetMonthlyAttendance() (types.MonthlyAttendanceListResponse, error) {
	var attendances []types.TodayAttendance

	// Get tolerance time from work_hours
	var toleranceTime string
	err := database.DB.QueryRow(`
		SELECT tolerance_time
		FROM work_hours
		ORDER BY id DESC
		LIMIT 1
	`).Scan(&toleranceTime)

	if err != nil {
		log.Printf("Error fetching work hours: %v", err)
		return types.MonthlyAttendanceListResponse{}, err
	}

	rows, err := database.DB.Query(`
		SELECT 
			u.id,
			u.name,
			u.email,
			d.name as department_name,
			u.position,
			at.created_at,
			at.token,
			at.is_used
		FROM attendance_tokens at
		JOIN users u ON at.user_id = u.id
		JOIN departments d ON u.department_id = d.id
		WHERE EXTRACT(MONTH FROM at.created_at) = EXTRACT(MONTH FROM CURRENT_DATE)
		  AND EXTRACT(YEAR FROM at.created_at) = EXTRACT(YEAR FROM CURRENT_DATE)
		  AND at.is_used = true
		ORDER BY at.created_at ASC
	`)

	if err != nil {
		log.Printf("Error fetching monthly attendance: %v", err)
		return types.MonthlyAttendanceListResponse{}, err
	}
	defer rows.Close()

	totalLate := 0

	for rows.Next() {
		var attendance types.TodayAttendance
		err := rows.Scan(
			&attendance.UserID,
			&attendance.UserName,
			&attendance.UserEmail,
			&attendance.DepartmentName,
			&attendance.Position,
			&attendance.CheckInTime,
			&attendance.Token,
			&attendance.IsUsed,
		)
		if err != nil {
			log.Printf("Error scanning attendance row: %v", err)
			continue
		}

		// Determine status (on-time or late)
		checkInTime := attendance.CheckInTime.Format("15:04:05")
		if checkInTime > toleranceTime {
			attendance.Status = "late"
			totalLate++
		} else {
			attendance.Status = "on-time"
		}

		attendances = append(attendances, attendance)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating attendance rows: %v", err)
		return types.MonthlyAttendanceListResponse{}, err
	}

	// Get absent users (active users who haven't attended at all this month)
	var absentUsers []types.AbsentUser
	absentRows, err := database.DB.Query(`
		SELECT 
			u.id,
			u.name,
			u.email,
			d.name as department_name,
			u.position
		FROM users u
		JOIN departments d ON u.department_id = d.id
		WHERE u.status = 'active'
		  AND u.id NOT IN (
			  SELECT DISTINCT user_id 
			  FROM attendance_tokens 
			  WHERE EXTRACT(MONTH FROM created_at) = EXTRACT(MONTH FROM CURRENT_DATE)
				AND EXTRACT(YEAR FROM created_at) = EXTRACT(YEAR FROM CURRENT_DATE)
				AND is_used = true
		  )
		ORDER BY u.name ASC
	`)

	if err != nil {
		log.Printf("Error fetching absent users: %v", err)
		// Continue even if absent users query fails
	} else {
		defer absentRows.Close()

		for absentRows.Next() {
			var absentUser types.AbsentUser
			err := absentRows.Scan(
				&absentUser.UserID,
				&absentUser.UserName,
				&absentUser.UserEmail,
				&absentUser.DepartmentName,
				&absentUser.Position,
			)
			if err != nil {
				log.Printf("Error scanning absent user row: %v", err)
				continue
			}
			absentUsers = append(absentUsers, absentUser)
		}

		if err = absentRows.Err(); err != nil {
			log.Printf("Error iterating absent user rows: %v", err)
		}
	}

	// Format current month and year
	now := time.Now()
	month := now.Format("01")  // MM format
	year := now.Format("2006") // YYYY format

	response := types.MonthlyAttendanceListResponse{
		Month:       month,
		Year:        year,
		TotalAttend: len(attendances),
		TotalLate:   totalLate,
		TotalAbsent: len(absentUsers),
		Attendances: attendances,
		AbsentUsers: absentUsers,
	}

	return response, nil
}

func GetEmployeeMonthlyAttendance(userID int, month int, year int) (types.EmployeeMonthlyAttendanceResponse, error) {
	// Get tolerance time
	var toleranceTime string
	err := database.DB.QueryRow(`
		SELECT tolerance_time
		FROM work_hours
		ORDER BY id DESC
		LIMIT 1
	`).Scan(&toleranceTime)

	if err != nil {
		log.Printf("Error fetching work hours: %v", err)
		return types.EmployeeMonthlyAttendanceResponse{}, err
	}

	// Get attendance records for the month
	rows, err := database.DB.Query(`
		SELECT 
			DATE(at.created_at) as attendance_date,
			at.created_at::time as check_in_time,
			at.is_used
		FROM attendance_tokens at
		WHERE at.user_id = $1
		  AND EXTRACT(MONTH FROM at.created_at) = $2
		  AND EXTRACT(YEAR FROM at.created_at) = $3
		  AND at.is_used = true
		ORDER BY at.created_at ASC
	`, userID, month, year)

	if err != nil {
		log.Printf("Error fetching employee attendance: %v", err)
		return types.EmployeeMonthlyAttendanceResponse{}, err
	}
	defer rows.Close()

	var attendances []types.EmployeeAttendance
	totalLateMinutes := 0

	for rows.Next() {
		var date string
		var checkInTime string
		var isUsed bool

		err := rows.Scan(&date, &checkInTime, &isUsed)
		if err != nil {
			log.Printf("Error scanning attendance row: %v", err)
			continue
		}

		status := "on-time"
		if checkInTime > toleranceTime {
			status = "late"
			// Calculate late minutes
			lateMinutes := calculateLateMinutes(checkInTime, toleranceTime)
			totalLateMinutes += lateMinutes
		}

		attendance := types.EmployeeAttendance{
			Date:        date,
			CheckInTime: checkInTime,
			Status:      status,
		}
		attendances = append(attendances, attendance)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error iterating attendance rows: %v", err)
		return types.EmployeeMonthlyAttendanceResponse{}, err
	}

	// Calculate total working days in the month (weekdays)
	totalWorkingDays := calculateWorkingDaysInMonth(year, month)

	// Total present is the number of attendance records
	totalPresent := len(attendances)

	// Total absent is working days minus present days
	totalAbsent := totalWorkingDays - totalPresent

	// Convert total late minutes to HH:MM format
	totalLateHours := formatMinutesToHHMM(totalLateMinutes)

	response := types.EmployeeMonthlyAttendanceResponse{
		Month:          fmt.Sprintf("%02d", month),
		Year:           fmt.Sprintf("%d", year),
		TotalPresent:   totalPresent,
		TotalAbsent:    totalAbsent,
		TotalLateHours: totalLateHours,
		Attendances:    attendances,
	}

	return response, nil
}

func calculateLateMinutes(checkInTime, toleranceTime string) int {
	// Parse times
	checkIn, err1 := time.Parse("15:04:05", checkInTime)
	tolerance, err2 := time.Parse("15:04:05", toleranceTime)

	if err1 != nil || err2 != nil {
		return 0
	}

	if checkIn.After(tolerance) {
		diff := checkIn.Sub(tolerance)
		return int(diff.Minutes())
	}

	return 0
}

func calculateWorkingDaysInMonth(year, month int) int {
	// Create time for first day of month
	firstDay := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	// Get last day of month
	lastDay := firstDay.AddDate(0, 1, -1)

	workingDays := 0
	for d := firstDay; d.Before(lastDay) || d.Equal(lastDay); d = d.AddDate(0, 0, 1) {
		// Skip weekends (Saturday = 6, Sunday = 0)
		if d.Weekday() != time.Saturday && d.Weekday() != time.Sunday {
			workingDays++
		}
	}

	return workingDays
}

func formatMinutesToHHMM(minutes int) string {
	hours := minutes / 60
	mins := minutes % 60
	return fmt.Sprintf("%02d:%02d", hours, mins)
}
