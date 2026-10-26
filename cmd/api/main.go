package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/joho/godotenv"

	"taskflow/internal/cache"
	"taskflow/internal/config"
	"taskflow/internal/handler"
	"taskflow/internal/middleware"
	"taskflow/internal/repository"
	"taskflow/internal/service"
	"taskflow/internal/worker"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env не найден, используются переменные окружения")
	}

	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Ошибка конфигурации: %v", err)
	}
	db := repository.Connect(cfg.DatabaseURL)

	redisDB, _ := strconv.Atoi(env("REDIS_DB", "0"))
	redisCache := cache.NewRedisCache(
		env("REDIS_ADDR", "localhost:6379"),
		env("REDIS_PASSWORD", ""),
		redisDB,
	)
	defer redisCache.Close()

	pingCtx, pingCancel := context.WithTimeout(context.Background(), 3*time.Second)
	if err := redisCache.Ping(pingCtx); err != nil {
		log.Printf("Redis недоступен, кэш будет работать без него: %v", err)
	}
	pingCancel()

	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	workerCount, _ := strconv.Atoi(env("WORKER_COUNT", "4"))
	queueSize, _ := strconv.Atoi(env("WORKER_QUEUE_SIZE", "100"))
	pool := worker.NewPool(rootCtx, workerCount, queueSize)
	defer pool.Close()
	_ = pool

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handler.NewAuthHandler(authService)

	taskRepo := repository.NewTaskRepository(db)
	taskService := service.NewTaskService(taskRepo, redisCache)
	taskHandler := handler.NewTaskHandler(taskService)

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.Login)

	tasks := router.Group("/tasks")
	tasks.Use(middleware.Auth(cfg.JWTSecret))
	{
		tasks.POST("", taskHandler.Create)
		tasks.GET("", taskHandler.GetAll)
		tasks.GET("/:id", taskHandler.GetByID)
		tasks.PUT("/:id", taskHandler.Update)
		tasks.DELETE("/:id", taskHandler.Delete)
	}

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		<-rootCtx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Printf("shutdown error: %v", err)
		}
	}()

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server failed: %v", err)
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
