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

	repository := user.NewRepository(db)

	service := user.NewService(repository)

	handler := user.NewHandler(service)

	router := gin.Default()

	user.RegisterRoutes(router, handler)

	router.GET("/health", func(c *gin.Context) {

		c.JSON(200, gin.H{
			"status": "ok",
		})

	})

	log.Printf("Server started on :%s", cfg.AppPort)

	router.Run(":" + cfg.AppPort)
}
