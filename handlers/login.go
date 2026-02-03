package handlers

import (
	"backend/database"
	"backend/types"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

// Store akan diinisialisasi setelah .env dimuat
var store *sessions.CookieStore

func init() {
	// Tunggu sampai SESSION_SECRET tersedia
	sessionSecret := os.Getenv("SESSION_SECRET")
	if sessionSecret == "" {
		log.Println("Warning: SESSION_SECRET tidak ditemukan, menggunakan default (tidak aman untuk production)")
		sessionSecret = "default-secret-key-change-this"
	}

	// Inisialisasi session store
	store = sessions.NewCookieStore([]byte(sessionSecret))

	// Konfigurasi session
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 24, // 24 jam
		HttpOnly: true,         // Mencegah JavaScript akses cookie
		Secure:   false,        // Set true jika pakai HTTPS
		SameSite: http.SameSiteLaxMode,
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginReq types.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validasi input
	if strings.TrimSpace(loginReq.Email) == "" {
		http.Error(w, "Email tidak boleh kosong", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(loginReq.Password) == "" {
		http.Error(w, "Password tidak boleh kosong", http.StatusBadRequest)
		return
	}
	if !strings.Contains(loginReq.Email, "@") {
		http.Error(w, "Format email tidak valid", http.StatusBadRequest)
		return
	}

	var userID int
	var hashedPassword string

	err := database.DB.QueryRow(`
		SELECT id, password_hash 
		FROM users 
		WHERE email = $1
	`, loginReq.Email).Scan(&userID, &hashedPassword)

	if err == sql.ErrNoRows {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Printf("DB error: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	// Verifikasi password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginReq.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Simpan sesi
	session, err := store.Get(r, "attendance-session")
	if err != nil {
		log.Printf("Session error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	session.Values["user_id"] = userID

	if err := session.Save(r, w); err != nil {
		log.Printf("Failed to save session: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}

// LogoutHandler menghapus session user
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	session, err := store.Get(r, "attendance-session")
	if err != nil {
		log.Printf("Session error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Hapus session dengan set MaxAge ke -1
	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		log.Printf("Failed to delete session: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Logout successful"})
}
