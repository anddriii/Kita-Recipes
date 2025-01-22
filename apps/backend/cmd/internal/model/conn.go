package model

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OpenConnection() *gorm.DB {

	// for docker
	// // Ambil environment variables
	// dbHost := os.Getenv("DB_HOST")
	// dbPort := os.Getenv("DB_PORT")
	// dbUser := os.Getenv("DB_USER")
	// dbPassword := os.Getenv("DB_PASSWORD")
	// dbName := os.Getenv("DB_NAME")

	// // Format DSN
	// dsn := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

	//for local
	dsn := "root:root@tcp(127.0.0.1:3306)/kita_recipes_test?charset=utf8mb4&parseTime=True&loc=Local"

	// Open connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Enable SQL logging
	})

	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	return db
}
