package controllers

import "database/sql"

type Department struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetDepartments(db *sql.DB) ([]Department, error) {
	rows, err := db.Query(`
		SELECT id, name 
		FROM departments 
		ORDER BY name ASC
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	departments := []Department{}
	for rows.Next() {
		var dept Department
		err := rows.Scan(&dept.ID, &dept.Name)
		if err != nil {
			return nil, err
		}
		departments = append(departments, dept)
	}

	return departments, nil
}
