package config

import (
	"fmt"
	"os"
	"task-tracker/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// 1. Load .env
	godotenv.Load()

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// ========================================================
	// TAHAP 1: Konek ke database 'postgres' (Database System)
	// Tujuannya cuma buat ngecek dan bikin database baru
	// ========================================================
	dsnRoot := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable TimeZone=Asia/Jakarta",
		host, user, password, port,
	)

	dbRoot, err := gorm.Open(postgres.Open(dsnRoot), &gorm.Config{})
	if err != nil {
		panic("Gagal konek ke Postgres System (Cek password/username di .env): " + err.Error())
	}

	// Cek apakah database sudah ada?
	var exists int
	dbRoot.Raw("SELECT 1 FROM pg_database WHERE datname = ?", dbName).Scan(&exists)

	if exists == 0 {
		// Kalau belum ada, kita buatkan!
		fmt.Println("⚠️  Database belum ada. Sedang membuat database:", dbName, "...")
		
		// Kita harus close koneksi root dulu atau pakai sql command langsung, 
		// tapi di GORM paling aman kita eksekusi command create
		// Note: CREATE DATABASE tidak bisa jalan di dalam transaction, jadi kita pakai Exec biasa
		createDbCommand := fmt.Sprintf("CREATE DATABASE %s", dbName)
		if result := dbRoot.Exec(createDbCommand); result.Error != nil {
			panic("Gagal membuat database: " + result.Error.Error())
		}
		fmt.Println("✅ Berhasil membuat database baru!")
	}

	// ========================================================
	// TAHAP 2: Konek ke database sesungguhnya (task_db)
	// ========================================================
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		host, user, password, dbName, port,
	)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Gagal koneksi ke database " + dbName + ": " + err.Error())
	}

	// Auto Migrate (Bikin Tabel)
	fmt.Println("🚀 Menjalankan Auto Migration...")
	database.AutoMigrate(&models.Task{})

	DB = database
	fmt.Println("✅ Koneksi & Migrasi Selesai!")
}