package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/s-usmonalizoda25/marketService/internal/handlers"
	"github.com/s-usmonalizoda25/marketService/internal/infrastructure/security"
	"github.com/s-usmonalizoda25/marketService/internal/repository"
	"github.com/s-usmonalizoda25/marketService/internal/service"
	"github.com/s-usmonalizoda25/marketService/pkg/cache"
	"github.com/s-usmonalizoda25/marketService/pkg/logger"
	"github.com/s-usmonalizoda25/marketService/router"
	"go.uber.org/zap"
)

func main() {
	err := godotenv.Load("config/config.env")
	if err != nil {
		log.Fatalf("failed to initialize env: %v", err)
	}

	mainLog := logger.New()
	mainLog.Info("Environment config loaded successfully")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		mainLog.Fatal("Failed to create pgx pool", zap.Error(err))
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		mainLog.Fatal("PostgreSQL is unreachable", zap.Error(err))
	}
	mainLog.Info("Successfully connected to PostgreSQL pool!")

	err = repository.RunMigration(ctx, pool)
	if err != nil {
		mainLog.Fatal("Failed to run table migration", zap.Error(err))
	}
	mainLog.Info("Migrations applied successfully!")

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		mainLog.Fatal("JWT_SECRET environment variable is required but not set")
	}

	jwtManager := security.NewJWTManager(jwtSecret, 15*time.Minute)
	hasher := security.NewBcryptHasher(10)

	userRepo := repository.NewPostgresUserRepo(pool)
	orderRepo := repository.NewPostgresOrderRepo(pool)

	userCache := cache.NewMemoryCache()

	userService := service.NewMyUserService(userRepo, hasher, jwtManager, userCache)
	orderService := service.NewMyOrderService(orderRepo)

	userHandler := handlers.NewUserHandler(userService, mainLog)
	orderHandler := handlers.NewOrderHandler(orderService, mainLog)

	appRouter := router.NewRouter(jwtManager, userHandler, orderHandler)

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		mainLog.Fatal("SERVER_PORT environment variable is required but not set")
	}

	notifyCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	server := &http.Server{
		Addr:         ":" + serverPort,
		Handler:      appRouter,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		mainLog.Info("Starting HTTP server...", zap.String("port", serverPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			mainLog.Fatal("HTTP server failed to start", zap.Error(err))
		}
	}()

	<-notifyCtx.Done()
	mainLog.Info("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		mainLog.Fatal("Server forced to shutdown:", zap.Error(err))
	}

	mainLog.Info("Server exiting")
}
