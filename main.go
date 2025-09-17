package main

import (
	"context"
	"log"
	"subscriptions/database"
	"subscriptions/repository"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Запуск приложения...")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки переменных окружения: %v", err)
	}

	db := database.ConnectDB()

	servicerepo := repository.NewServiceRepo(db)
	subscriptionrepo := repository.NewSubscriptionRepo(db)

	//проверка репозитория
	log.Println(servicerepo.GetAll(context.Background()))
	log.Println(subscriptionrepo.GetAll(context.Background()))

	log.Println("Приложение запущено")
}
