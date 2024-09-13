package config

import (
	"fmt"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB
var once sync.Once

func DatabaseConnection() *gorm.DB {

	once.Do(func() {
		dbInstance, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn(),
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err != nil {
			panic("failed to connect database")
		}
		db = dbInstance
	})

	return db
}

func dsn() string {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	result := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=allow TimeZone=UTC",
		host, user, password, dbname, port)

	return result
}
