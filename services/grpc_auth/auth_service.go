package grpc_auth

import (
	"awesomeProject1/ent"
	"awesomeProject1/ent/user"
	"awesomeProject1/internal/db"
	pb "awesomeProject1/services/grpc_auth/proto"
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceServer struct {
	client *ent.Client
	pb.UnimplementedAuthServiceServer
}

func NewAuthService(client *ent.Client) *AuthServiceServer {
	return &AuthServiceServer{client: client}
}

// Login handles user login
func (s *AuthServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	u, err := s.client.User.Query().Where(user.Email(req.Email)).Only(ctx)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Pwd), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	token, err := GenerateJWT(u.ID)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{Token: token}, nil
}

// Register handles user registration
func (s *AuthServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	u, err := s.client.User.Create().
		SetName(req.Name).
		SetEmail(req.Email).
		SetPwd(string(hashedPwd)).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{Id: u.ID}, nil
}

// ListUsers temp output all users
func (s *AuthServiceServer) ListUsers(ctx context.Context, req *pb.EmptyRequest) (*pb.ListUsersResponse, error) {
	users, err := db.Client.User.Query().All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}
	var userList []*pb.User
	for _, user := range users {
		userList = append(userList, &pb.User{
			Id:    fmt.Sprintf("%d", user.ID), // Convert ID to string if needed
			Name:  user.Name,
			Email: user.Email,
		})
	}
	return &pb.ListUsersResponse{Users: userList}, nil
}
