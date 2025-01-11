package auth

import (
	"awesomeProject1/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/login", handlers.loginHandler)
		api.POST("/register", handlers.createUser)
		api.GET("/users", handlers.listUsers)
	}
}
