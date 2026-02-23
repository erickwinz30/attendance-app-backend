# Attendance App - Backend

![Go](https://img.shields.io/badge/Go-1.25.5-00ADD8?style=flat&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-336791?style=flat&logo=postgresql)

Backend API untuk sistem manajemen absensi dengan autentikasi berbasis session dan token-based attendance.

## 🚀 Quick Start

**Prasyarat:** Go >= 1.20, PostgreSQL >= 13

```bash
# Clone & install dependencies
git clone <repository-url>
cd attendance-app/backend
go mod download

# Setup database
psql -U postgres -c "CREATE DATABASE attendance_db;"

# Buat file .env (lihat bagian Konfigurasi)
cp .env.example .env

# Jalankan server
go run main.go
```

Server berjalan di `http://localhost:8080`

## ⚙️ Konfigurasi

Buat file `.env`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=attendance_db
SESSION_SECRET=your_secret_key_here
PORT=8080
```

## 📦 Tech Stack

- **Go 1.25.5** dengan Gorilla Mux & Sessions
- **PostgreSQL** - Database utama
- **golang.org/x/crypto** - Password hashing
- **Rate Limiting** - 60 req/menit
- **CORS** enabled

## 🔌 API Endpoints

### Public Routes

- `GET /api/health` - Health check

### Authentication

- `POST /api/login` - Login dengan email & password
- `POST /api/logout` - Logout

### Protected Routes (Perlu Login)

- `GET /api/auth/check` - Cek status autentikasi
- `GET /api/attendance/token` - Generate token absensi
- `POST /api/attendance/token/check` - Validasi token
- `POST /api/attendance/submit` - Submit absensi
- `GET /api/work-hours` - Get jam kerja

### HR Only Routes

- `GET /api/users` - List semua user
- `GET /api/users/search` - Search user
- `POST /api/users` - Buat user baru
- `GET /api/users/{id}` - Detail user
- `PUT /api/users/{id}` - Update user
- `GET /api/departments` - List departemen
- `GET /api/attendance/today` - Absensi hari ini
- `GET /api/attendance/monthly` - Laporan bulanan
- `GET /api/attendance/employee/monthly` - Laporan karyawan per bulan

## 📁 Struktur Project

```
backend/
├── controllers/    # Business logic
├── database/       # Database connection
├── handlers/       # HTTP handlers & middleware
├── middleware/     # CORS & rate limiting
├── types/          # Data models
├── utils/          # Logger utilities
├── seed/           # Database seeding
└── main.go         # Entry point
```

## 🗄️ Database Tables

- **users** - Data karyawan (id, name, email, password, role, department_id)
- **departments** - Departemen & divisi
- **attendance** - Record absensi (check-in/out)
- **work_hours** - Jam kerja operasional
- **attendance_tokens** - Token sementara untuk submit absensi

## 🔐 Security Features

- Session-based authentication dengan gorilla/sessions
- Password hashing dengan bcrypt
- Rate limiting (60 requests/menit)
- Role-based access control (employee/hr)
- Protected routes dengan middleware
