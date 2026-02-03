// main.go
package main

import (
	"backend/database"
	"backend/handlers"
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
	route := mux.NewRouter()

	// Route publik (tidak perlu login)
	route.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}).Methods("GET")

	// Route autentikasi
	route.HandleFunc("/api/login", handlers.LoginHandler).Methods("POST")
	route.HandleFunc("/api/logout", handlers.LogoutHandler).Methods("POST")

	// Route yang memerlukan autentikasi
	protected := route.PathPrefix("/api").Subrouter()
	protected.Use(handlers.RequireAuth) // Middleware untuk cek login

	// User routes (butuh login)
	protected.HandleFunc("/users/search", handlers.SearchUsers()).Methods("GET")
	protected.HandleFunc("/users", handlers.GetUsers()).Methods("GET")
	protected.HandleFunc("/users", handlers.CreateUser()).Methods("POST")
	protected.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		handlers.GetUser(userID)(w, r)
	}).Methods("GET")

	// Attendance & Department routes (butuh login)
	protected.HandleFunc("/attendance/token", handlers.GenerateToken()).Methods("GET")
	protected.HandleFunc("/departments", handlers.GetDepartments()).Methods("GET")

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", route))
}
