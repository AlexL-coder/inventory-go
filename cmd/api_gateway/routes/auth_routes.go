package routes

import (
	"awesomeProject1/services/grpc_auth/proto"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func AuthRoutes(router *gin.Engine, authClient proto.AuthServiceClient) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", func(c *gin.Context) {
			// Login handler logic
		})

		auth.POST("/register", func(c *gin.Context) {
			// Bind JSON request body to the RegisterRequest struct
			var req struct {
				Name     string `json:"name" binding:"required"`
				Email    string `json:"email" binding:"required,email"`
				Password string `json:"password" binding:"required"`
			}

			if err := c.ShouldBindJSON(&req); err != nil {
				log.Printf("Invalid request body: %v", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
				return
			}

			// Call gRPC Register method
			grpcReq := &proto.RegisterRequest{
				Name:     req.Name,
				Email:    req.Email,
				Password: req.Password,
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			resp, err := authClient.Register(ctx, grpcReq)
			if err != nil {
				log.Printf("gRPC Register call failed: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})
				return
			}

			// Respond to the HTTP client
			c.JSON(http.StatusOK, gin.H{
				"message": "Registration successful",
				"user_id": resp.Id,
			})
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
