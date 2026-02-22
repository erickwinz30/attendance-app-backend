package controllers

import (
	"backend/database"
	"backend/types"
	"backend/utils"
	"database/sql"
	"fmt"
)

func GetAllUsers() ([]types.User, error) {
	rows, err := database.DB.Query(`
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

	var users []types.User
	for rows.Next() {
		var user types.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.Position,
			&user.DepartmentID,
			&user.DepartmentName,
			&user.Status,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func GetUser(userID int) (*types.User, error) {
	var user types.User
	err := database.DB.QueryRow(`
        SELECT 
            u.id, u.name, u.email, u.phone, u.position, 
            u.department_id, d.name as department_name,
            u.status, u.created_at 
        FROM users u
        INNER JOIN departments d ON u.department_id = d.id
        WHERE u.id = $1
    `, userID).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Phone,
		&user.Position,
		&user.DepartmentID,
		&user.DepartmentName,
		&user.Status,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateUser(req types.CreateUserRequest) (types.CreateUserRequest, error) {
	fmt.Println("Creating user with data (controller):", req)

	// cek apakah email sudah ada
	var existingUserID int
	err := database.DB.QueryRow(`
				SELECT id FROM users WHERE email = $1
		`, req.Email).Scan(&existingUserID)

	if err != nil && err != sql.ErrNoRows {
		// error unexpected (bukan "tidak ditemukan")
		return req, fmt.Errorf("gagal memeriksa email: %w", err)
	}
	if err == nil {
		return req, fmt.Errorf("email sudah terdaftar")
	}

	// cek apakah nomor telepon sudah ada
	err = database.DB.QueryRow(`
				SELECT id FROM users WHERE phone = $1
		`, req.Phone).Scan(&existingUserID)

	if err != nil && err != sql.ErrNoRows {
		// error unexpected (bukan "tidak ditemukan")
		return req, fmt.Errorf("gagal memeriksa nomor telepon: %w", err)
	}
	if err == nil {
		return req, fmt.Errorf("nomor telepon sudah terdaftar")
	}

	_, err = database.DB.Exec(`
		INSERT INTO users (name, email, phone, position, department_id, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`, req.Name, req.Email, req.Phone, req.Position, req.DepartmentID, req.Status)

	if err != nil {
		return req, fmt.Errorf("gagal insert user: %w", err)
	}

	return req, nil
}

func SearchUsers(query string) ([]types.User, error) {
	rows, err := database.DB.Query(`
        SELECT 
            u.id, u.name, u.email, u.phone, u.position, 
            u.department_id, d.name as department_name,
            u.status, u.created_at 
        FROM users u
        INNER JOIN departments d ON u.department_id = d.id
        WHERE u.name ILIKE $1 
          OR u.email ILIKE $1 
          OR u.phone ILIKE $1
          OR u.position ILIKE $1
        ORDER BY u.created_at DESC
    `, "%"+query+"%")

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []types.User
	for rows.Next() {
		var user types.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Phone,
			&user.Position,
			&user.DepartmentID,
			&user.DepartmentName,
			&user.Status,
			&user.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func EditUser(userID int, req types.EditUserRequest) (map[string]interface{}, error) {
	fmt.Println("Editing user with ID (controller):", userID)
	fmt.Println("Edit User data (controller):", req)

	// Buat map dari data yang ingin diedit
	editData := map[string]interface{}{
		"user_id":       userID,
		"name":          req.Name,
		"email":         req.Email,
		"phone":         req.Phone,
		"position":      req.Position,
		"department_id": req.DepartmentID,
		"status":        req.Status,
	}

	// Import utils untuk logging
	utils.LogEditUser(userID, editData)

	// ambil data user lama
	var oldData types.User
	err := database.DB.QueryRow(`
        SELECT id, name, email, phone, position, department_id, status
        FROM users
        WHERE id = $1
    `, userID).Scan(
		&oldData.ID,
		&oldData.Name,
		&oldData.Email,
		&oldData.Phone,
		&oldData.Position,
		&oldData.DepartmentID,
		&oldData.Status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.LogEditUser(userID, map[string]interface{}{
				"error": "user tidak ditemukan",
			})
			return nil, fmt.Errorf("user tidak ditemukan")
		}
		return nil, fmt.Errorf("gagal query user: %w", err)
	}

	// 2. Buat map untuk menyimpan field yang akan diupdate
	updateData := make(map[string]interface{})

	// 3. Cek setiap field, jika ada perubahan, tambahkan ke updateData
	if req.Name != "" && req.Name != oldData.Name {
		updateData["name"] = req.Name
	}
	if req.Email != "" && req.Email != oldData.Email {
		updateData["email"] = req.Email
	}
	if req.Phone != "" && req.Phone != oldData.Phone {
		updateData["phone"] = req.Phone
	}
	if req.Position != "" && req.Position != oldData.Position {
		updateData["position"] = req.Position
	}
	if req.DepartmentID != 0 && req.DepartmentID != oldData.DepartmentID {
		updateData["department_id"] = req.DepartmentID
	}
	if req.Status != "" && req.Status != oldData.Status {
		updateData["status"] = req.Status
	}

	// 4. Jika tidak ada perubahan, return pesan tidak ada perubahan
	if len(updateData) == 0 {
		utils.LogEditUser(userID, map[string]interface{}{
			"message": "tidak ada perubahan data",
		})
		return map[string]interface{}{
			"message": "tidak ada perubahan data",
			"user_id": userID,
		}, nil
	}

	// 5. Mulai transaction
	tx, err := database.DB.Begin()
	if err != nil {
		utils.LogEditUser(userID, map[string]interface{}{
			"error": fmt.Sprintf("gagal memulai transaction: %v", err),
		})
		return nil, fmt.Errorf("gagal memulai transaction: %w", err)
	}

	// 6. Build UPDATE query dinamis
	updateQuery := "UPDATE users SET "
	args := []interface{}{}
	argIndex := 1

	for field, value := range updateData {
		if argIndex > 1 {
			updateQuery += ", "
		}
		updateQuery += fmt.Sprintf("%s = $%d", field, argIndex)
		args = append(args, value)
		argIndex++
	}

	updateQuery += fmt.Sprintf(" WHERE id = $%d", argIndex)
	args = append(args, userID)

	// 7. Execute update
	result, err := tx.Exec(updateQuery, args...)
	if err != nil {
		tx.Rollback()
		utils.LogEditUser(userID, map[string]interface{}{
			"error": fmt.Sprintf("gagal update user: %v", err),
		})
		return nil, fmt.Errorf("gagal update user: %w", err)
	}

	// 8. Cek apakah ada row yang teraffected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("gagal cek rows affected: %w", err)
	}

	if rowsAffected == 0 {
		tx.Rollback()
		return nil, fmt.Errorf("user tidak ditemukan")
	}

	// 9. Commit transaction
	if err := tx.Commit(); err != nil {
		utils.LogEditUser(userID, map[string]interface{}{
			"error": fmt.Sprintf("gagal commit transaction: %v", err),
		})
		return nil, fmt.Errorf("gagal commit transaction: %w", err)
	}

	// 10. Log success dengan field yang berubah
	changedFields := make([]string, 0)
	for field := range updateData {
		changedFields = append(changedFields, field)
	}

	utils.LogEditUser(userID, map[string]interface{}{
		"message":        "user berhasil diupdate",
		"changed_fields": changedFields,
		"updated_data":   updateData,
	})

	return map[string]interface{}{
		"message":        "user berhasil diupdate",
		"user_id":        userID,
		"changed_fields": changedFields,
		"updated_data":   updateData,
	}, nil
}
