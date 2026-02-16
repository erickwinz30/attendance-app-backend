package handlers

import (
	"backend/controllers"
	"encoding/json"
	"log"
	"net/http"
)

func GetWorkHours() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		workHours, err := controllers.GetWorkHours()
		if err != nil {
			log.Printf("Error getting work hours: %v", err)
			http.Error(w, "Failed to get work hours", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(workHours); err != nil {
			log.Printf("JSON encoding error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}
