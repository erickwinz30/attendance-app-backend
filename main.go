// main.go
package main

import (
	"backend/database"
	"backend/handlers"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func init() {
	database.Init()
}

func main() {
	route := mux.NewRouter()

	route.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}).Methods("GET")

	// user routes - tidak perlu passing db lagi
	route.HandleFunc("/api/users", handlers.GetUsers()).Methods("GET")
	route.HandleFunc("/api/users", handlers.CreateUser()).Methods("POST")

	route.HandleFunc("/api/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		handlers.GetUser(userID)(w, r)
	}).Methods("GET")

	route.HandleFunc("/api/departments", handlers.GetDepartments()).Methods("GET")

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", route))
}
