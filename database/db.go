package database

import (
	"fmt"
	"log"
	"os"
	"sync"

	"subscriptions/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var once sync.Once

func ConnectDB() *gorm.DB {
	once.Do(func() {
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"))

		var err error
		DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Ошибка подключения к базе данных: %v", err)
		}

		log.Println("Подключение к базе данных установлено")

		err = DB.AutoMigrate(&models.Service{}, &models.Subscription{})
		if err != nil {
			log.Fatalf("Ошибка миграции базы данных: %v", err)
		}

		log.Println("Миграция базы данных выполнена")

	})

	return DB
}
