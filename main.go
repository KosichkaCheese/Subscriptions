package main

import (
	"log"
	"os"
	"subscriptions/database"
	_ "subscriptions/docs"
	"subscriptions/handlers"
	"subscriptions/repository"
	"subscriptions/routes"
	"subscriptions/services"

	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// @title Subscriptions API
// @version 1.0
// @description API for Subscriptions
// @BasePath /api
func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Ошибка инициализации логгера: %v", err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	sugar.Info("Запуск приложения...")

	err = godotenv.Load()
	if err != nil {
		sugar.Fatalf("Ошибка загрузки переменных окружения: %v", err)
	}

	db := database.ConnectDB(sugar)

	servicerepo := repository.NewServiceRepo(db)
	subscriptionrepo := repository.NewSubscriptionRepo(db)

	serviceservice := services.NewServiceService(servicerepo, sugar)
	subscriptionservice := services.NewSubscriptionService(subscriptionrepo, servicerepo, sugar)

	servicehandler := handlers.NewServiceHandler(serviceservice)
	subscriptionhandler := handlers.NewSubscriptionHandler(subscriptionservice)

	router := routes.SetupRouter(servicehandler, subscriptionhandler)
	router.GET("/swagger/*any", swagger.WrapHandler(swaggerFiles.Handler))
	err = router.Run(":" + os.Getenv("APP_PORT"))
	if err != nil {
		sugar.Fatalf("Ошибка запуска приложения:, %v", err)
	}
}
