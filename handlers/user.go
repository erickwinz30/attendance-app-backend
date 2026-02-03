package handlers

import (
	"backend/controllers"
	"backend/types"
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
