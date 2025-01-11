package routes

import (
	"awesomeProject1/services/grpc_auth/proto"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func AuthRoutes(router *gin.Engine, authClient proto.AuthServiceClient) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", func(c *gin.Context) {
			// Login handler logic
		})

		auth.POST("/register", func(c *gin.Context) {
			// Register handler logic
		})

		auth.GET("/list_users", func(c *gin.Context) {
			resp, err := authClient.ListUsers(context.Background(), &proto.EmptyRequest{})
			if err != nil {
				log.Printf("Failed to call ListUsers: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"users": resp.Users})
		})
	}
}
