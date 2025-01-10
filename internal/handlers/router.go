package handlers

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		api.POST("/login", loginHandler)
		api.POST("/user", createUser)
		api.GET("/users", listUsers)
	}

}
