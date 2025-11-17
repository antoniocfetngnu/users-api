package grpc

import (
	"context"
	"fmt"

	"github.com/antoniocfetngnu/users-api/database"
	"github.com/antoniocfetngnu/users-api/models"
	pb "github.com/antoniocfetngnu/users-api/proto"
)

// UsersServer implements the gRPC UsersService
type UsersServer struct {
	pb.UnimplementedUsersServiceServer
}

// GetUser returns a single user by ID
func (s *UsersServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	var user models.User

	if err := database.DB.First(&user, req.UserId).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &pb.UserResponse{
		Id:        uint32(user.ID),
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetUsers returns multiple users by IDs (batch request)
func (s *UsersServer) GetUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.UsersResponse, error) {
	if len(req.UserIds) == 0 {
		return &pb.UsersResponse{Users: []*pb.UserResponse{}}, nil
	}

	var users []models.User

	// Query users with IN clause
	if err := database.DB.Where("id IN ?", req.UserIds).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	// Convert to proto response
	userResponses := make([]*pb.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = &pb.UserResponse{
			Id:        uint32(user.ID),
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return &pb.UsersResponse{Users: userResponses}, nil
}
