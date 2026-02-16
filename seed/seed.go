// seed.go
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"

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

	// Migrate Fresh: Drop tables, create, dan seed ulang
	dropTables(db)
	createTables(db)

	// Seed data
	seedDepartments(db)
	seedUsers(db)
	seedWorkHours(db)

	fmt.Println("🌱 Migrate Fresh & Seeding selesai!")
}

func dropTables(db *sql.DB) {
	fmt.Println("🗑️  Menghapus tabel yang ada...")

	// Drop tables dalam urutan terbalik (karena foreign key constraints)
	tables := []string{"attendance_tokens", "users", "departments", "work_hours"}
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table))
		if err != nil {
			log.Printf("Gagal menghapus tabel %s: %v", table, err)
		}
	}

	fmt.Println("✅ Semua tabel berhasil dihapus")
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
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatal("Gagal membuat tabel users:", err)
	}

	// Tabel attendance_tokens
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS attendance_tokens (
			id SERIAL PRIMARY KEY,
			user_id INTEGER NOT NULL REFERENCES users(id),
			token TEXT UNIQUE NOT NULL,
			expired_at TIMESTAMP NOT NULL,
			is_used BOOLEAN DEFAULT false,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatal("Gagal membuat tabel attendance_tokens:", err)
	}

	// Tabel work_hours
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS work_hours (
			id SERIAL PRIMARY KEY,
			work_start_time TIME NOT NULL,
			work_end_time TIME NOT NULL,
			tolerance_time TIME NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		log.Fatal("Gagal membuat tabel work_hours:", err)
	}

	fmt.Println("✅ Semua tabel siap (departments, users, attendance_tokens, work_hours)")
}

func seedDepartments(db *sql.DB) {
	// Daftar departemen unik dari data kamu
	departments := []string{"IT", "Product", "Design", "HR", "Marketing", "Administrator"}

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
		Password   string
	}{
		{"Admin System", "admin@company.com", "+62 811-0000-0000", "System Administrator", "Administrator", "active", "admin123"},
		{"Ahmad Fauzi", "ahmad.fauzi@company.com", "+62 812-3456-7890", "Software Engineer", "IT", "active", "password123"},
		{"Siti Nurhaliza", "siti.nurhaliza@company.com", "+62 813-4567-8901", "Product Manager", "Product", "active", "password123"},
		{"Budi Santoso", "budi.santoso@company.com", "+62 814-5678-9012", "UI/UX Designer", "Design", "active", "password123"},
		{"Rina Wijaya", "rina.wijaya@company.com", "+62 815-6789-0123", "HR Manager", "HR", "inactive", "password123"},
		{"Dewi Lestari", "dewi.lestari@company.com", "+62 816-7890-1234", "Marketing Specialist", "Marketing", "active", "password123"},
		{"Andi Pratama", "andi.pratama@company.com", "+62 817-8901-2345", "Backend Developer", "IT", "active", "password123"},
	}

	for _, u := range users {
		var deptID int
		err := db.QueryRow("SELECT id FROM departments WHERE name = $1", u.Department).Scan(&deptID)
		if err != nil {
			log.Printf("Departemen tidak ditemukan: %s", u.Department)
			continue
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Gagal hash password untuk %s: %v", u.Name, err)
			continue
		}

		_, err = db.Exec(`
		INSERT INTO users (name, email, phone, position, department_id, status, password_hash)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (email) DO NOTHING;
	`, u.Name, u.Email, u.Phone, u.Position, deptID, u.Status, string(hashedPassword))

		if err != nil {
			log.Printf("Gagal menyisipkan user %s: %v", u.Name, err)
		}
	}
	fmt.Println("✅ Pengguna disisipkan")
}

func seedWorkHours(db *sql.DB) {
	// Set jam kerja global: 08:00 - 17:00 dengan toleransi sampai 08:15
	_, err := db.Exec(`
		INSERT INTO work_hours (work_start_time, work_end_time, tolerance_time)
		VALUES ('08:00:00', '17:00:00', '08:15:00');
	`)
	if err != nil {
		log.Printf("Gagal menyisipkan work_hours: %v", err)
		return
	}
	fmt.Println("✅ Work hours disisipkan (08:00 - 17:00, toleransi sampai 08:15)")
}
