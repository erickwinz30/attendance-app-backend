package handlers

import (
	"backend/controllers"
	"database/sql"
	"encoding/json"
	"net/http"
)

func GetUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users, err := controllers.GetAllUsers(db)
		if err != nil {
			http.Error(w, "Gagal mengambil data pengguna", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}
