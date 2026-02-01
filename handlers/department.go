package handlers

import (
	"backend/controllers"
	"encoding/json"
	"net/http"
)

func GetDepartments() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		departments, err := controllers.GetDepartments()

		if err != nil {
			http.Error(w, "Gagal mengambil data departemen", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(departments)
	}
}
