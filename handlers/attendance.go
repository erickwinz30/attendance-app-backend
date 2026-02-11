package handlers

import (
	"backend/controllers"
	"backend/types"
	"encoding/json"
	"net/http"
)

func GenerateToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// cek dulu session dan user_id
		session, err := store.Get(r, "attendance-session")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, ok := session.Values["user_id"]
		if !ok || userID == nil {
			http.Error(w, "Unauthorized - Please login first", http.StatusUnauthorized)
			return
		}

		user, err := controllers.CheckAuthentication(userID.(int))
		if err != nil || !user.Authenticated {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// generate token untuk user tersebut
		attendanceToken, err := controllers.GenerateUserAttendanceToken(userID.(int))

		if err != nil {
			http.Error(w, "Failed to generate attendance token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(attendanceToken); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func CheckAttendanceToken() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var checkReq types.CheckAttendanceToken

		if err := json.NewDecoder(r.Body).Decode(&checkReq); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		checkResp, err := controllers.CheckAttendanceToken(checkReq)

		if err != nil {
			http.Error(w, "Failed to check attendance token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(checkResp); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
