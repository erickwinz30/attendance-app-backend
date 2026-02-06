package handlers

import (
	"backend/controllers"
	"encoding/json"
	"log"
	"net/http"
)

// RequireAuth adalah middleware untuk memastikan user sudah login
// Middleware ini mengecek apakah ada session yang valid
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ambil session
		session, err := store.Get(r, "attendance-session")
		if err != nil {
			log.Printf("Session error: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Cek apakah user_id ada di session
		userID, ok := session.Values["user_id"]
		if !ok || userID == nil {
			http.Error(w, "Unauthorized - Please login first", http.StatusUnauthorized)
			return
		}

		// Jika sudah login, lanjutkan ke handler berikutnya
		next.ServeHTTP(w, r)
	})
}

func CheckAuthentication() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := store.Get(r, "attendance-session")
		if err != nil {
			log.Printf("Session error: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, ok := session.Values["user_id"]
		if !ok || userID == nil {
			http.Error(w, "Unauthorized - Please login first", http.StatusUnauthorized)
			return
		}

		users, err := controllers.CheckAuthentication(userID.(int))
		if err != nil {
			log.Printf("Authentication check error: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(users); err != nil {
			log.Printf("JSON encoding error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func RequireHR(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ambil session
		session, err := store.Get(r, "attendance-session")
		if err != nil {
			log.Printf("Session error: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// cek apakah ada session user_id
		userID, ok := session.Values["user_id"]
		if !ok || userID == nil {
			http.Error(w, "Unauthorized - Please login first", http.StatusUnauthorized)
			return
		}

		// Cek apakah user adalah HR
		authResponse, err := controllers.CheckAuthentication(userID.(int))
		if err != nil {
			log.Printf("Authentication check error: %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !authResponse.Authenticated || authResponse.User == nil || !authResponse.User.IsHRD {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(map[string]string{"error": "Forbidden - HR access only"})
			return
		}

		// jika user adalah dari department HR
		next.ServeHTTP(w, r)
	})
}
