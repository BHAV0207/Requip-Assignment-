package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	log.Println("Server started on :8080")

	router.Run(":8080")
}
