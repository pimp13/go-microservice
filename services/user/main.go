package main

import (
	"microservice_project/services/user/handler"
	"microservice_project/services/user/model"
	"microservice_project/services/user/repository"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Database connection
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=password dbname=pouya_test port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	// Auto migrate
	db.AutoMigrate(&model.User{})

	// Initialize service
	userRepo := repository.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepo)

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.GET("/health", userHandler.HealthCheck)
	// e.GET("/users", userService.getUsers)
	// e.GET("/users/:id", userService.getUser)
	e.POST("/users", userHandler.CreateUser)
	// e.PUT("/users/:id", userService.updateUser)
	// e.DELETE("/users/:id", userService.deleteUser)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	e.Logger.Printf("User Service starting on port %s", port)
	e.Start(":" + port)
}
