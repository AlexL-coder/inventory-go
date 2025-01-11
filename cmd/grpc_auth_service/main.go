package main

import (
	"awesomeProject1/internal/db"
	"awesomeProject1/services/grpc_auth"
	pb "awesomeProject1/services/grpc_auth/proto"
	_ "github.com/lib/pq" // PostgreSQL driver
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {
	db.InitDB()
	defer db.CloseDB()
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
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
