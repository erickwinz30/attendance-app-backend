package controllers

import (
	"database/sql"
	"time"
)

type User struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	Position       string    `json:"position"`
	DepartmentID   int       `json:"department_id"`
	DepartmentName string    `json:"department_name"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

// type Department struct {
// 	ID   int    `json:"id"`
// 	Name string `json:"name"`
// }

func GetAllUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query(`
        SELECT 
            u.id, u.name, u.email, u.phone, u.position, 
            u.department_id, d.name as department_name,
            u.status, u.created_at 
        FROM users u
        INNER JOIN departments d ON u.department_id = d.id
        ORDER BY u.created_at DESC
    `)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Position, &user.DepartmentID, &user.DepartmentName, &user.Status, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
