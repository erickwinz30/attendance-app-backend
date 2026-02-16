package types

import "time"

type WorkHours struct {
	ID            int       `json:"id" db:"id"`
	WorkStartTime string    `json:"work_start_time" db:"work_start_time"`
	WorkEndTime   string    `json:"work_end_time" db:"work_end_time"`
	ToleranceTime string    `json:"tolerance_time" db:"tolerance_time"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
