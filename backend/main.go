package main

import (
	"log"
	"os"

	"task-management-backend/config"
	"task-management-backend/routes"

	"github.com/joho/godotenv"
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

	router := routes.SetupRouter(db)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server berjalan di http://localhost:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}