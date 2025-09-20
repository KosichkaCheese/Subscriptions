package database

import (
	"fmt"
	"os"
	"sync"
	"time"

	"subscriptions/models"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var once sync.Once

func ConnectDB(logger *zap.SugaredLogger) *gorm.DB {
	once.Do(func() {
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"))

		var err error

		for i := 0; i < 10; i++ {
			DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

			if err == nil {
				break
			}

			logger.Warn("База данных недоступна. Повторная попытка подключения...")

			time.Sleep(2 * time.Second)
		}

		if err != nil {
			logger.Fatalf("Ошибка подключения к базе данных: %v", err)
		}

		logger.Info("Подключение к базе данных установлено")

		err = DB.AutoMigrate(&models.Service{}, &models.Subscription{})
		if err != nil {
			logger.Fatalf("Ошибка миграции базы данных: %v", err)
		}

		logger.Info("Миграция базы данных выполнена")

	})

	return DB
}
