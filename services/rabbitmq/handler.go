package rabbitmq

import (
	"awesomeProject1/services/grpc_auth"
	"awesomeProject1/services/grpc_auth/proto"
	"context"
	//"github.com/gin-gonic/gin"
	"log"
	//"net/http"
	"encoding/json"
	"time"
)

// Example of a handler function that processes messages based on the routing key
func messageHandler(msgBody []byte, queueName string, authService *grpc_auth.AuthServiceServer) error {
	// Bind JSON request body to the RegisterRequest struct
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	// Deserialize the msgBody into the req struct
	if err := json.Unmarshal(msgBody, &req); err != nil {
		log.Printf("Failed to unmarshal message body: %v", err)
		return err
	}
	log.Printf("Processing message queueName: %s, %s", queueName, string(msgBody))
	// Call gRPC Register method
	grpcReq := &proto.RegisterRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := authService.Register(ctx, grpcReq)
	if err != nil {
		log.Printf("gRPC Register call failed: %v", err)
		//c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})
		return err
	}
	log.Println("Registration successful", resp.Id)

	//// Respond to the HTTP client
	//c.JSON(http.StatusOK, gin.H{
	//	"message": "Registration successful",
	//	"user_id": resp.Id,
	//})
	return nil
}
