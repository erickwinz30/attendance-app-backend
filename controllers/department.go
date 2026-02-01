package controllers

import "backend/database"

type Department struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetDepartments() ([]Department, error) {
	rows, err := database.DB.Query(`
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
