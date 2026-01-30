// seed.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load .env dari parent directory
	envPath := filepath.Join("..", ".env")
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning: Gagal memuat .env dari %s: %v", envPath, err)
	}

	// Konfigurasi database dari environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Gagal membuka koneksi ke database:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Gagal ping database:", err)
	}

	fmt.Println("✅ Terhubung ke PostgreSQL")

	// Buat tabel jika belum ada
	createTables(db)

	// Seed data
	seedDepartments(db)
	seedUsers(db)

	fmt.Println("🌱 Seeding selesai!")
}

func createTables(db *sql.DB) {
	// Tabel departments
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS departments (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		);
	`)
	if err != nil {
		log.Fatal("Gagal membuat tabel departments:", err)
	}

	// Tabel users
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			phone TEXT,
			position TEXT,
			department_id INTEGER NOT NULL REFERENCES departments(id) ON DELETE RESTRICT,
			status TEXT NOT NULL CHECK (status IN ('active', 'inactive')),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatal("Gagal membuat tabel users:", err)
	}

	fmt.Println("✅ Tabel departments dan users siap")
}

func seedDepartments(db *sql.DB) {
	// Daftar departemen unik dari data kamu
	departments := []string{"IT", "Product", "Design", "HR", "Marketing"}

	for _, name := range departments {
		_, err := db.Exec(`
			INSERT INTO departments (name)
			VALUES ($1)
			ON CONFLICT (name) DO NOTHING;
		`, name)
		if err != nil {
			log.Printf("Gagal menyisipkan departemen %s: %v", name, err)
		}
	}
	fmt.Println("✅ Departemen disisipkan")
}

func seedUsers(db *sql.DB) {
	// Data pengguna (disesuaikan dengan struktur tabel)
	users := []struct {
		Name       string
		Email      string
		Phone      string
		Position   string
		Department string
		Status     string
	}{
		{"Ahmad Fauzi", "ahmad.fauzi@company.com", "+62 812-3456-7890", "Software Engineer", "IT", "active"},
		{"Siti Nurhaliza", "siti.nurhaliza@company.com", "+62 813-4567-8901", "Product Manager", "Product", "active"},
		{"Budi Santoso", "budi.santoso@company.com", "+62 814-5678-9012", "UI/UX Designer", "Design", "active"},
		{"Rina Wijaya", "rina.wijaya@company.com", "+62 815-6789-0123", "HR Manager", "HR", "inactive"},
		{"Dewi Lestari", "dewi.lestari@company.com", "+62 816-7890-1234", "Marketing Specialist", "Marketing", "active"},
		{"Andi Pratama", "andi.pratama@company.com", "+62 817-8901-2345", "Backend Developer", "IT", "active"},
	}

	for _, u := range users {
		// Dapatkan department_id dari nama departemen
		var deptID int
		err := db.QueryRow("SELECT id FROM departments WHERE name = $1", u.Department).Scan(&deptID)
		if err != nil {
			log.Printf("Departemen tidak ditemukan: %s", u.Department)
			continue
		}

		// Sisipkan user (hindari duplikat email)
		_, err = db.Exec(`
			INSERT INTO users (name, email, phone, position, department_id, status)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (email) DO NOTHING;
		`, u.Name, u.Email, u.Phone, u.Position, deptID, u.Status)

		if err != nil {
			log.Printf("Gagal menyisipkan user %s: %v", u.Name, err)
		}
	}
	fmt.Println("✅ Pengguna disisipkan")
}
