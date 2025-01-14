package main

import (
	"awesomeProject1/config"
	"awesomeProject1/internal/db"
	"awesomeProject1/services/grpc_auth"
	pb "awesomeProject1/services/grpc_auth/proto"
	"awesomeProject1/services/rabbitmq"
	_ "github.com/lib/pq" // PostgreSQL driver
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sync"
)

func main() {
	// Init logger
	config.InitLogging()
	// Initialize DB PostgreSQL
	db.InitDB(config.Log)
	defer db.CloseDB()
	// Initialize RabbitMQ
	maxRetries := 10
	rabbitConn, rabbitCh, err := config.InitializeRabbitMQ(maxRetries)
	if err != nil {
		log.Fatalf("Could not establish RabbitMQ connection after %d attempts: %v", maxRetries, err)
	}
	defer rabbitConn.Close()
	defer rabbitCh.Close()

	// Now, set up the gRPC server
	port := ":50051"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	authService := grpc_auth.NewAuthService(db.Client)

	// Register the AuthService with the gRPC server
	pb.RegisterAuthServiceServer(grpcServer, authService)
	reflection.Register(grpcServer)
	log.Printf("Starting gRPC Auth service on port %s...\n", port)
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// List of queues to consume from
	queueNames := []string{rabbitmq.AuthRegisterQueue, rabbitmq.EmailWelcomeQueue}
	// Create a wait group to synchronize the goroutines
	var wg sync.WaitGroup
	// Set a maximum number of concurrent consumer
	maxConsumers := 5
	sem := make(chan struct{}, maxConsumers)

	for _, queueName := range queueNames {
		sem <- struct{}{}
		wg.Add(1)

		go func(queueName string) {
			defer func() {
				// Release the slot in the semaphore once the consumer is done
				<-sem
				// Mark this goroutine as done in the wait group
				wg.Done()
			}()
			// Consume messages from the queue
			rabbitmq.ConsumeMessages(rabbitCh, queueName, &wg, authService)
		}(queueName) // Passing queueName to the goroutine
	}

	// Wait for all consumers to finish
	wg.Wait()
	log.Println("Service is running and waiting for messages...")
}
