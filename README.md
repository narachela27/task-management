# Task Management Backend API (Go + Gin + MySQL)

Backend REST API untuk aplikasi manajemen tugas dengan autentikasi JWT dan Role-Based Access Control (`USER` & `ADMIN`).

## Tech Stack

- **Language:** Go
- **Framework:** Gin (`github.com/gin-gonic/gin`)
- **Database:** MySQL (`database/sql` + `go-sql-driver/mysql`)
- **Auth & Security:** JWT (`golang-jwt/jwt/v5`), Bcrypt (`golang.org/x/crypto/bcrypt`)

## Cara Menjalankan

1. Import migrasi SQL di folder `backend/migrations/` ke database MySQL (`task_management`).
2. Salin `backend/.env.example` menjadi `backend/.env` dan sesuaikan konfigurasi database serta `JWT_SECRET`.
3. Jalankan seeder admin pertama:

   ```bash
   cd backend
   go run cmd/seed/main.go
