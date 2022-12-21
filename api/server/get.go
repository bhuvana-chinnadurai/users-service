package server

import (
	"context"
	"fmt"
	"time"

	pb "github.com/bhuvana-chinnadurai/users-service/api/proto"
	"github.com/bhuvana-chinnadurai/users-service/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	defaultLimit = 100
	minPage      = 1
)

func (s *Server) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequests) (*pb.GetAllUsersResponse, error) {

	pagination := model.Pagination{
		Limit: int(req.GetPagination().GetLimit()),
		Page:  int(req.GetPagination().GetPage()),
	}
	if (pagination.Limit) > defaultLimit {
		pagination.Limit = defaultLimit
	}
	if pagination.Page < minPage {
		pagination.Page = minPage
	}
	filter := model.Filter{
		Country: req.Country,
	}

	users, err := s.UsersRepository.GetAll(filter, pagination)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("error while querying users: %s", err.Error()),
		)
	}

	var usersResponse []*pb.User
	for _, u := range users {
		usersResponse = append(usersResponse, &pb.User{
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Nickname:  u.Nickname,
			Email:     u.Email,
			Country:   u.Country,
			CreatedAt: u.CreatedAt.Format(time.RFC1123),
			UpdatedAt: u.UpdatedAt.Format(time.RFC1123),
		})
	}
	return &pb.GetAllUsersResponse{Users: usersResponse}, nil
}
