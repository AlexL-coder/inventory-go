package routes

import (
	"awesomeProject1/services/grpc_auth/proto"
	"awesomeProject1/services/rabbitmq"
	"awesomeProject1/services/rabbitmq/handler"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"time"
)

func AuthRoutes(router *gin.Engine, authClient proto.AuthServiceClient, rabit_ch *amqp.Channel) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", func(c *gin.Context) {
			// Login handler logic
		})

		auth.GET("/login", func(c *gin.Context) {
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
			//// rabbit logic
			//// Convert request to JSON string
			//hashedPwd, errP := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
			//if errP != nil {
			//	return
			//}
			//message := fmt.Sprintf(`{"name": "%s", "email": "%s", "password": "%s"}`, req.Name, req.Email, string(hashedPwd))
			//// Publish the message to RabbitMQ
			//err := rabbitmq.PublishMessage(rabit_ch, "auth_exchange", "auth.register", message)
			//if err != nil {
			//	log.Printf("Failed to publish message: %v", err)
			//	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue registration request"})
			//	return
			//}
			//log.Printf("Message published: %s", message)
			//c.JSON(http.StatusAccepted, gin.H{"message": "User registration is finished. User was registered."})

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
				"email":   "User registration is finished. Welcome email was sent.",
			})
			// rabbit logic
			emailDetails := handler.EmailDetails{
				Recipient: req.Email,
				Subject:   "Welcome to Our Service!",
				Body:      fmt.Sprintf("Hello %s, your user ID is %s. Welcome!", req.Name, resp.Id),
			}
			// Publish the message to RabbitMQ
			err = rabbitmq.PublishMessage(rabit_ch, rabbitmq.AuthExchange, rabbitmq.EmailWelcomeRoutingKey, emailDetails)
			if err != nil {
				log.Printf("Failed to publish message: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue registration request"})
				return
			}
			log.Printf("Message published: %s", emailDetails)
			//c.JSON(http.StatusAccepted, gin.H{"message": "User registration is finished. Welcome email was sent."})

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
