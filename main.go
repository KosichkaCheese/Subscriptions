package main

import (
	"context"
	"log"
	"os"
	"subscriptions/database"
	_ "subscriptions/docs"
	"subscriptions/repository"
	"subscriptions/routes"
	"subscriptions/services"

	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
)

// @title Subscriptions API
// @version 1.0
// @description API for Subscriptions
// @BasePath /api
func main() {
	log.Println("Запуск приложения...")

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки переменных окружения: %v", err)
	}

	db := database.ConnectDB()

	servicerepo := repository.NewServiceRepo(db)
	subscriptionrepo := repository.NewSubscriptionRepo(db)

	serviceservice := services.NewServiceService(servicerepo)
	subscriptionservice := services.NewSubscriptionService(subscriptionrepo, servicerepo)

	//проверка репозитория
	log.Println(serviceservice.GetAll(context.Background()))
	log.Println(subscriptionservice.GetAll(context.Background()))

	router := routes.SetupRouter()
	router.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))
	err = router.Run(":" + os.Getenv("APP_PORT"))
	if err != nil {
		log.Fatalf("Ошибка запуска приложения: %v", err)
	}

	log.Println("Приложение запущено")
}
