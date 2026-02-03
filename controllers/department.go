package controllers

import (
	"backend/database"
	"backend/types"
)

func GetDepartments() ([]types.Department, error) {
	rows, err := database.DB.Query(`
		SELECT id, name 
		FROM departments 
		ORDER BY name ASC
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	departments := []types.Department{}
	for rows.Next() {
		var dept types.Department
		err := rows.Scan(&dept.ID, &dept.Name)
		if err != nil {
			return nil, err
		}
		departments = append(departments, dept)
	}

	return departments, nil
}
