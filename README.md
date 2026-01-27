# Attendance App - Backend

![Go](https://img.shields.io/badge/Go-1.25.5-00ADD8?style=flat&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Database-336791?style=flat&logo=postgresql)

Backend API server untuk sistem manajemen absensi berbasis QR Code yang dibangun dengan Go dan PostgreSQL.

## 📋 Daftar Isi

- [Fitur](#fitur)
- [Teknologi](#teknologi)
- [Prasyarat](#prasyarat)
- [Instalasi](#instalasi)
- [Konfigurasi](#konfigurasi)
- [Menjalankan Aplikasi](#menjalankan-aplikasi)
- [Struktur Database](#struktur-database)
- [API Endpoints](#api-endpoints)
- [Development](#development)

## ✨ Fitur

- 🔌 **REST API** - RESTful API untuk operasi CRUD
- 🗄️ **PostgreSQL Database** - Database relational untuk data persistence
- 🔐 **Environment Variables** - Konfigurasi dengan .env
- ⚡ **Lightweight** - Minimal dependencies, performance tinggi
- 🔄 **CORS Support** - Cross-Origin Resource Sharing enabled
- 📊 **Health Check** - Endpoint untuk monitoring

## 🛠 Teknologi

### Core

- **Go** 1.25.5 - Programming language
- **net/http** - Standard HTTP server
- **database/sql** - Database driver interface

### Database

- **PostgreSQL** - Relational database
- **lib/pq** 1.10.9 - PostgreSQL driver for Go

### Utilities

- **godotenv** 1.5.1 - Environment variable management

## 📦 Prasyarat

Pastikan Anda telah menginstall:

- **Go** >= 1.20
- **PostgreSQL** >= 13
- **Git**

## 🚀 Instalasi

1. **Clone repository**

   ```bash
   git clone <repository-url>
   cd attendance-app/backend
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Setup PostgreSQL Database**

   Buat database baru:

   ```bash
   psql -U postgres
   CREATE DATABASE attendance_db;
   ```

4. **Buat file .env**
   ```bash
   cp .env.example .env
   ```

## ⚙️ Konfigurasi

### File `.env`

Buat file `.env` di root folder backend:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password_here
DB_NAME=attendance_db

# Server Configuration
PORT=8080
```

### Konfigurasi Database

Update nilai environment variables sesuai dengan setup PostgreSQL Anda:

- `DB_HOST`: Host database (default: localhost)
- `DB_PORT`: Port PostgreSQL (default: 5432)
- `DB_USER`: Username PostgreSQL
- `DB_PASSWORD`: Password database
- `DB_NAME`: Nama database

## 🏃 Menjalankan Aplikasi

### Development Mode

```bash
go run db.go
```

Server akan berjalan di [http://localhost:8080](http://localhost:8080)

### Build untuk Production

```bash
go build -o attendance-server db.go
./attendance-server
```

## 🗄️ Struktur Database

### Table: `users`

Menyimpan data pengguna/karyawan

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    position VARCHAR(100),
    department VARCHAR(100),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Table: `attendance`

Menyimpan data absensi

```sql
CREATE TABLE attendance (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    date DATE NOT NULL,
    time TIME NOT NULL,
    status VARCHAR(20) NOT NULL, -- 'hadir', 'tidak-hadir'
    qr_code VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Table: `qr_codes`

Menyimpan QR code yang valid

```sql
CREATE TABLE qr_codes (
    id SERIAL PRIMARY KEY,
    code VARCHAR(255) UNIQUE NOT NULL,
    user_id INTEGER REFERENCES users(id),
    valid_from TIMESTAMP NOT NULL,
    valid_until TIMESTAMP NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🔌 API Endpoints

### Health Check

```http
GET /api/health
```

**Response:**

```
OK
```

### Users Endpoints (Planned)

#### Get All Users

```http
GET /api/users
```

#### Get User by ID

```http
GET /api/users/:id
```

#### Create User

```http
POST /api/users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "phone": "+62812345678",
  "position": "Software Engineer",
  "department": "IT"
}
```

#### Update User

```http
PUT /api/users/:id
Content-Type: application/json

{
  "name": "John Doe Updated",
  "position": "Senior Software Engineer"
}
```

#### Delete User

```http
DELETE /api/users/:id
```

### Attendance Endpoints (Planned)

#### Record Attendance (QR Scan)

```http
POST /api/attendance/scan
Content-Type: application/json

{
  "qr_code": "USER-12345",
  "timestamp": "2026-01-19T08:15:30Z"
}
```

#### Get Attendance History

```http
GET /api/attendance?date=2026-01-19&status=hadir
```

#### Get User Attendance

```http
GET /api/attendance/user/:user_id
```

### QR Code Endpoints (Planned)

#### Generate QR Code

```http
POST /api/qrcode/generate
Content-Type: application/json

{
  "user_id": 1,
  "valid_duration": "24h"
}
```

#### Verify QR Code

```http
POST /api/qrcode/verify
Content-Type: application/json

{
  "code": "USER-12345"
}
```

## 📁 Struktur File

```
backend/
├── db.go              # Main server file dengan database setup
├── go.mod             # Go module dependencies
├── go.sum             # Go module checksums
├── .env               # Environment variables (tidak di-commit)
├── .env.example       # Template environment variables
└── README.md          # Dokumentasi ini
```

## 🔧 Development

### Menambah Dependencies

```bash
go get <package-name>
go mod tidy
```

### Running Tests

```bash
go test ./...
```

### Code Formatting

```bash
go fmt ./...
```

### Linting

```bash
go vet ./...
```

## 🚀 Pengembangan Selanjutnya

- [ ] **API Endpoints** - Implementasi lengkap CRUD untuk Users, Attendance, QR Code
- [ ] **Authentication** - JWT-based authentication
- [ ] **Authorization** - Role-based access control (RBAC)
- [ ] **Middleware** - Logging, CORS, Rate limiting
- [ ] **QR Code Generation** - Generate QR dengan library Go
- [ ] **Validation** - Input validation dengan validator
- [ ] **Error Handling** - Standardized error responses
- [ ] **Pagination** - Pagination untuk list endpoints
- [ ] **Filtering & Sorting** - Query parameters untuk filter dan sort
- [ ] **WebSocket** - Real-time updates untuk attendance
- [ ] **Docker** - Containerization dengan Docker
- [ ] **Testing** - Unit tests dan integration tests
- [ ] **Documentation** - Swagger/OpenAPI documentation
- [ ] **Migrations** - Database migration tool
- [ ] **Seeding** - Database seeding untuk development

## 📝 Database Migration (Rencana)

Contoh migration script:

```sql
-- migrations/001_create_users_table.up.sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    position VARCHAR(100),
    department VARCHAR(100),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- migrations/001_create_users_table.down.sql
DROP TABLE IF EXISTS users;
```

## 🔐 Security Considerations

- ✅ Environment variables untuk credentials
- ⏳ Input validation (coming soon)
- ⏳ SQL injection prevention dengan prepared statements (coming soon)
- ⏳ Password hashing (coming soon)
- ⏳ HTTPS/TLS (production)
- ⏳ Rate limiting (coming soon)

## 🐛 Troubleshooting

### Database Connection Error

```
Error: failed to connect to database
```

**Solution:**

1. Pastikan PostgreSQL running: `sudo systemctl status postgresql`
2. Cek credentials di file `.env`
3. Pastikan database sudah dibuat
4. Test koneksi: `psql -U <username> -d <database_name>`

### Port Already in Use

```
Error: bind: address already in use
```

**Solution:**

```bash
# Cari process yang menggunakan port 8080
lsof -i :8080
# Kill process
kill -9 <PID>
```

## 📚 Resources

- [Go Documentation](https://golang.org/doc/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [lib/pq Driver](https://github.com/lib/pq)
- [godotenv](https://github.com/joho/godotenv)

## 🤝 Kontribusi

Silakan buat issue atau pull request untuk kontribusi atau perbaikan.

## 📄 License

MIT License - bebas digunakan untuk pembelajaran dan pengembangan.

---

**Dibuat dengan ❤️ menggunakan Go & PostgreSQL**
Backend untuk project absensi QR Code menggunakan bahasa pemrograman Go
