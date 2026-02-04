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
revoks-gym-backend/
│
├── 📂 cmd/                          # Entry point aplikasi
│   └── main.go                      # File utama untuk menjalankan server
│
├── 📂 internal/                     # Kode internal aplikasi
│   ├── 📂 config/                   # Konfigurasi (DB + seeding)
│   │   ├── database.go
│   │   └── seed.go
│   ├── 📂 dto/                      # Data Transfer Objects
│   │   ├── attendance_dto.go
│   │   ├── auth_dto.go
│   │   ├── common_dto.go
│   │   ├── membership_dto.go
│   │   ├── notification_dto.go
│   │   ├── qr_dto.go
│   │   ├── settings_dto.go
│   │   ├── target_dto.go
│   │   ├── template_dto.go
│   │   ├── trainer_dto.go
│   │   ├── user_dto.go
│   │   ├── workout_dto.go
│   │   ├── 📂 request/
│   │   │   └── user_request.go
│   │   └── 📂 response/
│   │       └── user_response.go
│   ├── 📂 handlers/                 # HTTP Handlers (Controllers)
│   │   ├── attendance_handler.go
│   │   ├── auth_handler.go
│   │   ├── handler.go
│   │   ├── membership_handler.go
│   │   ├── me_handler.go
│   │   ├── notification_handler.go
│   │   ├── qr_handler.go
│   │   ├── setting_handler.go
│   │   ├── target_handler.go
│   │   ├── template_handler.go
│   │   ├── trainer_handler.go
│   │   ├── user_handler.go
│   │   └── workout_handler.go
│   ├── 📂 middleware/               # Auth middleware
│   │   └── auth.go
│   ├── 📂 models/                   # Database Models (Entities)
│   │   ├── attendance.go
│   │   ├── membership.go
│   │   ├── notification.go
│   │   ├── setting.go
│   │   ├── target.go
│   │   ├── template.go
│   │   ├── trainer.go
│   │   ├── user.go
│   │   └── workout_progress.go
│   ├── 📂 repository/               # Data Access Layer
│   │   ├── attendance_repository.go
│   │   ├── membership_repository.go
│   │   ├── notification_repository.go
│   │   ├── repository.go
│   │   ├── setting_repository.go
│   │   ├── target_repository.go
│   │   ├── template_repository.go
│   │   ├── trainer_repository.go
│   │   ├── user_repository.go
│   │   └── workout_repository.go
│   ├── 📂 routes/                   # Route Definitions
│   │   └── routes.go
│   ├── 📂 services/                 # Business Logic Layer
│   │   ├── attendance_service.go
│   │   ├── auth_service.go
│   │   ├── membership_service.go
│   │   ├── notification_service.go
│   │   ├── qr_service.go
│   │   ├── service.go
│   │   ├── setting_service.go
│   │   ├── target_service.go
│   │   ├── template_service.go
│   │   ├── trainer_service.go
│   │   ├── user_service.go
│   │   └── workout_service.go
│   └── 📂 utils/                    # Helpers (JWT, dll)
│       └── jwt.go
│
├── 📂 scripts/                      # Script helper & testing
│   ├── run.sh
│   ├── smoke-test.ps1
│   └── smoke-test.sh
│
├── 📂 tmp/                          # Folder sementara (Air)
│   └── main
│
├── 📂 .git/                          # Metadata git (tidak perlu diedit)
│   └── (git internal files)
│
├── 📄 .air.toml                     # Konfigurasi hot reload (Air)
├── 📄 .env                          # Environment variables
├── 📄 .env.example                  # Contoh environment variables
│
├── 🐳 docker-compose.dev.yml        # Docker Compose untuk Development
├── 🐳 docker-compose.prod.yml       # Docker Compose untuk Production
├── 🐳 Dockerfile.dev                # Dockerfile untuk Development
├── 🐳 Dockerfile.prod               # Dockerfile untuk Production
│
├── 📄 go.mod                        # Go module dependencies
├── 📄 go.sum                        # Go module checksums
├── 📄 Makefile                      # Make commands
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
git clone <repo-url>
cd revoks-gym-backend
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
chmod +x ./scripts/run.sh
chmod +x ./scripts/smoke-test.sh
```

#### 4. Jalankan Development Environment
```bash
# Menggunakan shell script
./scripts/run.sh dev:up

# ATAU menggunakan Makefile
make dev-up

# ATAU menggunakan docker compose langsung
docker compose -f docker-compose.dev.yml up -d --build
```

#### 5. Verifikasi Aplikasi Berjalan
```bash
# Cek status container
./scripts/run.sh status

# Cek health endpoint
curl http://localhost:3000/health

# Lihat logs
./scripts/run.sh dev:logs
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
| POST | `/api/v1/auth/login` | Login dan mendapatkan access token (JWT) |
| GET | `/api/v1/me` | Ambil profile user yang sedang login |
| POST | `/api/v1/users` | Buat user baru |
| GET | `/api/v1/users` | Ambil semua users (dengan pagination) |
| GET | `/api/v1/users/:id` | Ambil user berdasarkan ID |
| PUT | `/api/v1/users/:id` | Update user |
| DELETE | `/api/v1/users/:id` | Hapus user |
| PUT | `/api/v1/users/:id/password` | Ganti password user |
| GET | `/api/v1/membership` | Status membership + history |
| POST | `/api/v1/membership/renew` | Renew membership (CTA Renew) |
| GET | `/api/v1/qr/code` | Generate QR token untuk akses gate |
| POST | `/api/v1/qr/scan` | Validasi hasil scan gate (Accepted/Rejected) |
| GET | `/api/v1/attendance/history` | Attendance history (IN/OUT) |
| GET | `/api/v1/workouts/progress` | Workout progress history |
| POST | `/api/v1/workouts/progress` | Tambah workout progress |
| GET | `/api/v1/templates/followed` | Followed templates |
| GET | `/api/v1/templates/:id` | Template detail |
| POST | `/api/v1/templates/follow` | Follow template |
| GET | `/api/v1/targets` | Targets (weekly/monthly) |
| GET | `/api/v1/targets/:id/progress` | Progress history target |
| POST | `/api/v1/targets/:id/progress` | Add/update progress target |
| GET | `/api/v1/pt` | PT list |
| GET | `/api/v1/pt/:id` | PT detail |
| GET | `/api/v1/pt/:id/schedule` | PT schedule |
| GET | `/api/v1/notifications` | Notifications list (bell icon) |
| GET | `/api/v1/notifications/:id` | Notification detail |
| PUT | `/api/v1/notifications/:id/read` | Mark notification as read |
| GET | `/api/v1/settings` | Settings |
| PUT | `/api/v1/settings` | Update settings |

### Query Parameters

| Parameter | Type | Default | Deskripsi |
|-----------|------|---------|-----------|
| `page` | int | 1 | Nomor halaman |
| `per_page` | int | 10 | Jumlah item per halaman |

---

## 📝 Contoh Penggunaan API

## 🧭 Alur Mobile App → Backend (sesuai gambar)

### 1) AUTH

- **Splash / Boot**
  - Mobile cek ada token tersimpan.
  - Jika ada: call `GET /api/v1/me`.
    - `200` → **authenticated** masuk ke MAIN.
    - `401` → **not authenticated** arahkan ke Login.
- **Login**
  - `POST /api/v1/auth/login` → simpan `access_token`.
  - Setelah **login success** → `GET /api/v1/me` untuk hydrate profile.

### 2) MAIN (Bottom Tab Bar)

- **Home Tab**
  - Home Dashboard (future agregasi)
  - PT Online List (preview) → `GET /api/v1/pt`
  - Membership Status (preview) → `GET /api/v1/membership`
  - CTA Renew → `POST /api/v1/membership/renew`

- **Activity Tab**
  - Attendance History (IN/OUT) → `GET /api/v1/attendance/history`
  - Workout Progress History → `GET /api/v1/workouts/progress`
  - Add/Update Progress → `POST /api/v1/workouts/progress`

- **QR Tab (QR ACCESS)**
  - My QR Code → `GET /api/v1/qr/code` (token pendek, aman untuk ditampilkan)
  - Gate scan device → `POST /api/v1/qr/scan`
    - Response: `accepted/rejected`
    - Setelah gate response → app bisa request QR lagi (show QR again)
  - Jika membership **EXPIRED** → backend akan `rejected` dengan reason `membership_required` / `membership_expired` (Locked: Membership Required)

- **Program Tab**
  - Followed Templates → `GET /api/v1/templates/followed`
  - Template Detail → `GET /api/v1/templates/:id`
  - Targets (Weekly/Monthly) → `GET /api/v1/targets?period=weekly|monthly`
  - Add/Update Progress → `POST /api/v1/targets/:id/progress`

- **Profile Tab**
  - My Profile → `GET /api/v1/me`
  - Membership (Status + History + Renew) → `GET /api/v1/membership` + `POST /api/v1/membership/renew`
  - PT List / PT Detail / PT Schedule → `GET /api/v1/pt`, `GET /api/v1/pt/:id`, `GET /api/v1/pt/:id/schedule`
  - Settings → `GET/PUT /api/v1/settings`

### 3) NOTIFICATIONS (global)

- Bell icon dari tab manapun:
  - Notifications List → `GET /api/v1/notifications`
  - Notification Detail → `GET /api/v1/notifications/:id`
  - Mark read → `PUT /api/v1/notifications/:id/read`

### State / Flow rules yang dipakai backend

- **not authenticated**: request ke endpoint protected akan `401 Unauthorized`.
- **authenticated**: request wajib pakai header `Authorization: Bearer <token>`.
- **membership ACTIVE/EXPIRED**:
  - Membership dianggap **ACTIVE** jika record terakhir `status=active` dan `ends_at` masih di masa depan.
  - QR scan akan **Rejected** jika membership tidak aktif.

### Smoke Test Otomatis (PowerShell)

Script ini akan menjalankan flow end-to-end: health → (optional create user) → login → me → membership → (optional renew) → QR → attendance → workouts → settings → notifications → PT.

Jalankan dari root project:

```powershell
./scripts/smoke-test.ps1
```

Opsi flags (`true/false`, `1/0`, termasuk `$true/$false`):

```powershell
# Skip renew membership (tetap lanjut walau QR rejected)
./scripts/smoke-test.ps1 -RenewMembership 0

# Pakai existing user (wajib isi Email/Password)
./scripts/smoke-test.ps1 -CreateUser 0 -Email "you@example.com" -Password "password123"
```

### Smoke Test Otomatis (macOS / Linux)

Script Bash dengan `curl` + `python3` (tanpa `jq`). Pastikan `curl` dan `python3` tersedia di PATH. Jalankan dari root project:

```bash
chmod +x ./scripts/smoke-test.sh
./scripts/smoke-test.sh
```

Contoh opsi:

```bash
# Skip renew membership
./scripts/smoke-test.sh --renew-membership 0

# Pakai existing user (wajib isi email/password)
./scripts/smoke-test.sh --create-user 0 --email you@example.com --password password123
```


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

### Shell Script (scripts/run.sh)

| Perintah | Fungsi |
|----------|--------|
| `./scripts/run.sh dev:up` | Start development dengan hot reload |
| `./scripts/run.sh dev:down` | Stop development |
| `./scripts/run.sh dev:logs` | Lihat semua logs development |
| `./scripts/run.sh dev:logs:app` | Lihat logs app saja |
| `./scripts/run.sh dev:restart` | Restart development |
| `./scripts/run.sh dev:rebuild` | Rebuild dan restart development |
| `./scripts/run.sh prod:up` | Start production |
| `./scripts/run.sh prod:down` | Stop production |
| `./scripts/run.sh status` | Lihat status container |
| `./scripts/run.sh clean:all` | Hapus semua resources Docker |
| `./scripts/run.sh help` | Tampilkan bantuan |

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
