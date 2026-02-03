package handlers

import (
	"backend/controllers"
	"net/http"
)

func GenerateToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := controllers.GenerateToken()

		if err != nil {
			http.Error(w, "Gagal menghasilkan token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("Token: " + token))
	}
}
