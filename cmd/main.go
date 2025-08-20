package main

import (
	"fmt"
	"log"

	"subscription-service/config"
	_ "subscription-service/docs"
	"subscription-service/internal/handlers"
	"subscription-service/internal/models"
	"subscription-service/internal/repository"
	"subscription-service/internal/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title           Subscriptions API
// @version         1.0
// @description     REST-сервис для агрегации онлайн-подписок пользователей.
// @BasePath        /

// @host      localhost:8080
// @schemes   http

func main() {
	cfg := config.LoadConfig()

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	sugar.Infow("starting server",
		"port", cfg.AppPort,
		"db_host", cfg.DBHost,
		"db_name", cfg.DBName,
	)

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&models.Subscription{}); err != nil {
		log.Fatal("failed to migrate:", err)
	}
	sugar.Info("auto-migrations done")
	repo := repository.NewSubscriptionRepository(db)
	svc := services.NewSubscriptionService(repo, sugar)
	h := handlers.NewSubscriptionHandler(svc, sugar)

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/subscriptions", h.Create)
	r.GET("/subscriptions/:id", h.Get)
	r.PUT("/subscriptions/:id", h.Update)
	r.DELETE("/subscriptions/:id", h.Delete)
	r.GET("/subscriptions", h.List)

	r.GET("/subscriptions/total", h.Total)

	r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	if err := r.Run(":" + cfg.AppPort); err != nil {
		sugar.Fatalw("server failed", "error", err)
	}
}
