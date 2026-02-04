# 🏋️ Revoks Gym Backend Application

Backend API untuk aplikasi manajemen gym menggunakan **Go (Golang)** dengan clean architecture.

---

## 📋 Daftar Isi

1. [Tentang Project](#tentang-project)
2. [Tech Stack & Versi](#tech-stack--versi)
3. [Arsitektur Layer](#arsitektur-layer)
4. [Penjelasan Repository-Service Pattern](#penjelasan-repository-service-pattern)
5. [Struktur Project](#struktur-project)
6. [Penjelasan Docker Setup](#penjelasan-docker-setup)
7. [Cara Setup & Menjalankan](#cara-setup--menjalankan)
8. [API Endpoints](#api-endpoints)
9. [Contoh Penggunaan API](#contoh-penggunaan-api)

---

## 📖 Tentang Project

Project ini adalah backend REST API untuk sistem manajemen gym yang dibangun dengan prinsip **Clean Architecture**. Clean Architecture memisahkan kode menjadi beberapa layer yang independen sehingga:

- **Mudah di-maintain** - Setiap layer memiliki tanggung jawab sendiri
- **Mudah di-test** - Setiap layer dapat diuji secara terpisah
- **Mudah dikembangkan** - Menambah fitur baru tidak mengganggu fitur yang sudah ada
- **Loosely coupled** - Perubahan di satu layer tidak mempengaruhi layer lain

---

## 🛠️ Tech Stack & Versi

| Technology | Version | Deskripsi |
|------------|---------|-----------|
| **Go (Golang)** | 1.21 | Bahasa pemrograman utama |
| **Fiber** | v2.52.0 | Web framework (mirip Express.js untuk Node.js) |
| **GORM** | v1.25.5 | ORM (Object Relational Mapping) untuk database |
| **PostgreSQL** | 16-alpine | Database relational |
| **pgAdmin** | 4 (latest) | GUI untuk manajemen PostgreSQL |
| **Docker** | 24.x | Container platform |
| **Docker Compose** | 2.x | Multi-container orchestration |
| **Air** | v1.49.0 | Hot reload untuk development |

### Dependensi Go (go.mod)

```go
require (
    github.com/gofiber/fiber/v2 v2.52.0    // Web framework
    github.com/joho/godotenv v1.5.1        // Load .env file
    golang.org/x/crypto v0.17.0            // Password hashing (bcrypt)
    gorm.io/driver/postgres v1.5.4         // PostgreSQL driver untuk GORM
    gorm.io/gorm v1.25.5                   // ORM library
)
```

---

## 🏗️ Arsitektur Layer

```
┌─────────────────────────────────────────────────────────────┐
│                      HTTP Request                            │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     ROUTES (routes.go)                       │
│         Mendefinisikan endpoint dan method HTTP              │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   HANDLERS (user_handler.go)                 │
│    Menerima request, validasi input, memanggil service       │
│              Menggunakan DTO untuk request/response          │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   SERVICES (user_service.go)                 │
│         Business logic, validasi bisnis, transformasi data   │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                 REPOSITORY (user_repository.go)              │
│            Operasi database (CRUD), query builder            │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    MODELS (user.go)                          │
│              Struktur tabel database (Entity)                │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      DATABASE                                │
│                     (PostgreSQL)                             │
└─────────────────────────────────────────────────────────────┘
```

### Penjelasan Setiap Layer:

| Layer | File | Fungsi |
|-------|------|--------|
| **Routes** | `routes/routes.go` | Mendefinisikan URL endpoint dan method HTTP (GET, POST, PUT, DELETE) |
| **Handlers** | `handlers/user_handler.go` | Menerima HTTP request, parsing body/params, validasi input, mengembalikan response |
| **DTO** | `dto/user_dto.go` | Data Transfer Object - struktur data untuk request dan response (bukan entity database) |
| **Services** | `services/user_service.go` | Berisi business logic, validasi aturan bisnis, koordinasi antar repository |
| **Repository** | `repository/user_repository.go` | Abstraksi akses database, semua query SQL/GORM ada di sini |
| **Models** | `models/user.go` | Representasi tabel database dalam bentuk struct Go |
| **Config** | `config/database.go` | Konfigurasi koneksi database |

---

## 🔍 Penjelasan Repository-Service Pattern

### Apa itu Repository Pattern?

**Repository Pattern** adalah design pattern yang memisahkan logika akses data dari logika bisnis. Repository bertindak sebagai **"perantara"** antara aplikasi dan database.

```go
// Repository hanya bertugas CRUD ke database
type UserRepository interface {
    Create(user *models.User) error
    FindAll(page, perPage int) ([]models.User, int64, error)
    FindByID(id uint) (*models.User, error)
    FindByEmail(email string) (*models.User, error)
    Update(user *models.User) error
    Delete(id uint) error
}
```

**Kenapa menggunakan Repository Pattern?**

| Alasan | Penjelasan |
|--------|------------|
| **Separation of Concerns** | Logika database terpisah dari logika bisnis |
| **Testability** | Mudah di-mock untuk unit testing |
| **Maintainability** | Perubahan database hanya di satu tempat |
| **Flexibility** | Mudah ganti database (PostgreSQL → MySQL) tanpa ubah service |
| **Single Responsibility** | Setiap class punya satu tugas saja |

### Apa itu Service Layer?

**Service Layer** berisi **business logic** aplikasi. Service memanggil repository untuk operasi database dan menerapkan aturan bisnis.

```go
// Service bertugas menjalankan business logic
type UserService interface {
    CreateUser(req *dto.CreateUserRequest) (*dto.UserResponse, error)
    GetAllUsers(page, perPage int) (*dto.PaginatedResponse, error)
    GetUserByID(id uint) (*dto.UserResponse, error)
    UpdateUser(id uint, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
    DeleteUser(id uint) error
    ChangePassword(id uint, req *dto.ChangePasswordRequest) error
}
```

**Kenapa menggunakan Service Layer?**

| Alasan | Penjelasan |
|--------|------------|
| **Business Logic Centralized** | Semua aturan bisnis di satu tempat |
| **Reusability** | Service bisa dipanggil dari berbagai handler |
| **Transaction Management** | Mengelola transaksi database yang kompleks |
| **Validation** | Validasi aturan bisnis sebelum ke database |
| **Data Transformation** | Mengubah data dari DTO ke Model dan sebaliknya |

### Alur Data (Flow)

```
Request → Handler → Service → Repository → Database
                ↓
Response ← Handler ← Service ← Repository ← Database
```

**Contoh Alur Create User:**

1. **Handler** menerima HTTP POST request dengan JSON body
2. **Handler** parsing JSON ke `CreateUserRequest` DTO
3. **Handler** memanggil `service.CreateUser(dto)`
4. **Service** validasi bisnis (email sudah ada?)
5. **Service** hash password
6. **Service** convert DTO ke Model
7. **Service** memanggil `repository.Create(model)`
8. **Repository** insert ke database menggunakan GORM
9. **Repository** return hasil ke Service
10. **Service** convert Model ke Response DTO
11. **Handler** return JSON response

---

## 📁 Struktur Project

```
fiber-gorm-app/
│
├── 📂 cmd/                          # Entry point aplikasi
│   └── main.go                      # File utama untuk menjalankan server
│
├── 📂 internal/                     # Kode internal aplikasi
│   │
│   ├── 📂 config/                   # Konfigurasi
│   │   └── database.go              # Setup koneksi PostgreSQL
│   │
│   ├── 📂 dto/                      # Data Transfer Objects
│   │   ├── common_dto.go            # DTO umum (Response, Pagination)
│   │   ├── user_dto.go              # DTO untuk User
│   │   ├── 📂 request/              # Request DTOs
│   │   │   └── user_request.go
│   │   └── 📂 response/             # Response DTOs
│   │       └── user_response.go
│   │
│   ├── 📂 handlers/                 # HTTP Handlers (Controllers)
│   │   ├── handler.go               # Aggregator semua handlers
│   │   └── user_handler.go          # Handler untuk User endpoints
│   │
│   ├── 📂 models/                   # Database Models (Entities)
│   │   └── user.go                  # Model User (representasi tabel)
│   │
│   ├── 📂 repository/               # Data Access Layer
│   │   ├── repository.go            # Aggregator semua repositories
│   │   └── user_repository.go       # Repository untuk User
│   │
│   ├── 📂 routes/                   # Route Definitions
│   │   └── routes.go                # Semua endpoint API
│   │
│   └── 📂 services/                 # Business Logic Layer
│       ├── service.go               # Aggregator semua services
│       └── user_service.go          # Service untuk User
│
├── 📄 .air.toml                     # Konfigurasi hot reload (Air)
├── 📄 .dockerignore                 # File yang diabaikan Docker
├── 📄 .env                          # Environment variables
├── 📄 .env.example                  # Contoh environment variables
├── 📄 .gitignore                    # File yang diabaikan Git
│
├── 🐳 docker-compose.dev.yml        # Docker Compose untuk Development
├── 🐳 docker-compose.prod.yml       # Docker Compose untuk Production
├── 🐳 Dockerfile.dev                # Dockerfile untuk Development
├── 🐳 Dockerfile.prod               # Dockerfile untuk Production
│
├── 📄 go.mod                        # Go module dependencies
├── 📄 Makefile                      # Make commands
├── 📄 run.sh                        # Shell script shortcuts
└── 📄 README.md                     # Dokumentasi ini
```

---

## 🐳 Penjelasan Docker Setup

### Kenapa Docker Dipisah Development dan Production?

| Aspek | Development | Production | Alasan |
|-------|-------------|------------|--------|
| **User** | `root` | `non-root (appuser)` | Production lebih aman dengan user terbatas |
| **Hot Reload** | ✅ Ya (Air) | ❌ Tidak | Development butuh reload otomatis saat coding |
| **Volume Mount** | ✅ Ya | ❌ Tidak | Development perlu sync kode lokal |
| **DB Port Exposed** | ✅ 5432 | ❌ Internal | Production database tidak boleh diakses dari luar |
| **Build** | Single stage | Multi-stage | Production image lebih kecil |
| **Security** | Minimal | Enhanced | Production butuh keamanan maksimal |
| **Debug Tools** | ✅ Ada | ❌ Tidak | Development butuh tools debugging |

### Dockerfile.dev (Development)

```dockerfile
# Base image Go versi 1.21 dengan Alpine Linux (ringan)
FROM golang:1.21-alpine

# Install tools yang dibutuhkan
RUN apk add --no-cache gcc musl-dev git

# Install Air untuk hot reload
RUN go install github.com/cosmtrek/air@v1.49.0

# Set working directory
WORKDIR /app

# Copy dan download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy semua source code
COPY . .

# Expose port 3000
EXPOSE 3000

# Jalankan dengan Air (hot reload)
CMD ["air", "-c", ".air.toml"]
```

**Kenapa Development pakai Root User?**
- Lebih mudah untuk debugging
- Tidak ada masalah permission saat install package
- Volume mount berfungsi dengan baik
- Tidak perlu khawatir keamanan karena hanya di lokal

### Dockerfile.prod (Production)

```dockerfile
# ============ STAGE 1: Build ============
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache gcc musl-dev git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build binary dengan optimasi
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# ============ STAGE 2: Runtime ============
FROM alpine:3.19

# Install ca-certificates untuk HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Buat user non-root untuk keamanan
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /app

# Copy binary dari stage builder
COPY --from=builder /app/main .

# Set ownership ke user non-root
RUN chown -R appuser:appgroup /app

# Gunakan user non-root
USER appuser

EXPOSE 3000

CMD ["./main"]
```

**Kenapa Production pakai Non-Root User?**
- **Security Best Practice** - Jika ada vulnerability, attacker tidak punya akses root
- **Principle of Least Privilege** - Aplikasi hanya punya permission yang dibutuhkan
- **Container Isolation** - Mengurangi risiko container escape
- **Compliance** - Memenuhi standar keamanan (ISO, SOC2, dll)

### Docker Compose Development vs Production

#### docker-compose.dev.yml
```yaml
services:
  app:
    build:
      dockerfile: Dockerfile.dev      # Pakai Dockerfile development
    volumes:
      - .:/app                        # Mount kode lokal untuk hot reload
    ports:
      - "3000:3000"
    depends_on:
      - postgres

  postgres:
    image: postgres:16-alpine
    ports:
      - "5432:5432"                   # Port exposed untuk tools database
    volumes:
      - postgres_dev_data:/var/lib/postgresql/data

  pgadmin:
    image: dpage/pgadmin4:latest
    ports:
      - "5050:80"
```

#### docker-compose.prod.yml
```yaml
services:
  app:
    build:
      dockerfile: Dockerfile.prod     # Pakai Dockerfile production
    read_only: true                   # Filesystem read-only (keamanan)
    security_opt:
      - no-new-privileges:true        # Tidak bisa escalate privilege
    ports:
      - "3000:3000"

  postgres:
    image: postgres:16-alpine
    # Port TIDAK di-expose ke luar (keamanan)
    volumes:
      - postgres_prod_data:/var/lib/postgresql/data
```

---

## 🚀 Cara Setup & Menjalankan

### Prerequisites

Pastikan sudah terinstall:
- Docker Desktop (versi 24.x atau lebih baru)
- Git

### Langkah-langkah Setup

#### 1. Clone Repository
```bash
cd /Users/rasya2121/Documents/code/revoks-gym/revoks-gym-backend
cd fiber-gorm-app
```

#### 2. Setup Environment Variables
```bash
# Copy file contoh environment
cp .env.example .env

# Edit jika diperlukan (opsional)
nano .env
```

#### 3. Beri Permission Script
```bash
chmod +x run.sh
```

#### 4. Jalankan Development Environment
```bash
# Menggunakan shell script
./run.sh dev:up

# ATAU menggunakan Makefile
make dev-up

# ATAU menggunakan docker compose langsung
docker compose -f docker-compose.dev.yml up -d --build
```

#### 5. Verifikasi Aplikasi Berjalan
```bash
# Cek status container
./run.sh status

# Cek health endpoint
curl http://localhost:3000/health

# Lihat logs
./run.sh dev:logs
```

### 📍 Access Points

| Service | URL | Keterangan |
|---------|-----|------------|
| API | http://localhost:3000 | REST API endpoint |
| pgAdmin | http://localhost:5050 | Database management GUI |

### Login pgAdmin

| Field | Value |
|-------|-------|
| Email | admin@admin.com |
| Password | admin |

### Koneksi pgAdmin ke PostgreSQL

1. Buka http://localhost:5050
2. Login dengan kredensial di atas
3. Klik kanan **Servers** → **Create** → **Server**
4. Tab **General**: Name = `fiber-db`
5. Tab **Connection**:
   - Host: `postgres`
   - Port: `5432`
   - Username: `postgres`
   - Password: `postgres`

---

## 📡 API Endpoints

### Base URL
```
http://localhost:3000
```

### Daftar Endpoints

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| GET | `/health` | Health check server |
| POST | `/api/v1/users` | Buat user baru |
| GET | `/api/v1/users` | Ambil semua users (dengan pagination) |
| GET | `/api/v1/users/:id` | Ambil user berdasarkan ID |
| PUT | `/api/v1/users/:id` | Update user |
| DELETE | `/api/v1/users/:id` | Hapus user |
| PUT | `/api/v1/users/:id/password` | Ganti password user |

### Query Parameters

| Parameter | Type | Default | Deskripsi |
|-----------|------|---------|-----------|
| `page` | int | 1 | Nomor halaman |
| `per_page` | int | 10 | Jumlah item per halaman |

---

## 📝 Contoh Penggunaan API

### Health Check
```bash
curl http://localhost:3000/health
```

Response:
```json
{
  "status": "success",
  "message": "Server is running"
}
```

### Create User
```bash
curl -X POST http://localhost:3000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "password123",
    "phone": "08123456789",
    "address": "Jakarta, Indonesia"
  }'
```

Response:
```json
{
  "status": "success",
  "message": "User created successfully",
  "data": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "phone": "08123456789",
    "address": "Jakarta, Indonesia",
    "created_at": "2026-02-04T10:00:00Z",
    "updated_at": "2026-02-04T10:00:00Z"
  }
}
```

### Get All Users (dengan Pagination)
```bash
curl "http://localhost:3000/api/v1/users?page=1&per_page=10"
```

Response:
```json
{
  "status": "success",
  "message": "Users retrieved successfully",
  "data": [...],
  "meta": {
    "current_page": 1,
    "per_page": 10,
    "total": 25,
    "total_pages": 3
  }
}
```

### Get User by ID
```bash
curl http://localhost:3000/api/v1/users/1
```

### Update User
```bash
curl -X PUT http://localhost:3000/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Updated",
    "phone": "08987654321"
  }'
```

### Delete User
```bash
curl -X DELETE http://localhost:3000/api/v1/users/1
```

### Change Password
```bash
curl -X PUT http://localhost:3000/api/v1/users/1/password \
  -H "Content-Type: application/json" \
  -d '{
    "old_password": "password123",
    "new_password": "newpassword456"
  }'
```

---

## 🔧 Perintah Shortcut

### Shell Script (run.sh)

| Perintah | Fungsi |
|----------|--------|
| `./run.sh dev:up` | Start development dengan hot reload |
| `./run.sh dev:down` | Stop development |
| `./run.sh dev:logs` | Lihat semua logs development |
| `./run.sh dev:logs:app` | Lihat logs app saja |
| `./run.sh dev:restart` | Restart development |
| `./run.sh dev:rebuild` | Rebuild dan restart development |
| `./run.sh prod:up` | Start production |
| `./run.sh prod:down` | Stop production |
| `./run.sh status` | Lihat status container |
| `./run.sh clean:all` | Hapus semua resources Docker |
| `./run.sh help` | Tampilkan bantuan |

### Makefile

| Perintah | Fungsi |
|----------|--------|
| `make dev-up` | Start development |
| `make dev-down` | Stop development |
| `make dev-logs` | Lihat logs development |
| `make prod-up` | Start production |
| `make prod-down` | Stop production |
| `make status` | Lihat status container |
| `make clean` | Hapus semua resources Docker |

---

## 📚 Referensi

- [Fiber Documentation](https://docs.gofiber.io/)
- [GORM Documentation](https://gorm.io/docs/)
- [Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Repository Pattern](https://martinfowler.com/eaaCatalog/repository.html)

---

## 👨‍💻 Author

Project ini dibuat untuk keperluan ujian backend development.

---

*Terakhir diupdate: 4 Februari 2026*
