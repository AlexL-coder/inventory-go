package main

import (
	"awesomeProject1/cmd/api_gateway/routes"
	"awesomeProject1/services/grpc_auth/proto"
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func main() {
	// Initialize Gin
	router := gin.Default()

	// Create a context with timeout for the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Dial with credentials - using insecure for dev only
	conn, err := grpc.DialContext(
		ctx,
		"grpc_auth:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(), // Optionally block until the connection is established or times out
	)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	log.Println("Connected to gRPC server!")

	authClient := proto.NewAuthServiceClient(conn)

	// Setup Auth routes
	routes.AuthRoutes(router, authClient)

	// Start the server
	log.Println("Starting API Gateway on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run API Gateway: %v", err)
	}
}
