package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"task-management-backend/config"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan, menggunakan environment system")
	}

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Gagal koneksi database: %v", err)
	}
	defer db.Close()

	adminName := os.Getenv("ADMIN_NAME")
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminName == "" || adminEmail == "" || adminPassword == "" {
		log.Fatal("ADMIN_NAME, ADMIN_EMAIL, dan ADMIN_PASSWORD wajib diisi di .env")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existingID int64
	err = db.QueryRowContext(ctx, "SELECT id FROM users WHERE email = ?", adminEmail).Scan(&existingID)
	if err == nil {
		log.Printf("Akun Admin dengan email %s sudah ada (ID: %d), seeder dilewati.", adminEmail, existingID)
		return
	}
	if err != sql.ErrNoRows {
		log.Fatalf("Error saat mengecek user admin: %v", err)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Gagal hash password: %v", err)
	}

	query := `INSERT INTO users (name, email, password, role, is_active, created_at, updated_at) 
	          VALUES (?, ?, ?, 'ADMIN', 1, NOW(), NOW())`

	res, err := db.ExecContext(ctx, query, adminName, adminEmail, string(hashedPassword))
	if err != nil {
		log.Fatalf("Gagal membuat akun Admin: %v", err)
	}

	insertedID, _ := res.LastInsertId()
	log.Printf("Berhasil membuat akun Admin pertama! ID: %d, Email: %s", insertedID, adminEmail)
}