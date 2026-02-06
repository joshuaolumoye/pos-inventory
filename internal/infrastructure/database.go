package infrastructure

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitializeDB() *gorm.DB {
	dbType := os.Getenv("DB_TYPE") // "mysql", "postgres", or "sqlite"
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	var dsn string
	var dialector gorm.Dialector

	// Default to mysql if not specified
	if dbType == "" {
		dbType = "mysql"
	}

	switch dbType {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			user, password, host, port, dbname)
		dialector = mysql.Open(dsn)
	case "postgres":
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
		dialector = postgres.Open(dsn)
	case "sqlite":
		sqliteFile := dbname
		if sqliteFile == "" {
			sqliteFile = "pos.db"
		}
		dialector = sqlite.Open(sqliteFile)
	default:
		log.Fatalf("unsupported database type: %s", dbType)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Printf("Connected to %s database successfully!", dbType)
	return db
}

func AutoMigrate(db *gorm.DB) error {
	log.Println("Starting auto-migration...")

	err := db.AutoMigrate(
		&Business{},
		&Branch{},
		&Staff{},
		&Product{},
		&Sale{},
		&SaleItem{},
		&Notification{},
	)

	if err != nil {
		log.Printf("Migration failed: %v", err)
		return err
	}

	log.Println("Auto-migration completed successfully!")
	return nil
}
