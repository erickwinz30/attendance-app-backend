package handlers

import (
	"backend/controllers"
	"backend/types"
	"backend/utils"
	"encoding/json"
	"fmt"
	"net/http"
)

func GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := controllers.GetAllUsers()
		if err != nil {
			http.Error(w, "Gagal mengambil data pengguna", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func GetUser(userID int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := controllers.GetUser(userID)
		if err != nil {
			http.Error(w, "Gagal mengambil data pengguna", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}

func CreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUser types.CreateUserRequest
		// gunakan json decoder untuk parsing body
		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			http.Error(w, fmt.Sprintf("Gagal memproses data JSON: %v", err), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		fmt.Println("Parsed form (handlers):", newUser)

		// Validasi department_id
		if newUser.DepartmentID == 0 {
			http.Error(w, "department_id is required", http.StatusBadRequest)
			return
		}

		// Validasi field lainnya jika perlu
		if newUser.Name == "" || newUser.Email == "" {
			http.Error(w, "name and email are required", http.StatusBadRequest)
			return
		}

		// Panggil controller untuk membuat user baru
		result, err := controllers.CreateUser(newUser)
		if err != nil {
			http.Error(w, fmt.Sprintf("Gagal membuat pengguna baru: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Println("User created successfully (handlers):", result)

		w.WriteHeader(http.StatusCreated)
		// w.Write([]byte("Pengguna baru berhasil dibuat"))
		// json.NewEncoder(w).Encode(result)

	}
}

func SearchUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Implementasi pencarian pengguna berdasarkan query parameter

		query := r.URL.Query().Get("q")
		users, err := controllers.SearchUsers(query)
		if err != nil {
			http.Error(w, "Gagal mencari pengguna", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func EditUser(userID int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var editRequest types.EditUserRequest

		// Parse JSON dari request body
		err := json.NewDecoder(r.Body).Decode(&editRequest)
		if err != nil {
			http.Error(w, fmt.Sprintf("Gagal memproses data JSON: %v", err), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Validasi data kosong
		if editRequest.Name == "" && editRequest.Email == "" && editRequest.Phone == "" &&
			editRequest.Position == "" && editRequest.DepartmentID == 0 && editRequest.Status == "" {
			http.Error(w, "Minimal ada satu field yang harus diisi untuk update", http.StatusBadRequest)
			return
		}

		// Jika name kosong, return error
		if editRequest.Name == "" {
			http.Error(w, "field name tidak boleh kosong", http.StatusBadRequest)
			return
		}

		// Jika email kosong, return error
		if editRequest.Email == "" {
			http.Error(w, "field email tidak boleh kosong", http.StatusBadRequest)
			return
		}

		// Jika phone kosong, return error
		if editRequest.Phone == "" {
			http.Error(w, "field phone tidak boleh kosong", http.StatusBadRequest)
			return
		}

		// Jika position kosong, return error
		if editRequest.Position == "" {
			http.Error(w, "field position tidak boleh kosong", http.StatusBadRequest)
			return
		}

		// Jika department_id kosong/0, return error
		if editRequest.DepartmentID == 0 {
			http.Error(w, "field department_id tidak boleh kosong", http.StatusBadRequest)
			return
		}

		// Jika status kosong, return error
		if editRequest.Status == "" {
			http.Error(w, "field status tidak boleh kosong", http.StatusBadRequest)
			return
		}

		fmt.Println("Edit User Request (handlers):", editRequest)

		// Panggil controller untuk edit user
		result, err := controllers.EditUser(userID, editRequest)
		if err != nil {
			http.Error(w, fmt.Sprintf("Gagal mengedit pengguna: %v", err), http.StatusInternalServerError)
			return
		}

		utils.LogEditUser(userID, result)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}
