package controllers

import (
	"backend/database"
	"backend/types"
	"log"
)

func GetWorkHours() (types.WorkHours, error) {
	var workHours types.WorkHours

	err := database.DB.QueryRow(`
		SELECT id, work_start_time, work_end_time, tolerance_time, created_at, updated_at
		FROM work_hours
		ORDER BY id DESC
		LIMIT 1
	`).Scan(
		&workHours.ID,
		&workHours.WorkStartTime,
		&workHours.WorkEndTime,
		&workHours.ToleranceTime,
		&workHours.CreatedAt,
		&workHours.UpdatedAt,
	)

	if err != nil {
		log.Printf("Error fetching work hours: %v", err)
		return types.WorkHours{}, err
	}

	return workHours, nil
}
