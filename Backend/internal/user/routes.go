package user

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine, handler *Handler) {

	users := router.Group("/api/v1/users")

	{
		users.POST("", handler.Create)
	}
}