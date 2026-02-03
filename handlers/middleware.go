package handlers

import (
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
