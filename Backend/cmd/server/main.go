package main

import (
	"log"

	"github.com/BHAV0207/Requip-assignment/internal/config"
	"github.com/BHAV0207/Requip-assignment/internal/database"
	"github.com/BHAV0207/Requip-assignment/internal/user"
	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.Load()

	db := database.Connect(cfg)

	// Auto Migrate
	db.AutoMigrate(&user.User{})

	router := gin.Default()

	api := router.Group("/api/v1")

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	user.RegisterRoutes(api, userHandler)

	router.GET("/health", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"status": "ok",
		})

	})

	log.Printf("Server started on :%s", cfg.AppPort)

	router.Run(":" + cfg.AppPort)
}
