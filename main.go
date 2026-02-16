// main.go
package main

import (
	"backend/database"
	"backend/handlers"
	"backend/middleware"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	// Muat .env terlebih dahulu
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// Inisialisasi database
	database.Init()

	// Pastikan SESSION_SECRET ada
	if os.Getenv("SESSION_SECRET") == "" {
		log.Println("Warning: SESSION_SECRET tidak diset di .env")
	}
}

func main() {
	r := mux.NewRouter()
	r.Use(middleware.CORSMiddleware)

	// Route publik (tidak perlu login)
	r.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Route autentikasi
	r.HandleFunc("/api/login", handlers.LoginHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/logout", handlers.LogoutHandler).Methods("POST", "OPTIONS")

	// Route yang memerlukan autentikasi
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(handlers.RequireAuth) // Middleware untuk cek login

	// route untuk check session
	protected.HandleFunc("/auth/check", handlers.CheckAuthentication()).Methods("GET")
	// route untuk generate attendance token
	protected.HandleFunc("/attendance/token", handlers.GenerateToken()).Methods("GET")
	protected.HandleFunc("/attendance/token/check", handlers.CheckAttendanceToken()).Methods("POST", "OPTIONS")

	// route untuk proses absensi
	protected.HandleFunc("/attendance/submit", handlers.SubmitAttendance()).Methods("POST", "OPTIONS")

	// route untuk work hours
	protected.HandleFunc("/work-hours", handlers.GetWorkHours()).Methods("GET")

	// buat route khusus HR
	hrOnly := r.PathPrefix("/api").Subrouter()
	hrOnly.Use(handlers.RequireAuth)
	hrOnly.Use(handlers.RequireHR)

	// User routes (butuh login)
	hrOnly.HandleFunc("/users/search", handlers.SearchUsers()).Methods("GET")
	hrOnly.HandleFunc("/users", handlers.GetUsers()).Methods("GET")
	hrOnly.HandleFunc("/users", handlers.CreateUser()).Methods("POST")
	hrOnly.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		handlers.GetUser(userID)(w, r)
	}).Methods("GET")

	// Attendance & Department routes (butuh login)
	// hrOnly.HandleFunc("/attendance/token", handlers.GenerateToken()).Methods("GET")
	hrOnly.HandleFunc("/departments", handlers.GetDepartments()).Methods("GET")
	hrOnly.HandleFunc("/attendance/today", handlers.GetTodayAttendance()).Methods("GET")
	hrOnly.HandleFunc("/attendance/monthly", handlers.GetMonthlyAttendance()).Methods("GET")

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
